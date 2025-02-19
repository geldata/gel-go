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

package gelcfg

import (
	"fmt"
	"math"
	"time"

	"github.com/edgedb/edgedb-go/internal/snc"
)

var rnd = snc.NewRand()

// RetryBackoff returns the duration to wait after the nth attempt
// before making the next attempt when retrying a transaction.
type RetryBackoff func(n int) time.Duration

func defaultBackoff(attempt int) time.Duration {
	backoff := math.Pow(2.0, float64(attempt)) * 100.0
	jitter := rnd.Float64() * 100.0
	return time.Duration(backoff+jitter) * time.Millisecond
}

// NewRetryRule returns the default RetryRule value.
func NewRetryRule() RetryRule {
	return RetryRule{
		fromFactory: true,
		attempts:    3,
		backoff:     defaultBackoff,
	}
}

// RetryRule determines how transactions should be retried when run in Tx()
// methods. See Client.Tx() for details.
type RetryRule struct {
	// fromFactory indicates that a RetryOptions value was created using
	// NewRetryOptions() and not created directly. Requiring users to use the
	// factory function allows for nonzero default values.
	fromFactory bool

	// Total number of times to attempt a transaction.
	// attempts <= 0 indicate that a default value should be used.
	attempts int

	// backoff determines how long to wait between transaction attempts.
	// nil indicates that a default function should be used.
	backoff RetryBackoff
}

// Attempts retruns the number of retry attempts allowed.
func (r RetryRule) Attempts() int { return r.attempts }

// Backoff returns the RetryBackoff.
func (r RetryRule) Backoff() RetryBackoff { return r.backoff }

// WithAttempts sets the rule's attempts. attempts must be greater than zero.
func (r RetryRule) WithAttempts(attempts int) RetryRule {
	if attempts < 1 {
		panic(fmt.Sprintf(
			"RetryRule attempts must be greater than 0, got %v",
			attempts,
		))
	}

	r.attempts = attempts
	return r
}

// WithBackoff returns a copy of the RetryRule with backoff set to fn.
func (r RetryRule) WithBackoff(fn RetryBackoff) RetryRule {
	if fn == nil {
		panic("the backoff function must not be nil")
	}

	r.backoff = fn
	return r
}
