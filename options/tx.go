package options

import (
	"errors"

	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// Constants for the Tx struct fields.
const (
	NameTxChainID            = "TxChainID"
	NameTxFeeGranterAddr     = "TxFeeGranterAddr"
	NameTxFees               = "TxFees"
	NameTxFromName           = "TxFromName"
	NameTxGas                = "TxGas"
	NameTxGasAdjustment      = "TxGasAdjustment"
	NameTxGasPrices          = "TxGasPrices"
	NameTxMemo               = "TxMemo"
	NameTxSimulateAndExecute = "TxSimulateAndExecute"
	NameTxTimeoutHeight      = "TxTimeoutHeight"
)

// Default values for the Tx fields.
const (
	DefaultTxChainID            = "sentinelhub-2"
	DefaultTxFeeGranterAddr     = ""
	DefaultTxFees               = ""
	DefaultTxFromName           = ""
	DefaultTxGas                = 200_000
	DefaultTxGasAdjustment      = 1.0 + (1.0 / 6.0)
	DefaultTxGasPrices          = "0.1udvpn"
	DefaultTxMemo               = ""
	DefaultTxSimulateAndExecute = true
	DefaultTxTimeoutHeight      = 0
)

// Flags for command-line options for Tx.
const (
	FlagTxChainID            = "tx.chain-id"
	FlagTxFeeGranterAddr     = "tx.fee-granter-addr"
	FlagTxFees               = "tx.fees"
	FlagTxFromName           = "tx.from-name"
	FlagTxGas                = "tx.gas"
	FlagTxGasAdjustment      = "tx.gas-adjustment"
	FlagTxGasPrices          = "tx.gas-prices"
	FlagTxMemo               = "tx.memo"
	FlagTxSimulateAndExecute = "tx.simulate-and-execute"
	FlagTxTimeoutHeight      = "tx.timeout-height"
)

// init function sets the default values for transaction-related parameters at package initialization.
func init() {
	SetDefault(NameTxChainID, DefaultTxChainID)
	SetDefault(NameTxFeeGranterAddr, DefaultTxFeeGranterAddr)
	SetDefault(NameTxFees, DefaultTxFees)
	SetDefault(NameTxFromName, DefaultTxFromName)
	SetDefault(NameTxGas, DefaultTxGas)
	SetDefault(NameTxGasAdjustment, DefaultTxGasAdjustment)
	SetDefault(NameTxGasPrices, DefaultTxGasPrices)
	SetDefault(NameTxMemo, DefaultTxMemo)
	SetDefault(NameTxSimulateAndExecute, DefaultTxSimulateAndExecute)
	SetDefault(NameTxTimeoutHeight, DefaultTxTimeoutHeight)
}

// Tx defines a structure for holding transaction-related parameters.
type Tx struct {
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

// NewTx creates a new Tx instance with default values.
func NewTx() *Tx {
	return &Tx{
		ChainID:            cast.ToString(GetDefault(NameTxChainID)),
		FeeGranterAddr:     cast.ToString(GetDefault(NameTxFeeGranterAddr)),
		Fees:               cast.ToString(GetDefault(NameTxFees)),
		FromName:           cast.ToString(GetDefault(NameTxFromName)),
		Gas:                cast.ToUint64(GetDefault(NameTxGas)),
		GasAdjustment:      cast.ToFloat64(GetDefault(NameTxGasAdjustment)),
		GasPrices:          cast.ToString(GetDefault(NameTxGasPrices)),
		Memo:               cast.ToString(GetDefault(NameTxMemo)),
		SimulateAndExecute: cast.ToBool(GetDefault(NameTxSimulateAndExecute)),
		TimeoutHeight:      cast.ToUint64(GetDefault(NameTxTimeoutHeight)),
	}
}

// WithChainID sets the ChainID field and returns the modified Tx instance.
func (t *Tx) WithChainID(v string) *Tx {
	t.ChainID = v
	return t
}

// WithFeeGranterAddr sets the FeeGranterAddr field and returns the modified Tx instance.
func (t *Tx) WithFeeGranterAddr(v cosmossdk.AccAddress) *Tx {
	t.FeeGranterAddr = v.String()
	return t
}

// WithFees sets the Fees field and returns the modified Tx instance.
func (t *Tx) WithFees(v cosmossdk.Coins) *Tx {
	t.Fees = v.String()
	return t
}

// WithFromName sets the FromName field and returns the modified Tx instance.
func (t *Tx) WithFromName(v string) *Tx {
	t.FromName = v
	return t
}

// WithGas sets the Gas field and returns the modified Tx instance.
func (t *Tx) WithGas(v uint64) *Tx {
	t.Gas = v
	return t
}

// WithGasAdjustment sets the GasAdjustment field and returns the modified Tx instance.
func (t *Tx) WithGasAdjustment(v float64) *Tx {
	t.GasAdjustment = v
	return t
}

// WithGasPrices sets the GasPrices field and returns the modified Tx instance.
func (t *Tx) WithGasPrices(v cosmossdk.DecCoins) *Tx {
	t.GasPrices = v.String()
	return t
}

// WithMemo sets the Memo field and returns the modified Tx instance.
func (t *Tx) WithMemo(v string) *Tx {
	t.Memo = v
	return t
}

// WithSimulateAndExecute sets the SimulateAndExecute field and returns the modified Tx instance.
func (t *Tx) WithSimulateAndExecute(v bool) *Tx {
	t.SimulateAndExecute = v
	return t
}

// WithTimeoutHeight sets the TimeoutHeight field and returns the modified Tx instance.
func (t *Tx) WithTimeoutHeight(v uint64) *Tx {
	t.TimeoutHeight = v
	return t
}

// GetChainID returns the blockchain network identifier for the transaction.
func (t *Tx) GetChainID() string {
	return t.ChainID
}

// GetFeeGranterAddr returns the fee granter address for the transaction.
func (t *Tx) GetFeeGranterAddr() cosmossdk.AccAddress {
	if t.FeeGranterAddr == "" {
		return nil
	}

	v, err := cosmossdk.AccAddressFromBech32(t.FeeGranterAddr)
	if err != nil {
		panic(err)
	}

	return v
}

// GetFees returns the transaction fees.
func (t *Tx) GetFees() cosmossdk.Coins {
	v, err := cosmossdk.ParseCoinsNormalized(t.Fees)
	if err != nil {
		panic(err)
	}

	return v
}

// GetFromName returns the sender name for the transaction.
func (t *Tx) GetFromName() string {
	return t.FromName
}

// GetGas returns the gas limit for the transaction.
func (t *Tx) GetGas() uint64 {
	return t.Gas
}

// GetGasAdjustment returns the gas adjustment factor for the transaction.
func (t *Tx) GetGasAdjustment() float64 {
	return t.GasAdjustment
}

// GetGasPrices returns the gas prices for the transaction.
func (t *Tx) GetGasPrices() cosmossdk.DecCoins {
	v, err := cosmossdk.ParseDecCoins(t.GasPrices)
	if err != nil {
		panic(err)
	}

	return v
}

// GetMemo returns the memo for the transaction.
func (t *Tx) GetMemo() string {
	return t.Memo
}

// GetSimulateAndExecute returns whether to simulate and execute the transaction.
func (t *Tx) GetSimulateAndExecute() bool {
	return t.SimulateAndExecute
}

// GetTimeoutHeight returns the timeout height for the transaction.
func (t *Tx) GetTimeoutHeight() uint64 {
	return t.TimeoutHeight
}

// ValidateTxChainID validates the ChainID field.
func ValidateTxChainID(v string) error {
	if v == "" {
		return errors.New("chain_id must not be empty")
	}

	return nil
}

// ValidateTxFeeGranterAddr validates the FeeGranterAddr field.
func ValidateTxFeeGranterAddr(v string) error {
	if v == "" {
		return nil
	}
	if _, err := cosmossdk.AccAddressFromBech32(v); err != nil {
		return errors.New("fee_granter_addr must be a valid address")
	}

	return nil
}

// ValidateTxFees validates the Fees field.
func ValidateTxFees(v string) error {
	if _, err := cosmossdk.ParseCoinsNormalized(v); err != nil {
		return errors.New("fees must be a valid coins format")
	}

	return nil
}

// ValidateTxFromName validates the FromName field.
func ValidateTxFromName(v string) error {
	if v == "" {
		return errors.New("from_name must not be empty")
	}

	return nil
}

// ValidateTxGas validates the Gas field.
func ValidateTxGas(v uint64) error {
	if v == 0 {
		return errors.New("gas must be greater than zero")
	}

	return nil
}

// ValidateTxGasAdjustment validates the GasAdjustment field.
func ValidateTxGasAdjustment(v float64) error {
	if v <= 0 {
		return errors.New("gas_adjustment must be greater than zero")
	}

	return nil
}

// ValidateTxGasPrices validates the GasPrices field.
func ValidateTxGasPrices(v string) error {
	if _, err := cosmossdk.ParseDecCoins(v); err != nil {
		return errors.New("gas_prices must be a valid decimal coins format")
	}

	return nil
}

// Validate validates all the fields of the Tx struct.
func (t *Tx) Validate() error {
	if err := ValidateTxChainID(t.ChainID); err != nil {
		return err
	}
	if err := ValidateTxFeeGranterAddr(t.FeeGranterAddr); err != nil {
		return err
	}
	if err := ValidateTxFees(t.Fees); err != nil {
		return err
	}
	if err := ValidateTxFromName(t.FromName); err != nil {
		return err
	}
	if err := ValidateTxGas(t.Gas); err != nil {
		return err
	}
	if err := ValidateTxGasAdjustment(t.GasAdjustment); err != nil {
		return err
	}
	if err := ValidateTxGasPrices(t.GasPrices); err != nil {
		return err
	}

	return nil
}

// GetTxChainIDFromCmd retrieves the blockchain network identifier from the command-line flags.
func GetTxChainIDFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagTxChainID)
}

// GetTxFeeGranterAddrFromCmd retrieves the fee granter address from the command-line flags.
func GetTxFeeGranterAddrFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagTxFeeGranterAddr)
}

// GetTxFeesFromCmd retrieves the transaction fees from the command-line flags.
func GetTxFeesFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagTxFees)
}

// GetTxFromNameFromCmd retrieves the sender name from the command-line flags.
func GetTxFromNameFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagTxFromName)
}

// GetTxGasFromCmd retrieves the gas limit from the command-line flags.
func GetTxGasFromCmd(cmd *cobra.Command) (uint64, error) {
	return cmd.Flags().GetUint64(FlagTxGas)
}

// GetTxGasAdjustmentFromCmd retrieves the gas adjustment factor from the command-line flags.
func GetTxGasAdjustmentFromCmd(cmd *cobra.Command) (float64, error) {
	return cmd.Flags().GetFloat64(FlagTxGasAdjustment)
}

// GetTxGasPricesFromCmd retrieves the gas prices from the command-line flags.
func GetTxGasPricesFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagTxGasPrices)
}

// GetTxMemoFromCmd retrieves the memo from the command-line flags.
func GetTxMemoFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagTxMemo)
}

// GetTxSimulateAndExecuteFromCmd retrieves the simulate and execute flag from the command-line flags.
func GetTxSimulateAndExecuteFromCmd(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool(FlagTxSimulateAndExecute)
}

// GetTxTimeoutHeightFromCmd retrieves the timeout height from the command-line flags.
func GetTxTimeoutHeightFromCmd(cmd *cobra.Command) (uint64, error) {
	return cmd.Flags().GetUint64(FlagTxTimeoutHeight)
}

// SetFlagTxChainID sets the flag for the blockchain network identifier in the given command.
func SetFlagTxChainID(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameTxChainID))
	cmd.Flags().String(FlagTxChainID, value, "Blockchain network identifier.")
}

// SetFlagTxFeeGranterAddr sets the flag for the fee granter address in the given command.
func SetFlagTxFeeGranterAddr(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameTxFeeGranterAddr))
	cmd.Flags().String(FlagTxFeeGranterAddr, value, "Address of the entity granting fees for the transaction.")
}

// SetFlagTxFees sets the flag for the transaction fees in the given command.
func SetFlagTxFees(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameTxFees))
	cmd.Flags().String(FlagTxFees, value, "Transaction fees to be paid.")
}

// SetFlagTxFromName sets the flag for the sender name in the given command.
func SetFlagTxFromName(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameTxFromName))
	cmd.Flags().String(FlagTxFromName, value, "Name of the sender's account in the keyring.")
}

// SetFlagTxGas sets the flag for the gas limit in the given command.
func SetFlagTxGas(cmd *cobra.Command) {
	value := cast.ToUint64(GetDefault(NameTxGas))
	cmd.Flags().Uint64(FlagTxGas, value, "Gas limit set for the transaction.")
}

// SetFlagTxGasAdjustment sets the flag for the gas adjustment factor in the given command.
func SetFlagTxGasAdjustment(cmd *cobra.Command) {
	value := cast.ToFloat64(GetDefault(NameTxGasAdjustment))
	cmd.Flags().Float64(FlagTxGasAdjustment, value, "Gas adjustment factor for the transaction.")
}

// SetFlagTxGasPrices sets the flag for the gas prices in the given command.
func SetFlagTxGasPrices(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameTxGasPrices))
	cmd.Flags().String(FlagTxGasPrices, value, "Gas prices for the transaction execution.")
}

// SetFlagTxMemo sets the flag for the memo in the given command.
func SetFlagTxMemo(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameTxMemo))
	cmd.Flags().String(FlagTxMemo, value, "Memo text attached to the transaction.")
}

// SetFlagTxSimulateAndExecute sets the flag for simulate and execute in the given command.
func SetFlagTxSimulateAndExecute(cmd *cobra.Command) {
	value := cast.ToBool(GetDefault(NameTxSimulateAndExecute))
	cmd.Flags().Bool(FlagTxSimulateAndExecute, value, "Indicates whether to simulate and execute the transaction.")
}

// SetFlagTxTimeoutHeight sets the flag for the timeout height in the given command.
func SetFlagTxTimeoutHeight(cmd *cobra.Command) {
	value := cast.ToUint64(GetDefault(NameTxTimeoutHeight))
	cmd.Flags().Uint64(FlagTxTimeoutHeight, value, "Block height at which the transaction times out.")
}

// SetTxFlags sets all transaction-related flags for the command.
func SetTxFlags(cmd *cobra.Command) {
	SetFlagTxChainID(cmd)
	SetFlagTxFeeGranterAddr(cmd)
	SetFlagTxFees(cmd)
	SetFlagTxFromName(cmd)
	SetFlagTxGas(cmd)
	SetFlagTxGasAdjustment(cmd)
	SetFlagTxGasPrices(cmd)
	SetFlagTxMemo(cmd)
	SetFlagTxSimulateAndExecute(cmd)
	SetFlagTxTimeoutHeight(cmd)
}

// NewTxFromCmd creates a new Tx object from the command-line flags.
func NewTxFromCmd(cmd *cobra.Command) (*Tx, error) {
	// Retrieve chain ID from the command flags
	chainID, err := GetTxChainIDFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve fee granter address from the command flags
	feeGranterAddr, err := GetTxFeeGranterAddrFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve fees from the command flags
	fees, err := GetTxFeesFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve sender name from the command flags
	fromName, err := GetTxFromNameFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve gas limit from the command flags
	gas, err := GetTxGasFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve gas adjustment factor from the command flags
	gasAdjustment, err := GetTxGasAdjustmentFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve gas prices from the command flags
	gasPrices, err := GetTxGasPricesFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve memo from the command flags
	memo, err := GetTxMemoFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve simulate and execute flag from the command flags
	simulateAndExecute, err := GetTxSimulateAndExecuteFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve timeout height from the command flags
	timeoutHeight, err := GetTxTimeoutHeightFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Create a new Tx object with the retrieved values
	return &Tx{
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
