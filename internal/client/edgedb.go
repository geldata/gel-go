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
	"fmt"
	"log"
	"time"

	types "github.com/geldata/gel-go/geltypes"
	"github.com/geldata/gel-go/internal"
	"github.com/geldata/gel-go/internal/buff"
	"github.com/geldata/gel-go/internal/cache"
	"github.com/geldata/gel-go/internal/codecs"
	"github.com/geldata/gel-go/internal/gelerr"
	"github.com/geldata/gel-go/internal/snc"
	"github.com/geldata/gel-go/internal/soc"
)

type systemConfig struct {
	ID                 types.OptionalUUID     `gel:"id"`
	SessionIdleTimeout types.OptionalDuration `gel:"session_idle_timeout"`
}

type cacheCollection struct {
	ServerSettings    *snc.ServerSettings
	typeIDCache       *cache.Cache
	inCodecCache      *cache.Cache
	outCodecCache     *cache.Cache
	capabilitiesCache *cache.Cache // nolint:structcheck
}

type protocolConnection struct {
	soc                 *autoClosingSocket
	writeMemory         [1024]byte
	acquireReaderSignal chan struct{}
	readerChan          chan *buff.Reader

	protocolVersion internal.ProtocolVersion
	cacheCollection

	SystemConfig systemConfig
	stateCodec   codecs.Encoder
}

// connectWithTimeout makes a single attempt to connect to `addr`.
func connectWithTimeout(
	ctx context.Context,
	cfg *connConfig,
	caches cacheCollection,
) (*protocolConnection, error) {
	socket, err := connectAutoClosingSocket(ctx, cfg)
	if err != nil {
		return nil, err
	}

	deadline, _ := ctx.Deadline()
	err = socket.SetDeadline(deadline)
	if err != nil {
		return nil, err
	}

	conn := &protocolConnection{
		soc:                 socket,
		acquireReaderSignal: make(chan struct{}, 1),
		readerChan:          make(chan *buff.Reader, 1),
		cacheCollection:     caches,
	}

	toBeDeserialized := make(chan *soc.Data, 2)
	go soc.Read(socket, soc.NewMemPool(4, 256*1024), toBeDeserialized)
	r := buff.NewReader(toBeDeserialized)

	err = conn.connect(r, cfg)
	if err != nil {
		return nil, err
	}

	err = socket.SetDeadline(time.Time{})
	if err != nil {
		return nil, err
	}

	return conn, conn.releaseReader(r)
}

// closeAcquireTimeout bounds close()'s wait for the reader. Pool.Release
// calls close() on the caller's goroutine, so a connection whose reader
// hand-off is broken must still close instead of parking the caller. A
// healthy hand-off takes microseconds; a second is already generous.
var closeAcquireTimeout = time.Second

func (c *protocolConnection) acquireReader(
	ctx context.Context,
) (*buff.Reader, error) {
	if c.isClosed() {
		return nil, gelerr.NewClientConnectionClosedError("", nil)
	}

	// A ctx that is already over never touches the hand-off channels. The
	// signal used to be sent unconditionally and the select below could
	// then leave on ctx.Done(), stranding the connection: the signal
	// stayed buffered with nobody to consume it, or the reader was parked
	// in readerChan with no loop goroutine left. The next acquire on the
	// connection (a retry of the same dead query, or close() from
	// Pool.Release) blocked forever on the full 1-slot signal channel.
	// Nothing was touched here, so the pool keeps the connection.
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("gel: %w", err)
	}

	select {
	case c.acquireReaderSignal <- struct{}{}:
	case <-ctx.Done():
		// ctx ended while the send waited (the channel was full), or
		// both were ready and the select picked ctx. The hand-off state
		// is not known from here, so retire the connection.
		_ = c.soc.Close()
		return nil, wrapNetError(ctx.Err())
	}

	select {
	case r := <-c.readerChan:
		if r.Err != nil {
			return nil, gelerr.NewClientConnectionClosedError("", r.Err)
		}
		if c.soc.Closed() {
			return nil, gelerr.NewClientConnectionClosedError("", nil)
		}
		return r, nil
	case <-ctx.Done():
		// The signal is out and the reader will be parked for nobody.
		// Retire the connection: every later acquire (close() included)
		// fails fast with ClientConnectionClosedError and the pool
		// reconnects; the loop goroutine exits on the signal or on the
		// closed socket.
		_ = c.soc.Close()
		return nil, wrapNetError(ctx.Err())
	}
}

func (c *protocolConnection) releaseReader(r *buff.Reader) error {
	if c.isClosed() {
		return gelerr.NewClientConnectionClosedError("", nil)
	}

	if err := c.soc.SetDeadline(time.Time{}); err != nil {
		return err
	}

	go func() {
		for r.Next(c.acquireReaderSignal) {
			switch Message(r.MsgType) {
			case ErrorResponse:
				err := decodeErrorResponseMsg(r, "", "")
				log.Println("from background:", err)
				r.Err = wrapAll(r.Err, err)
			default:
				if e := c.fallThrough(r); e != nil {
					log.Println("fall through in background:", e)
					r.Err = wrapAll(r.Err, e)
					err := c.soc.Close()
					if err != nil {
						log.Println("error closing socket:", err)
					}
					c.readerChan <- r
					return
				}
			}
		}

		if r.Err != nil {
			err := c.soc.Close()
			if err != nil {
				log.Println("error closing socket:", err)
			}
		}
		c.readerChan <- r
	}()

	return nil
}

// Close the db connection
func (c *protocolConnection) close() error {
	if c.soc == nil {
		return gelerr.NewInterfaceError(
			"connection closed more than once",
			nil,
		)
	}

	// Bounded: Pool.Release calls close() on the caller's goroutine, and
	// the socket closes whether or not the reader ever comes back.
	ctx, cancel := context.WithTimeout(
		context.Background(), closeAcquireTimeout)
	defer cancel()

	if _, err := c.acquireReader(ctx); err != nil {
		return FirstError(err, c.soc.Close())
	}

	if err := c.terminate(); err != nil {
		return err
	}

	return c.soc.Close()
}

func (c *protocolConnection) isClosed() bool {
	if c.soc == nil || c.soc.Closed() {
		return true
	}

	return false
}

func (c *protocolConnection) ScriptFlow(ctx context.Context, q *query) error {
	if q.lang == SQL && c.protocolVersion.LT(protocolVersion3p0) {
		return gelerr.NewUnsupportedFeatureError(
			"the server does not support SQL queries, "+
				"upgrade to 6.0 or newer",
			nil,
		)
	}

	r, err := c.acquireReader(ctx)
	if err != nil {
		return err
	}

	deadline, _ := ctx.Deadline()
	err = c.soc.SetDeadline(deadline)
	if err != nil {
		return err
	}

	switch {
	case c.protocolVersion.GTE(protocolVersion2p0):
		err = c.execGranularFlow2pX(r, q)
	default:
		err = c.execGranularFlow1pX(r, q)
	}

	return FirstError(err, c.releaseReader(r))
}

func (c *protocolConnection) granularFlow(
	ctx context.Context,
	q *query,
) error {
	if q.lang == SQL && c.protocolVersion.LT(protocolVersion3p0) {
		return gelerr.NewUnsupportedFeatureError(
			"the server does not support SQL queries, "+
				"upgrade to 6.0 or newer",
			nil,
		)
	}

	r, err := c.acquireReader(ctx)
	if err != nil {
		return err
	}

	deadline, _ := ctx.Deadline()
	err = c.soc.SetDeadline(deadline)
	if err != nil {
		return err
	}

	switch {
	case c.protocolVersion.GTE(protocolVersion2p0):
		err = c.execGranularFlow2pX(r, q)
	default:
		err = c.execGranularFlow1pX(r, q)
	}

	return FirstError(err, c.releaseReader(r))
}
