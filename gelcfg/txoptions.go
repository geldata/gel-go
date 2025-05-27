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

import "fmt"

// IsolationLevel documentation can be found [here]
//
// [here]: https://docs.geldata.com/reference/stdlib/sys#type::sys::TransactionIsolation
type IsolationLevel string

const (
	// Serializable isolation level
	Serializable IsolationLevel = "Serializable"

	// RepeatableRead isolation level (supported in read-only transactions)
	RepeatableRead IsolationLevel = "RepeatableRead"

	// PreferRepeatableRead use RepeatableRead isolation level and fall back to
	// Serializable if the transaction is read-only.
	PreferRepeatableRead IsolationLevel = "PreferRepeatableRead"
)

// NewTxOptions returns the default TxOptions value.
func NewTxOptions() TxOptions {
	return TxOptions{
		fromFactory: true,
	}
}

// TxOptions configures how transactions behave.
//
// See [github.com/geldata/gel-go.Client.WithTxOptions] for an example.
type TxOptions struct {
	// fromFactory indicates that a TxOptions value was created using
	// NewTxOptions() and not created directly with TxOptions{}.
	// Requiring users to use the factory function allows for nonzero
	// default values.
	fromFactory bool

	readOnly   string
	deferrable bool
	isolation  IsolationLevel
}

// ReadOnly returns the access mode name or "" if unset.
//
// This method is intended for internal use only and is not subject to semantic
// visioning guarantees.
func (o TxOptions) ReadOnly() string { return o.readOnly }

// Deferrable returns true if deferrable mode is set.
//
// This method is intended for internal use only and is not subject to semantic
// visioning guarantees.
func (o TxOptions) Deferrable() bool { return o.deferrable }

// IsolationLevel returns the TxOptions IsolationLevel setting.
//
// This method is intended for internal use only and is not subject to semantic
// visioning guarantees.
func (o TxOptions) IsolationLevel() IsolationLevel { return o.isolation }

// IsValid returns true if the TxOptions was created with NewTxOptions().
//
// This method is intended for internal use only and is not subject to semantic
// visioning guarantees.
func (o TxOptions) IsValid() bool { return o.fromFactory }

// WithIsolation returns a copy of the TxOptions
// with the isolation level set to i.
func (o TxOptions) WithIsolation(i IsolationLevel) TxOptions {
	switch i {
	case Serializable:
	case RepeatableRead:
	case PreferRepeatableRead:
	default:
		panic(fmt.Sprintf("unknown isolation level: %q", i))
	}

	o.isolation = i
	return o
}

// WithReadOnly returns a copy of the TxOptions with the transaction read only
// access mode set to r.
func (o TxOptions) WithReadOnly(r bool) TxOptions {
	if r {
		o.readOnly = "ReadOnly"
	} else {
		o.readOnly = "ReadWrite"
	}
	return o
}

// WithDeferrable returns a copy of the TxOptions with the transaction
// deferrable mode set to d.
func (o TxOptions) WithDeferrable(d bool) TxOptions {
	o.deferrable = d
	return o
}
