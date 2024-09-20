package options

import (
	"errors"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// Constants for the log fields.
const (
	NameLogFormat = "LogFormat"
	NameLogLevel  = "LogLevel"
)

// Default values for the log fields.
const (
	DefaultLogFormat = "text"
	DefaultLogLevel  = "info"
)

// Flags for command-line options for log configuration.
const (
	FlagLogFormat = "log.format"
	FlagLogLevel  = "log.level"
)

// init function sets the default values for the log-related parameters at package initialization.
func init() {
	SetDefault(NameLogFormat, DefaultLogFormat)
	SetDefault(NameLogLevel, DefaultLogLevel)
}

// Log defines a structure for holding log-related parameters.
type Log struct {
	Format string `json:"format" toml:"format"` // Format defines the format of the logs (e.g., text, json).
	Level  string `json:"level" toml:"level"`   // Level defines the log level (e.g., debug, info, warn, error).
}

// NewLog creates a new Log instance with default values.
func NewLog() *Log {
	return &Log{
		Format: cast.ToString(GetDefault(NameLogFormat)),
		Level:  cast.ToString(GetDefault(NameLogLevel)),
	}
}

// WithFormat sets the log format and returns the updated Log instance.
func (l *Log) WithFormat(v string) *Log {
	l.Format = v
	return l
}

// WithLevel sets the log level and returns the updated Log instance.
func (l *Log) WithLevel(v string) *Log {
	l.Level = v
	return l
}

// GetFormat returns the log format from the Log.
func (l *Log) GetFormat() string {
	return l.Format
}

// GetLevel returns the log level from the Log.
func (l *Log) GetLevel() string {
	return l.Level
}

// ValidateLogFormat validates the log format.
func ValidateLogFormat(v string) error {
	allowedFormats := map[string]bool{
		"json":  true,
		"plain": true,
	}

	if v == "" {
		return errors.New("format must be non-empty")
	}
	if _, ok := allowedFormats[v]; !ok {
		return errors.New("format must be one of: json, plain")
	}

	return nil
}

// ValidateLogLevel validates the log level.
func ValidateLogLevel(v string) error {
	validLevels := map[string]bool{
		"debug": true,
		"error": true,
		"fatal": true,
		"info":  true,
		"panic": true,
		"warn":  true,
	}

	if v == "" {
		return errors.New("level must be non-empty")
	}
	if _, ok := validLevels[v]; !ok {
		return errors.New("level must be one of: debug, error, fatal, info, panic, warn")
	}

	return nil
}

// Validate checks all fields of the Log struct for validity.
func (l *Log) Validate() error {
	if err := ValidateLogFormat(l.Format); err != nil {
		return err
	}
	if err := ValidateLogLevel(l.Level); err != nil {
		return err
	}

	return nil
}

// GetLogFormatFromCmd retrieves the log format from the command-line flags.
func GetLogFormatFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagLogFormat)
}

// GetLogLevelFromCmd retrieves the log level from the command-line flags.
func GetLogLevelFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagLogLevel)
}

// SetFlagLogFormat sets the flag for the log format in the given command.
func SetFlagLogFormat(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameLogFormat))
	cmd.Flags().String(FlagLogFormat, value, "Specify the log format (json or plain).")
}

// SetFlagLogLevel sets the flag for the log level in the given command.
func SetFlagLogLevel(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameLogLevel))
	cmd.Flags().String(FlagLogLevel, value, "Specify the log level (debug, info, warn, error, fatal, panic).")
}

// SetLogFlags sets all log-related flags for the command.
func SetLogFlags(cmd *cobra.Command) {
	SetFlagLogFormat(cmd)
	SetFlagLogLevel(cmd)
}

// NewLogFromCmd creates a new Log object from the command-line flags.
func NewLogFromCmd(cmd *cobra.Command) (*Log, error) {
	// Retrieve log format from the command flags
	format, err := GetLogFormatFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve log level from the command flags
	level, err := GetLogLevelFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Return a new Log object with the retrieved values
	return &Log{
		Format: format,
		Level:  level,
	}, nil
}
