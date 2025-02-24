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
	"math/rand"
	"strings"
	"testing"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/geltypes"
	types "github.com/geldata/gel-go/geltypes"
	gel "github.com/geldata/gel-go/internal/client"
	"github.com/geldata/gel-go/internal/snc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var rnd = snc.NewRand()

func TestTxRollesBack(t *testing.T) {
	ctx := context.Background()
	err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		query := "INSERT TxTest {name := 'Test Roll Back'};"
		if e := tx.Execute(ctx, query); e != nil {
			return e
		}

		return tx.Execute(ctx, "SELECT 1 / 0;")
	})

	var edbErr gelerr.Error
	require.True(t, errors.As(err, &edbErr), "wrong error: %v", err)
	require.True(
		t,
		edbErr.Category(gelerr.DivisionByZeroError),
		"wrong error: %v",
		err,
	)

	query := `
		SELECT (
			SELECT TxTest {name}
			FILTER .name = 'Test Roll Back'
		).name
		LIMIT 1
	`

	var testNames []string
	err = client.Query(ctx, query, &testNames)

	require.NoError(t, err)
	require.Equal(t, 0, len(testNames), "The transaction wasn't rolled back")
}

func TestTxRollesBackOnUserError(t *testing.T) {
	ctx := context.Background()
	err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		query := "INSERT TxTest {name := 'Test Roll Back'};"
		if e := tx.Execute(ctx, query); e != nil {
			return e
		}

		return errors.New("user defined error")
	})

	require.Equal(t, err, errors.New("user defined error"))

	query := `
		SELECT (
			SELECT TxTest {name}
			FILTER .name = 'Test Roll Back'
		).name
		LIMIT 1
	`

	var testNames []string
	err = client.Query(ctx, query, &testNames)

	require.NoError(t, err)
	require.Equal(t, 0, len(testNames), "The transaction wasn't rolled back")
}

func TestTxCommits(t *testing.T) {
	ctx := context.Background()
	err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		return tx.Execute(ctx, "INSERT TxTest {name := 'Test Commit'};")
	})
	require.NoError(t, err)

	query := `
		SELECT (
			SELECT TxTest {name}
			FILTER .name = 'Test Commit'
		).name
		LIMIT 1
	`

	var testNames []string
	err = client.Query(ctx, query, &testNames)

	require.NoError(t, err)
	require.Equal(
		t,
		[]string{"Test Commit"},
		testNames,
		"The transaction wasn't commited",
	)
}

func newTxOpts(
	level gelcfg.IsolationLevel,
	readOnly,
	deferrable bool,
) gelcfg.TxOptions {
	return gelcfg.NewTxOptions().
		WithIsolation(level).
		WithReadOnly(readOnly).
		WithDeferrable(deferrable)
}

func TestTxKinds(t *testing.T) {
	ctx := context.Background()

	combinations := []gelcfg.TxOptions{
		newTxOpts(gelcfg.Serializable, true, true),
		newTxOpts(gelcfg.Serializable, true, false),
		newTxOpts(gelcfg.Serializable, false, true),
		newTxOpts(gelcfg.Serializable, false, false),
		gelcfg.NewTxOptions().
			WithIsolation(gelcfg.Serializable).
			WithReadOnly(true),
		gelcfg.NewTxOptions().
			WithIsolation(gelcfg.Serializable).
			WithReadOnly(false),
		gelcfg.NewTxOptions().
			WithIsolation(gelcfg.Serializable).
			WithDeferrable(true),
		gelcfg.NewTxOptions().
			WithIsolation(gelcfg.Serializable).
			WithDeferrable(false),
		gelcfg.NewTxOptions().WithReadOnly(true).WithDeferrable(true),
		gelcfg.NewTxOptions().WithReadOnly(true).WithDeferrable(false),
		gelcfg.NewTxOptions().WithReadOnly(false).WithDeferrable(true),
		gelcfg.NewTxOptions().WithReadOnly(false).WithDeferrable(false),
		gelcfg.NewTxOptions().WithIsolation(gelcfg.Serializable),
		gelcfg.NewTxOptions().WithReadOnly(true),
		gelcfg.NewTxOptions().WithReadOnly(false),
		gelcfg.NewTxOptions().WithDeferrable(true),
		gelcfg.NewTxOptions().WithDeferrable(false),
	}

	noOp := func(ctx context.Context, tx geltypes.Tx) error { return nil }

	for _, opts := range combinations {
		name := fmt.Sprintf("%#v", opts)

		t.Run(name, func(t *testing.T) {
			p := client.WithTxOptions(opts)
			require.NoError(t, p.Tx(ctx, noOp))
		})
	}
}

func TestWithConfigInTx(t *testing.T) {
	if protocolVersion.LT(gel.ProtocolVersion1p0) {
		t.Skip()
	}

	ctx := context.Background()

	err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		var id types.UUID
		_, e := rnd.Read(id[:])
		assert.NoError(t, e)

		e = tx.Execute(ctx, `insert User { id := <uuid>$0 }`, id)
		assert.True(t, strings.HasPrefix(
			e.Error(),
			"gel.QueryError: cannot assign to property 'id'",
		))

		return errors.New("rollback")
	})
	assert.EqualError(t, err, "rollback")

	c := client.WithConfig(map[string]interface{}{
		"allow_user_specified_id": true,
	})

	var id types.UUID
	_, e := rnd.Read(id[:])
	assert.NoError(t, e)

	// todo: remove this Execute query after
	// https://github.com/edgedb/edgedb/issues/4816
	// is resolved
	e = c.Execute(ctx, `insert User { id := <uuid>$0 }`, id)
	assert.NoError(t, e)

	err = c.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		var id types.UUID
		_, e := rnd.Read(id[:])
		assert.NoError(t, e)

		e = tx.Execute(ctx, `insert User { id := <uuid>$0 }`, id)
		assert.NoError(t, e)

		return errors.New("rollback")
	})
	assert.EqualError(t, err, "rollback")
}

func serverVersion(t *testing.T) int64 {
	ctx := context.Background()

	var version int64
	err := client.QuerySingle(ctx, "SELECT sys::get_version().major", &version)
	assert.NoError(t, err)

	return version
}

func TestSQLTx(t *testing.T) {
	ctx := context.Background()
	rollback := errors.New("rollback")

	if serverVersion(t) >= 6 {
		typename := "ExecuteSQL_01"
		query := fmt.Sprintf("select %s.prop1 limit 1", typename)

		err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
			if e := tx.Execute(ctx, fmt.Sprintf(`
	 		  CREATE TYPE %s {
	 			  CREATE REQUIRED PROPERTY prop1 -> std::str;
	 		  };
			`, typename)); e != nil {
				return e
			}

			if e := tx.ExecuteSQL(ctx, fmt.Sprintf(`
			  INSERT INTO "%s" (prop1) VALUES (123);
			`, typename)); e != nil {
				return e
			}

			var res string
			if e := tx.QuerySingle(ctx, query, &res); e != nil {
				return e
			}
			assert.Equal(t, "123", res)

			if e := tx.ExecuteSQL(ctx, fmt.Sprintf(`
				UPDATE "%s" SET prop1 = '345';
			`, typename)); e != nil {
				return e
			}

			var res2 string
			if e := tx.QuerySingle(ctx, query, &res2); e != nil {
				return e
			}
			assert.Equal(t, "345", res2)

			var updateRes []struct {
				Prop1 string `edgedb:"prop1"`
			}
			if e := tx.QuerySQL(ctx, fmt.Sprintf(`
		    UPDATE "%s" SET prop1 = '567' RETURNING prop1;
		  `, typename), &updateRes); e != nil {
				return e
			}
			assert.Equal(t, "567", updateRes[0].Prop1)

			return rollback
		})
		assert.Equal(t, rollback, err)
	} else {
		err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
			if e := tx.ExecuteSQL(ctx, "SELECT 1"); e != nil {
				return e
			}

			return rollback
		})
		assert.EqualError(
			t, err,
			"gel.UnsupportedFeatureError: "+
				"the server does not support SQL queries, "+
				"upgrade to 6.0 or newer",
		)
	}
}

func selectInTx(t *testing.T, cb func(context.Context, geltypes.Tx, string)) {
	name := fmt.Sprintf("test%v", rand.Intn(10_000_000))
	ctx := context.Background()
	err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		e := tx.Execute(ctx, "INSERT TxTest {name := <str>$0};", name)
		require.NoError(t, e)

		cb(ctx, tx, name)
		return nil
	})
	require.NoError(t, err)

	query := `
		WITH removed := (DELETE TxTest FILTER .name = <str>$0)
		SELECT removed.name LIMIT 1
	`
	var testNames []string
	err = client.Query(ctx, query, &testNames, name)

	require.NoError(t, err)
	require.Equal(
		t,
		[]string{name},
		testNames,
		"The transaction wasn't commited",
	)
}

func TestTxExerciseQuery(t *testing.T) {
	selectInTx(t, func(ctx context.Context, tx geltypes.Tx, name string) {
		var result []string
		query := "SELECT name := TxTest.name FILTER name = <str>$0"
		err := tx.Query(ctx, query, &result, name)
		require.NoError(t, err)
		assert.Equal(t, []string{name}, result)
	})
}

func TestTxExerciseQueryJSON(t *testing.T) {
	selectInTx(t, func(ctx context.Context, tx geltypes.Tx, name string) {
		var result []byte
		query := "SELECT name := TxTest.name FILTER name = <str>$0"
		err := tx.QueryJSON(ctx, query, &result, name)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`["%s"]`, name), string(result))
	})
}

func TestTxExerciseQuerySQL(t *testing.T) {
	if serverVersion(t) < 6 {
		t.Skipf("server version is too old to support SQL")
	}

	selectInTx(t, func(ctx context.Context, tx geltypes.Tx, name string) {
		var result []struct {
			Name string `gel:"name"`
		}
		query := `SELECT name FROM "TxTest" WHERE name = $1`
		err := tx.QuerySQL(ctx, query, &result, name)
		require.NoError(t, err)
		require.Equal(t, 1, len(result))
		assert.Equal(t, name, result[0].Name)
	})
}

func TestTxQuerySQLMalformedQuery(t *testing.T) {
	t.Skip("see https://github.com/geldata/gel-go/issues/335")
	if serverVersion(t) < 6 {
		t.Skip("server version is too old to support SQL")
	}

	selectInTx(t, func(ctx context.Context, tx geltypes.Tx, name string) {
		var result []string
		err := tx.QuerySQL(ctx, `malformed query`, &result, name)
		require.NoError(t, err)
		assert.Equal(t, []string{name}, result)
	})
}

func TestTxExerciseQuerySingle(t *testing.T) {
	if serverVersion(t) < 6 {
		t.Skipf("server version is too old to support SQL")
	}

	selectInTx(t, func(ctx context.Context, tx geltypes.Tx, name string) {
		var result string
		query := "SELECT name := TxTest.name FILTER name = <str>$0 LIMIT 1"
		err := tx.QuerySingle(ctx, query, &result, name)
		require.NoError(t, err)
		assert.Equal(t, name, result)
	})
}

func TestTxExerciseQuerySingleJSON(t *testing.T) {
	if serverVersion(t) < 6 {
		t.Skipf("server version is too old to support SQL")
	}

	selectInTx(t, func(ctx context.Context, tx geltypes.Tx, name string) {
		var result []byte
		query := "SELECT name := TxTest.name FILTER name = <str>$0 LIMIT 1"
		err := tx.QuerySingleJSON(ctx, query, &result, name)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(`"%s"`, name), string(result))
	})
}
