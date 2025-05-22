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

package gel

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/internal"
	gelint "github.com/geldata/gel-go/internal/client"
	"github.com/geldata/gel-go/internal/testserver"
)

var (
	client          *Client
	once            sync.Once
	opts            gelcfg.Options
	protocolVersion internal.ProtocolVersion
)

func execOrFatal(command string) {
	ctx := context.Background()
	err := client.Execute(ctx, command)
	if err != nil {
		testserver.Fatal(err)
	}
}

func initServer() {
	defer log.Println("test server is ready for use")

	initClient()
	initProtocolVersion()

	log.Println("configuring instance")
	execOrFatal(`
		CONFIGURE INSTANCE SET session_idle_timeout := <duration>'1s';
	`)

	ctx := context.Background()
	err := client.Execute(ctx, `
		CREATE SUPERUSER ROLE user_with_password {
			SET password := 'secret';
		};
	`)
	if err != nil &&
		!strings.Contains(err.Error(), "'user_with_password' already exists") {
		testserver.Fatal(err)
	}
	execOrFatal(`
		CONFIGURE INSTANCE RESET Auth;
	`)
	execOrFatal(`
		CONFIGURE INSTANCE INSERT Auth {
			comment := "no password",
			priority := 1,
			method := (INSERT Trust),
			user := {'*'},
		};
	`)
	execOrFatal(`
		CONFIGURE INSTANCE INSERT Auth {
			comment := "password required",
			priority := 0,
			method := (INSERT SCRAM),
			user := {'user_with_password'}
		};
	`)

	log.Println("running migration")
	execOrFatal(`
		START MIGRATION TO {
			module default {
				global global_id -> uuid;
				required global global_str -> str {
					default := "default";
				};
				global global_bytes -> bytes;
				global global_int16 -> int16;
				global global_int32 -> int32;
				global global_int64 -> int64;
				global global_float32 -> float32;
				global global_float64 -> float64;
				global global_bool -> bool;
				global global_datetime -> datetime;
				global global_duration -> duration;
				global global_json -> json;
				global global_local_datetime -> cal::local_datetime;
				global global_local_date -> cal::local_date;
				global global_local_time -> cal::local_time;
				global global_bigint -> bigint;
				global global_relative_duration -> cal::relative_duration;
				global global_date_duration -> cal::date_duration;
				global global_memory -> cfg::memory;

				type User {
					property name -> str;
				}
				type TxTest {
					required property name -> str;
				}
				type Counter {
					name: str {
						constraint exclusive;
					};
					value: int32 {
						default := 0;
					}
				}
			}
		};
		POPULATE MIGRATION;
		COMMIT MIGRATION;
	`)
}

func initClient() {
	opts = testserver.Options()

	log.Println("initializing testserver.Client")
	var err error
	client, err = CreateClient(opts)
	if err != nil {
		testserver.Fatal(err)
	}
}

func initProtocolVersion() {
	log.Println("initializing testserver.ProtocolVersion")
	var err error
	protocolVersion, err = gelint.ProtocolVersion(
		context.Background(),
		client.pool,
	)
	if err != nil {
		testserver.Fatal(err)
	}
	log.Printf("using Protocol Version: %v", protocolVersion)
}
