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
	"errors"
	"strings"
	"time"

	"github.com/geldata/gel-go/gelerr"
	types "github.com/geldata/gel-go/geltypes"
	gelerrint "github.com/geldata/gel-go/internal/gelerr"
)

type transactableConn struct {
	*reconnectingConn
}

// sleepBackoff waits out a retry back-off unless ctx ends first. A query
// whose ctx is already over cannot succeed on retry; retrying it only
// re-entered the connection with a dead ctx, which is how the reader
// hand-off used to get stranded (see acquireReader).
func sleepBackoff(ctx context.Context, d time.Duration) bool {
	if ctx.Err() != nil {
		return false
	}

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func (c *transactableConn) withRetries(
	ctx context.Context,
	q *query,
	cb func(context.Context, *query) error,
) error {
	var (
		err    error
		edbErr gelerr.Error
	)

	for i := 1; true; i++ {
		if errors.As(err, &edbErr) && c.conn.soc.Closed() {
			err = c.reconnect(ctx, true)
			if err != nil {
				goto Error
			}
		}

		err = cb(ctx, q)

	Error:
		// q is a read only query if it has no capabilities
		// i.e. capabilities == 0. Read only queries are always
		// retryable, mutation queries are retryable if the
		// error explicitly indicates a transaction conflict.
		capabilities, ok := c.getCachedCapabilities(q)
		if ok &&
			errors.As(err, &edbErr) &&
			edbErr.HasTag(gelerr.ShouldRetry) &&
			(capabilities == 0 ||
				edbErr.Category(gelerr.TransactionConflictError)) {
			rule, e := q.cfg.RetryOptions.RuleForException(edbErr)
			if e != nil {
				return e
			}

			if i >= rule.Attempts() {
				return err
			}

			if !sleepBackoff(ctx, rule.Backoff()(i)) {
				return err
			}
			continue
		}

		return err
	}

	return gelerrint.NewClientError("unreachable", nil)
}

func (c *transactableConn) granularFlow(ctx context.Context, q *query) error {
	return c.withRetries(ctx, q, c.reconnectingConn.granularFlow)
}

func (c *transactableConn) ScriptFlow(ctx context.Context, q *query) error {
	return c.withRetries(ctx, q, c.reconnectingConn.ScriptFlow)
}

// abandonTx retires the connection when a commit or rollback did not get
// a server answer: a gelerr.Error means the server answered and passes
// through untouched; anything else means the statement was never sent (a
// dead ctx) or its answer could not be read, so the transaction may still
// be open on the server. Closing the socket makes the server roll it back
// and keeps a mid-transaction session out of the pool. A ctx error is
// returned as ClientConnectionClosedError so Pool.Release drops the
// connection at once; any other error keeps its type (the outcome is
// unknown, a caller must not take it for a retryable connection loss).
func abandonTx(conn *protocolConnection, err error) error {
	var edbErr gelerr.Error
	if err == nil || errors.As(err, &edbErr) {
		return err
	}

	_ = conn.soc.Close()
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return gelerrint.NewClientConnectionClosedError("", err)
	}

	return err
}

func (c *transactableConn) Tx(
	ctx context.Context,
	action types.TxBlock,
	state map[string]interface{},
	cfg *QueryConfig,
) (err error) {
	conn, err := c.borrow("transaction")
	if err != nil {
		return err
	}
	defer func() { err = FirstError(err, c.unborrow()) }()

	optimisticRepeatableRead := true

	var edbErr gelerr.Error
	for i := 1; true; i++ {
		if errors.As(err, &edbErr) && c.conn.soc.Closed() {
			err = c.reconnect(ctx, true)
			if err != nil {
				goto Error
			}
			// get the newly connected protocolConnection
			conn = c.conn
		}

		{
			tx := &Tx{
				borrowableConn: borrowableConn{conn: conn},
				txState:        &txState{},
				state:          state,
				cfg:            *cfg,
			}
			err = tx.start(ctx, optimisticRepeatableRead)
			if err != nil {
				goto Error
			}

			err = action(ctx, tx)
			if err == nil {
				err = abandonTx(conn, tx.commit(ctx))
				if errors.As(err, &edbErr) &&
					edbErr.Category(gelerr.TransactionError) &&
					edbErr.HasTag(gelerr.ShouldRetry) {
					goto Error
				}
				return err
			} else if isClientConnectionError(err) && conn.isClosed() {
				// This socket is gone; the server dropped the transaction
				// with it. A connection error from elsewhere (the callback
				// used another client) leaves this transaction open, so it
				// still gets the rollback below.
				goto Error
			}

			if e := tx.rollback(ctx); e != nil && !errors.As(e, &edbErr) {
				return abandonTx(conn, e)
			}
		}

	Error:
		if errors.As(err, &edbErr) &&
			edbErr.Category(gelerr.CapabilityError) &&
			strings.Contains(err.Error(), "REPEATABLE READ") {
			if !optimisticRepeatableRead {
				return err
			}

			optimisticRepeatableRead = false
			i--
			continue
		}

		if errors.As(err, &edbErr) && edbErr.HasTag(gelerr.ShouldRetry) {
			rule, e := cfg.RetryOptions.RuleForException(edbErr)
			if e != nil {
				return e
			}

			if i >= rule.Attempts() {
				return err
			}

			if !sleepBackoff(ctx, rule.Backoff()(i)) {
				return err
			}
			continue
		}

		return err
	}

	return gelerrint.NewClientError("unreachable", nil)
}
