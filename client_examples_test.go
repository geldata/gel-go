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
	"fmt"
	"log"

	gel "github.com/geldata/gel-go"
)

func ExampleCreateClient() {
	client, err := gel.CreateClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	var result int64
	err = client.QuerySingle(ctx, `SELECT 1`, &result)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
	// Output: 1
}

func ExampleClient_Execute() {
	err := client.Execute(ctx, "INSERT Product")
	if err != nil {
		log.Fatal(err)
	}

	// Output:
}

func ExampleClient_ExecuteSQL() {
	err := client.ExecuteSQL(ctx, `INSERT INTO "Product" DEFAULT VALUES`)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
}

func ExampleClient_Query() {
	var output []struct {
		Result int64 `gel:"result"`
	}

	err := client.Query(ctx, `SELECT {result := 2 + 2}`, &output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(output)
	// Output: [{4}]
}

func ExampleClient_QueryJSON() {
	var output []byte
	err := client.QueryJSON(ctx, `SELECT {result := 2 + 2}`, &output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
	// Output: [{"result" : 4}]
}

func ExampleClient_QuerySQL() {
	var output []struct {
		Result int32 `gel:"result"`
	}

	err := client.QuerySQL(ctx, `SELECT 2 + 2 AS result`, &output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
	// Output: [{4}]
}

func ExampleClient_QuerySingle() {
	var output struct {
		Result int64 `gel:"result"`
	}

	err := client.QuerySingle(ctx, `SELECT {result := 2 + 2}`, &output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(output)
	// Output: {4}
}

func ExampleClient_QuerySingleJSON() {
	var output []byte
	err := client.QuerySingleJSON(ctx, `SELECT {result := 2 + 2}`, &output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
	// Output: {"result" : 4}
}
