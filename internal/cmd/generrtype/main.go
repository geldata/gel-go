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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/edgedb/edgedb-go/internal/errgen"
)

func printError(errType *errgen.Type) {
	fmt.Printf(`

func New%[1]s(msg string, err error) error {
	return &%[1]s{msg, err}
}

type %[1]v struct {
	msg string
	err error
}

func (e *%[1]v) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.%[1]v: " + msg
}

func (e *%[1]v) Unwrap() error { return e.err }
`, errType.Name)

	fmt.Printf(`

func (e *%v) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.%v:
		return true`, errType.Name, errType.Name)

	for _, ancestor := range errType.Ancestors {
		fmt.Printf(`
	case gelerr.%v:
		return true`, ancestor)
	}

	fmt.Print(`
	default:
		return false
	}
}
`)
	for _, ancestor := range errType.Ancestors {
		fmt.Printf(`
func (e *%v) isEdgeDB%v() {}
`, errType.Name, ancestor)
	}

	fmt.Printf(`
func (e *%v) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {`, errType.Name)

	for _, tag := range errType.Tags {
		fmt.Printf(`
	case gelerr.%v:
		return true`, tag.Identifyer())
	}

	fmt.Printf(`
	default:
		return false
	}
}`)
}

func printErrors(types []*errgen.Type) {
	for _, typ := range types {
		printError(typ)
	}
}

func printCodeMap(types []*errgen.Type) {
	fmt.Print(`

func ErrorFromCode(code uint32, msg string) error {
	switch code {`)

	for _, typ := range types {
		fmt.Printf(`
	case 0x%02x_%02x_%02x_%02x:
		return &%v{msg: msg}`,
			typ.Code[0], typ.Code[1], typ.Code[2], typ.Code[3],
			typ.Name,
		)
	}
	code := `
	default:
		return &UnexpectedMessageError{
			msg: fmt.Sprintf(
				"invalid error code 0x%x with message %q", code, msg,
			),
		}
	}
}`
	fmt.Print(code)
}

//nolint:typecheck
func main() {
	var data [][]interface{}
	if e := json.NewDecoder(os.Stdin).Decode(&data); e != nil {
		log.Fatal(e)
	}

	types := errgen.ParseTypes(data)

	fmt.Print(`// This source file is part of the EdgeDB open source project.
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
// run 'make errors' to regenerate

// internal/cmd/export should ignore this file
//go:build !export
`)

	fmt.Println()
	fmt.Println("package gelerr")
	fmt.Println()
	fmt.Print(`import (
	"fmt"
	"github.com/edgedb/edgedb-go/gelerr"
	)`)
	printErrors(types)
	printCodeMap(types)
}
