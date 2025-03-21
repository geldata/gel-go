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

package gel_test

import (
	"context"
	"log"

	gel "github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
	"github.com/geldata/gel-go/internal/testserver"
)

var (
	// Define identifiers that are used in examples, but keep them in their own
	// file. This way the definition is not included in the example. This keeps
	// examples concise.
	ctx    context.Context
	client *gel.Client
	opts   gelcfg.Options
	id     geltypes.UUID
)

func init() {
	ctx = context.Background()
	opts = testserver.Options()
	var err error
	client, err = gel.CreateClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Execute(ctx, `
		START MIGRATION TO {
			module default {
				# ExampleClient_WithGlobals
				global used_everywhere: int64;

				type User {
					required name: str {
						default := 'default';
					};
				};
				type Product {};
			}
		};
		POPULATE MIGRATION;
		COMMIT MIGRATION;
	`)
	if err != nil {
		log.Fatal(err)
	}

	// The server sends warnings in response to both parse and execute.  We
	// don't want the same warning to show up twice in
	// ExampleClient_WithWarningHandler, so make sure this query is cached
	// before the example is run.
	err = client.
		// Don't log the warning.
		WithWarningHandler(func(_ []error) error { return nil }).
		Execute(ctx, `SELECT _warn_on_call()`)
	if err != nil {
		log.Fatal(err)
	}
}
