package client

import (
	"context"

	core "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// gRPC methods for simulating the transaction
	methodSimulate = "/cosmos.tx.v1beta1.Service/Simulate"
)

// calculateFees computes transaction fees based on gas prices and gas limit.
func calculateFees(gasPrices cosmossdk.DecCoins, gasLimit uint64) cosmossdk.Coins {
	fees := make(cosmossdk.Coins, len(gasPrices))
	for i, price := range gasPrices {
		fee := price.Amount.MulInt64(int64(gasLimit))
		fees[i] = cosmossdk.NewCoin(price.Denom, fee.Ceil().RoundInt())
	}

	return fees
}

// Simulate simulates the execution of a transaction before broadcasting it.
// Takes transaction bytes as input and returns the simulation response or an error.
func (c *Client) Simulate(ctx context.Context, buf []byte) (*tx.SimulateResponse, error) {
	var (
		resp tx.SimulateResponse
		req  = &tx.SimulateRequest{TxBytes: buf}
	)

	// Perform a gRPC query to simulate the transaction.
	if err := c.QueryGRPC(ctx, methodSimulate, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// simulateTx calculates the gas usage of a transaction.
// Returns the gas used and any error encountered.
func (c *Client) simulateTx(ctx context.Context, txb client.TxBuilder) (uint64, error) {
	// Encode the transaction into bytes.
	buf, err := c.txConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return 0, err
	}

	// Simulate the transaction execution.
	res, err := c.Simulate(ctx, buf)
	if err != nil {
		return 0, err
	}

	// Calculate the adjusted gas usage.
	return uint64(c.txGasAdjustment * float64(res.GasInfo.GasUsed)), nil
}

// broadcastTxSync broadcasts a transaction synchronously.
// Returns the broadcast result or an error if the operation fails.
func (c *Client) broadcastTxSync(ctx context.Context, txb client.TxBuilder) (*core.ResultBroadcastTx, error) {
	// Encode the transaction into bytes.
	buf, err := c.txConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, err
	}

	// Get the HTTP client for broadcasting.
	http, err := c.HTTP()
	if err != nil {
		return nil, err
	}

	// Broadcast the transaction synchronously.
	return http.BroadcastTxSync(ctx, buf)
}

// signTx signs a transaction using the provided key and account information.
// Returns an error if the signing process fails.
func (c *Client) signTx(txb client.TxBuilder, key *keyring.Record, account auth.AccountI) error {
	// Prepare the single signature data.
	singleSignatureData := txsigning.SingleSignatureData{
		SignMode:  txsigning.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}

	// Retrieve the public key from the key record.
	pubKey, err := key.GetPubKey()
	if err != nil {
		return err
	}

	// Create the signature information.
	signature := txsigning.SignatureV2{
		PubKey:   pubKey,
		Data:     &singleSignatureData,
		Sequence: account.GetSequence(),
	}

	// Set the initial signature in the transaction builder.
	if err := txb.SetSignatures(signature); err != nil {
		return err
	}

	// Prepare the signer data for signing the transaction.
	signerData := authsigning.SignerData{
		ChainID:       c.chainID,
		AccountNumber: account.GetAccountNumber(),
		Sequence:      account.GetSequence(),
	}

	// Get the bytes to be signed.
	buf, err := c.txConfig.SignModeHandler().GetSignBytes(singleSignatureData.SignMode, signerData, txb.GetTx())
	if err != nil {
		return err
	}

	// Sign the transaction bytes.
	buf, _, err = c.Sign(c.txFromName, buf)
	if err != nil {
		return err
	}

	// Update the signature data with the signed bytes.
	singleSignatureData.Signature = buf
	signature.Data = &singleSignatureData

	// Set the updated signature in the transaction builder.
	if err := txb.SetSignatures(signature); err != nil {
		return err
	}

	return nil
}

// prepareTx prepares a transaction for broadcasting by setting fees, gas, and other parameters.
// Returns the transaction builder and any error encountered.
func (c *Client) prepareTx(ctx context.Context, key *keyring.Record, account auth.AccountI, msgs []cosmossdk.Msg) (client.TxBuilder, error) {
	// Create a new transaction builder.
	txb := c.txConfig.NewTxBuilder()
	if err := txb.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	// Set transaction parameters.
	txb.SetFeeAmount(c.txFees)
	txb.SetFeeGranter(c.txFeeGranterAddr)
	txb.SetGasLimit(c.txGas)
	txb.SetMemo(c.txMemo)
	txb.SetTimeoutHeight(c.txTimeoutHeight)

	if !c.txGasPrices.IsZero() {
		fees := calculateFees(c.txGasPrices, c.txGas)
		txb.SetFeeAmount(fees)
	}

	// Retrieve the public key from the key record.
	pubKey, err := key.GetPubKey()
	if err != nil {
		return nil, err
	}

	// Set an initial signature in the transaction builder.
	signature := txsigning.SignatureV2{
		PubKey: pubKey,
		Data: &txsigning.SingleSignatureData{
			SignMode: txsigning.SignMode_SIGN_MODE_DIRECT,
		},
		Sequence: account.GetSequence(),
	}

	if err := txb.SetSignatures(signature); err != nil {
		return nil, err
	}

	// Simulate the transaction to calculate gas usage if required.
	if c.txSimulateAndExecute {
		gasLimit, err := c.simulateTx(ctx, txb)
		if err != nil {
			return nil, err
		}

		txb.SetGasLimit(gasLimit)
		if !c.txGasPrices.IsZero() {
			fees := calculateFees(c.txGasPrices, c.txGas)
			txb.SetFeeAmount(fees)
		}
	}

	return txb, nil
}

// BroadcastTx broadcasts a signed transaction and returns the broadcast result or an error.
func (c *Client) BroadcastTx(ctx context.Context, msgs []cosmossdk.Msg) (*core.ResultBroadcastTx, error) {
	// Retrieve the signing key.
	key, err := c.Key(c.txFromName)
	if err != nil {
		return nil, err
	}

	// Get the sender's address from the key.
	accAddr, err := key.GetAddress()
	if err != nil {
		return nil, err
	}

	// Retrieve the sender's account information.
	account, err := c.Account(ctx, accAddr)
	if err != nil {
		return nil, err
	}

	// Prepare the transaction for broadcasting.
	txb, err := c.prepareTx(ctx, key, account, msgs)
	if err != nil {
		return nil, err
	}

	// Sign the transaction.
	if err := c.signTx(txb, key, account); err != nil {
		return nil, err
	}

	// Broadcast the signed transaction synchronously.
	return c.broadcastTxSync(ctx, txb)
}

// Tx retrieves a transaction from the blockchain using its hash.
// Returns the transaction result or an error.
func (c *Client) Tx(ctx context.Context, hash []byte) (*core.ResultTx, error) {
	// Get the HTTP client for querying the blockchain.
	http, err := c.HTTP()
	if err != nil {
		return nil, err
	}

	// Perform the query using the transaction hash.
	return http.Tx(ctx, hash, c.queryProve)
}
