package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/cometbft/cometbft/rpc/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/v2fly/v2ray-core/v5/common/retry"
)

// ABCIQueryWithOptions performs an ABCI query with configurable options.
// It retries the query in case of failures based on the Client's retry configuration.
// Returns the ABCI query response or an error.
func (c *Client) ABCIQueryWithOptions(ctx context.Context, path string, data bytes.HexBytes) (*abcitypes.ResponseQuery, error) {
	var result *coretypes.ResultABCIQuery

	// Define the function to perform the ABCI query.
	fn := func() error {
		// Get the RPC client for querying.
		http, err := c.HTTP()
		if err != nil {
			return err
		}

		// Configure the query options.
		opts := client.ABCIQueryOptions{
			Height: c.queryHeight,
			Prove:  c.queryProve,
		}

		// Perform the query and store the result.
		result, err = http.ABCIQueryWithOptions(ctx, path, data, opts)
		if err != nil {
			return err
		}

		return nil
	}

	// Retry the query using the configured maximum retries and delay.
	retryDelay := uint32(c.queryRetryDelay / time.Millisecond)
	if err := retry.Timed(c.queryMaxRetries, retryDelay).On(fn); err != nil {
		return nil, err
	}

	// Return nil if no result was produced.
	if result == nil {
		return nil, nil
	}

	return &result.Response, nil
}

// QueryKey performs an ABCI query for a specific key in a store.
// Constructs the query path and delegates the query to ABCIQueryWithOptions.
// Returns the query response or an error.
func (c *Client) QueryKey(ctx context.Context, store string, data bytes.HexBytes) (*abcitypes.ResponseQuery, error) {
	// Construct the path for querying the key.
	path := fmt.Sprintf("/store/%s/key", store)

	// Perform the query.
	return c.ABCIQueryWithOptions(ctx, path, data)
}

// QuerySubspace performs an ABCI query for a subspace in a store.
// Constructs the query path and delegates the query to ABCIQueryWithOptions.
// Returns the query response or an error.
func (c *Client) QuerySubspace(ctx context.Context, store string, data bytes.HexBytes) (*abcitypes.ResponseQuery, error) {
	// Construct the path for querying the subspace.
	path := fmt.Sprintf("/store/%s/subspace", store)

	// Perform the query.
	return c.ABCIQueryWithOptions(ctx, path, data)
}

// QueryGRPC performs a gRPC query using ABCI with configurable options.
// Marshals the request, queries via ABCI, and unmarshals the response.
// Returns an error if any step fails.
func (c *Client) QueryGRPC(ctx context.Context, method string, req, resp codec.ProtoMarshaler) error {
	// Marshal the request into bytes.
	data, err := c.Marshal(req)
	if err != nil {
		return err
	}

	// Perform the query using ABCIQueryWithOptions.
	reply, err := c.ABCIQueryWithOptions(ctx, method, data)
	if err != nil {
		return err
	}

	// Check for a nil reply.
	if reply == nil {
		return errors.New("nil reply")
	}

	// Unmarshal the response value into the provided response object.
	if err := c.Unmarshal(reply.Value, resp); err != nil {
		return err
	}

	return nil
}
