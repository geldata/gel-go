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

package geltypes_test

import (
	"fmt"
	"log"

	"github.com/geldata/gel-go/geltypes"
)

func ExampleOptional() {
	type User struct {
		geltypes.Optional
		Name string `gel:"name"`
	}

	var result User
	query := `
		SELECT User { name }
		FILTER .name = "doesn't exist"
		LIMIT 1
	`
	err := client.QuerySingle(ctx, query, &result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Missing())

	err = client.QuerySingle(ctx, `SELECT User { name } LIMIT 1`, &result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Missing())

	// Output:
	// true
	// false
}
