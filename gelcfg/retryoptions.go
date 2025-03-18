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

package gelcfg

import (
	"fmt"

	"github.com/geldata/gel-go/gelerr"
	gelerrint "github.com/geldata/gel-go/internal/gelerr"
)

// RetryCondition represents scenarios that can cause a transaction
// or a single query run outside a geltypes.TxBlock to be retried.
type RetryCondition int

const (
	// TxConflict indicates that the server could not complete a transaction
	// because it encountered a deadlock or serialization error.
	TxConflict RetryCondition = iota

	// NetworkError indicates that the transaction was interupted
	// by a network error.
	NetworkError
)

// NewRetryOptions returns the default RetryOptions.
func NewRetryOptions() RetryOptions {
	return RetryOptions{fromFactory: true}.WithDefault(NewRetryRule())
}

// RetryOptions configures how failed transactions or queries outside of
// transactions are retried.  Use [NewRetryOptions] to get a default
// RetryOptions value instead of creating one yourself.
//
// See [github.com/geldata/gel-go.Client.WithRetryOptions] for an example.
type RetryOptions struct {
	fromFactory bool
	txConflict  RetryRule
	network     RetryRule
}

// IsValid returns true if o was created with NewRetryOptions()
func (o RetryOptions) IsValid() bool { return o.fromFactory }

// WithDefault sets the rule for all conditions to rule.
func (o RetryOptions) WithDefault(rule RetryRule) RetryOptions { // nolint:gocritic,lll
	if !rule.fromFactory {
		panic("RetryRule not created with NewRetryRule() is not valid")
	}

	o.txConflict = rule
	o.network = rule
	return o
}

// WithCondition sets the retry rule for the specified condition.
func (o RetryOptions) WithCondition( // nolint:gocritic
	condition RetryCondition,
	rule RetryRule,
) RetryOptions {
	if !rule.fromFactory {
		panic("RetryRule not created with NewRetryRule() is not valid")
	}

	switch condition {
	case TxConflict:
		o.txConflict = rule
	case NetworkError:
		o.network = rule
	default:
		panic(fmt.Sprintf("unexpected condition: %v", condition))
	}

	return o
}

// RuleForException returns the RetryRule to be applied for err.
func (o RetryOptions) RuleForException(err gelerr.Error) (RetryRule, error) {
	switch {
	case err.Category(gelerr.TransactionConflictError):
		return o.txConflict, nil
	case err.Category(gelerr.ClientError):
		return o.network, nil
	default:
		return RetryRule{}, gelerrint.NewClientError(
			fmt.Sprintf("unexpected error type: %T", err),
			nil,
		)
	}
}
