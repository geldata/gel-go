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
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/geltypes"
	gel "github.com/geldata/gel-go/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectClient(t *testing.T) {
	ctx := context.Background()
	p, err := CreateClient(opts)
	require.NoError(t, err)

	var result string
	err = p.QuerySingle(ctx, "SELECT 'hello';", &result)
	assert.NoError(t, err)
	assert.Equal(t, "hello", result)

	p2 := p.WithTxOptions(gelcfg.NewTxOptions())

	err = p.Close()
	assert.NoError(t, err)

	// Client should not be closeable a second time.
	err = p.Close()
	assert.EqualError(t, err, "gel.InterfaceError: client closed")

	// Copied clients should be closed if a different copy is closed.
	err = p2.Close()
	assert.EqualError(t, err, "gel.InterfaceError: client closed")
}

func TestClientRejectsTransaction(t *testing.T) {
	ctx := context.Background()
	p, err := CreateClient(opts)
	require.NoError(t, err)

	expected := "gel.DisabledCapabilityError: " +
		"cannot execute transaction control commands.*"

	err = p.Execute(ctx, "START TRANSACTION")
	assert.Regexp(t, expected, err)

	var result []byte
	err = p.Query(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)

	err = p.QueryJSON(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)

	err = p.QuerySingle(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)

	err = p.QuerySingleJSON(ctx, "START TRANSACTION", &result)
	assert.Regexp(t, expected, err)

	err = p.Close()
	assert.NoError(t, err)
}

func TestConnectClientZeroConcurrency(t *testing.T) {
	o := opts
	o.Concurrency = 0

	ctx := context.Background()
	p, err := CreateClient(o)
	require.NoError(t, err)
	require.NoError(t, p.EnsureConnected(ctx))

	expected := client.pool.Cfg.ServerSettings.
		Get("suggested_pool_concurrency").(int)
	if err != nil {
		expected = gel.DefaultConcurrency
	}
	require.Equal(t, expected, p.pool.Concurrency)

	var result string
	err = p.QuerySingle(ctx, "SELECT 'hello';", &result)
	assert.NoError(t, err)
	assert.Equal(t, "hello", result)

	err = p.Close()
	assert.NoError(t, err)
}

func TestCloseClientConcurently(t *testing.T) {
	p, err := CreateClient(opts)
	require.NoError(t, err)

	errs := make(chan error)
	go func() { errs <- p.Close() }()
	go func() { errs <- p.Close() }()

	assert.NoError(t, <-errs)
	var edbErr gelerr.Error
	require.True(t, errors.As(<-errs, &edbErr), "wrong error: %v", err)
	assert.True(
		t,
		edbErr.Category(gelerr.InterfaceError),
		"wrong error: %v",
		err,
	)
}

func TestClientTx(t *testing.T) {
	ctx := context.Background()

	p, err := CreateClient(opts)
	require.NoError(t, err)
	defer p.Close() // nolint:errcheck

	var result int64
	err = p.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		return tx.QuerySingle(ctx, "SELECT 33*21", &result)
	})

	require.NoError(t, err)
	require.Equal(t, int64(693), result, "Client.Tx() failed")
}

func TestQuerySingleMissingResult(t *testing.T) {
	ctx := context.Background()

	var result string
	err := client.QuerySingle(ctx, "SELECT <str>{}", &result)
	assert.EqualError(t, err, "gel.NoDataError: zero results")

	optionalResult := geltypes.NewOptionalStr("this should be set to missing")
	err = client.QuerySingle(ctx, "SELECT <str>{}", &optionalResult)
	assert.NoError(t, err)
	assert.Equal(t, geltypes.OptionalStr{}, optionalResult)

	var objectResult struct {
		Name string `gel:"name"`
	}
	err = client.QuerySingle(ctx,
		"SELECT sys::Database { name } FILTER .name = 'does not exist'",
		&objectResult,
	)
	assert.EqualError(t, err, "gel.NoDataError: zero results")

	var optionalObjectResult struct {
		geltypes.Optional
		Name string `gel:"name"`
	}
	optionalObjectResult.SetMissing(false)
	err = client.QuerySingle(ctx,
		"SELECT sys::Database { name } FILTER .name = 'does not exist'",
		&optionalObjectResult,
	)
	assert.NoError(t, err)
	assert.Equal(t, "", optionalObjectResult.Name)
	assert.True(t, optionalObjectResult.Missing())
}

func TestQuerySingleJSONMissingResult(t *testing.T) {
	ctx := context.Background()

	var result []byte
	err := client.QuerySingleJSON(ctx, "SELECT <str>{}", &result)
	assert.EqualError(t, err, "gel.NoDataError: zero results")

	optionalResult := geltypes.NewOptionalBytes(
		[]byte("this should be set to missing"),
	)
	err = client.QuerySingleJSON(ctx, "SELECT <str>{}", &optionalResult)
	assert.NoError(t, err)
	assert.Equal(t, geltypes.OptionalBytes{}, optionalResult)

	var wrongType string
	err = client.QuerySingleJSON(ctx, "SELECT <str>{}", &wrongType)
	assert.EqualError(t, err, "gel.InterfaceError: "+
		"the \"out\" argument must be *[]byte or *OptionalBytes, got *string")
}

func TestQuerySQL(t *testing.T) {
	ctx := context.Background()

	if serverVersionGTE(t, 6, 0) {
		err := client.ExecuteSQL(ctx, "select 1")
		assert.NoError(t, err)

		var result []struct {
			Col1 int32 `edgedb:"col~1"`
		}
		err = client.QuerySQL(ctx, "select 1", &result)
		assert.NoError(t, err)
		assert.Equal(t, int32(1), result[0].Col1)

		type res2 struct {
			Foo int32 `edgedb:"foo"`
			Bar int32 `edgedb:"bar"`
		}
		var result2 []res2
		err = client.QuerySQL(ctx, "select 1 AS foo, 2 AS bar", &result2)
		assert.NoError(t, err)
		assert.Equal(t, []res2{
			{
				Foo: 1,
				Bar: 2,
			},
		}, result2)

		var result3 []struct {
			Col1 int64 `edgedb:"col~1"`
		}
		err = client.QuerySQL(ctx, "select 1 + $1::int8", &result3, int64(41))
		assert.NoError(t, err)
		assert.Equal(t, int64(42), result3[0].Col1)
	} else {
		var res []interface{}
		err := client.QuerySQL(ctx, "select 1", &res)
		assert.EqualError(
			t, err,
			"gel.UnsupportedFeatureError: "+
				"the server does not support SQL queries, "+
				"upgrade to 6.0 or newer",
		)

		err = client.ExecuteSQL(ctx, "select 1")
		assert.EqualError(
			t, err,
			"gel.UnsupportedFeatureError: "+
				"the server does not support SQL queries, "+
				"upgrade to 6.0 or newer",
		)
	}
}

func TestSessionIdleTimeout(t *testing.T) {
	ctx := context.Background()
	p, err := gel.NewPool("", opts)
	require.NoError(t, err)

	con, err := p.Acquire(ctx)
	require.NoError(t, err)

	var result geltypes.Duration
	err = gel.RunQuery(
		ctx,
		con,
		"QuerySingle",
		"SELECT assert_single(cfg::Config.session_idle_timeout)",
		&result,
		[]any{},
		p.State,
		&p.QueryConfig,
		false,
	)
	require.NoError(t, p.Release(con, err))
	require.NoError(t, err)
	require.Equal(t, geltypes.Duration(1_000_000), result)

	// The client keeps one connection in the pool.
	// Get a reference to that connection.
	con1, err := p.Acquire(ctx)
	require.NoError(t, err)
	require.NotNil(t, con1)

	err = p.Release(con1, nil)
	require.NoError(t, err)

	// After releasing we should get the same connection
	// back again on acquire.
	con2, err := p.Acquire(ctx)
	require.NoError(t, err)
	require.NotNil(t, con2)
	assert.Equal(t, con1, con2)

	err = p.Release(con2, nil)
	require.NoError(t, err)

	// If the pooled connection is not used for longer than the
	// session_idle_timeout then the next acquired connection
	// should be a new connection.
	time.Sleep(1_200 * time.Millisecond)

	con3, err := p.Acquire(ctx)
	require.NoError(t, err)
	require.NotNil(t, con3)
	assert.NotEqual(t, unsafe.Pointer(con1), unsafe.Pointer(con3))

	err = p.Release(con3, nil)
	assert.NoError(t, err)
}

// Try to trigger race conditions
func TestConcurentClientUsage(t *testing.T) {
	ctx := context.Background()
	var done sync.WaitGroup

	for i := 0; i < 2; i++ {
		done.Add(1)
		go func() {
			var result int64
			for j := 0; j < 10; j++ {
				_ = client.QuerySingle(ctx, "SELECT 1", &result)
			}
			done.Done()
		}()
	}

	done.Wait()
}
