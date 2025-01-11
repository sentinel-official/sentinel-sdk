package log

import (
	"io"
	"time"

	"cosmossdk.io/log"
	"github.com/cometbft/cometbft/config"
	"github.com/rs/zerolog"
)

var logger = log.NewNopLogger() // Default to a no-op logger

// Info logs an informational message.
func Info(msg string, keyVals ...any) {
	logger.Info(msg, keyVals...)
}

// Warn logs a warning message.
func Warn(msg string, keyVals ...any) {
	logger.Warn(msg, keyVals...)
}

// Error logs an error message.
func Error(msg string, keyVals ...any) {
	logger.Error(msg, keyVals...)
}

// Debug logs a debug message.
func Debug(msg string, keyVals ...any) {
	logger.Debug(msg, keyVals...)
}

// With returns a logger with additional context.
func With(keyVals ...any) log.Logger {
	return logger.With(keyVals...)
}

// Impl returns the underlying implementation of the logger.
func Impl() any {
	return logger.Impl()
}

// SetLogger sets the global logger instance.
func SetLogger(l log.Logger) {
	if l != nil {
		logger = l
	}
}

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
		log.TimeFormatOption(time.RFC3339),
	}

	// Set log format based on the provided format string
	if format == config.LogFormatJSON {
		opts = append(opts, log.OutputJSONOption())
	}

	// Create and return the logger with the specified options
	return log.NewLogger(w, opts...), nil
}
