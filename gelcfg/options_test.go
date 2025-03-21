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
	"fmt"
	"log"
	"time"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
)

func ExampleOptions() {
	opts := gelcfg.Options{
		ConnectTimeout:     60 * time.Second,
		WaitUntilAvailable: 5 * time.Second,
		WarningHandler:     gelcfg.WarningsAsErrors,
	}

	client, err := gel.CreateClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	var message string
	err = client.QuerySingle(ctx, `SELECT "hello Gel"`, &message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(message)
	// Output: hello Gel
}
