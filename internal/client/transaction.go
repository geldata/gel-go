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

	"github.com/edgedb/edgedb-go/gelcfg"
	"github.com/edgedb/edgedb-go/internal/gelerr"
)

type txStatus int

const (
	newTx txStatus = iota
	startedTx
	committedTx
	rolledBackTx
	failedTx
)

type txState struct {
	txStatus txStatus
}

// assertNotDone returns an error if the transaction is in a done state.
func (s *txState) assertNotDone(opName string) error {
	switch s.txStatus {
	case committedTx:
		return gelerr.NewInterfaceError(fmt.Sprintf(
			"cannot %v; the transaction is already committed", opName,
		), nil)
	case rolledBackTx:
		return gelerr.NewInterfaceError(fmt.Sprintf(
			"cannot %v; the transaction is already rolled back", opName,
		), nil)
	case failedTx:
		return gelerr.NewInterfaceError(fmt.Sprintf(
			"cannot %v; the transaction is in error state", opName,
		), nil)
	default:
		return nil
	}
}

// assertStarted returns an error if the transaction is not in Started state.
func (s *txState) assertStarted(opName string) error {
	switch s.txStatus {
	case startedTx:
		return nil
	case newTx:
		return gelerr.NewInterfaceError(fmt.Sprintf(
			"cannot %v; the transaction is not yet started", opName,
		), nil)
	default:
		return s.assertNotDone(opName)
	}
}

// Tx is a transaction. Use Client.Tx() to get a transaction.
type Tx struct {
	borrowableConn
	*txState
	options        gelcfg.TxOptions
	state          map[string]interface{}
	warningHandler gelcfg.WarningHandler
}

func (t *Tx) execute(
	ctx context.Context,
	cmd string,
	sucessState txStatus,
) error {
	q, err := NewQuery(
		"Execute",
		cmd,
		nil,
		txCapabilities,
		t.state,
		nil,
		false,
		t.warningHandler,
	)
	if err != nil {
		return err
	}

	err = t.borrowableConn.ScriptFlow(ctx, q)

	switch err {
	case nil:
		t.txStatus = sucessState
	default:
		t.txStatus = failedTx
	}

	return err
}

func startTxQuery(o gelcfg.TxOptions) string { // nolint:gocritic
	query := "START TRANSACTION"

	switch o.IsolationLevel() {
	case gelcfg.Serializable:
		query += " ISOLATION SERIALIZABLE"
	default:
		panic(fmt.Sprintf("unknown isolation level: %q", o.IsolationLevel()))
	}

	if o.ReadOnly() {
		query += ", READ ONLY"
	} else {
		query += ", READ WRITE"
	}

	if o.Deferrable() {
		query += ", DEFERRABLE"
	} else {
		query += ", NOT DEFERRABLE"
	}

	query += ";"
	return query
}

func (t *Tx) start(ctx context.Context) error {
	if e := t.assertNotDone("start"); e != nil {
		return e
	}

	if t.txStatus == startedTx {
		return gelerr.NewInterfaceError(
			"cannot start; the transaction is already started",
			nil,
		)
	}

	query := startTxQuery(t.options)
	return t.execute(ctx, query, startedTx)
}

func (t *Tx) commit(ctx context.Context) error {
	if e := t.assertStarted("commit"); e != nil {
		return e
	}

	return t.execute(ctx, "COMMIT;", committedTx)
}

func (t *Tx) rollback(ctx context.Context) error {
	if e := t.assertStarted("rollback"); e != nil {
		return e
	}

	return t.execute(ctx, "ROLLBACK;", rolledBackTx)
}

func (t *Tx) scriptFlow(ctx context.Context, q *query) error {
	if e := t.assertStarted("Execute"); e != nil {
		return e
	}

	return t.borrowableConn.ScriptFlow(ctx, q)
}

func (t *Tx) granularFlow(ctx context.Context, q *query) error {
	if e := t.assertStarted(q.method); e != nil {
		return e
	}

	return t.borrowableConn.granularFlow(ctx, q)
}

// Execute an EdgeQL command (or commands).
func (t *Tx) Execute(
	ctx context.Context,
	cmd string,
	args ...interface{},
) error {
	q, err := NewQuery(
		"Execute",
		cmd,
		args,
		t.Capabilities1pX(),
		t.state,
		nil,
		true,
		t.warningHandler,
	)
	if err != nil {
		return err
	}

	return t.scriptFlow(ctx, q)
}

// Query runs a query and returns the results.
func (t *Tx) Query(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	return RunQuery(
		ctx,
		t,
		"Query",
		cmd,
		out,
		args,
		t.state,
		t.warningHandler,
	)
}

// QuerySingle runs a singleton-returning query and returns its element.
// If the query executes successfully but doesn't return a result
// a NoDataError is returned. If the out argument is an optional type the out
// argument will be set to missing instead of returning a NoDataError.
func (t *Tx) QuerySingle(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	return RunQuery(
		ctx,
		t,
		"QuerySingle",
		cmd,
		out,
		args,
		t.state,
		t.warningHandler,
	)
}

// QueryJSON runs a query and return the results as JSON.
func (t *Tx) QueryJSON(
	ctx context.Context,
	cmd string,
	out *[]byte,
	args ...interface{},
) error {
	return RunQuery(
		ctx,
		t,
		"QueryJSON",
		cmd,
		out,
		args,
		t.state,
		t.warningHandler,
	)
}

// QuerySingleJSON runs a singleton-returning query.
// If the query executes successfully but doesn't have a result
// a NoDataError is returned.
func (t *Tx) QuerySingleJSON(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	return RunQuery(
		ctx,
		t,
		"QuerySingleJSON",
		cmd,
		out,
		args,
		t.state,
		t.warningHandler,
	)
}

// ExecuteSQL executes a SQL command (or commands).
func (t *Tx) ExecuteSQL(
	ctx context.Context,
	cmd string,
	args ...interface{},
) error {
	q, err := NewQuery(
		"ExecuteSQL",
		cmd,
		args,
		t.Capabilities1pX(),
		t.state,
		nil,
		true,
		t.warningHandler,
	)
	if err != nil {
		return err
	}

	return t.scriptFlow(ctx, q)
}

// QuerySQL runs a SQL query and returns the results.
func (t *Tx) QuerySQL(
	ctx context.Context,
	cmd string,
	out interface{},
	args ...interface{},
) error {
	return RunQuery(
		ctx,
		t,
		"QuerySQL",
		cmd,
		out,
		args,
		t.state,
		t.warningHandler,
	)
}
