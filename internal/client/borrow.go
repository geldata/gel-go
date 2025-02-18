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
	"fmt"

	gelerrint "github.com/edgedb/edgedb-go/internal/gelerr"
)

type borrowableConn struct {
	conn   *protocolConnection
	reason string
}

func (c *borrowableConn) borrow(reason string) (*protocolConnection, error) {
	switch c.reason {
	case "":
		// this is the expected value
	case "transaction":
		return nil, gelerrint.NewInterfaceError(
			"The connection is borrowed for a transaction. "+
				"Use the methods on the transaction object instead.",
			nil,
		)
	case "subtransaction":
		return nil, gelerrint.NewInterfaceError(
			"The transaction is borrowed for a subtransaction. "+
				"Use the methods on the subtransaction object instead.",
			nil,
		)
	default:
		return nil, gelerrint.NewInterfaceError(fmt.Sprintf(
			"existing borrow reason is unexpected: %q", c.reason), nil)
	}

	switch reason {
	case "transaction", "subtransaction":
		c.reason = reason
		return c.conn, nil
	default:
		return nil, gelerrint.NewInterfaceError(fmt.Sprintf(
			"unexpected borrow reason: %q", c.reason), nil)
	}
}

func (c *borrowableConn) unborrow() error {
	if c.reason == "" {
		return gelerrint.NewInterfaceError(
			"not currently borrowed, cannot unborrow",
			nil,
		)
	}

	c.reason = ""
	return nil
}

func (c *borrowableConn) assertUnborrowed() error {
	switch c.reason {
	case "":
		return nil
	case "transaction":
		return gelerrint.NewInterfaceError(
			"The connection is borrowed for a transaction. "+
				"Use the methods on the transaction object instead.",
			nil,
		)
	case "subtransaction":
		return gelerrint.NewInterfaceError(
			"The transaction is borrowed for a subtransaction. "+
				"Use the methods on the subtransaction object instead.",
			nil,
		)
	default:
		return gelerrint.NewInterfaceError(fmt.Sprintf(
			"existing borrow reason is unexpected: %q", c.reason), nil)
	}
}

func (c *borrowableConn) capabilities1pX() uint64 {
	return userCapabilities
}

func (c *borrowableConn) scriptFlow(ctx context.Context, q *query) error {
	if e := c.assertUnborrowed(); e != nil {
		return e
	}

	return c.conn.scriptFlow(ctx, q)
}

func (c *borrowableConn) granularFlow(ctx context.Context, q *query) error {
	if e := c.assertUnborrowed(); e != nil {
		return e
	}

	return c.conn.granularFlow(ctx, q)
}
