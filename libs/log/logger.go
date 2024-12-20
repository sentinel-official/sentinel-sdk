package log

import (
	"io"

	"cosmossdk.io/log"
	"github.com/cometbft/cometbft/config"
	"github.com/rs/zerolog"
)

// NewLogger creates a new logger instance with the specified output writer, format, and log level.
func NewLogger(w io.Writer, format, level string) (log.Logger, error) {
	// Parse the log level from the string
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	// Prepare options for logger
	opts := []log.Option{
		log.LevelOption(logLevel),
	}

	// Set log format based on the provided format string
	if format == config.LogFormatJSON {
		opts = append(opts, log.OutputJSONOption())
	}

	// Create and return the logger with the specified options
	return log.NewLogger(w, opts...), nil
}
