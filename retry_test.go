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
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/geldata/gel-go/geltypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomName() string {
	return fmt.Sprintf("test%v", rand.Intn(10_000_000))
}

func assertRetry(
	t *testing.T,
	cb func(context.Context, string, string, int64, int64) error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := client.Execute(ctx, `SELECT 1`)
	require.NoError(t, err)

	query := `
		SELECT (
			INSERT Counter {
				name := <str>$0,
				value := 1,
			} UNLESS CONFLICT ON .name
			ELSE (
				UPDATE Counter
				SET { value := .value + 1 }
			)
		).value
		ORDER BY sys::_sleep(<int64>$1)
		THEN <int64>$2
	`

	name := randomName()
	errors := make(chan error, 3)

	go func() {
		errors <- cb(ctx, query, name, 0, 0)
	}()

	go func() {
		errors <- cb(ctx, query, name, 5, 1)
	}()

	go func() {
		errors <- cb(ctx, query, name, 5, 2)
	}()

	for i := 0; i < 3; i++ {
		assert.NoError(t, <-errors)
	}

	var value int32
	query = "SELECT (SELECT Counter FILTER .name = <str>$0).value"
	err = client.QuerySingle(ctx, query, &value, name)
	assert.NoError(t, err)
	assert.Equal(t, int32(3), value)
}

func TestRetryExecute(t *testing.T) {
	assertRetry(t, func(
		ctx context.Context,
		query, name string,
		sleep, nonce int64,
	) error {
		return client.Execute(ctx, query, name, sleep, nonce)
	})
}

func TestRetryQuery(t *testing.T) {
	assertRetry(t, func(
		ctx context.Context,
		query, name string,
		sleep, nonce int64,
	) error {
		var value []int32
		return client.Query(ctx, query, &value, name, sleep, nonce)
	})
}

func TestRetryQueryJSON(t *testing.T) {
	assertRetry(t, func(
		ctx context.Context,
		query, name string,
		sleep, nonce int64,
	) error {
		var value []byte
		return client.QueryJSON(ctx, query, &value, name, sleep, nonce)
	})
}

func TestRetryQuerySingle(t *testing.T) {
	assertRetry(t, func(
		ctx context.Context,
		query, name string,
		sleep, nonce int64,
	) error {
		var value int32
		return client.QuerySingle(ctx, query, &value, name, sleep, nonce)
	})
}

func TestRetryQuerySingleJSON(t *testing.T) {
	assertRetry(t, func(
		ctx context.Context,
		query, name string,
		sleep, nonce int64,
	) error {
		var value []byte
		return client.QuerySingleJSON(ctx, query, &value, name, sleep, nonce)
	})
}

func assertRetryTx(
	t *testing.T,
	cb func(context.Context, geltypes.Tx, string, string) error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	name := randomName()
	iterMtx := sync.Mutex{}
	iterations := 0
	barrier := sync.WaitGroup{}
	barrier.Add(2)
	errors := make(chan error, 2)

	for i := 0; i < 2; i++ {
		go func() {
			errors <- client.Tx(
				ctx,
				func(ctx context.Context, tx geltypes.Tx) error {
					iterMtx.Lock()
					iterations++
					iterMtx.Unlock()

					require.NoError(t, tx.Execute(ctx, "SELECT 1"))

					barrier.Done()
					barrier.Wait()

					query := `
					SELECT (
						INSERT Counter {
							name := <str>$0,
							value := 1,
						}
						UNLESS CONFLICT ON .name
						ELSE (
							UPDATE Counter
							SET { value := .value + 1 }
						)
					).value
				`

					err := cb(ctx, tx, query, name)
					if err != nil {
						barrier.Add(1)
					}
					return err
				},
			)
		}()
	}

	for i := 0; i < 2; i++ {
		require.NoError(t, <-errors)
	}

	assert.Equal(t, 3, iterations)

	var value int32
	query := "SELECT (SELECT Counter FILTER .name = <str>$0).value"
	err := client.QuerySingle(ctx, query, &value, name)
	assert.NoError(t, err)
	assert.Equal(t, int32(2), value)
}

func TestRetryTxQuery(t *testing.T) {
	assertRetryTx(
		t,
		func(ctx context.Context, tx geltypes.Tx, query, name string) error {
			var value []int32
			return tx.Query(ctx, query, &value, name)
		},
	)
}

func TestRetryTxQueryJSON(t *testing.T) {
	assertRetryTx(
		t,
		func(ctx context.Context, tx geltypes.Tx, query, name string) error {
			var value []byte
			return tx.QueryJSON(ctx, query, &value, name)
		},
	)
}

func TestRetryTxQuerySingle(t *testing.T) {
	assertRetryTx(
		t,
		func(ctx context.Context, tx geltypes.Tx, query, name string) error {
			var value int32
			return tx.QuerySingle(ctx, query, &value, name)
		},
	)
}

func TestRetryTxQuerySingleJSON(t *testing.T) {
	assertRetryTx(
		t,
		func(ctx context.Context, tx geltypes.Tx, query, name string) error {
			var value []byte
			return tx.QuerySingleJSON(ctx, query, &value, name)
		},
	)
}
