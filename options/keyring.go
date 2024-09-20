package options

import (
	"errors"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// Constants for the keyring struct fields.
const (
	NameKeyringBackend = "KeyringBackend"
	NameKeyringHomeDir = "KeyringHomeDir"
	NameKeyringName    = "KeyringName"
)

// Default values for the keyring fields.
const (
	DefaultKeyringBackend = "test"
	DefaultKeyringHomeDir = ""
	DefaultKeyringName    = "sentinel"
)

// Flags for command-line options for keyring.
const (
	FlagKeyringBackend = "keyring.backend"
	FlagKeyringHomeDir = "keyring.home-dir"
	FlagKeyringName    = "keyring.name"
)

// init function sets the default values for the keyring-related parameters at package initialization.
func init() {
	SetDefault(NameKeyringBackend, DefaultKeyringBackend)
	SetDefault(NameKeyringHomeDir, DefaultKeyringHomeDir)
	SetDefault(NameKeyringName, DefaultKeyringName)
}

// Keyring defines a structure for holding keyring-related parameters.
type Keyring struct {
	Backend string    `json:"backend" toml:"backend"`   // Backend specifies the backend type for the keyring.
	HomeDir string    `json:"home_dir" toml:"home_dir"` // HomeDir specifies the home directory for storing keyring data.
	Input   io.Reader `json:"-" toml:"-"`               // Input reader, typically used for user interaction
	Name    string    `json:"name" toml:"name"`         // Name represents the name of the keyring.
}

// NewKeyring creates a new Keyring instance with default values.
func NewKeyring() *Keyring {
	return &Keyring{
		Backend: cast.ToString(GetDefault(NameKeyringBackend)),
		HomeDir: cast.ToString(GetDefault(NameKeyringHomeDir)),
		Input:   os.Stdin,
		Name:    cast.ToString(GetDefault(NameKeyringName)),
	}
}

// WithBackend sets the Backend field and returns the updated Keyring instance.
func (k *Keyring) WithBackend(v string) *Keyring {
	k.Backend = v
	return k
}

// WithHomeDir sets the HomeDir field and returns the updated Keyring instance.
func (k *Keyring) WithHomeDir(v string) *Keyring {
	k.HomeDir = v
	return k
}

// WithInput sets the Input field and returns the updated Keyring instance.
func (k *Keyring) WithInput(v io.Reader) *Keyring {
	k.Input = v
	return k
}

// WithName sets the Name field and returns the updated Keyring instance.
func (k *Keyring) WithName(v string) *Keyring {
	k.Name = v
	return k
}

// GetBackend returns the backend type from the Keyring.
func (k *Keyring) GetBackend() string {
	return k.Backend
}

// GetHomeDir returns the home directory from the Keyring.
func (k *Keyring) GetHomeDir() string {
	return k.HomeDir
}

// GetInput returns the input from the Keyring.
func (k *Keyring) GetInput() io.Reader {
	return k.Input
}

// GetName returns the name from the Keyring.
func (k *Keyring) GetName() string {
	return k.Name
}

// Keystore creates and returns a new keyring based on the provided options.
func (k *Keyring) Keystore(cdc codec.Codec) (keyring.Keyring, error) {
	return keyring.New(k.GetName(), k.GetBackend(), k.GetHomeDir(), k.GetInput(), cdc)
}

// ValidateKeyringBackend validates the backend type for the keyring.
func ValidateKeyringBackend(v string) error {
	allowedBackends := map[string]bool{
		"file":    true,
		"kwallet": true,
		"memory":  true,
		"os":      true,
		"pass":    true,
		"test":    true,
	}

	if v == "" {
		return errors.New("backend must be non-empty")
	}
	if _, ok := allowedBackends[v]; !ok {
		return errors.New("backend must be one of: file, kwallet, memory, os, pass, test")
	}

	return nil
}

// ValidateKeyringHomeDir validates the home directory for the keyring.
func ValidateKeyringHomeDir(v string) error {
	if v == "" {
		return errors.New("home directory must be non-empty")
	}

	return nil
}

// ValidateKeyringName validates the name for the keyring.
func ValidateKeyringName(v string) error {
	if v == "" {
		return errors.New("name must be non-empty")
	}

	return nil
}

// Validate validates all fields of the Keyring struct.
func (k *Keyring) Validate() error {
	if err := ValidateKeyringBackend(k.Backend); err != nil {
		return err
	}
	if err := ValidateKeyringHomeDir(k.HomeDir); err != nil {
		return err
	}
	if err := ValidateKeyringName(k.Name); err != nil {
		return err
	}

	return nil
}

// GetKeyringBackendFromCmd retrieves the backend type from the command-line flags.
func GetKeyringBackendFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagKeyringBackend)
}

// GetKeyringHomeDirFromCmd retrieves the home directory from the command-line flags.
func GetKeyringHomeDirFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagKeyringHomeDir)
}

// GetKeyringNameFromCmd retrieves the name from the command-line flags.
func GetKeyringNameFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagKeyringName)
}

// SetFlagKeyringBackend sets the flag for the backend type in the given command.
func SetFlagKeyringBackend(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameKeyringBackend))
	cmd.Flags().String(FlagKeyringBackend, value, "Keyring backend to use.")
}

// SetFlagKeyringHomeDir sets the flag for the home directory in the given command.
func SetFlagKeyringHomeDir(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameKeyringHomeDir))
	cmd.Flags().String(FlagKeyringHomeDir, value, "Home directory for keyring to store the keys.")
}

// SetFlagKeyringName sets the flag for the name in the given command.
func SetFlagKeyringName(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameKeyringName))
	cmd.Flags().String(FlagKeyringName, value, "Name for keyring.")
}

// SetKeyringFlags sets all keyring-related flags for the command.
func SetKeyringFlags(cmd *cobra.Command) {
	SetFlagKeyringBackend(cmd)
	SetFlagKeyringHomeDir(cmd)
	SetFlagKeyringName(cmd)
}

// NewKeyringFromCmd creates a new Keyring object from the command-line flags.
func NewKeyringFromCmd(cmd *cobra.Command) (*Keyring, error) {
	// Retrieve backend type from the command flags
	backend, err := GetKeyringBackendFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve home directory from the command flags
	homeDir, err := GetKeyringHomeDirFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve name from the command flags
	name, err := GetKeyringNameFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Return a new Keyring object with the retrieved values
	return &Keyring{
		Backend: backend,
		HomeDir: homeDir,
		Input:   cmd.InOrStdin(),
		Name:    name,
	}, nil
}
