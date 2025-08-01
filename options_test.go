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
	"fmt"
	"log"
	"time"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
)

//nolint:lll
func ExampleClient_WithTxOptions() {
	opts := gelcfg.NewTxOptions().WithReadOnly(true)
	configured := client.WithTxOptions(opts)

	err := configured.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		return tx.Execute(ctx, "INSERT User")
	})
	fmt.Println(err)
	// Output:
	// gel.TransactionError: Modifications not allowed in a read-only transaction
}

func ExampleClient_WithRetryOptions() {
	linearBackoff := func(n int) time.Duration { return time.Second }

	rule := gelcfg.NewRetryRule().
		WithAttempts(5).
		WithBackoff(linearBackoff)
	opts := gelcfg.NewRetryOptions().WithDefault(rule)
	configured := client.WithRetryOptions(opts)

	err := configured.Execute(
		ctx,
		"INSERT Product { name := 'shiny and new' }",
	)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_WithQueryTag() {
	tag := "my-app/backend"
	configured, err := client.WithQueryTag(tag)
	if err != nil {
		log.Fatal(err)
	}

	err = configured.Execute(ctx, "SELECT User { ** }")
	if err != nil {
		log.Fatal(err)
	}

	var query string
	err = client.QuerySingle(
		ctx,
		`
	SELECT (
		SELECT sys::QueryStats
		FILTER .tag = <str>$0
		LIMIT 1
	).query
	`,
		&query,
		tag,
	)
	if err != nil {
		log.Fatal(err)
	}

	// sys::QueryStats reformats queries
	fmt.Println(query)
	// Output:
	// select
	//     User {
	//         **
	//     }
}

func ExampleClient_WithoutQueryTag() {
	tag := "my-app/api"
	configured, err := client.WithQueryTag(tag)
	if err != nil {
		log.Fatal(err)
	}

	unconfigured := configured.WithoutQueryTag()

	err = unconfigured.Execute(ctx, "SELECT Product { ** }")
	if err != nil {
		log.Fatal(err)
	}

	var queries []string
	err = client.Query(
		ctx,
		`
	SELECT (
		SELECT sys::QueryStats
		FILTER .tag = <str>$0
	).query
	`,
		&queries,
		tag,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(queries)
	// Output: []
}

func ExampleClient_WithConfig() {
	configured := client.WithConfig(map[string]any{
		"allow_user_specified_id": true,
	})

	err := configured.Execute(ctx, `INSERT USER { id := <uuid>$0 }`, id)
	if err != nil {
		log.Fatal(err)
	}
}

//nolint:lll
func ExampleClient_WithoutConfig() {
	configured := client.WithConfig(map[string]any{
		"allow_user_specified_id": true,
	})

	unconfigured := configured.WithoutConfig("allow_user_specified_id")

	err := unconfigured.Execute(ctx, `INSERT User { id := <uuid>$0 }`, id)
	fmt.Println(err)
	// Output:
	// gel.QueryError: cannot assign to property 'id'
	// query:1:15
	//
	// INSERT User { id := <uuid>$0 }
	//               ^ consider enabling the "allow_user_specified_id" configuration parameter to allow setting custom object ids
}

func ExampleClient_WithGlobals() {
	configured := client.WithGlobals(map[string]any{
		"used_everywhere": int64(42),
	})

	var result int64
	err := configured.QuerySingle(
		ctx,
		"SELECT GLOBAL used_everywhere",
		&result,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// Output: 42
}

func ExampleClient_WithoutGlobals() {
	configured := client.WithGlobals(map[string]any{
		"used_everywhere": int64(42),
	})

	unconfigured := configured.WithoutGlobals("used_everywhere")

	var result int64
	err := unconfigured.QuerySingle(
		ctx,
		`SELECT GLOBAL used_everywhere`,
		&result,
	)
	fmt.Println(err)
	// Output: gel.NoDataError: zero results
}

func ExampleClient_WithModuleAliases() {
	configured := client.WithModuleAliases(
		gelcfg.ModuleAlias{
			Module: "math",
			Alias:  "m",
		},
	)

	var result int64
	err := configured.QuerySingle(ctx, "SELECT m::abs(-42)", &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// Output: 42
}

func ExampleClient_WithoutModuleAliases() {
	configured := client.WithModuleAliases(
		gelcfg.ModuleAlias{
			Module: "math",
			Alias:  "m",
		},
	)

	unconfigured := configured.WithoutModuleAliases("m")

	var result int64
	err := unconfigured.QuerySingle(ctx, "SELECT m::abs(-42)", &result)
	fmt.Println(err)
	// Output:
	// gel.InvalidReferenceError: function 'm::abs' does not exist
	// query:1:8
	//
	// SELECT m::abs(-42)
	//        ^ error
}

//nolint:lll
func ExampleClient_WithQueryOptions() {
	opts := gelcfg.NewQueryOptions().WithReadOnly(true)
	configured := client.WithQueryOptions(opts)

	err := configured.Execute(ctx, "INSERT User")
	fmt.Println(err)
	// Output:
	// gel.DisabledCapabilityError: cannot execute data modification queries: disabled by the client
}

func ExampleClient_WithWarningHandler() {
	handler := func(warnings []error) error {
		for _, warning := range warnings {
			fmt.Println(warning)
		}
		return nil
	}

	configured := client.WithWarningHandler(handler)
	err := configured.Execute(ctx, `SELECT _warn_on_call()`)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// gel.QueryError: Test warning please ignore
	// query:1:8
	//
	// SELECT _warn_on_call()
	//        ^ error
}
