package client

import (
	"time"

	"cosmossdk.io/log"
	"github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
)

// contextKey is a custom type used as a key to store the Client struct in a context.
type contextKey byte

// ContextKey is the constant key value used for storing and retrieving the Client struct from a context.
const ContextKey contextKey = 0

// Client contains all necessary components for transaction handling, query management, and configuration settings.
type Client struct {
	chainID              string                    // The chain ID used to identify the blockchain network
	keyring              keyring.Keyring           // Keyring for managing private keys and signatures
	logger               log.Logger                // Embedded logger
	protoCodec           codec.ProtoCodecMarshaler // Used for marshaling and unmarshaling protobuf data
	queryHeight          int64                     // Query height for blockchain data
	queryMaxRetries      int                       // Maximum number of retries for queries
	queryProve           bool                      // Flag indicating whether to prove queries
	queryRetryDelay      time.Duration             // Delay between query retries
	rpcAddr              string                    // RPC server address
	rpcTimeout           time.Duration             // RPC timeout duration
	txConfig             client.TxConfig           // Configuration related to transactions (e.g., signing modes)
	txFeeGranterAddr     types.AccAddress          // Address that grants transaction fees
	txFees               types.Coins               // Fees for transactions
	txFromName           string                    // Sender name for transactions
	txGasAdjustment      float64                   // Adjustment factor for gas estimation
	txGasPrices          types.DecCoins            // Gas price settings for transactions
	txGas                uint64                    // Gas limit for transactions
	txMemo               string                    // Memo attached to transactions
	txSimulateAndExecute bool                      // Flag for simulating and executing transactions
	txTimeoutHeight      uint64                    // Transaction timeout height
}

// New initializes a new Client instance.
func New() *Client {
	return &Client{}
}

// WithChainID sets the blockchain chain ID and returns the updated Client.
func (c *Client) WithChainID(chainID string) *Client {
	c.chainID = chainID
	return c
}

// WithKeyring assigns the keyring to the Client and returns the updated Client.
func (c *Client) WithKeyring(keyring keyring.Keyring) *Client {
	c.keyring = keyring
	return c
}

// WithLogger assigns a logger instance to the Client and returns the updated Client.
func (c *Client) WithLogger(logger log.Logger) *Client {
	c.logger = logger
	return c
}

// WithProtoCodec sets the protobuf codec and returns the updated Client.
func (c *Client) WithProtoCodec(protoCodec codec.ProtoCodecMarshaler) *Client {
	c.protoCodec = protoCodec
	return c
}

// WithQueryMaxRetries sets the maximum number of retries for queries and returns the updated Client.
func (c *Client) WithQueryMaxRetries(maxRetries int) *Client {
	c.queryMaxRetries = maxRetries
	return c
}

// WithQueryProve sets the prove flag for queries and returns the updated Client.
func (c *Client) WithQueryProve(prove bool) *Client {
	c.queryProve = prove
	return c
}

// WithQueryRetryDelay sets the retry delay duration for queries and returns the updated Client.
func (c *Client) WithQueryRetryDelay(delay time.Duration) *Client {
	c.queryRetryDelay = delay
	return c
}

// WithRPCAddr sets the RPC server address and returns the updated Client.
func (c *Client) WithRPCAddr(rpcAddr string) *Client {
	c.rpcAddr = rpcAddr
	return c
}

// WithRPCTimeout sets the RPC timeout duration and returns the updated Client.
func (c *Client) WithRPCTimeout(timeout time.Duration) *Client {
	c.rpcTimeout = timeout
	return c
}

// WithTxConfig sets the transaction configuration and returns the updated Client.
func (c *Client) WithTxConfig(txConfig client.TxConfig) *Client {
	c.txConfig = txConfig
	return c
}

// WithTxFeeGranterAddr sets the transaction fee granter address and returns the updated Client.
func (c *Client) WithTxFeeGranterAddr(addr types.AccAddress) *Client {
	c.txFeeGranterAddr = addr
	return c
}

// WithTxFees assigns transaction fees and returns the updated Client.
func (c *Client) WithTxFees(fees types.Coins) *Client {
	c.txFees = fees
	return c
}

// WithTxFromName sets the "from" name for transactions and returns the updated Client.
func (c *Client) WithTxFromName(name string) *Client {
	c.txFromName = name
	return c
}

// WithTxGasAdjustment sets the gas adjustment factor for transactions and returns the updated Client.
func (c *Client) WithTxGasAdjustment(adjustment float64) *Client {
	c.txGasAdjustment = adjustment
	return c
}

// WithTxGasPrices sets the gas prices for transactions and returns the updated Client.
func (c *Client) WithTxGasPrices(prices types.DecCoins) *Client {
	c.txGasPrices = prices
	return c
}

// WithTxGas sets the gas limit for transactions and returns the updated Client.
func (c *Client) WithTxGas(gas uint64) *Client {
	c.txGas = gas
	return c
}

// WithTxMemo sets the memo for transactions and returns the updated Client.
func (c *Client) WithTxMemo(memo string) *Client {
	c.txMemo = memo
	return c
}

// WithTxSimulateAndExecute sets the simulate and execute flag and returns the updated Client.
func (c *Client) WithTxSimulateAndExecute(simulate bool) *Client {
	c.txSimulateAndExecute = simulate
	return c
}

// WithTxTimeoutHeight sets the timeout height for transactions and returns the updated Client.
func (c *Client) WithTxTimeoutHeight(height uint64) *Client {
	c.txTimeoutHeight = height
	return c
}

// HTTP creates an HTTP client for the given RPC address and timeout configuration.
// Returns the HTTP client or an error if initialization fails.
func (c *Client) HTTP() (*http.HTTP, error) {
	timeout := uint(c.rpcTimeout / time.Second)
	return http.NewWithTimeout(c.rpcAddr, "/websocket", timeout)
}
