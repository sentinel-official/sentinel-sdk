package options

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/cometbft/cometbft/rpc/client"
	"github.com/cometbft/cometbft/rpc/client/http"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/sentinel-go-sdk/utils"
)

// Constants for the Query struct fields.
const (
	NameQueryHeight     = "QueryHeight"
	NameQueryMaxRetries = "QueryMaxRetries"
	NameQueryProve      = "QueryProve"
	NameQueryRetryDelay = "QueryRetryDelay"
	NameQueryRPCAddr    = "QueryRPCAddr"
	NameQueryTimeout    = "QueryTimeout"
)

// Default values for the Query fields.
const (
	DefaultQueryHeight     = 0
	DefaultQueryMaxRetries = 5
	DefaultQueryProve      = false
	DefaultQueryRetryDelay = "1s"
	DefaultQueryRPCAddr    = "https://rpc.sentinel.co:443"
	DefaultQueryTimeout    = "10s"
)

// Flags for command-line options for Query.
const (
	FlagQueryHeight     = "query.height"
	FlagQueryMaxRetries = "query.max-retries"
	FlagQueryProve      = "query.prove"
	FlagQueryRetryDelay = "query.retry-delay"
	FlagQueryRPCAddr    = "query.rpc-addr"
	FlagQueryTimeout    = "query.timeout"
)

// init function sets the default values for query-related parameters at package initialization.
func init() {
	SetDefault(NameQueryHeight, DefaultQueryHeight)
	SetDefault(NameQueryMaxRetries, DefaultQueryMaxRetries)
	SetDefault(NameQueryProve, DefaultQueryProve)
	SetDefault(NameQueryRetryDelay, DefaultQueryRetryDelay)
	SetDefault(NameQueryRPCAddr, DefaultQueryRPCAddr)
	SetDefault(NameQueryTimeout, DefaultQueryTimeout)
}

// Query defines a structure for holding query-related parameters.
type Query struct {
	Height     int64  `json:"height" toml:"height"`           // Height is the block height at which the query is to be performed.
	MaxRetries int    `json:"max_retries" toml:"max_retries"` // MaxRetries is the maximum number of retries for the query.
	Prove      bool   `json:"prove" toml:"prove"`             // Prove indicates whether to include proof in query results.
	RetryDelay string `json:"retry_delay" toml:"retry_delay"` // RetryDelay is the delay between query retries.
	RPCAddr    string `json:"rpc_addr" toml:"rpc_addr"`       // RPCAddr is the address of the RPC server.
	Timeout    string `json:"timeout" toml:"timeout"`         // Timeout is the maximum duration for the query to be executed.
}

// NewQuery creates a new Query instance with default values.
func NewQuery() *Query {
	return &Query{
		Height:     cast.ToInt64(GetDefault(NameQueryHeight)),
		MaxRetries: cast.ToInt(GetDefault(NameQueryMaxRetries)),
		Prove:      cast.ToBool(GetDefault(NameQueryProve)),
		RetryDelay: cast.ToString(GetDefault(NameQueryRetryDelay)),
		RPCAddr:    cast.ToString(GetDefault(NameQueryRPCAddr)),
		Timeout:    cast.ToString(GetDefault(NameQueryTimeout)),
	}
}

// WithHeight sets the Height field and returns the modified Query instance.
func (q *Query) WithHeight(v int64) *Query {
	q.Height = v
	return q
}

// WithMaxRetries sets the MaxRetries field and returns the modified Query instance.
func (q *Query) WithMaxRetries(v int) *Query {
	q.MaxRetries = v
	return q
}

// WithProve sets the Prove field and returns the modified Query instance.
func (q *Query) WithProve(v bool) *Query {
	q.Prove = v
	return q
}

// WithRetryDelay sets the RetryDelay field and returns the modified Query instance.
func (q *Query) WithRetryDelay(v time.Duration) *Query {
	q.RetryDelay = v.String()
	return q
}

// WithRPCAddr sets the RPCAddr field and returns the modified Query instance.
func (q *Query) WithRPCAddr(v string) *Query {
	q.RPCAddr = v
	return q
}

// WithTimeout sets the Timeout field and returns the modified Query instance.
func (q *Query) WithTimeout(v time.Duration) *Query {
	q.Timeout = v.String()
	return q
}

// GetHeight returns the block height for the query.
func (q *Query) GetHeight() int64 {
	return q.Height
}

// GetMaxRetries returns the maximum number of retries for the query.
func (q *Query) GetMaxRetries() int {
	return q.MaxRetries
}

// GetProve returns whether to include proof in query results.
func (q *Query) GetProve() bool {
	return q.Prove
}

// GetRetryDelay returns the delay between retries for the query.
func (q *Query) GetRetryDelay() time.Duration {
	v, err := time.ParseDuration(q.RetryDelay)
	if err != nil {
		panic(err)
	}

	return v
}

// GetRPCAddr returns the address of the RPC server.
func (q *Query) GetRPCAddr() string {
	return q.RPCAddr
}

// GetTimeout returns the maximum duration for the query.
func (q *Query) GetTimeout() time.Duration {
	v, err := time.ParseDuration(q.Timeout)
	if err != nil {
		panic(err)
	}

	return v
}

// ValidateQueryHeight validates the Height field.
func ValidateQueryHeight(v int64) error {
	if v < 0 {
		return errors.New("height must be non-negative")
	}

	return nil
}

// ValidateQueryMaxRetries validates the MaxRetries field.
func ValidateQueryMaxRetries(v int) error {
	if v < 0 {
		return errors.New("max_retries must be non-negative")
	}

	return nil
}

// ValidateQueryRetryDelay validates the RetryDelay field.
func ValidateQueryRetryDelay(v string) error {
	duration, err := time.ParseDuration(v)
	if err != nil {
		return errors.New("retry_delay must be a valid duration")
	}
	if duration < 0 {
		return errors.New("retry_delay must not be negative")
	}

	return nil
}

// ValidateQueryRPCAddr validates the RPCAddr field.
func ValidateQueryRPCAddr(v string) error {
	if v == "" {
		return errors.New("rpc_addr must not be empty")
	}

	// Parse the URL
	addr, err := url.Parse(v)
	if err != nil {
		return errors.New("rpc_addr must be a valid URL")
	}

	// Check if the URL scheme is set
	if addr.Scheme == "" {
		return errors.New("rpc_addr must have a valid scheme (e.g., http, https)")
	}

	// Check if the port is a valid number
	port, err := strconv.Atoi(addr.Port())
	if err != nil {
		return errors.New("rpc_addr must include a valid port number")
	}

	// Check if the port number is within the valid range
	if port < 1 || port > 65535 {
		return errors.New("rpc_addr must include a port number between 1 and 65535")
	}

	return nil
}

// ValidateQueryTimeout validates the Timeout field.
func ValidateQueryTimeout(v string) error {
	duration, err := time.ParseDuration(v)
	if err != nil {
		return errors.New("timeout must be a valid duration")
	}
	if duration < 0 {
		return errors.New("timeout must not be negative")
	}

	return nil
}

// Validate validates all the fields of the Query struct.
func (q *Query) Validate() error {
	if err := ValidateQueryHeight(q.Height); err != nil {
		return err
	}
	if err := ValidateQueryMaxRetries(q.MaxRetries); err != nil {
		return err
	}
	if err := ValidateQueryRetryDelay(q.RetryDelay); err != nil {
		return err
	}
	if err := ValidateQueryRPCAddr(q.RPCAddr); err != nil {
		return err
	}
	if err := ValidateQueryTimeout(q.Timeout); err != nil {
		return err
	}

	return nil
}

// ABCIQueryOptions converts Query to ABCIQueryOptions.
func (q *Query) ABCIQueryOptions() client.ABCIQueryOptions {
	return client.ABCIQueryOptions{
		Height: q.GetHeight(),
		Prove:  q.GetProve(),
	}
}

// Client creates a new HTTP client with the configured options.
func (q *Query) Client() (*http.HTTP, error) {
	timeout := utils.UIntSecondsFromDuration(q.GetTimeout())
	return http.NewWithTimeout(q.GetRPCAddr(), "/websocket", timeout)
}

// GetQueryHeightFromCmd retrieves the block height from the command-line flags.
func GetQueryHeightFromCmd(cmd *cobra.Command) (int64, error) {
	return cmd.Flags().GetInt64(FlagQueryHeight)
}

// GetQueryMaxRetriesFromCmd retrieves the maximum number of retries from the command-line flags.
func GetQueryMaxRetriesFromCmd(cmd *cobra.Command) (int, error) {
	return cmd.Flags().GetInt(FlagQueryMaxRetries)
}

// GetQueryProveFromCmd retrieves whether to include proof from the command-line flags.
func GetQueryProveFromCmd(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool(FlagQueryProve)
}

// GetQueryRetryDelayFromCmd retrieves the retry delay from the command-line flags.
func GetQueryRetryDelayFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagQueryRetryDelay)
}

// GetQueryRPCAddrFromCmd retrieves the RPC server address from the command-line flags.
func GetQueryRPCAddrFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagQueryRPCAddr)
}

// GetQueryTimeoutFromCmd retrieves the query timeout from the command-line flags.
func GetQueryTimeoutFromCmd(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(FlagQueryTimeout)
}

// SetFlagQueryHeight sets the flag for the block height in the given command.
func SetFlagQueryHeight(cmd *cobra.Command) {
	value := cast.ToInt64(GetDefault(NameQueryHeight))
	cmd.Flags().Int64(FlagQueryHeight, value, "Block height at which the query is to be performed.")
}

// SetFlagQueryMaxRetries sets the flag for the maximum number of retries in the given command.
func SetFlagQueryMaxRetries(cmd *cobra.Command) {
	value := cast.ToInt(GetDefault(NameQueryMaxRetries))
	cmd.Flags().Int(FlagQueryMaxRetries, value, "Maximum number of retries for the query.")
}

// SetFlagQueryProve sets the flag for including proof in the query results in the given command.
func SetFlagQueryProve(cmd *cobra.Command) {
	value := cast.ToBool(GetDefault(NameQueryProve))
	cmd.Flags().Bool(FlagQueryProve, value, "Include proof in query results.")
}

// SetFlagQueryRetryDelay sets the flag for the delay between retries in the given command.
func SetFlagQueryRetryDelay(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameQueryRetryDelay))
	cmd.Flags().String(FlagQueryRetryDelay, value, "Delay between retries for the query.")
}

// SetFlagQueryRPCAddr sets the flag for the RPC server address in the given command.
func SetFlagQueryRPCAddr(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameQueryRPCAddr))
	cmd.Flags().String(FlagQueryRPCAddr, value, "Address of the RPC server.")
}

// SetFlagQueryTimeout sets the flag for the query timeout in the given command.
func SetFlagQueryTimeout(cmd *cobra.Command) {
	value := cast.ToString(GetDefault(NameQueryTimeout))
	cmd.Flags().String(FlagQueryTimeout, value, "Maximum duration for the query to be executed.")
}

// SetQueryFlags sets all query-related flags for the command.
func SetQueryFlags(cmd *cobra.Command) {
	SetFlagQueryHeight(cmd)
	SetFlagQueryMaxRetries(cmd)
	SetFlagQueryProve(cmd)
	SetFlagQueryRetryDelay(cmd)
	SetFlagQueryRPCAddr(cmd)
	SetFlagQueryTimeout(cmd)
}

// NewQueryFromCmd creates a new Query object from the command-line flags.
func NewQueryFromCmd(cmd *cobra.Command) (*Query, error) {
	// Retrieve block height from the command flags
	height, err := GetQueryHeightFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve maximum number of retries from the command flags
	maxRetries, err := GetQueryMaxRetriesFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve whether to include proof from the command flags
	prove, err := GetQueryProveFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve retry delay from the command flags
	retryDelay, err := GetQueryRetryDelayFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve RPC server address from the command flags
	rpcAddr, err := GetQueryRPCAddrFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Retrieve query timeout from the command flags
	timeout, err := GetQueryTimeoutFromCmd(cmd)
	if err != nil {
		return nil, err
	}

	// Return a new Query object with the retrieved values
	return &Query{
		Height:     height,
		MaxRetries: maxRetries,
		Prove:      prove,
		RetryDelay: retryDelay,
		RPCAddr:    rpcAddr,
		Timeout:    timeout,
	}, nil
}
