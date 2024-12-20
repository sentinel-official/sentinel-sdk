package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// writeOutputJSON formats the output as JSON and writes it to the provided writer.
func writeOutputJSON(w io.Writer, v interface{}) error {
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s", buf)
	return err
}

// writeOutputText formats the output as YAML and writes it to the provided writer.
func writeOutputText(w io.Writer, v interface{}) error {
	buf, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s", buf)
	return err
}

// writeOutput formats the output according to the specified format and writes it to the provided writer.
func writeOutput(w io.Writer, v interface{}, format string) error {
	switch format {
	case keys.OutputFormatJSON:
		return writeOutputJSON(w, v)
	case keys.OutputFormatText:
		return writeOutputText(w, v)
	default:
		return fmt.Errorf("invalid output format %s", format)
	}
}

// writeOutputToCmd writes the formatted output to the command's output and adds a newline.
func writeOutputToCmd(cmd *cobra.Command, v interface{}, format string) error {
	if err := writeOutput(cmd.OutOrStderr(), v, format); err != nil {
		return err
	}

	cmd.Println() // Adding a newline after output
	return nil
}

// writeMnemonicWarningToCmd prints a formatted warning message to save the mnemonic securely.
func writeMnemonicWarningToCmd(cmd *cobra.Command) {
	cmd.Printf("\n")
	cmd.Printf("####################################################################\n")
	cmd.Printf("WARNING: YOU MUST SAVE THE FOLLOWING MNEMONIC SECURELY!\n")
	cmd.Printf("THIS MNEMONIC IS REQUIRED TO RECOVER YOUR KEY.\n")
	cmd.Printf("IF YOU LOSE THIS MNEMONIC, YOU WILL NOT BE ABLE TO RECOVER YOUR KEY.\n")
	cmd.Printf("####################################################################\n")
	cmd.Printf("\n")
}
