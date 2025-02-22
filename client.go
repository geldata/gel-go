// This source file is part of the EdgeDB open source project.
//
// Copyright EdgeDB Inc. and the EdgeDB authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gel

import (
	"context"

	"github.com/edgedb/edgedb-go/gelcfg"
	"github.com/edgedb/edgedb-go/geltypes"
	gel "github.com/edgedb/edgedb-go/internal/client"
)

// CreateClient returns a new client. The client connects lazily. Call
// Client.EnsureConnected() to force a connection.
func CreateClient(ctx context.Context, opts gelcfg.Options) (*Client, error) { // nolint:gocritic,lll
	return CreateClientDSN(ctx, "", opts)
}

// CreateClientDSN returns a new client. See also CreateClient.
//
// dsn is either an instance name
// https://www.edgedb.com/docs/clients/connection
// or it specifies a single string in the following format:
//
//	gel://user:password@host:port/database?option=value.
//
// The following options are recognized: host, port, user, database, password.
func CreateClientDSN(_ context.Context, dsn string, opts gelcfg.Options) (*Client, error) { // nolint:gocritic,lll
	pool, err := gel.NewPool(dsn, opts)
	if err != nil {
		return nil, err
	}

	warningHandler := gelcfg.LogWarnings
	if opts.WarningHandler != nil {
		warningHandler = opts.WarningHandler
	}

	p := &Client{
		pool:           pool,
		warningHandler: warningHandler,
	}

	return p, nil
}

// Client is a connection pool and is safe for concurrent use.
type Client struct {
	pool           *gel.Pool
	warningHandler gelcfg.WarningHandler
}

// EnsureConnected forces the client to connect if it hasn't already.
func (c *Client) EnsureConnected(ctx context.Context) error {
	return c.pool.EnsureConnected(ctx)
}

// Close closes all connections in the client.
// Calling close blocks until all acquired connections have been released,
// and returns an error if called more than once.
func (c *Client) Close() error { return c.pool.Close() }

// Execute an EdgeQL command (or commands).
func (c *Client) Execute(
	ctx context.Context,
	cmd string,
	args ...interface{},
) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	q, err := gel.NewQuery(
		"Execute",
		cmd,
		args,
		conn.Capabilities1pX(),
		gel.CopyState(c.pool.State),
		nil,
		true,
		c.warningHandler,
	)
	if err != nil {
		return err
	}

	err = conn.ScriptFlow(ctx, q)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// Query runs a query and returns the results.
func (c *Client) Query(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = gel.RunQuery(
		ctx, conn, "Query", cmd, out, args, c.pool.State, c.warningHandler)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// QuerySingle runs a singleton-returning query and returns its element.
// If the query executes successfully but doesn't return a result
// a NoDataError is returned. If the out argument is an optional type the out
// argument will be set to missing instead of returning a NoDataError.
func (c *Client) QuerySingle(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = gel.RunQuery(
		ctx,
		conn,
		"QuerySingle",
		cmd,
		out,
		args,
		c.pool.State,
		c.warningHandler,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// QueryJSON runs a query and return the results as JSON.
func (c *Client) QueryJSON(
	ctx context.Context,
	cmd string,
	out *[]byte,
	args ...interface{},
) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = gel.RunQuery(
		ctx,
		conn,
		"QueryJSON",
		cmd,
		out,
		args,
		c.pool.State,
		c.warningHandler,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// QuerySingleJSON runs a singleton-returning query.
// If the query executes successfully but doesn't have a result
// a NoDataError is returned.
func (c *Client) QuerySingleJSON(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = gel.RunQuery(
		ctx,
		conn,
		"QuerySingleJSON",
		cmd,
		out,
		args,
		c.pool.State,
		c.warningHandler,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// QuerySQL runs a SQL query and returns the results.
func (c *Client) QuerySQL(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = gel.RunQuery(
		ctx, conn, "QuerySQL", cmd, out, args, c.pool.State, c.warningHandler)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// ExecuteSQL executes a SQL command (or commands).
func (c *Client) ExecuteSQL(
	ctx context.Context,
	cmd string,
	args ...interface{},
) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	q, err := gel.NewQuery(
		"ExecuteSQL",
		cmd,
		args,
		conn.Capabilities1pX(),
		gel.CopyState(c.pool.State),
		nil,
		true,
		c.warningHandler,
	)
	if err != nil {
		return err
	}

	err = conn.ScriptFlow(ctx, q)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// Tx runs an action in a transaction retrying failed actions
// if they might succeed on a subsequent attempt.
//
// Retries are governed by retry rules.
// The default rule can be set with WithRetryRule().
// For more fine grained control a retry rule can be set
// for each defined RetryCondition using WithRetryCondition().
// When a transaction fails but is retryable
// the rule for the failure condition is used to determine if the transaction
// should be tried again based on RetryRule.Attempts and the amount of time
// to wait before retrying is determined by RetryRule.Backoff.
// If either field is unset (see RetryRule) then the default rule is used.
// If the object's default is unset the fall back is 3 attempts
// and exponential backoff.
func (c *Client) Tx(ctx context.Context, action geltypes.TxBlock) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = conn.Tx(ctx, action, c.pool.State, c.warningHandler)
	return gel.FirstError(err, c.pool.Release(conn, err))
}
