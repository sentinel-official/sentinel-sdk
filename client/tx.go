package client

import (
	"context"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	txsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Simulate simulates the execution of a transaction before broadcasting it.
// Takes transaction bytes as input and returns the simulation response or an error.
func (c *Client) Simulate(ctx context.Context, buf []byte) (*txtypes.SimulateResponse, error) {
	var (
		resp   txtypes.SimulateResponse
		method = "/cosmos.tx.v1beta1.Service/Simulate"
		req    = &txtypes.SimulateRequest{TxBytes: buf}
	)

	// Perform a gRPC query to simulate the transaction.
	if err := c.QueryGRPC(ctx, method, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// simulateTx calculates the gas usage of a transaction.
// Returns the gas used and any error encountered.
func (c *Client) simulateTx(ctx context.Context, txb client.TxBuilder) (uint64, error) {
	// Encode the transaction into bytes.
	buf, err := c.TxEncoder()(txb.GetTx())
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
func (c *Client) broadcastTxSync(ctx context.Context, txb client.TxBuilder) (*coretypes.ResultBroadcastTx, error) {
	// Encode the transaction into bytes.
	buf, err := c.TxEncoder()(txb.GetTx())
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
func (c *Client) signTx(txb client.TxBuilder, key *keyring.Record, account authtypes.AccountI) error {
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
	buf, err := c.SignModeHandler().GetSignBytes(singleSignatureData.SignMode, signerData, txb.GetTx())
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
func (c *Client) prepareTx(ctx context.Context, key *keyring.Record, account authtypes.AccountI, msgs []sdk.Msg) (client.TxBuilder, error) {
	// Create a new transaction builder.
	txb := c.NewTxBuilder()
	if err := txb.SetMsgs(msgs...); err != nil {
		return nil, err
	}

	// Set transaction parameters.
	txb.SetFeeAmount(c.txFees)
	txb.SetFeeGranter(c.txFeeGranterAddr)
	txb.SetGasLimit(c.txGas)
	txb.SetMemo(c.txMemo)
	txb.SetTimeoutHeight(c.txTimeoutHeight)

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
	}

	return txb, nil
}

// BroadcastTx broadcasts a signed transaction and returns the broadcast result or an error.
func (c *Client) BroadcastTx(ctx context.Context, msgs []sdk.Msg) (*coretypes.ResultBroadcastTx, error) {
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
func (c *Client) Tx(ctx context.Context, hash []byte) (*coretypes.ResultTx, error) {
	// Get the HTTP client for querying the blockchain.
	http, err := c.HTTP()
	if err != nil {
		return nil, err
	}

	// Perform the query using the transaction hash.
	return http.Tx(ctx, hash, c.queryProve)
}
