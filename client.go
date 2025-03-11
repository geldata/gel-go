// This source file is part of the Gel open source project.
//
// Copyright Gel Data Inc. and the Gel authors.
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

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/geltypes"
	gel "github.com/geldata/gel-go/internal/client"
)

// CreateClient returns a new client. The client connects lazily. Call
// [Client.EnsureConnected] to force a connection.
//
// Instead of providing connection details directly, we recommend connecting
// using projects locally, and environment variables for remote instances,
// providing an empty [gelcfg.Options] struct here.
// [Read more] about the recommended ways to configure client connections.
//
// [Read more]: https://docs.geldata.com/reference/clients/connection
func CreateClient(opts gelcfg.Options) (*Client, error) { // nolint:gocritic,lll
	return CreateClientDSN("", opts)
}

// CreateClientDSN returns a new Client. We recommend using [CreateClient] for
// most use cases.
//
// dsn is either an instance name or a [DSN].
//
// [DSN]: https://docs.geldata.com/reference/reference/dsn#ref-dsn
func CreateClientDSN(dsn string, opts gelcfg.Options) (*Client, error) { // nolint:gocritic,lll
	pool, err := gel.NewPool(dsn, opts)
	if err != nil {
		return nil, err
	}

	p := &Client{pool: pool}

	return p, nil
}

// Client is a connection pool and is safe for concurrent use.
type Client struct {
	pool *gel.Pool
}

// EnsureConnected forces the client to connect if it hasn't already. This can
// be used to ensure that your program will fail early in the case that the
// [configured connection parameters] are not correct.
//
// [configured connection parameters]: https://docs.geldata.com/reference/clients/connection
func (c *Client) EnsureConnected(ctx context.Context) error {
	return c.pool.EnsureConnected(ctx)
}

// Close closes all connections in the client.
// Calling Close() blocks until all acquired connections have been released,
// and returns an error if called more than once.
func (c *Client) Close() error { return c.pool.Close() }

// Execute an EdgeQL command (or commands).
func (c *Client) Execute(ctx context.Context, cmd string, args ...interface{}) error { //nolint:lll
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
		&c.pool.QueryConfig,
	)
	if err != nil {
		return err
	}

	err = conn.ScriptFlow(ctx, q)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// Query runs a query and returns the results.
func (c *Client) Query(ctx context.Context, cmd string, out interface{}, args ...interface{}) error { //nolint:lll
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = gel.RunQuery(
		ctx,
		conn,
		"Query",
		cmd,
		out,
		args,
		c.pool.State,
		&c.pool.QueryConfig,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// needed for dock link [gelerr.NoDataError]
var _ = gelerr.NoDataError

// QuerySingle runs a singleton-returning query and returns its element.  If
// the query executes successfully but doesn't return a result a
// [gelerr.NoDataError] is returned.  If the out
// argument is an optional type the out argument will be set to missing instead
// of returning a NoDataError.
func (c *Client) QuerySingle(ctx context.Context, cmd string, out interface{}, args ...interface{}) error { //nolint:lll
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
		&c.pool.QueryConfig,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// QueryJSON runs a query and returns the results as JSON.
func (c *Client) QueryJSON(ctx context.Context, cmd string, out *[]byte, args ...interface{}) error { //nolint:lll
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
		&c.pool.QueryConfig,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// QuerySingleJSON runs a singleton-returning query.
// If the query executes successfully but doesn't have a result
// a [gelerr.NoDataError] is returned.
func (c *Client) QuerySingleJSON(ctx context.Context, cmd string, out interface{}, args ...interface{}) error { //nolint:lll
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
		&c.pool.QueryConfig,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// QuerySQL runs a SQL query and returns the results.
func (c *Client) QuerySQL(ctx context.Context, cmd string, out interface{}, args ...interface{}) error { //nolint:lll
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = gel.RunQuery(
		ctx,
		conn,
		"QuerySQL",
		cmd,
		out,
		args,
		c.pool.State,
		&c.pool.QueryConfig,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// ExecuteSQL executes a SQL command (or commands).
func (c *Client) ExecuteSQL(ctx context.Context, cmd string, args ...interface{}) error { //nolint:lll
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
		&c.pool.QueryConfig,
	)
	if err != nil {
		return err
	}

	err = conn.ScriptFlow(ctx, q)
	return gel.FirstError(err, c.pool.Release(conn, err))
}

// Tx runs action in a transaction retrying failed attempts. Queries must be
// executed on the [geltypes.Tx] that is passed to action. Queries executed on
// the client in a geltypes.TxBlock will not run in the transaction and will be
// applied immediately.
//
// The geltypes.TxBlock may be re-run if any of the queries fail in a way that
// might succeed on subsequent attempts. Retries are governed by
// [gelcfg.RetryOptions] and [gelcfg.RetryRule].  Retry options can be set
// using [Client.WithRetryOptions].  See [gelcfg.RetryRule] for more details on
// how they work.
func (c *Client) Tx(ctx context.Context, action geltypes.TxBlock) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = conn.Tx(
		ctx,
		action,
		c.pool.State,
		&c.pool.QueryConfig,
	)
	return gel.FirstError(err, c.pool.Release(conn, err))
}
