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

package gelerr

// ErrorTag is the argument type to Error.HasTag().
type ErrorTag string

// ErrorCategory values represent Gel's error types.
type ErrorCategory string

// Error is the error type returned from gel.
type Error interface {
	Error() string
	Unwrap() error

	// HasTag returns true if the error is marked with the supplied tag.
	HasTag(ErrorTag) bool

	// Category returns true if the error is in the provided category.
	Category(ErrorCategory) bool
}
