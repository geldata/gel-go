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

// This file is auto generated. Do not edit!
// run 'go generate ./...' to regenerate

package gel

import (
	gel "github.com/edgedb/edgedb-go/internal/client"
)

type (
	// Tx is a transaction. Use Client.Tx() to get a transaction.
	Tx = gel.Tx

	// TxBlock is work to be done in a transaction.
	TxBlock = gel.TxBlock
)

var (
	// LogWarnings is an gel.WarningHandler that logs warnings.
	LogWarnings = gel.LogWarnings

	// WarningsAsErrors is an gel.WarningHandler that returns warnings as
	// errors.
	WarningsAsErrors = gel.WarningsAsErrors
)
