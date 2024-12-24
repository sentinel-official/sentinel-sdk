package client

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/go-bip39"
)

// Key retrieves key information from the keyring based on the provided name.
// Returns the key record or an error if the key cannot be found.
func (c *Client) Key(name string) (*keyring.Record, error) {
	return c.keyring.Key(name)
}

// Sign signs the provided data using the key from the keyring identified by the given name.
// Returns the signed bytes, the public key, and any error encountered.
func (c *Client) Sign(name string, buf []byte) ([]byte, types.PubKey, error) {
	return c.keyring.Sign(name, buf)
}

// Keys retrieves a list of all keys from the keyring.
// Returns the list of key records or an error if the operation fails.
func (c *Client) Keys() ([]*keyring.Record, error) {
	return c.keyring.List()
}

// DeleteKey removes a key from the keyring based on the provided name.
// Returns an error if the key cannot be deleted.
func (c *Client) DeleteKey(name string) error {
	return c.keyring.Delete(name)
}

// NewMnemonic generates a new mnemonic phrase using bip39 with 256 bits of entropy.
// Returns the mnemonic or an error if the operation fails.
func (c *Client) NewMnemonic() (string, error) {
	// Generate new entropy for the mnemonic.
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	// Create a new mnemonic phrase from the entropy.
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
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
			return "", nil, err
		}
	}

	// Create an HD path for the key.
	hdPath := hd.CreateHDPath(coinType, account, index)
	signAlgo := hd.Secp256k1

	// Create a new key in the keyring.
	key, err := c.keyring.NewAccount(name, mnemonic, bip39Pass, hdPath.String(), signAlgo)
	if err != nil {
		return "", nil, err
	}

	return mnemonic, key, nil
}
