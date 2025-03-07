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

	"github.com/geldata/gel-go/geltypes"
)

func ExampleClient_Tx() {
	err := client.Tx(ctx, func(ctx context.Context, tx geltypes.Tx) error {
		return tx.Execute(ctx, "INSERT User { name := 'Don' }")
	})
	if err != nil {
		log.Println(err)
	}
}
