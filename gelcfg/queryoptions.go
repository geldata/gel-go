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

func NewQueryOptions() QueryOptions {
	return QueryOptions{fromFactory: true}
}

type QueryOptions struct {
	fromFactory bool

	readOnly      bool
	implicitLimit uint64
}

func (o QueryOptions) WithReadOnly(readOnly bool) QueryOptions {
	o.readOnly = readOnly
	return o
}

func (o QueryOptions) WithImplicitLimit(limit uint64) QueryOptions {
	o.implicitLimit = limit
	return o
}

func (o QueryOptions) ReadOnly() bool        { return o.readOnly }
func (o QueryOptions) ImplicitLimit() uint64 { return o.implicitLimit }
