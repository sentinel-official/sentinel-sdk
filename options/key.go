package options

import (
	cryptohd "github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// Constants for the key struct fields.
const (
	NameKeyAccount  = "KeyAccount"
	NameKeyCoinType = "KeyCoinType"
	NameKeyIndex    = "KeyIndex"
)

// Default values for key fields.
const (
	DefaultKeyAccount  = 0
	DefaultKeyCoinType = 118
	DefaultKeyIndex    = 0
)

// Flags for command-line options for key creation.
const (
	FlagKeyAccount  = "key.account"
	FlagKeyCoinType = "key.coin-type"
	FlagKeyIndex    = "key.index"
)

// init function sets the default values for the key-related parameters at package initialization.
func init() {
	SetDefault(NameKeyAccount, DefaultKeyAccount)
	SetDefault(NameKeyCoinType, DefaultKeyCoinType)
	SetDefault(NameKeyIndex, DefaultKeyIndex)
}

// Key defines a structure for holding key-related parameters.
type Key struct {
	Account  uint32 `json:"account" toml:"account"`     // Account represents the account number in the key.
	CoinType uint32 `json:"coin_type" toml:"coin_type"` // CoinType represents the type of coin used.
	Index    uint32 `json:"index" toml:"index"`         // Index represents the specific key index.
}

// NewKey creates a new Key instance with default values.
func NewKey() *Key {
	return &Key{
		Account:  cast.ToUint32(GetDefault(NameKeyAccount)),
		CoinType: cast.ToUint32(GetDefault(NameKeyCoinType)),
		Index:    cast.ToUint32(GetDefault(NameKeyIndex)),
	}
}

// WithAccount sets the Account field and returns the updated Key instance.
func (k *Key) WithAccount(v uint32) *Key {
	k.Account = v
	return k
}

// WithCoinType sets the CoinType field and returns the updated Key instance.
func (k *Key) WithCoinType(v uint32) *Key {
	k.CoinType = v
	return k
}

// WithIndex sets the Index field and returns the updated Key instance.
func (k *Key) WithIndex(v uint32) *Key {
	k.Index = v
	return k
}

// GetAccount returns the account number from the Key.
func (k *Key) GetAccount() uint32 {
	return k.Account
}

// GetCoinType returns the coin type from the Key.
func (k *Key) GetCoinType() uint32 {
	return k.CoinType
}

// GetIndex returns the index number from the Key.
func (k *Key) GetIndex() uint32 {
	return k.Index
}

// ValidateKeyAccount validates the account number.
func ValidateKeyAccount(_ uint32) error {
	return nil
}

// ValidateKeyCoinType validates the coin type.
func ValidateKeyCoinType(_ uint32) error {
	return nil
}

// ValidateKeyIndex validates the index.
func ValidateKeyIndex(_ uint32) error {
	return nil
}

// Validate validates all fields of the Key struct.
func (k *Key) Validate() error {
	if err := ValidateKeyAccount(k.Account); err != nil {
		return err
	}
	if err := ValidateKeyCoinType(k.CoinType); err != nil {
		return err
	}
	if err := ValidateKeyIndex(k.Index); err != nil {
		return err
	}

	return nil
}

// HDPath returns the hierarchical deterministic (HD) path for the key.
func (k *Key) HDPath() string {
	v := cryptohd.CreateHDPath(k.GetCoinType(), k.GetAccount(), k.GetIndex())
	return v.String()
}

// SignatureAlgo returns the signature algorithm to be used (secp256k1).
func (k *Key) SignatureAlgo() keyring.SignatureAlgo {
	return cryptohd.Secp256k1
}

// GetKeyAccountFromCmd retrieves the account number from the command-line flags.
func GetKeyAccountFromCmd(cmd *cobra.Command) (uint32, error) {
	return cmd.Flags().GetUint32(FlagKeyAccount)
}

// GetKeyCoinTypeFromCmd retrieves the coin type from the command-line flags.
func GetKeyCoinTypeFromCmd(cmd *cobra.Command) (uint32, error) {
	return cmd.Flags().GetUint32(FlagKeyCoinType)
}

// GetKeyIndexFromCmd retrieves the index from the command-line flags.
func GetKeyIndexFromCmd(cmd *cobra.Command) (uint32, error) {
	return cmd.Flags().GetUint32(FlagKeyIndex)
}

// SetFlagKeyAccount sets the flag for the account number in the given command.
func SetFlagKeyAccount(cmd *cobra.Command) {
	value := cast.ToUint32(GetDefault(NameKeyAccount))
	cmd.Flags().Uint32(FlagKeyAccount, value, "Account number for key creation.")
}

// SetFlagKeyCoinType sets the flag for the coin type in the given command.
func SetFlagKeyCoinType(cmd *cobra.Command) {
	value := cast.ToUint32(GetDefault(NameKeyCoinType))
	cmd.Flags().Uint32(FlagKeyCoinType, value, "Coin type for key creation.")
}

// SetFlagKeyIndex sets the flag for the index in the given command.
func SetFlagKeyIndex(cmd *cobra.Command) {
	value := cast.ToUint32(GetDefault(NameKeyIndex))
	cmd.Flags().Uint32(FlagKeyIndex, value, "Index for key creation.")
}

// SetKeyFlags sets all key-related flags for the command.
func SetKeyFlags(cmd *cobra.Command) {
	SetFlagKeyAccount(cmd)
	SetFlagKeyCoinType(cmd)
	SetFlagKeyIndex(cmd)
}

// NewKeyFromCmd creates a new Key object from the command-line flags.
func NewKeyFromCmd(cmd *cobra.Command) (*Key, error) {
	// Retrieve account from the command flags
	account, err := GetKeyAccountFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve coin type from the command flags
	coinType, err := GetKeyCoinTypeFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve index from the command flags
	index, err := GetKeyIndexFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Return a new Key object with the retrieved values
	return &Key{
		Account:  account,
		CoinType: coinType,
		Index:    index,
	}, nil
}
