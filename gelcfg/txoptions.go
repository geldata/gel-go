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
	// Serializable is the only isolation level
	Serializable IsolationLevel = "serializable"
)

// NewTxOptions returns the default TxOptions value.
func NewTxOptions() TxOptions {
	return TxOptions{
		fromFactory: true,
		isolation:   Serializable,
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

	readOnly   bool
	deferrable bool
	isolation  IsolationLevel
}

// ReadOnly returns true if the read only access mode is set.
func (o TxOptions) ReadOnly() bool { return o.readOnly }

// Deferrable returns true if deferrable mode is set.
func (o TxOptions) Deferrable() bool { return o.deferrable }

// IsolationLevel returns the TxOptions IsolationLevel setting.
func (o TxOptions) IsolationLevel() IsolationLevel { return o.isolation }

// IsValid returns true if the TxOptions was created with NewTxOptions().
func (o TxOptions) IsValid() bool { return o.fromFactory }

// WithIsolation returns a copy of the TxOptions
// with the isolation level set to i.
func (o TxOptions) WithIsolation(i IsolationLevel) TxOptions {
	if i != Serializable {
		panic(fmt.Sprintf("unknown isolation level: %q", i))
	}

	o.isolation = i
	return o
}

// WithReadOnly returns a copy of the TxOptions with the transaction read only
// access mode set to r.
func (o TxOptions) WithReadOnly(r bool) TxOptions {
	o.readOnly = r
	return o
}

// WithDeferrable returns a copy of the TxOptions with the transaction
// deferrable mode set to d.
func (o TxOptions) WithDeferrable(d bool) TxOptions {
	o.deferrable = d
	return o
}
