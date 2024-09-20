package options

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/spf13/cobra"
)

// SetFlagOutputFormat sets the flag for the output format in the given command.
func SetFlagOutputFormat(cmd *cobra.Command) {
	cmd.Flags().String("output-format", keys.OutputFormatText, "Specify the output format (json or text).")
}

// GetOutputFormatFromCmd retrieves the output format from the command-line flags.
func GetOutputFormatFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString("output-format")
}
