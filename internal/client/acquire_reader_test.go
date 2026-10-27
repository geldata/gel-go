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
	"net"
	"testing"
	"time"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/internal/buff"
	"github.com/geldata/gel-go/internal/cache"
	gelerrint "github.com/geldata/gel-go/internal/gelerr"
	"github.com/geldata/gel-go/internal/soc"
)

// newIdlePipeConn builds a protocolConnection over an in-memory pipe with
// the reader loop running: the state a pooled connection rests in between
// queries. Nothing ever writes to the server end, so no message arrives,
// exactly like an idle connection.
func newIdlePipeConn(t *testing.T) *protocolConnection {
	t.Helper()
	client, server := net.Pipe()
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})

	c := &protocolConnection{
		soc:                 &autoClosingSocket{conn: client},
		acquireReaderSignal: make(chan struct{}, 1),
		readerChan:          make(chan *buff.Reader, 1),
	}
	r := buff.NewReader(make(chan *soc.Data, 2))
	if err := c.releaseReader(r); err != nil {
		t.Fatal(err)
	}
	return c
}

// acquireWithin reports whether acquireReader returned at all within limit.
func acquireWithin(
	ctx context.Context,
	c *protocolConnection,
	limit time.Duration,
) (returned bool, r *buff.Reader, err error) {
	type result struct {
		r   *buff.Reader
		err error
	}
	done := make(chan result, 1)
	go func() {
		r, err := c.acquireReader(ctx)
		done <- result{r, err}
	}()
	select {
	case x := <-done:
		return true, x.r, x.err
	case <-time.After(limit):
		return false, nil, nil
	}
}

// TestAcquireReaderSurvivesCancelledWaits reproduces a production hang seen
// on 2026-10-24: acquires with an already-expired ctx (a query retried after
// its deadline) left the reader hand-off in a state no later caller could
// get past, and that caller parked forever in `chan send` on
// acquireReaderSignal.
func TestAcquireReaderSurvivesCancelledWaits(t *testing.T) {
	c := newIdlePipeConn(t)
	dead, cancel := context.WithCancel(context.Background())
	cancel()

	// Drive the connection with dead-ctx acquires the way withRetries does
	// for a retryable error. A reader that does come back goes through
	// releaseReader, keeping the protocol honest. Stop once the stale state
	// is reached: a buffered signal nobody will consume plus a parked
	// reader with no loop goroutine left.
	for i := 0; i < 50; i++ {
		r, err := c.acquireReader(dead)
		if err == nil {
			if e := c.releaseReader(r); e != nil {
				break
			}
			continue
		}
		if c.isClosed() {
			break
		}
		time.Sleep(2 * time.Millisecond)
		if len(c.acquireReaderSignal) == 1 && len(c.readerChan) == 1 {
			break
		}
	}

	live, cancelLive := context.WithTimeout(
		context.Background(), time.Second)
	defer cancelLive()
	returned, r, err := acquireWithin(live, c, 3*time.Second)
	if !returned {
		t.Fatal("acquireReader wedged: a live caller never returned " +
			"after cancelled waits on the same connection")
	}
	if err == nil {
		_ = c.releaseReader(r)
	}
}

// TestWithRetriesStopsOnDeadContext: a query whose ctx already expired used
// to be retried twice more, with the back-off sleeps in between, because the
// ctx exit surfaces as ClientConnectionTimeoutError, which is retryable.
// Those retries are what re-entered acquireReader with a dead ctx.
func TestWithRetriesStopsOnDeadContext(t *testing.T) {
	pc := newIdlePipeConn(t)
	c := &transactableConn{reconnectingConn: &reconnectingConn{
		borrowableConn: borrowableConn{conn: pc},
		cacheCollection: cacheCollection{
			typeIDCache:       cache.New(10),
			inCodecCache:      cache.New(10),
			outCodecCache:     cache.New(10),
			capabilitiesCache: cache.New(10),
		},
	}}
	q := &query{
		cmd:    "select 1",
		method: "Query",
		cfg:    QueryConfig{RetryOptions: gelcfg.NewRetryOptions()},
	}
	// a read-only query (no capabilities) is always retryable
	c.capabilitiesCache.Put(makeKey(q), uint64(0))

	dead, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	start := time.Now()
	err := c.withRetries(dead, q, func(ctx context.Context, _ *query) error {
		calls++
		return wrapNetError(ctx.Err())
	})
	if err == nil {
		t.Fatal("expected the ctx error, got nil")
	}
	if calls != 1 {
		t.Fatalf("callback ran %d times on a dead ctx, want 1", calls)
	}
	if d := time.Since(start); d > 100*time.Millisecond {
		t.Fatalf("withRetries took %v on a dead ctx, want no back-off", d)
	}
}

// TestCloseReturnsOnBrokenHandoff: close() acquires the reader before the
// Terminate message. It used to wait with context.Background(); on a
// connection whose hand-off was broken it never returned, and Pool.Release,
// which calls it, parked the request's goroutine with it.
func TestCloseReturnsOnBrokenHandoff(t *testing.T) {
	c := newIdlePipeConn(t)
	r, err := c.acquireReader(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// Break the hand-off by hand: park the reader with no loop goroutine
	// and leave a signal nobody will consume.
	c.readerChan <- r
	c.acquireReaderSignal <- struct{}{}

	old := closeAcquireTimeout
	closeAcquireTimeout = 200 * time.Millisecond
	t.Cleanup(func() { closeAcquireTimeout = old })

	done := make(chan error, 1)
	go func() { done <- c.close() }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("close() returned nil on a broken hand-off")
		}
		if !c.isClosed() {
			t.Fatal("close() returned but left the socket open")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("close() parked on the broken hand-off")
	}
}

// TestTxAbandonsOpenTransactionOnDeadCtx: with the ctx already over,
// rollback (like commit) fails before anything reaches the server, so the
// transaction stays open there. Tx must retire the connection instead of
// letting Pool.Release put a mid-transaction session back into the pool.
func TestTxAbandonsOpenTransactionOnDeadCtx(t *testing.T) {
	pc := newIdlePipeConn(t)
	tx := &Tx{
		borrowableConn: borrowableConn{conn: pc},
		txState:        &txState{txStatus: startedTx},
	}
	dead, cancel := context.WithCancel(context.Background())
	cancel()

	e := tx.rollback(dead)
	if e == nil {
		t.Fatal("rollback on a dead ctx returned nil")
	}
	if isClientConnectionError(e) || pc.isClosed() {
		t.Fatalf("rollback alone must not touch the connection: %v", e)
	}

	e = abandonTx(pc, e)
	if !isClientConnectionError(e) {
		t.Fatalf("abandonTx returned %T, want a client connection error", e)
	}
	if !pc.isClosed() {
		t.Fatal("abandonTx left the socket open with a transaction on it")
	}

	// A server answer passes through untouched.
	pc2 := newIdlePipeConn(t)
	serverErr := gelerrint.NewTransactionSerializationError("", nil)
	if got := abandonTx(pc2, serverErr); got != serverErr || pc2.isClosed() {
		t.Fatalf("abandonTx changed a server error (%v) or closed the socket",
			got)
	}

	// An unknown outcome (not a ctx error, no server answer) closes the
	// socket but keeps its error type: nobody may read it as retryable.
	pc3 := newIdlePipeConn(t)
	other := errors.New("decode failed")
	if got := abandonTx(pc3, other); got != other || !pc3.isClosed() {
		t.Fatalf("abandonTx on an unknown outcome: got %v closed=%v",
			got, pc3.isClosed())
	}
}

// TestAcquireReaderRetiresConnectionOnCancelledWait: the loop goroutine is
// mid-message (a header announcing a body that never arrives), so the
// hand-off cannot complete. A caller whose ctx ends while it waits must
// retire the connection, and the loop goroutine must still exit.
func TestAcquireReaderRetiresConnectionOnCancelledWait(t *testing.T) {
	client, server := net.Pipe()
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})
	c := &protocolConnection{
		soc:                 &autoClosingSocket{conn: client},
		acquireReaderSignal: make(chan struct{}, 1),
		readerChan:          make(chan *buff.Reader, 1),
	}
	// The real socket reader feeds the loop through a forwarder so the test
	// can see when a chunk has reached the reader's channel.
	raw := make(chan *soc.Data, 2)
	toBeDeserialized := make(chan *soc.Data, 2)
	pushed := make(chan struct{}, 8)
	go soc.Read(c.soc, soc.NewMemPool(4, 256*1024), raw)
	go func() {
		for d := range raw {
			toBeDeserialized <- d
			pushed <- struct{}{}
			if d.Err != nil {
				return
			}
		}
	}()
	if err := c.releaseReader(buff.NewReader(toBeDeserialized)); err != nil {
		t.Fatal(err)
	}
	// A ParameterStatus header announcing 96 more bytes that never come.
	// Once the loop goroutine has taken the chunk it is committed to the
	// body read and no longer watches the signal; wait for that moment,
	// else the loop may still answer the signal and hand the reader over.
	header := []byte{uint8(ParameterStatus), 0, 0, 0, 100}
	if _, err := server.Write(header); err != nil {
		t.Fatal(err)
	}
	<-pushed
	for len(toBeDeserialized) > 0 {
		time.Sleep(time.Millisecond)
	}

	// Cancel only once the signal is out, so the test exercises the wait
	// after the send and never the early return.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := c.acquireReader(ctx)
		result <- err
	}()
	for len(c.acquireReaderSignal) == 0 {
		time.Sleep(time.Millisecond)
	}
	cancel()
	err := <-result
	if !isClientConnectionError(err) {
		t.Fatalf("got %v, want a client connection error", err)
	}
	if !c.isClosed() {
		t.Fatal("the cancelled wait left the connection open")
	}
	deadline := time.Now().Add(time.Second)
	for len(c.readerChan) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("the loop goroutine did not exit after the socket closed")
		}
		time.Sleep(5 * time.Millisecond)
	}
}
