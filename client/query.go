package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/sentinel-official/sentinel-go-sdk/client/options"
)

// ABCIQueryWithOptions performs an ABCI query with configurable options.
// It retries the query according to the specified maximum number of retries.
func (c *Client) ABCIQueryWithOptions(ctx context.Context, path string, data bytes.HexBytes, opts *options.Options) (*abcitypes.ResponseQuery, error) {
	// Get the RPC client from the provided options.
	rpcClient, err := opts.Client()
	if err != nil {
		return nil, err
	}

	// Retry the query for the specified number of times.
	for t := 0; t < opts.MaxRetries; t++ {
		// Perform the ABCI query with options.
		result, err := rpcClient.ABCIQueryWithOptions(ctx, path, data, opts.ABCIQueryOptions())
		if err != nil {
			// Retry on specific errors, such as EOF or invalid character.
			if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "invalid character '<' looking for beginning of value") {
				continue
			}

			// Return other errors.
			return nil, err
		}

		// If the result is nil, return nil.
		if result == nil {
			return nil, nil
		}

		// Return the response from the successful query.
		return &result.Response, nil
	}

	// Return an error if the maximum retry limit is reached.
	return nil, errors.New("reached max retry limit")
}

// QueryKey performs an ABCI query for a specific key in a store.
func (c *Client) QueryKey(ctx context.Context, store string, data bytes.HexBytes, opts *options.Options) (*abcitypes.ResponseQuery, error) {
	// Construct the path for querying a key in the store.
	path := fmt.Sprintf("/store/%s/key", store)

	// Delegate the ABCI query to ABCIQueryWithOptions.
	return c.ABCIQueryWithOptions(ctx, path, data, opts)
}

// QuerySubspace performs an ABCI query for a subspace in a store.
func (c *Client) QuerySubspace(ctx context.Context, store string, data bytes.HexBytes, opts *options.Options) (*abcitypes.ResponseQuery, error) {
	// Construct the path for querying a subspace in the store.
	path := fmt.Sprintf("/store/%s/subspace", store)

	// Delegate the ABCI query to ABCIQueryWithOptions.
	return c.ABCIQueryWithOptions(ctx, path, data, opts)
}

// QueryGRPC performs a gRPC query using ABCI with configurable options.
// It marshals the request, queries with ABCI, and unmarshals the response.
func (c *Client) QueryGRPC(ctx context.Context, method string, req, resp codec.ProtoMarshaler, opts *options.Options) error {
	// Marshal the gRPC request.
	data, err := c.Marshal(req)
	if err != nil {
		return err
	}

	// Perform ABCI query with options.
	reply, err := c.ABCIQueryWithOptions(ctx, method, data, opts)
	if err != nil {
		return err
	}

	// Check for a nil reply.
	if reply == nil {
		return errors.New("nil reply")
	}

	// Unmarshal the ABCI response value into the provided response object.
	if err := c.Unmarshal(reply.Value, resp); err != nil {
		return err
	}

	// Return nil on success.
	return nil
}
