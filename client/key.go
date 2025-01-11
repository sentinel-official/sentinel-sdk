package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/go-bip39"
)

// Key retrieves key information from the keyring based on the provided name.
// Returns the key record or an error if the key cannot be found.
func (c *Client) Key(name string) (*keyring.Record, error) {
	key, err := c.keyring.Key(name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key: %w", err)
	}

	return key, nil
}

// Sign signs the provided data using the key from the keyring identified by the given name.
// Returns the signed bytes, the public key, and any error encountered.
func (c *Client) Sign(name string, buf []byte) ([]byte, types.PubKey, error) {
	sig, pubKey, err := c.keyring.Sign(name, buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign data: %w", err)
	}

	return sig, pubKey, nil
}

// Keys retrieves a list of all keys from the keyring.
// Returns the list of key records or an error if the operation fails.
func (c *Client) Keys() ([]*keyring.Record, error) {
	keys, err := c.keyring.List()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve keys: %w", err)
	}

	return keys, nil
}

// DeleteKey removes a key from the keyring based on the provided name.
// Returns an error if the key cannot be deleted.
func (c *Client) DeleteKey(name string) error {
	if err := c.keyring.Delete(name); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// NewMnemonic generates a new mnemonic phrase using bip39 with 256 bits of entropy.
// Returns the mnemonic or an error if the operation fails.
func (c *Client) NewMnemonic() (string, error) {
	// Generate new entropy for the mnemonic.
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	// Create a new mnemonic phrase from the entropy.
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// CreateKey generates and stores a new key in the keyring with the provided name, mnemonic, and options.
// If no mnemonic is provided, it generates a new one.
// Returns the mnemonic, the created key record, and any error encountered.
func (c *Client) CreateKey(name, mnemonic, bip39Pass string, coinType, account, index uint32) (s string, k *keyring.Record, err error) {
	// Generate a new mnemonic if none is provided.
	if mnemonic == "" {
		mnemonic, err = c.NewMnemonic()
		if err != nil {
			return "", nil, fmt.Errorf("failed to generate new mnemonic: %w", err)
		}
	}

	// Create an HD path for the key.
	hdPath := hd.CreateHDPath(coinType, account, index)
	signAlgo := hd.Secp256k1

	// Create a new key in the keyring.
	key, err := c.keyring.NewAccount(name, mnemonic, bip39Pass, hdPath.String(), signAlgo)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create new account: %w", err)
	}

	return mnemonic, key, nil
}
