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

// Package gel is the official Go driver for [Gel]. Additionally,
// [github.com/geldata/gel-go/cmd/edgeql-go] is a code generator that
// generates go functions from edgeql files.
//
// Typical client usage looks like this:
//
//	package main
//
//	import (
//	    "context"
//	    "log"
//
//	    "github.com/geldata/gel-go"
//	    "github.com/geldata/gel-go/gelcfg"
//	)
//
//	func main() {
//	    ctx := context.Background()
//	    client, err := gel.CreateClient(gelcfg.Options{})
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer client.Close()
//
//	    var (
//	        age   int64 = 21
//	        users []struct {
//	            ID   geltypes.UUID `gel:"id"`
//	            Name string   `gel:"name"`
//	        }
//	    )
//
//	    query := "SELECT User{name} FILTER .age = <int64>$0"
//	    err = client.Query(ctx, query, &users, age)
//	    ...
//	}
//
// We recommend using environment variables for connection parameters. See the
// [client connection docs] for more information.
//
// You may also connect to a database using a DSN:
//
//	dsn := "gel://edgedb@localhost/edgedb"
//	client, err := gel.CreateClientDSN(dsn, opts)
//
// Or you can use Option fields.
//
//	opts := gelcfg.Options{
//	    Database:    "edgedb",
//	    User:        "edgedb",
//	    Concurrency: 4,
//	}
//
//	client, err := gel.CreateClient(opts)
//
// # Errors
//
// gel never returns underlying errors directly.
// If you are checking for things like context expiration
// use [errors.Is] or [errors.As].
//
//	err := client.Query(...)
//	if errors.Is(err, context.Canceled) { ... }
//
// Most errors returned by the gel package will satisfy the gelerr.Error
// interface which has methods for introspecting.
//
//	err := client.Query(...)
//
//	var gelErr gelerr.Error
//	if errors.As(err, &gelErr) && gelErr.Category(gelcfg.NoDataError){
//	    ...
//	}
//
// # Datatypes
//
// The following list shows the marshal/unmarshal
// mapping between Gel types and go types. See also [geltypes]:
//
//	Gel                      Go
//	---------                ---------
//	Set                      []anytype
//	array<anytype>           []anytype
//	tuple                    struct
//	named tuple              struct
//	Object                   struct
//	bool                     bool, geltypes.OptionalBool
//	bytes                    []byte, geltypes.OptionalBytes
//	str                      string, geltypes.OptionalStr
//	anyenum                  string, geltypes.OptionalStr
//	datetime                 time.Time, geltypes.OptionalDateTime
//	cal::local_datetime      geltypes.LocalDateTime,
//	                         geltypes.OptionalLocalDateTime
//	cal::local_date          geltypes.LocalDate, geltypes.OptionalLocalDate
//	cal::local_time          geltypes.LocalTime, geltypes.OptionalLocalTime
//	duration                 geltypes.Duration, geltypes.OptionalDuration
//	cal::relative_duration   geltypes.RelativeDuration,
//	                         geltypes.OptionalRelativeDuration
//	float32                  float32, geltypes.OptionalFloat32
//	float64                  float64, geltypes.OptionalFloat64
//	int16                    int16, geltypes.OptionalFloat16
//	int32                    int32, geltypes.OptionalInt16
//	int64                    int64, geltypes.OptionalInt64
//	uuid                     geltypes.UUID, geltypes.OptionalUUID
//	json                     []byte, geltypes.OptionalBytes
//	bigint                   *big.Int, geltypes.OptionalBigInt
//
//	decimal                  user defined (see Custom Marshalers)
//
// Note that Gel's std::duration type is represented in int64 microseconds
// while go's time.Duration type is int64 nanoseconds. It is incorrect to cast
// one directly to the other.
//
// Shape fields that are not required must use optional types for receiving
// query results. The [gelcfg.Optional] struct can be embedded to make structs
// optional.
//
//	type User struct {
//	    gelcfg.Optional
//	    Email string `gel:"email"`
//	}
//
//	var result User
//	err := client.QuerySingle(ctx, `SELECT User { email } LIMIT 0`, $result)
//	fmt.Println(result.Missing())
//	// Output: true
//
//	err := client.QuerySingle(ctx, `SELECT User { email } LIMIT 1`, $result)
//	fmt.Println(result.Missing())
//	// Output: false
//
// Not all types listed above are valid query parameters.  To pass a slice of
// scalar values use array in your query. Gel doesn't currently support
// using sets as parameters.
//
//	query := `select User filter .id in array_unpack(<array<uuid>>$1)`
//	client.QuerySingle(ctx, query, $user, []geltypes.UUID{...})
//
// Nested structures are also not directly allowed but you can use [json]
// instead.
//
// By default Gel will ignore embedded structs when marshaling/unmarshaling.
// To treat an embedded struct's fields as part of the parent struct's fields,
// tag the embedded struct with `gel:"$inline"`.
//
//	type Object struct {
//	    ID geltypes.UUID
//	}
//
//	type User struct {
//	    Object `gel:"$inline"`
//	    Name string
//	}
//
// # Custom Marshalers
//
// Interfaces for user defined marshaler/unmarshalers  are documented in the
// [github.com/geldata/gel-go/internal/marshal] package.
//
// [Gel]: https://www.geldata.com
// [json]: https://docs.geldata.com/reference/edgeql/insert#bulk-inserts
// [client connection docs]: https://docs.geldata.com/learn/clients#connection
package gel
