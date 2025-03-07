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

// edgeql-go is a tool to generate go functions from edgeql queries.
//
// When run in a Gel project directory or subdirectory an *_edgeql.go source
// file will be generated for each *.edgeql file.  The generated go will have a
// query and queryJSON function with typed arguments and return value matching
// the query's arguments and result shape.  For example if a directory contains
// a get_user.edgeql file, edgeql-go will create a get_user_edgeql.go file
// with a getUser(...) and getUserJSON(...) function.
//
// # Install
//
// For go 1.24 and above:
//
//	go get -tool github.com/geldata/gel-go/cmd/edgeql-go@latest
//
// For go 1.23 and below:
//
//	go install github.com/geldata/gel-go/cmd/edgeql-go@latest
//
// # Usage
//
// Typically this process would be run using [go generate].
// For go 1.24 and above:
//
//	//go:generate go tool edgeql-go -pubfuncs -pubtypes -mixedcaps
//
// For go 1.23 and below:
//
//	//go:generate edgeql-go -pubfuncs -pubtypes -mixedcaps
//
// For a complete list of options:
//
//	edgeql-go -help
//
// [go generate]: https://go.dev/blog/generate
package main
