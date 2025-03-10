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

//go:build tools

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/geldata/gel-go/internal/errgen"
)

func printCategories(types []*errgen.Type) {
	fmt.Print(`

const (`)

	for _, typ := range types {
		fmt.Printf(`
	%[1]v ErrorCategory = "errors::%[1]v"`, typ.Name)
	}

	fmt.Print(`
)`)
}

func printTags(tags []errgen.Tag) {
	fmt.Print(`

const (`)

	for _, tag := range tags {
		fmt.Printf(`
	%[1]v ErrorTag = %[2]q`, tag.Identifyer(), tag)
	}

	fmt.Print(`
)`)
}

//nolint:typecheck
func main() {
	var data [][]interface{}
	if e := json.NewDecoder(os.Stdin).Decode(&data); e != nil {
		log.Fatal(e)
	}

	types := errgen.ParseTypes(data)
	tags := errgen.ParseTags(data)

	fmt.Print(`// This source file is part of the Gel open source project.
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

// This file is auto generated. Do not edit!
// run 'make errors' to regenerate

`)

	fmt.Println()
	fmt.Println("package gelerr")
	fmt.Println()
	printTags(tags)
	printCategories(types)
}
