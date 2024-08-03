package options

import (
	"github.com/spf13/cobra"

	"github.com/sentinel-official/sentinel-go-sdk/cmd/flags"
)

// TxOptions represents options for transactions.
type TxOptions struct {
	ChainID            string  `json:"chain_id" toml:"chain_id"`                         // ChainID is the identifier of the blockchain network.
	FeeGranterAddr     string  `json:"fee_granter_addr" toml:"fee_granter_addr"`         // FeeGranterAddr is the address of the entity granting fees.
	Fees               string  `json:"fees" toml:"fees"`                                 // Fees is the transaction fees.
	FromName           string  `json:"from_name" toml:"from_name"`                       // FromName is the name of the sender.
	Gas                uint64  `json:"gas" toml:"gas"`                                   // Gas is the gas limit for the transaction.
	GasAdjustment      float64 `json:"gas_adjustment" toml:"gas_adjustment"`             // GasAdjustment is the adjustment factor for gas estimation.
	GasPrices          string  `json:"gas_prices" toml:"gas_prices"`                     // GasPrices is the gas prices for transaction execution.
	Memo               string  `json:"memo" toml:"memo"`                                 // Memo is a memo attached to the transaction.
	SimulateAndExecute bool    `json:"simulate_and_execute" toml:"simulate_and_execute"` // SimulateAndExecute indicates whether to simulate and execute the transaction.
	TimeoutHeight      uint64  `json:"timeout_height" toml:"timeout_height"`             // TimeoutHeight is the block height at which the transaction times out.
}

// NewDefaultTx creates a new TxOptions instance with default values.
func NewDefaultTx() *TxOptions {
	return &TxOptions{
		ChainID:            flags.DefaultTxChainID,
		Gas:                flags.DefaultTxGas,
		GasAdjustment:      flags.DefaultTxGasAdjustment,
		GasPrices:          flags.DefaultTxGasPrices,
		SimulateAndExecute: flags.DefaultTxSimulateAndExecute,
	}
}

// WithChainID sets the ChainID field and returns the modified TxOptions instance.
func (t *TxOptions) WithChainID(v string) *TxOptions {
	t.ChainID = v
	return t
}

// WithFeeGranterAddr sets the FeeGranterAddr field and returns the modified TxOptions instance.
func (t *TxOptions) WithFeeGranterAddr(v string) *TxOptions {
	t.FeeGranterAddr = v
	return t
}

// WithFees sets the Fees field and returns the modified TxOptions instance.
func (t *TxOptions) WithFees(v string) *TxOptions {
	t.Fees = v
	return t
}

// WithFromName sets the FromName field and returns the modified TxOptions instance.
func (t *TxOptions) WithFromName(v string) *TxOptions {
	t.FromName = v
	return t
}

// WithGas sets the Gas field and returns the modified TxOptions instance.
func (t *TxOptions) WithGas(v uint64) *TxOptions {
	t.Gas = v
	return t
}

// WithGasAdjustment sets the GasAdjustment field and returns the modified TxOptions instance.
func (t *TxOptions) WithGasAdjustment(v float64) *TxOptions {
	t.GasAdjustment = v
	return t
}

// WithGasPrices sets the GasPrices field and returns the modified TxOptions instance.
func (t *TxOptions) WithGasPrices(v string) *TxOptions {
	t.GasPrices = v
	return t
}

// WithMemo sets the Memo field and returns the modified TxOptions instance.
func (t *TxOptions) WithMemo(v string) *TxOptions {
	t.Memo = v
	return t
}

// WithSimulateAndExecute sets the SimulateAndExecute field and returns the modified TxOptions instance.
func (t *TxOptions) WithSimulateAndExecute(v bool) *TxOptions {
	t.SimulateAndExecute = v
	return t
}

// WithTimeoutHeight sets the TimeoutHeight field and returns the modified TxOptions instance.
func (t *TxOptions) WithTimeoutHeight(v uint64) *TxOptions {
	t.TimeoutHeight = v
	return t
}

// AddTxFlagsToCmd configures all transaction-related flags for the given command.
func AddTxFlagsToCmd(cmd *cobra.Command) {
	flags.SetFlagTxChainID(cmd)
	flags.SetFlagTxFeeGranterAddr(cmd)
	flags.SetFlagTxFees(cmd)
	flags.SetFlagTxFromName(cmd)
	flags.SetFlagTxGas(cmd)
	flags.SetFlagTxGasAdjustment(cmd)
	flags.SetFlagTxGasPrices(cmd)
	flags.SetFlagTxMemo(cmd)
	flags.SetFlagTxSimulateAndExecute(cmd)
	flags.SetFlagTxTimeoutHeight(cmd)
}

// NewTxOptionsFromCmd creates and returns TxOptions from the given cobra command's flags.
func NewTxOptionsFromCmd(cmd *cobra.Command) (*TxOptions, error) {
	// Retrieve the chain ID flag value from the command.
	chainID, err := flags.GetTxChainIDFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the fee granter address flag value from the command.
	feeGranterAddr, err := flags.GetTxFeeGranterAddrFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the fees flag value from the command.
	fees, err := flags.GetTxFeesFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the from name flag value from the command.
	fromName, err := flags.GetTxFromNameFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the gas flag value from the command.
	gas, err := flags.GetTxGasFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the gas adjustment flag value from the command.
	gasAdjustment, err := flags.GetTxGasAdjustmentFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the gas prices flag value from the command.
	gasPrices, err := flags.GetTxGasPricesFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the memo flag value from the command.
	memo, err := flags.GetTxMemoFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the simulate and execute flag value from the command.
	simulateAndExecute, err := flags.GetTxSimulateAndExecuteFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve the timeout height flag value from the command.
	timeoutHeight, err := flags.GetTxTimeoutHeightFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Return a new TxOptions instance populated with the retrieved flag values.
	return &TxOptions{
		ChainID:            chainID,
		FeeGranterAddr:     feeGranterAddr,
		Fees:               fees,
		FromName:           fromName,
		Gas:                gas,
		GasAdjustment:      gasAdjustment,
		GasPrices:          gasPrices,
		Memo:               memo,
		SimulateAndExecute: simulateAndExecute,
		TimeoutHeight:      timeoutHeight,
	}, nil
}
