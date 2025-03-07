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

// NewQueryOptions returns the default QueryOptions value with readOnly set to
// false, and implicitLimit set to 0.
func NewQueryOptions() QueryOptions {
	return QueryOptions{fromFactory: true}
}

// QueryOptions controls limitations that the server can impose on
// queries.
type QueryOptions struct {
	fromFactory bool

	readOnly      bool
	implicitLimit uint64
}

// WithReadOnly enables read only mode if readOnly is true.
func (o QueryOptions) WithReadOnly(readOnly bool) QueryOptions {
	o.readOnly = readOnly
	return o
}

// WithImplicitLimit sets the max number of results that the server will
// return. If set to 0 the server will return all results. Defaults to 0.
func (o QueryOptions) WithImplicitLimit(limit uint64) QueryOptions {
	o.implicitLimit = limit
	return o
}

// ReadOnly returns true if read only mode is enabled.
func (o QueryOptions) ReadOnly() bool { return o.readOnly }

// ImplicitLimit  returns the configured implicit limit.
func (o QueryOptions) ImplicitLimit() uint64 { return o.implicitLimit }
