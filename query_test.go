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
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/geltypes"
	types "github.com/geldata/gel-go/geltypes"
	gel "github.com/geldata/gel-go/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryCachingIncludesOutType(t *testing.T) {
	ctx := context.Background()

	err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		var result struct {
			Val OptionalTuple `gel:"val"`
		}
		e := tx.Execute(ctx, `
			CREATE TYPE Sample {
				CREATE PROPERTY val -> tuple<int64, int64>;
			};
		`)
		assert.NoError(t, e)

		// Run a query with a particular out type
		// that can later be run with a different out type.
		return tx.QuerySingle(ctx, `SELECT Sample { val } LIMIT 1`, &result)
	})
	assert.EqualError(t, err, "gel.NoDataError: zero results")

	err = client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		var result struct {
			Val OptionalNamedTuple `gel:"val"`
		}

		e := tx.Execute(ctx, `
			CREATE TYPE Sample {
				CREATE PROPERTY val -> tuple<a: int64, b: int64>;
			};
		`)
		assert.NoError(t, e)

		// Run the same query string again with a different out type.
		// There should not be any errors complaining about the out type.
		return tx.QuerySingle(ctx, `SELECT Sample { val } LIMIT 1`, &result)
	})
	assert.EqualError(t, err, "gel.NoDataError: zero results")
}

func TestObjectWithoutID(t *testing.T) {
	ctx := context.Background()

	type Function struct {
		Name string `gel:"name"`
	}

	var result Function
	err := client.QuerySingle(
		ctx, `
		SELECT schema::Function{ name }
		FILTER .name = 'std::str_trim'
		LIMIT 1`,
		&result,
	)
	assert.NoError(t, err)
	assert.Equal(t, "std::str_trim", result.Name)
}

func TestWrongNumberOfArguments(t *testing.T) {
	var result string
	ctx := context.Background()
	err := client.QuerySingle(ctx, `SELECT <str>$0`, &result)
	assert.EqualError(t, err,
		"gel.InvalidArgumentError: expected 1 arguments got 0")
}

func TestConnRejectsTransactions(t *testing.T) {
	expected := "gel.DisabledCapabilityError: " +
		"cannot execute transaction control commands.*"

	ctx := context.Background()
	err := client.Execute(ctx, "START TRANSACTION")
	assert.Regexp(t, expected, err)

	var result []byte
	err = client.Query(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)

	err = client.QueryJSON(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)

	err = client.QuerySingle(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)

	err = client.QuerySingleJSON(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)
}

func TestMissmatchedCardinality(t *testing.T) {
	ctx := context.Background()

	var result []int64
	err := client.QuerySingle(ctx, "SELECT {1, 2, 3}", &result)

	expected := "gel.ResultCardinalityMismatchError: " +
		"the query has cardinality AT_LEAST_ONE " +
		"which does not match the expected cardinality ONE"
	assert.EqualError(t, err, expected)
}

func TestMissmatchedResultType(t *testing.T) {
	type C struct { // nolint:unused
		z int // nolint:structcheck
	}

	type B struct { // nolint:unused
		y C // nolint:structcheck
	}

	type A struct {
		x B // nolint:structcheck,unused
	}

	var result A

	ctx := context.Background()
	err := client.QuerySingle(ctx, "SELECT (x := (y := (z := 7)))", &result)

	expected := "gel.InvalidArgumentError: " +
		"the \"out\" argument does not match query schema: " +
		"expected gel.A.x.y.z to be int64 or geltypes.OptionalInt64 got int"
	assert.EqualError(t, err, expected)
}

// The client should read all messages through ReadyForCommand
// before returning from a QueryX()
func TestParseAllMessagesAfterError(t *testing.T) {
	t.Skip()
	ctx := context.Background()

	// cause error during prepare
	var number float64
	err := client.QuerySingle(ctx, "SELECT 1 / <str>$0", &number, int64(5))

	// nolint:lll
	expected := `gel.InvalidTypeError: operator '/' cannot be applied to operands of type 'std::int64' and 'std::str'
query:1:8

SELECT 1 / <str>$0
       ^ Consider using an explicit type cast or a conversion function.`
	assert.EqualError(t, err, expected)

	// cause erroy during execute
	err = client.QuerySingle(ctx, "SELECT 1 / 0", &number)
	assert.EqualError(t, err, "gel.DivisionByZeroError: division by zero")

	// cache query so that it is run optimistically next time
	err = client.QuerySingle(ctx, "SELECT 1 / <int64>$0", &number, int64(3))
	assert.NoError(t, err)

	// cause error during optimistic execute
	err = client.QuerySingle(ctx, "SELECT 1 / <int64>$0", &number, int64(0))
	assert.EqualError(t, err, "gel.DivisionByZeroError: division by zero")
}

func TestArgumentTypeMissmatch(t *testing.T) {
	type Tuple struct {
		first int16 `gel:"0"` // nolint:unused,structcheck
	}

	var res Tuple
	ctx := context.Background()
	err := client.QuerySingle(ctx,
		"SELECT (<int16>$0 + <int16>$1,)", &res, 1, 1111)

	require.NotNil(t, err)
	assert.EqualError(t, err,
		"gel.InvalidArgumentError: expected args[0] to be int16, "+
			"geltypes.OptionalInt16 or Int16Marshaler got int")
}

func TestNamedQueryArguments(t *testing.T) {
	ctx := context.Background()
	var result [][]int64
	err := client.Query(
		ctx,
		"SELECT [<int64>$first, <int64>$second]",
		&result,
		map[string]interface{}{
			"first":  int64(5),
			"second": int64(8),
		},
	)

	require.NoError(t, err)
	assert.Equal(t, [][]int64{{5, 8}}, result)
}

func TestNumberedQueryArguments(t *testing.T) {
	ctx := context.Background()
	result := [][]int64{}
	err := client.Query(
		ctx,
		"SELECT [<int64>$0, <int64>$1]",
		&result,
		int64(5),
		int64(8),
	)

	assert.NoError(t, err)
	assert.Equal(t, [][]int64{{5, 8}}, result)
}

func TestQueryJSON(t *testing.T) {
	ctx := context.Background()
	var result []byte
	err := client.QueryJSON(
		ctx,
		"SELECT {(a := 0, b := <int64>$0), (a := 42, b := <int64>$1)}",
		&result,
		int64(1),
		int64(2),
	)

	// casting to string makes error message more helpful
	// when this test fails
	actual := string(result)

	require.NoError(t, err)
	assert.Equal(
		t,
		"[{\"a\" : 0, \"b\" : 1}, {\"a\" : 42, \"b\" : 2}]",
		actual,
	)
}

func TestQuerySingleJSON(t *testing.T) {
	ctx := context.Background()
	var result []byte
	err := client.QuerySingleJSON(
		ctx,
		"SELECT (a := 0, b := <int64>$0)",
		&result,
		int64(42),
	)

	// casting to string makes error messages more helpful
	// when this test fails
	actual := string(result)

	assert.NoError(t, err)
	assert.Equal(t, "{\"a\" : 0, \"b\" : 42}", actual)
}

func TestQuerySingleJSONZeroResults(t *testing.T) {
	ctx := context.Background()
	var result []byte
	err := client.QuerySingleJSON(ctx, "SELECT <int64>{}", &result)

	require.Equal(t, err, gel.ErrZeroResults)
	assert.Equal(t, []byte(nil), result)
}

func TestQuerySingle(t *testing.T) {
	ctx := context.Background()
	var result int64
	err := client.QuerySingle(ctx, "SELECT 42", &result)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), result)
}

func TestQuerySingleZeroResults(t *testing.T) {
	ctx := context.Background()
	var result int64
	err := client.QuerySingle(ctx, "SELECT <int64>{}", &result)

	assert.Equal(t, gel.ErrZeroResults, err)
}

func TestQuerySingleNestedSlice(t *testing.T) {
	ctx := context.Background()
	type IDField struct {
		ID types.UUID `gel:"id"`
	}
	type NameField struct {
		Name types.OptionalStr `gel:"name"`
	}
	type UserModel struct {
		IDField   `gel:"$inline"`
		NameField `gel:"$inline"`
	}
	type UsersField struct {
		Users []UserModel `gel:"users"`
	}
	type ViewModel struct {
		UsersField `gel:"$inline"`
	}
	result := ViewModel{}
	err := client.QuerySingle(
		ctx,
		`
with a := (INSERT User { name := 'a' }), b := (INSERT User { name := 'b' })
SELECT { users := (SELECT { a, b } { id, name }) }`,
		&result,
	)
	require.NoError(t, err)

	require.Equal(
		t, 2, len(result.Users),
		"wrong number of users, expected 2 got %v", len(result.Users))
	assert.NotEqual(t, result.Users[0].ID, types.UUID{})
	a, _ := result.Users[0].Name.Get()
	assert.Equal(t, a, "a")

	assert.NotEqual(t, result.Users[1].ID, types.UUID{})
	assert.NotEqual(t, result.Users[0].ID, result.Users[1].ID)
	b, _ := result.Users[1].Name.Get()
	assert.Equal(t, b, "b")
}

func TestError(t *testing.T) {
	ctx := context.Background()
	err := client.Execute(ctx, "malformed query;")
	assert.EqualError(
		t,
		err,
		`gel.EdgeQLSyntaxError: Unexpected 'malformed'
query:1:1

malformed query;
^ error`,
	)

	var expected gelerr.Error
	assert.True(t, errors.As(err, &expected))
}

func TestQueryTimesOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	defer cancel()

	var r int64
	err := client.QuerySingle(ctx, "SELECT 1;", &r)
	require.True(
		t,
		errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, os.ErrDeadlineExceeded),
		err,
	)
	require.Equal(t, int64(0), r)

	err = client.QuerySingle(context.Background(), "SELECT 2;", &r)
	require.NoError(t, err)
	assert.Equal(t, int64(2), r)
}

func TestNilResultValue(t *testing.T) {
	ctx := context.Background()
	err := client.Query(ctx, "SELECT 1", nil)
	assert.EqualError(t, err, "gel.InterfaceError: "+
		"the \"out\" argument must be a pointer, got untyped nil")
}

func TestExecutWithArgs(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	err := client.Execute(ctx, "select <int64>$0; select <int64>$0;", int64(1))
	assert.NoError(t, err)

	err = client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		err = tx.Execute(ctx, "select <int64>$0; select <int64>$0;", int64(1))
		assert.NoError(t, err)

		return nil
	})
	assert.NoError(t, err)
}

func TestClientRejectsSessionConfig(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	expected := "gel.DisabledCapabilityError: " +
		"cannot execute session configuration queries.*"

	ctx := context.Background()
	err := client.Execute(ctx, "SET ALIAS bar AS MODULE std")
	assert.Regexp(t, expected, err)

	var result []byte
	err = client.Query(ctx, "SET ALIAS bar AS MODULE std", &result)
	assert.Regexp(t, expected, err)

	err = client.QueryJSON(ctx, "SET ALIAS bar AS MODULE std", &result)
	assert.Regexp(t, expected, err)

	err = client.QuerySingle(ctx, "SET ALIAS bar AS MODULE std", &result)
	assert.Regexp(t, expected, err)

	err = client.QuerySingleJSON(ctx, "SET ALIAS bar AS MODULE std", &result)
	assert.Regexp(t, expected, err)
}

func TestWithConfig(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	query := "SELECT assert_single(cfg::Config.query_execution_timeout)"

	var result types.Duration
	err := client.QuerySingle(ctx, query, &result)
	require.NoError(t, err)
	assert.Equal(t, types.Duration(0), result)

	a := client.WithConfig(map[string]interface{}{
		"query_execution_timeout": types.Duration(65_432_000),
	})
	err = a.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(65_432_000), result)

	err = client.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(0), result)

	b := a.WithConfig(map[string]interface{}{
		"query_execution_timeout": types.Duration(32_100_000),
	})
	err = b.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(32_100_000), result)

	err = a.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(65_432_000), result)

	err = client.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(0), result)
}

func TestInvalidWithConfig(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result int64

	err := client.QuerySingle(ctx, "select 1", &result)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result)

	c := client.WithConfig(map[string]interface{}{"hello": "world"})
	err = c.QuerySingle(ctx, "select 1", &result)
	assert.EqualError(t, err, "gel.BinaryProtocolError: "+
		"invalid connection state: "+
		"found unknown state value state.config.hello")

	err = client.QuerySingle(ctx, "select 1", &result)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result)

	c = client.WithConfig(map[string]interface{}{
		"query_execution_timeout": "this should be Duration not string",
	})
	err = c.QuerySingle(ctx, "select 1", &result)
	assert.EqualError(t, err, "gel.BinaryProtocolError: "+
		"invalid connection state: "+
		"expected state.config.query_execution_timeout to be "+
		"geltypes.Duration, geltypes.OptionalDuration or DurationMarshaler "+
		"got string")

	err = client.QuerySingle(ctx, "select 1", &result)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result)
}

func TestDefaultIsolation(t *testing.T) {
	skipIfServerVersionLT(t, 6, 0)

	{
		c := client.WithConfig(map[string]any{
			"default_transaction_isolation": gelcfg.RepeatableRead,
		})

		var result []int64
		err := c.Query(context.Background(), `SELECT 1; SELECT 2;`, &result)
		require.NoError(t, err)
		require.Equal(t, []int64{2}, result)
	}

	{
		txOps := gelcfg.NewTxOptions().
			WithIsolation(gelcfg.RepeatableRead)
		c := client.WithTxOptions(txOps)

		query := `select sys::get_transaction_isolation()`
		var result string
		err := c.QuerySingle(context.Background(), query, &result)
		require.NoError(t, err)
		require.Equal(t, "RepeatableRead", result)
	}

	{
		txOps := gelcfg.NewTxOptions().
			WithIsolation(gelcfg.PreferRepeatableRead)
		c := client.WithTxOptions(txOps)

		query := `select sys::get_transaction_isolation()`
		var result string
		err := c.QuerySingle(context.Background(), query, &result)
		require.NoError(t, err)
		require.Equal(t, string(gelcfg.RepeatableRead), result)
	}

	t.Run(
		"TxOptions.IsolationLevel takes precedence over config",
		func(t *testing.T) {
			{
				c := client.WithConfig(map[string]any{
					"default_transaction_isolation": gelcfg.RepeatableRead,
				}).WithTxOptions(
					gelcfg.NewTxOptions().WithIsolation(gelcfg.Serializable),
				)

				query := `select sys::get_transaction_isolation()`
				var result string
				err := c.QuerySingle(context.Background(), query, &result)
				require.NoError(t, err)
				require.Equal(t, string(gelcfg.Serializable), result)
			}
		},
	)

	t.Run(
		"TxOptions.ReadOnly takes precedence over config",
		func(t *testing.T) {
			{
				c := client.WithConfig(map[string]any{
					"default_transaction_access_mode": "ReadOnly",
				}).WithTxOptions(
					gelcfg.NewTxOptions().WithReadOnly(false),
				)

				query := `SELECT (INSERT User { name := "test" }).name`
				var result string
				err := c.QuerySingle(context.Background(), query, &result)
				require.NoError(t, err)
				require.Equal(t, "test", result)
			}
		},
	)

	{
		txOps := gelcfg.NewTxOptions().WithReadOnly(true)
		c := client.WithTxOptions(txOps)

		query := `SELECT (INSERT User { name := "test" }).name`
		var result string
		err := c.QuerySingle(context.Background(), query, &result)

		var edbErr gelerr.Error
		require.Error(t, err)
		require.True(t, errors.As(err, &edbErr))
		require.True(t, edbErr.Category(gelerr.TransactionError), err)
	}
}

func conflictTypesInDB(t *testing.T) {
	ctx := context.Background()

	err := client.Execute(ctx, `
		CREATE MODULE test;
		CREATE TYPE test::Tmp {
			CREATE REQUIRED PROPERTY tmp -> std::str;
		};
		CREATE TYPE test::TmpConflict {
			CREATE REQUIRED PROPERTY tmp -> std::str {
				CREATE CONSTRAINT exclusive;
			}
		};
		CREATE TYPE test::TmpConflictChild extending test::TmpConflict;
	`)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := client.Execute(ctx, `
			DROP TYPE test::TmpConflictChild;
			DROP TYPE test::TmpConflict;
			DROP TYPE test::Tmp;
			DROP MODULE test;
		`)
		assert.NoError(t, err)
	})
}

func TestPreferRepeatableRead(t *testing.T) {
	skipIfServerVersionLT(t, 6, 5)
	conflictTypesInDB(t)
	ctx := context.Background()

	{
		c := client.WithConfig(map[string]any{
			"default_transaction_isolation": gelcfg.PreferRepeatableRead,
		})

		var result struct {
			Ins struct {
				ID types.UUID `gel:"id"`
			} `gel:"ins"`
			Level string `gel:"level"`
		}

		// This query can run in RepeatableRead mode
		query := `
			select {
				ins := (insert test::Tmp { tmp := "test" }),
				level := sys::get_transaction_isolation(),
			}
		`
		err := c.QuerySingle(ctx, query, &result)
		require.NoError(t, err)
		require.Equal(t, string(gelcfg.RepeatableRead), result.Level)

		// This one can cannot.
		query = `
			select {
				ins := (insert test::TmpConflict { tmp := "test" }),
				level := sys::get_transaction_isolation(),
			}
		`
		err = c.QuerySingle(ctx, query, &result)
		require.NoError(t, err)
		require.Equal(t, string(gelcfg.Serializable), result.Level)
	}

	{
		var result struct {
			Level string `gel:"level"`
		}

		query := `
			select {
				level := sys::get_transaction_isolation(),
			}
		`
		err := client.WithTxOptions(gelcfg.NewTxOptions().
			WithIsolation(gelcfg.PreferRepeatableRead)).
			QuerySingle(ctx, query, &result)
		require.NoError(t, err)
		require.Equal(t, string(gelcfg.RepeatableRead), result.Level)
	}
}

func TestWithoutConfig(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	query := "SELECT assert_single(cfg::Config.query_execution_timeout)"

	var result types.Duration
	err := client.QuerySingle(ctx, query, &result)
	require.NoError(t, err)
	assert.Equal(t, types.Duration(0), result)

	a := client.WithConfig(map[string]interface{}{
		"query_execution_timeout": types.Duration(65_432_000),
	})
	err = a.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(65_432_000), result)

	b := a.WithoutConfig("query_execution_timeout")
	err = b.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(0), result)

	err = a.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(65_432_000), result)

	err = client.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(0), result)

	b = client.WithoutConfig("some", "crazy", "names")
	err = b.QuerySingle(ctx, query, &result)
	assert.NoError(t, err)
	assert.Equal(t, types.Duration(0), result)
}

func TestWithModuleAliases(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result int64

	err := client.QuerySingle(
		ctx,
		"SELECT <my_new_name_for_std::int64>1",
		&result,
	)
	assert.EqualError(t, err, "gel.InvalidReferenceError: "+
		"type 'my_new_name_for_std::int64' does not exist\n"+
		"query:1:9\n\n"+
		"SELECT <my_new_name_for_std::int64>1\n"+
		"        ^ error")

	a := client.WithModuleAliases(gelcfg.ModuleAlias{
		Alias:  "my_new_name_for_std",
		Module: "std",
	})

	err = a.QuerySingle(ctx, "SELECT <my_new_name_for_std::int64>2", &result)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), result)

	err = client.QuerySingle(
		ctx,
		"SELECT <my_new_name_for_std::int64>3",
		&result,
	)
	assert.EqualError(t, err, "gel.InvalidReferenceError: "+
		"type 'my_new_name_for_std::int64' does not exist\n"+
		"query:1:9\n\n"+
		"SELECT <my_new_name_for_std::int64>3\n"+
		"        ^ error")

	b := a.WithModuleAliases(gelcfg.ModuleAlias{
		Alias:  "my_new_name_for_std",
		Module: "math",
	})
	err = b.QuerySingle(ctx, "SELECT <my_new_name_for_std::int64>4", &result)
	assert.EqualError(t, err, "gel.InvalidReferenceError: "+
		"type 'my_new_name_for_std::int64' does not exist\n"+
		"query:1:9\n\n"+
		"SELECT <my_new_name_for_std::int64>4\n"+
		"        ^ error")

	err = a.QuerySingle(ctx, "SELECT <my_new_name_for_std::int64>5", &result)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), result)

	err = client.QuerySingle(
		ctx,
		"SELECT <my_new_name_for_std::int64>6",
		&result,
	)
	assert.EqualError(t, err, "gel.InvalidReferenceError: "+
		"type 'my_new_name_for_std::int64' does not exist\n"+
		"query:1:9\n\n"+
		"SELECT <my_new_name_for_std::int64>6\n"+
		"        ^ error")
}

func TestInvalidWithModuleAliases(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result int64

	a := client.WithModuleAliases(gelcfg.ModuleAlias{
		Alias:  "my_alias",
		Module: "this_module_doesnt_exist",
	})

	err := a.QuerySingle(ctx, "SELECT <my_alias::int64>1", &result)
	assert.EqualError(t, err, "gel.InvalidReferenceError: "+
		"type 'my_alias::int64' does not exist\n"+
		"query:1:9\n\n"+
		"SELECT <my_alias::int64>1\n"+
		"        ^ error")
}

func TestWithoutModuleAliases(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result int64

	a := client.WithModuleAliases(gelcfg.ModuleAlias{
		Alias:  "my_new_name_for_std",
		Module: "std",
	})
	b := a.WithoutModuleAliases("my_new_name_for_std")

	err := b.QuerySingle(ctx, "SELECT <my_new_name_for_std::int64>4", &result)
	assert.EqualError(t, err, "gel.InvalidReferenceError: "+
		"type 'my_new_name_for_std::int64' does not exist\n"+
		"query:1:9\n\n"+
		"SELECT <my_new_name_for_std::int64>4\n"+
		"        ^ error")

	err = a.QuerySingle(ctx, "SELECT <my_new_name_for_std::int64>5", &result)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), result)

	err = client.QuerySingle(
		ctx,
		"SELECT <my_new_name_for_std::int64>6",
		&result,
	)
	assert.EqualError(t, err, "gel.InvalidReferenceError: "+
		"type 'my_new_name_for_std::int64' does not exist\n"+
		"query:1:9\n\n"+
		"SELECT <my_new_name_for_std::int64>6\n"+
		"        ^ error")
}

func TestWithGlobals(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result string

	err := client.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "default", result)

	a := client.WithGlobals(map[string]interface{}{
		"default::global_str": "first",
	})
	err = a.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "first", result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "default", result)

	b := a.WithGlobals(map[string]interface{}{
		"default::global_str": "second",
	})
	err = b.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "second", result)

	err = a.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "first", result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "default", result)
}

func TestWithGlobalUUID(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var id types.UUID
	err := client.
		WithGlobals(map[string]interface{}{"global_id": types.UUID{1, 2, 3}}).
		QuerySingle(ctx, "SELECT GLOBAL global_id", &id)
	require.NoError(t, err)
	assert.Equal(t, types.UUID{1, 2, 3}, id)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_id", &id)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalBytes(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var bytes []byte
	err := client.
		WithGlobals(map[string]interface{}{"global_bytes": []byte{1, 2, 3}}).
		QuerySingle(ctx, "SELECT GLOBAL global_bytes", &bytes)
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, bytes)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_bytes", &bytes)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalInt16(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var val int16
	err := client.
		WithGlobals(map[string]interface{}{"global_int16": int16(7)}).
		QuerySingle(ctx, "SELECT GLOBAL global_int16", &val)
	require.NoError(t, err)
	assert.Equal(t, int16(7), val)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_int16", &val)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalInt32(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var val int32
	err := client.
		WithGlobals(map[string]interface{}{"global_int32": int32(7)}).
		QuerySingle(ctx, "SELECT GLOBAL global_int32", &val)
	require.NoError(t, err)
	assert.Equal(t, int32(7), val)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_int32", &val)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalInt64(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var val int64
	err := client.
		WithGlobals(map[string]interface{}{"global_int64": int64(7)}).
		QuerySingle(ctx, "SELECT GLOBAL global_int64", &val)
	require.NoError(t, err)
	assert.Equal(t, int64(7), val)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_int64", &val)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalFloat32(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var val float32
	err := client.
		WithGlobals(map[string]interface{}{"global_float32": float32(7)}).
		QuerySingle(ctx, "SELECT GLOBAL global_float32", &val)
	require.NoError(t, err)
	assert.Equal(t, float32(7), val)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_float32", &val)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalFloat64(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result float64
	err := client.
		WithGlobals(map[string]interface{}{"global_float64": float64(7)}).
		QuerySingle(ctx, "SELECT GLOBAL global_float64", &result)
	require.NoError(t, err)
	assert.Equal(t, float64(7), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_float64", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalBool(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result bool
	err := client.
		WithGlobals(map[string]interface{}{"global_bool": true}).
		QuerySingle(ctx, "SELECT GLOBAL global_bool", &result)
	require.NoError(t, err)
	assert.True(t, result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_bool", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalDateTime(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result time.Time
	err := client.
		WithGlobals(map[string]interface{}{
			"global_datetime": time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		}).
		QuerySingle(ctx, "SELECT GLOBAL global_datetime", &result)
	require.NoError(t, err)
	assert.Equal(t, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_datetime", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalDuration(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result types.Duration
	err := client.
		WithGlobals(
			map[string]interface{}{"global_duration": types.Duration(7)},
		).
		QuerySingle(ctx, "SELECT GLOBAL global_duration", &result)
	require.NoError(t, err)
	assert.Equal(t, types.Duration(7), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_duration", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalJSON(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result []byte
	err := client.
		WithGlobals(map[string]interface{}{"global_json": []byte("7")}).
		QuerySingle(ctx, "SELECT GLOBAL global_json", &result)
	require.NoError(t, err)
	assert.Equal(t, []byte("7"), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_json", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalLocalDateTime(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result types.LocalDateTime
	err := client.
		WithGlobals(map[string]interface{}{
			"global_local_datetime": types.NewLocalDateTime(
				1970,
				1,
				1,
				0,
				0,
				0,
				0,
			),
		}).
		QuerySingle(ctx, "SELECT GLOBAL global_local_datetime", &result)
	require.NoError(t, err)
	assert.Equal(t, types.NewLocalDateTime(1970, 1, 1, 0, 0, 0, 0), result)

	err = client.QuerySingle(
		ctx,
		"SELECT GLOBAL global_local_datetime",
		&result,
	)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalLocalDate(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result types.LocalDate
	err := client.
		WithGlobals(map[string]interface{}{
			"global_local_date": types.NewLocalDate(1970, 1, 1),
		}).
		QuerySingle(ctx, "SELECT GLOBAL global_local_date", &result)
	require.NoError(t, err)
	assert.Equal(t, types.NewLocalDate(1970, 1, 1), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_local_date", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalLocalTime(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result types.LocalTime
	err := client.
		WithGlobals(map[string]interface{}{
			"global_local_time": types.NewLocalTime(1, 2, 3, 4),
		}).
		QuerySingle(ctx, "SELECT GLOBAL global_local_time", &result)
	require.NoError(t, err)
	assert.Equal(t, types.NewLocalTime(1, 2, 3, 4), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_local_time", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalBigInt(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result *big.Int
	err := client.
		WithGlobals(map[string]interface{}{"global_bigint": big.NewInt(7)}).
		QuerySingle(ctx, "SELECT GLOBAL global_bigint", &result)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(7), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_bigint", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalRelativeDuration(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result types.RelativeDuration
	err := client.
		WithGlobals(map[string]interface{}{
			"global_relative_duration": types.NewRelativeDuration(1, 2, 3),
		}).
		QuerySingle(ctx, "SELECT GLOBAL global_relative_duration", &result)
	require.NoError(t, err)
	assert.Equal(t, types.NewRelativeDuration(1, 2, 3), result)

	err = client.QuerySingle(
		ctx,
		"SELECT GLOBAL global_relative_duration",
		&result,
	)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalDateDuration(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result types.DateDuration
	err := client.
		WithGlobals(map[string]interface{}{
			"global_date_duration": types.NewDateDuration(1, 2),
		}).
		QuerySingle(ctx, "SELECT GLOBAL global_date_duration", &result)
	require.NoError(t, err)
	assert.Equal(t, types.NewDateDuration(1, 2), result)

	err = client.QuerySingle(
		ctx,
		"SELECT GLOBAL global_date_duration",
		&result,
	)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestWithGlobalMemory(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result types.Memory
	err := client.
		WithGlobals(map[string]interface{}{"global_memory": types.Memory(7)}).
		QuerySingle(ctx, "SELECT GLOBAL global_memory", &result)
	require.NoError(t, err)
	assert.Equal(t, types.Memory(7), result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_memory", &result)
	require.Equal(t, gel.ErrZeroResults, err)
}

func TestInvalidWithGlobals(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result int64

	a := client.WithGlobals(map[string]interface{}{
		"default::this": "thing donesnt exist",
	})

	err := a.QuerySingle(ctx, "SELECT GLOBAL this", &result)
	assert.EqualError(t, err, "gel.BinaryProtocolError: "+
		"invalid connection state: "+
		"found unknown state value state.globals.default::this")

	b := client.WithGlobals(map[string]interface{}{
		"default::global_str": 27,
	})

	err = b.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	assert.EqualError(t, err, "gel.BinaryProtocolError: "+
		"invalid connection state: "+
		"expected state.globals.default::global_str to be "+
		"string, geltypes.OptionalStr or StrMarshaler "+
		"got int")
}

func TestWithoutGlobals(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()
	var result string

	err := client.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "default", result)

	a := client.WithGlobals(map[string]interface{}{
		"default::global_str": "first",
	})

	b := a.WithoutGlobals("default::global_str")
	err = b.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "default", result)

	err = a.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "first", result)

	err = client.QuerySingle(ctx, "SELECT GLOBAL global_str", &result)
	require.NoError(t, err)
	assert.Equal(t, "default", result)
}

func TestWithConfigWrongServerVersion(t *testing.T) {
	if protocolVersion.GTE(gel.ProtocolVersion1p0) {
		t.Skip()
	}
	ctx := context.Background()
	var result int64

	a := client.WithGlobals(map[string]interface{}{
		"default::global_str": "first",
	})
	err := a.QuerySingle(ctx, "SELECT 1", &result)
	require.EqualError(t, err, "gel.InterfaceError: "+
		"client methods WithConfig, WithGlobals, and WithModuleAliases "+
		"are not supported by the server. "+
		"Upgrade your server to version 2.0 or greater to use these features.")

	b := client.WithModuleAliases(gelcfg.ModuleAlias{
		Alias:  "other_math",
		Module: "math",
	})
	err = b.QuerySingle(ctx, "SELECT 1", &result)
	require.EqualError(t, err, "gel.InterfaceError: "+
		"client methods WithConfig, WithGlobals, and WithModuleAliases "+
		"are not supported by the server. "+
		"Upgrade your server to version 2.0 or greater to use these features.")

	c := client.WithConfig(map[string]interface{}{
		"query_execution_timeout": types.Duration(65_432_000),
	})
	err = c.QuerySingle(ctx, "SELECT 1", &result)
	require.EqualError(t, err, "gel.InterfaceError: "+
		"client methods WithConfig, WithGlobals, and WithModuleAliases "+
		"are not supported by the server. "+
		"Upgrade your server to version 2.0 or greater to use these features.")
}

func TestWithWarningHandler(t *testing.T) {
	var hasWarnOnCall bool
	ctx := context.Background()
	err := client.QuerySingle(
		ctx,
		`
		SELECT EXISTS (
			SELECT schema::Function { id }
			FILTER .name = 'std::_warn_on_call'
		)
		`,
		&hasWarnOnCall,
	)
	require.NoError(t, err)

	if !hasWarnOnCall {
		t.Skip("missing std::_warn_on_call")
	}

	// The default warning handler logs warnings. Make sure we never log
	// warnings in this test so that there isn't noise in the test run output.
	a := client.WithWarningHandler(func(_ []error) error { return nil })

	seen := []error{}
	b := a.WithWarningHandler(func(warnings []error) error {
		seen = append(seen, warnings...)
		return nil
	})

	// The server sends warnings in response to both parse and execute.  Use a
	// query that we know is not cached so that we can assert exactly how many
	// warnings are expected.
	name := randomName()
	query := fmt.Sprintf(`SELECT (%s := _warn_on_call()).%s`, name, name)

	err = b.Execute(ctx, query)
	require.NoError(t, err)
	require.Equal(t, 2, len(seen))

	var resultMany []int64
	seen = []error{}
	err = b.Query(ctx, query, &resultMany)
	require.NoError(t, err)
	require.Equal(t, 2, len(seen))

	var resultJSON []byte
	seen = []error{}
	err = b.QueryJSON(ctx, query, &resultJSON)
	require.NoError(t, err)
	require.Equal(t, 2, len(seen))

	var resultSingle int64
	seen = []error{}
	err = b.QuerySingle(ctx, query, &resultSingle)
	require.NoError(t, err)
	require.Equal(t, 2, len(seen))

	seen = []error{}
	err = b.QuerySingleJSON(ctx, query, &resultJSON)
	require.NoError(t, err)
	require.Equal(t, 2, len(seen))

	// Assert that the client config was not changed.
	seen = []error{}
	err = a.QuerySingle(ctx, query, &resultSingle)
	require.NoError(t, err)
	require.Equal(t, 0, len(seen))
}

func TestWithQueryOptionsReadonly(t *testing.T) {
	ctx := context.Background()

	err := client.Execute(ctx, `create type QueryOptsTest {
		create property name -> str;
	};`)
	assert.NoError(t, err)
	defer func() {
		err = client.Execute(ctx, `drop type QueryOptsTest;`)
		assert.NoError(t, err)
	}()

	var res struct {
		ID types.UUID `edgedb:"id"`
	}
	err = client.QuerySingle(ctx,
		"insert QueryOptsTest {name := 'abc'}", &res)
	assert.NoError(t, err)

	readOnly := gelcfg.NewQueryOptions().WithReadOnly(true)
	readonlyClient := client.WithQueryOptions(readOnly)

	err = readonlyClient.QuerySingle(ctx,
		"insert QueryOptsTest {name := 'def'}", &res)
	assert.EqualError(t, err,
		"gel.DisabledCapabilityError: "+
			"cannot execute data modification queries: disabled by the client",
	)

	err = client.QuerySingle(ctx,
		"select QueryOptsTest {id} limit 1", &res)
	assert.NoError(t, err)

	// check we didn't modify the original client
	err = client.QuerySingle(ctx,
		"insert QueryOptsTest {name := 'abc'}", &res)
	assert.NoError(t, err)
}

func TestWithQueryOptionsImplicitLimit(t *testing.T) {
	ctx := context.Background()

	var res []struct {
		Name string `edgedb:"name"`
	}
	err := client.Query(ctx,
		"select schema::ObjectType {name}", &res)
	assert.NoError(t, err)
	assert.Greater(t, len(res), 10)

	limit := gelcfg.NewQueryOptions().WithImplicitLimit(10)
	limitClient := client.WithQueryOptions(limit)

	var limitRes []struct {
		Name string `edgedb:"name"`
	}
	err = limitClient.Query(ctx,
		"select schema::ObjectType {name}", &limitRes)
	assert.NoError(t, err)
	assert.Equal(t, len(limitRes), 10)

	// check we didn't modify the original client
	var res2 []struct {
		Name string `edgedb:"name"`
	}
	err = client.Query(ctx,
		"select schema::ObjectType {name}", &res2)
	assert.NoError(t, err)
	assert.Equal(t, len(res2), len(res))
}
