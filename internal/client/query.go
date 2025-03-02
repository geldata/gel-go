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
	"errors"
	"fmt"
	"reflect"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	types "github.com/geldata/gel-go/geltypes"
	gelerrint "github.com/geldata/gel-go/internal/gelerr"
	"github.com/geldata/gel-go/internal/introspect"
)

type query struct {
	out          reflect.Value
	outType      reflect.Type
	method       string
	lang         Language
	cmd          string
	fmt          Format
	expCard      Cardinality
	args         []interface{}
	capabilities uint64
	state        map[string]interface{}
	parse        bool
	cfg          QueryConfig
}

func (q *query) getCapabilities() uint64 {
	capabilities := q.capabilities
	if q.cfg.QueryOptions.ReadOnly() {
		capabilities &^= capabilitiesModifications
	}
	return capabilities
}

func (q *query) flat() bool {
	if q.expCard != Many {
		return true
	}

	if q.fmt == JSON {
		return true
	}

	return false
}

// NewQuery returns a new granular flow query.
func NewQuery(
	method, cmd string,
	args []interface{},
	capabilities uint64,
	state map[string]interface{},
	out interface{},
	parse bool,
	cfg *QueryConfig,
) (*query, error) { // nolint:revive
	var (
		expCard Cardinality
		frmt    Format
	)

	lang := EdgeQL

	switch method {
	case "Execute", "ExecuteSQL":
		if method == "ExecuteSQL" {
			lang = SQL
		}
		return &query{
			method:       method,
			lang:         lang,
			cmd:          cmd,
			fmt:          Null,
			expCard:      Many,
			args:         args,
			capabilities: capabilities,
			state:        state,
			cfg:          *cfg,
			parse:        parse,
		}, nil
	case "Query":
		expCard = Many
		frmt = Binary
	case "QuerySingle":
		expCard = AtMostOne
		frmt = Binary
	case "QueryJSON":
		expCard = Many
		frmt = JSON
	case "QuerySingleJSON":
		expCard = AtMostOne
		frmt = JSON
	case "QuerySQL":
		lang = SQL
		expCard = Many
		frmt = Binary
	default:
		return nil, fmt.Errorf("unknown query method %q", method)
	}

	q := query{
		method:       method,
		lang:         lang,
		cmd:          cmd,
		fmt:          frmt,
		expCard:      expCard,
		args:         args,
		capabilities: capabilities,
		state:        state,
		cfg:          *cfg,
		parse:        parse,
	}

	var err error

	if frmt == JSON || expCard == AtMostOne {
		q.out, err = introspect.ValueOf(out)
	} else {
		q.out, err = introspect.ValueOfSlice(out)
		if err == nil {
			q.out.SetLen(0)
		}
	}

	if err != nil {
		return &query{}, gelerrint.NewInterfaceError("", err)
	}

	q.outType = q.out.Type()
	if !q.flat() {
		q.outType = q.outType.Elem()
	}

	return &q, nil
}

type queryable interface {
	Capabilities1pX() uint64
	granularFlow(context.Context, *query) error
}

type unseter interface {
	Unset()
}

// QueryConfig controls how a query is executed
type QueryConfig struct {
	WarningHandler gelcfg.WarningHandler
	QueryOptions   gelcfg.QueryOptions
	TxOptions      gelcfg.TxOptions
	RetryOptions   gelcfg.RetryOptions
	Annotations    map[string]string
}

// RunQuery runs a query.
func RunQuery(
	ctx context.Context,
	c queryable,
	method, cmd string,
	out interface{},
	args []interface{},
	state map[string]interface{},
	cfg *QueryConfig,
) error {
	if method == "QuerySingleJSON" {
		switch out.(type) {
		case *[]byte, *types.OptionalBytes:
		default:
			return gelerrint.NewInterfaceError(fmt.Sprintf(
				`the "out" argument must be *[]byte or *OptionalBytes, got %T`,
				out), nil)
		}
	}

	q, err := NewQuery(
		method,
		cmd,
		args,
		c.Capabilities1pX(),
		state,
		out,
		true,
		cfg,
	)
	if err != nil {
		return err
	}

	err = c.granularFlow(ctx, q)

	var edbErr gelerr.Error
	if errors.As(err, &edbErr) &&
		edbErr.Category(gelerr.NoDataError) &&
		(q.method == "QuerySingle" || q.method == "QuerySingleJSON") {
		if opt, ok := out.(unseter); ok {
			opt.Unset()
			return nil
		}
	}

	return err
}

// CopyState makes a copy of the state.
func CopyState(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for k, v := range in {
		switch val := v.(type) {
		case map[string]interface{}:
			out[k] = CopyState(val)
		case []interface{}:
			out[k] = copyStateSlice(val)
		default:
			out[k] = val
		}
	}

	return out
}

func copyStateSlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))

	for i, v := range in {
		switch val := v.(type) {
		case map[string]interface{}:
			out[i] = CopyState(val)
		case []interface{}:
			out[i] = copyStateSlice(val)
		default:
			out[i] = val
		}
	}

	return out
}
