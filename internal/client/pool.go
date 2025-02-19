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
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/edgedb/edgedb-go/gelcfg"
	"github.com/edgedb/edgedb-go/internal/cache"
	gelerrint "github.com/edgedb/edgedb-go/internal/gelerr"
)

const defaultIdleConnectionTimeout = 30 * time.Second

// DefaultConcurrency is used if no other concurrency setting is found.
var DefaultConcurrency = max(4, runtime.NumCPU())

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// NewPool makes a new pool.
func NewPool(dsn string, opts gelcfg.Options) (*Pool, error) { // nolint:gocritic,lll
	cfg, err := parseConnectDSNAndArgs(dsn, &opts, newCfgPaths())
	if err != nil {
		return nil, err
	}

	False := false
	p := &Pool{
		isClosed:             &False,
		isClosedMutex:        &sync.RWMutex{},
		Cfg:                  cfg,
		TxOpts:               gelcfg.NewTxOptions(),
		Concurrency:          int(opts.Concurrency),
		freeConns:            make(chan func() *transactableConn, 1),
		potentialConnsMutext: &sync.Mutex{},
		RetryOpts:            gelcfg.NewRetryOptions(),
		cacheCollection: cacheCollection{
			ServerSettings:    cfg.ServerSettings,
			typeIDCache:       cache.New(1_000),
			inCodecCache:      cache.New(1_000),
			outCodecCache:     cache.New(1_000),
			capabilitiesCache: cache.New(1_000),
		},
		State: make(map[string]interface{}),
	}

	return p, nil
}

// Pool is a connection pool.
type Pool struct {
	isClosed      *bool
	isClosedMutex *sync.RWMutex // locks isClosed

	// A buffered channel of structs representing unconnected capacity.
	// This field remains nil until the first connection is acquired.
	potentialConns       chan struct{}
	potentialConnsMutext *sync.Mutex

	// A buffered channel of connections ready for use.
	freeConns chan func() *transactableConn

	TxOpts    gelcfg.TxOptions
	RetryOpts gelcfg.RetryOptions

	Cfg             *connConfig
	cacheCollection cacheCollection
	State           map[string]interface{}

	Concurrency int
}

func (p *Pool) newConn(ctx context.Context) (*transactableConn, error) {
	conn := transactableConn{
		txOpts:    p.TxOpts,
		retryOpts: p.RetryOpts,
		reconnectingConn: &reconnectingConn{
			Cfg:             p.Cfg,
			cacheCollection: p.cacheCollection,
		},
	}

	if err := conn.reconnect(ctx, false); err != nil {
		return nil, err
	}

	return &conn, nil
}

// Acquire gets a connection from the pool.
func (p *Pool) Acquire(
	ctx context.Context,
) (*transactableConn, error) { // nolint:revive
	p.isClosedMutex.RLock()
	defer p.isClosedMutex.RUnlock()

	if *p.isClosed {
		return nil, gelerrint.NewInterfaceError("client closed", nil)
	}

	p.potentialConnsMutext.Lock()
	if p.potentialConns == nil {
		conn, err := p.newConn(ctx)
		if err != nil {
			p.potentialConnsMutext.Unlock()
			return nil, err
		}

		if p.Concurrency == 0 {
			// The user did not set Concurrency in provided Options.
			// See if the server sends a suggested max size.
			suggested, ok := conn.Cfg.ServerSettings.
				GetOk("suggested_pool_concurrency")
			if ok {
				p.Concurrency = suggested.(int)
			} else {
				p.Concurrency = DefaultConcurrency
			}
		}

		p.potentialConns = make(chan struct{}, p.Concurrency)
		for i := 0; i < p.Concurrency-1; i++ {
			p.potentialConns <- struct{}{}
		}

		p.potentialConnsMutext.Unlock()
		return conn, nil
	}
	p.potentialConnsMutext.Unlock()

	// force do nothing if context is expired
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("gel: %w", ctx.Err())
	default:
	}

	// force using an existing connection over connecting a new socket.
	select {
	case acquireIfNotTimedout := <-p.freeConns:
		conn := acquireIfNotTimedout()
		if conn != nil {
			return conn, nil
		}
	default:
	}

	for {
		select {
		case acquireIfNotTimedout := <-p.freeConns:
			conn := acquireIfNotTimedout()
			if conn != nil {
				return conn, nil
			}
			continue
		case <-p.potentialConns:
			conn, err := p.newConn(ctx)
			if err != nil {
				p.potentialConns <- struct{}{}
				return nil, err
			}
			return conn, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("gel: %w", ctx.Err())
		}
	}
}

// Release puts a connection back in the pool.
func (p *Pool) Release(conn *transactableConn, err error) error {
	if isClientConnectionError(err) {
		p.potentialConns <- struct{}{}
		return conn.Close()
	}

	timeout := defaultIdleConnectionTimeout
	if t, ok := conn.conn.SystemConfig.SessionIdleTimeout.Get(); ok {
		timeout = time.Duration(1_000 * t)
	}

	// 0 or less disables the idle timeout
	if timeout <= 0 {
		select {
		case p.freeConns <- func() *transactableConn { return conn }:
			return nil
		default:
			// we have MinConns idle so no need to keep this connection.
			p.potentialConns <- struct{}{}
			return conn.Close()
		}
	}

	cancel := make(chan struct{}, 1)
	connChan := make(chan *transactableConn, 1)

	acquireIfNotTimedout := func() *transactableConn {
		cancel <- struct{}{}
		return <-connChan
	}

	select {
	case p.freeConns <- acquireIfNotTimedout:
		go func() {
			select {
			case <-cancel:
				connChan <- conn
			case <-time.After(timeout):
				connChan <- nil
				p.potentialConns <- struct{}{}
				if e := conn.Close(); e != nil {
					log.Println("error while closing idle connection:", e)
				}
			}
		}()
	default:
		// we have MinConns idle so no need to keep this connection.
		p.potentialConns <- struct{}{}
		return conn.Close()
	}

	return nil
}

// EnsureConnected forces the pool to connect if it hasn't already.
func (p *Pool) EnsureConnected(ctx context.Context) error {
	conn, err := p.Acquire(ctx)
	if err != nil {
		return err
	}

	return p.Release(conn, nil)
}

// Close closes all connections in the pool.
// Calling close blocks until all acquired connections have been released,
// and returns an error if called more than once.
func (p *Pool) Close() error {
	p.isClosedMutex.Lock()
	defer p.isClosedMutex.Unlock()

	if *p.isClosed {
		return gelerrint.NewInterfaceError("client closed", nil)
	}
	*p.isClosed = true

	p.potentialConnsMutext.Lock()
	if p.potentialConns == nil {
		// The client never made any connections.
		p.potentialConnsMutext.Unlock()
		return nil
	}
	p.potentialConnsMutext.Unlock()

	wg := sync.WaitGroup{}
	errs := make([]error, p.Concurrency)
	for i := 0; i < p.Concurrency; i++ {
		select {
		case acquireIfNotTimedout := <-p.freeConns:
			wg.Add(1)
			go func(i int) {
				conn := acquireIfNotTimedout()
				if conn != nil {
					errs[i] = conn.Close()
				}
				wg.Done()
			}(i)
		case <-p.potentialConns:
		}
	}

	wg.Wait()
	return wrapAll(errs...)
}
