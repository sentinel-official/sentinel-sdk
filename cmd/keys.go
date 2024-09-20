package cmd

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/sentinel-go-sdk/client"
	"github.com/sentinel-official/sentinel-go-sdk/client/input"
	"github.com/sentinel-official/sentinel-go-sdk/options"
)

// KeysCmd returns a new Cobra command for key management sub-commands.
func KeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Sub-commands for managing keys.",
	}

	cmd.AddCommand(
		keysAdd(),
		keysDelete(),
		keysList(),
		keysShow(),
	)

	return cmd
}

// keysAdd creates a new key with the specified name, mnemonic, and bip39 passphrase.
func keysAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add a new key with the specified name and optional mnemonic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := client.NewOptions()
			if _, err := opts.WithKeyFromCmd(cmd); err != nil {
				return err
			}
			if _, err := opts.WithKeyringFromCmd(cmd); err != nil {
				return err
			}

			outputFormat, err := options.GetOutputFormatFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(opts.Input)

			// Prompt for mnemonic
			mnemonic, err := input.GetString("Enter your bip39 mnemonic, or hit enter to generate one.\n", reader)
			if err != nil {
				return err
			}

			if mnemonic != "" && !bip39.IsMnemonicValid(mnemonic) {
				return fmt.Errorf("invalid mnemonic")
			}

			// Prompt for bip39 passphrase
			bip39Pass, err := input.GetPassword("Enter your bip39 passphrase, or hit enter to use the default:", reader)
			if err != nil {
				return err
			}

			// Confirm passphrase if provided
			if bip39Pass != "" {
				confirmPass, err := input.GetPassword("Confirm bip39 passphrase:", reader)
				if err != nil {
					return err
				}

				if bip39Pass != confirmPass {
					return errors.New("bip39 passphrase does not match")
				}
			}

			// Initialize the Client
			c := client.NewDefault()

			// Check if the key already exists
			if _, err := c.Key(args[0], opts); err == nil {
				return fmt.Errorf("key with name '%s' already exists", args[0])
			}

			// Create the key
			newMnemonic, key, err := c.CreateKey(args[0], mnemonic, bip39Pass, opts)
			if err != nil {
				return err
			}

			output, err := keyring.MkAccKeyOutput(key)
			if err != nil {
				return err
			}

			if newMnemonic != mnemonic {
				writeMnemonicWarningToCmd(cmd)
				output.Mnemonic = newMnemonic
			}

			// Output the key information
			if err := writeOutputToCmd(cmd, output, outputFormat); err != nil {
				return err
			}

			cmd.Println("Key created successfully.")
			return nil
		},
	}

	options.SetKeyFlags(cmd)
	options.SetKeyringFlags(cmd)
	options.SetFlagOutputFormat(cmd)

	return cmd
}

// keysDelete removes the key with the specified name.
func keysDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete the key with the specified name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := client.NewOptions()
			if _, err := opts.WithKeyringFromCmd(cmd); err != nil {
				return err
			}

			reader := bufio.NewReader(opts.Input)

			confirm, err := input.GetConfirmation("Are you sure you want to delete this key? [y/N]:", reader)
			if err != nil {
				return err
			}
			if !confirm {
				return errors.New("deletion aborted")
			}

			// Initialize the Client
			c := client.NewDefault()

			// Delete the key
			if err := c.DeleteKey(args[0], opts); err != nil {
				return err
			}

			cmd.Println("Key deleted successfully.")
			return nil
		},
	}

	options.SetKeyringFlags(cmd)

	return cmd
}

// keysList lists all the available keys.
func keysList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := client.NewOptions()
			if _, err := opts.WithKeyringFromCmd(cmd); err != nil {
				return err
			}

			outputFormat, err := options.GetOutputFormatFromCmd(cmd)
			if err != nil {
				return err
			}

			// Initialize the Client
			c := client.NewDefault()

			// Fetch the list of keys
			keys, err := c.Keys(opts)
			if err != nil {
				return err
			}

			output, err := keyring.MkAccKeysOutput(keys)
			if err != nil {
				return err
			}

			// Output the key list
			if err := writeOutputToCmd(cmd, output, outputFormat); err != nil {
				return err
			}

			return nil
		},
	}

	options.SetKeyringFlags(cmd)
	options.SetFlagOutputFormat(cmd)

	return cmd
}

// keysShow displays details of the key with the specified name.
func keysShow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Show details of the key with the specified name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := client.NewOptions()
			if _, err := opts.WithKeyringFromCmd(cmd); err != nil {
				return err
			}

			outputFormat, err := options.GetOutputFormatFromCmd(cmd)
			if err != nil {
				return err
			}

			// Initialize the Client
			c := client.NewDefault()

			// Fetch the key details
			key, err := c.Key(args[0], opts)
			if err != nil {
				return err
			}

			output, err := keyring.MkAccKeyOutput(key)
			if err != nil {
				return err
			}

			// Output the key details
			if err := writeOutputToCmd(cmd, output, outputFormat); err != nil {
				return err
			}

			return nil
		},
	}

	options.SetKeyringFlags(cmd)
	options.SetFlagOutputFormat(cmd)

	return cmd
}
