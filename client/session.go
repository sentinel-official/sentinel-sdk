package client

import (
	"context"

	cosmossdk "github.com/cosmos/cosmos-sdk/types"
	base "github.com/sentinel-official/hub/v12/types"
	"github.com/sentinel-official/hub/v12/x/session/types/v2"

	"github.com/sentinel-official/sentinel-go-sdk/v1/client/options"
)

const (
	// gRPC methods for querying session information
	methodQuerySession                           = "/sentinel.session.v2.QueryService/QuerySession"
	methodQuerySessions                          = "/sentinel.session.v2.QueryService/QuerySessions"
	methodQuerySessionsForAccount                = "/sentinel.session.v2.QueryService/QuerySessionsForAccount"
	methodQuerySessionsForNode                   = "/sentinel.session.v2.QueryService/QuerySessionsForNode"
	methodQuerySessionsForSubscription           = "/sentinel.session.v2.QueryService/QuerySessionsForSubscription"
	methodQuerySessionsForSubscriptionAllocation = "/sentinel.session.v2.QueryService/QuerySessionsForAllocation"
)

// Session queries and returns information about a specific session based on the provided session ID.
// It uses gRPC to send a request to the "/sentinel.session.v2.QueryService/QuerySession" endpoint.
// The result is a pointer to v2.Session and an error if the query fails.
func (c *Context) Session(ctx context.Context, id uint64, opts *options.QueryOptions) (res *v2.Session, err error) {
	// Initialize variables for the query.
	var (
		resp v2.QuerySessionResponse
		req  = &v2.QuerySessionRequest{
			Id: id,
		}
	)

	// Send a gRPC query using the provided context, method, request, response, and options.
	if err := c.QueryGRPC(ctx, methodQuerySession, req, &resp, opts); err != nil {
		return nil, err
	}

	// Return a pointer to the session and a nil error.
	return &resp.Session, nil
}

// Sessions queries and returns a list of sessions based on the provided options.
// It uses gRPC to send a request to the "/sentinel.session.v2.QueryService/QuerySessions" endpoint.
// The result is a slice of v2.Session and an error if the query fails.
func (c *Context) Sessions(ctx context.Context, opts *options.QueryOptions) (res []v2.Session, err error) {
	// Initialize variables for the query.
	var (
		resp v2.QuerySessionsResponse
		req  = &v2.QuerySessionsRequest{
			Pagination: opts.PageRequest(),
		}
	)

	// Send a gRPC query using the provided context, method, request, response, and options.
	if err := c.QueryGRPC(ctx, methodQuerySessions, req, &resp, opts); err != nil {
		return nil, err
	}

	// Return the list of sessions and a nil error.
	return resp.Sessions, nil
}

// SessionsForAccount queries and returns a list of sessions associated with a specific account
// based on the provided account address and options.
// It uses gRPC to send a request to the "/sentinel.session.v2.QueryService/QuerySessionsForAccount" endpoint.
// The result is a slice of v2.Session and an error if the query fails.
func (c *Context) SessionsForAccount(ctx context.Context, accAddr cosmossdk.AccAddress, opts *options.QueryOptions) (res []v2.Session, err error) {
	// Initialize variables for the query.
	var (
		resp v2.QuerySessionsForAccountResponse
		req  = &v2.QuerySessionsForAccountRequest{
			Address:    accAddr.String(),
			Pagination: opts.PageRequest(),
		}
	)

	// Send a gRPC query using the provided context, method, request, response, and options.
	if err := c.QueryGRPC(ctx, methodQuerySessionsForAccount, req, &resp, opts); err != nil {
		return nil, err
	}

	// Return the list of sessions and a nil error.
	return resp.Sessions, nil
}

// SessionsForNode queries and returns a list of sessions associated with a specific node
// based on the provided node address and options.
// It uses gRPC to send a request to the "/sentinel.session.v2.QueryService/QuerySessionsForNode" endpoint.
// The result is a slice of v2.Session and an error if the query fails.
func (c *Context) SessionsForNode(ctx context.Context, nodeAddr base.NodeAddress, opts *options.QueryOptions) (res []v2.Session, err error) {
	// Initialize variables for the query.
	var (
		resp v2.QuerySessionsForNodeResponse
		req  = &v2.QuerySessionsForNodeRequest{
			Address:    nodeAddr.String(),
			Pagination: opts.PageRequest(),
		}
	)

	// Send a gRPC query using the provided context, method, request, response, and options.
	if err := c.QueryGRPC(ctx, methodQuerySessionsForNode, req, &resp, opts); err != nil {
		return nil, err
	}

	// Return the list of sessions and a nil error.
	return resp.Sessions, nil
}

// SessionsForSubscription queries and returns a list of sessions associated with a specific subscription
// based on the provided subscription ID and options.
// It uses gRPC to send a request to the "/sentinel.session.v2.QueryService/QuerySessionsForSubscription" endpoint.
// The result is a slice of v2.Session and an error if the query fails.
func (c *Context) SessionsForSubscription(ctx context.Context, id uint64, opts *options.QueryOptions) (res []v2.Session, err error) {
	// Initialize variables for the query.
	var (
		resp v2.QuerySessionsForSubscriptionResponse
		req  = &v2.QuerySessionsForSubscriptionRequest{
			Id:         id,
			Pagination: opts.PageRequest(),
		}
	)

	// Send a gRPC query using the provided context, method, request, response, and options.
	if err := c.QueryGRPC(ctx, methodQuerySessionsForSubscription, req, &resp, opts); err != nil {
		return nil, err
	}

	// Return the list of sessions and a nil error.
	return resp.Sessions, nil
}

// SessionsForSubscriptionAllocation queries and returns a list of sessions associated with a specific subscription allocation
// based on the provided subscription ID, account address, and options.
// It uses gRPC to send a request to the "/sentinel.session.v2.QueryService/QuerySessionsForAllocation" endpoint.
// The result is a slice of v2.Session and an error if the query fails.
func (c *Context) SessionsForSubscriptionAllocation(ctx context.Context, id uint64, accAddr cosmossdk.AccAddress, opts *options.QueryOptions) (res []v2.Session, err error) {
	// Initialize variables for the query.
	var (
		resp v2.QuerySessionsForAllocationResponse
		req  = &v2.QuerySessionsForAllocationRequest{
			Id:         id,
			Address:    accAddr.String(),
			Pagination: opts.PageRequest(),
		}
	)

	// Send a gRPC query using the provided context, method, request, response, and options.
	if err := c.QueryGRPC(ctx, methodQuerySessionsForSubscriptionAllocation, req, &resp, opts); err != nil {
		return nil, err
	}

	// Return the list of sessions and a nil error.
	return resp.Sessions, nil
}
