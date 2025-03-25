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

package gelcfg_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/geldata/gel-go/internal/testserver"
)

var (
	// Define identifiers that are used in examples, but keep them in their own
	// file. This way the definition is not included in the example. This keeps
	// examples concise.
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	opts := testserver.Options()
	err := os.Setenv("GEL_DSN", testserver.AsDSN(opts))
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}
