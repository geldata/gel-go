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

package geltypes

import "context"

// TxBlock is work to be done in a transaction.
type TxBlock func(context.Context, Tx) error

// Executor is a common interface between *gel.Client and Tx,
// that can run queries on a Gel database.
type Executor interface {
	Execute(context.Context, string, ...any) error
	ExecuteSQL(context.Context, string, ...any) error
	Query(context.Context, string, any, ...any) error
	QueryJSON(context.Context, string, *[]byte, ...any) error
	QuerySQL(context.Context, string, any, ...any) error
	QuerySingle(context.Context, string, any, ...any) error
	QuerySingleJSON(context.Context, string, any, ...any) error
}

// Tx is a transaction. Use gel.Client.Tx() to get a transaction.
type Tx interface {
	Executor
}
