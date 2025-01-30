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

//go:build gendocs

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"
)

var lintMode = flag.Bool("lint", false, "Instead of writing output files, "+
	"check if contents of existing files match")

func main() {
	flag.Parse()

	if err := os.Mkdir("rstdocs", 0750); err != nil && !os.IsExist(err) {
		panic(err)
	}

	renderIndexPage()
	apiTypesInfo, datatypesInfo := gatherAPITypes()
	renderAPIPage(apiTypesInfo)
	renderDatatypesPage(datatypesInfo)
	renderCodegenPage()
}

func readAndParseFile(fset *token.FileSet, filename string) *ast.File {
	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	ast, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return ast
}

func writeFile(filename string, content string) {
	if *lintMode {
		if file, err := os.ReadFile(filename); err != nil ||
			string(file) != content {
			panic("Content of " + filename + " does not match generated " +
				"docs, Run 'make gendocs' to update docs")
		}
	} else {
		if err := os.WriteFile(
			filename, []byte(content), 0666); err != nil {
			panic(err)
		}
	}
}

func renderDecl(decl any, fset *token.FileSet) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, decl); err != nil {
		panic(err)
	}
	return buf.String()
}

func renderDoc(p *doc.Package, doc string) string {
	return strings.TrimSpace(string(printRST(p.Printer(), p.Parser().Parse(doc))))
}

type TypesInfo struct {
	FuncNames []string
	TypeNames []string
}

func gatherAPITypes() (TypesInfo, TypesInfo) {
	fset := token.NewFileSet()
	files := []*ast.File{
		readAndParseFile(fset, "export.go"),
	}

	p, err := doc.NewFromFiles(fset, files, "github.com/edgedb/edgedb-go")
	if err != nil {
		panic(err)
	}

	apiTypesInfo := TypesInfo{}
	datatypesInfo := TypesInfo{}

	for _, v := range p.Vars {
		lines := strings.Split(renderDecl(v.Decl, fset), "\n")
		for _, line := range lines[1 : len(lines)-1] {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "//") {
				continue
			}
			typeParts := strings.Split(strings.Split(trimmed, "=")[1], ".")

			switch strings.TrimSpace(typeParts[0]) {
			case "edgedb":
				apiTypesInfo.FuncNames = append(apiTypesInfo.FuncNames, strings.TrimSpace(typeParts[1]))
			case "edgedbtypes":
				datatypesInfo.FuncNames = append(datatypesInfo.FuncNames, strings.TrimSpace(typeParts[1]))
			default:
				panic(fmt.Errorf("unknown internal module %s", typeParts[0]))
			}
		}
	}

	for _, t := range p.Types {
		decl := strings.TrimPrefix(renderDecl(t.Decl, fset), "type ")
		typeParts := strings.Split(strings.Split(decl, "=")[1], ".")

		switch strings.TrimSpace(typeParts[0]) {
		case "edgedb":
			apiTypesInfo.TypeNames = append(apiTypesInfo.TypeNames, strings.TrimSpace(typeParts[1]))
		case "edgedbtypes":
			datatypesInfo.TypeNames = append(datatypesInfo.TypeNames, strings.TrimSpace(typeParts[1]))
		default:
			panic(fmt.Errorf("unknown internal module %s", typeParts[0]))
		}
	}

	return apiTypesInfo, datatypesInfo
}

func renderAPIPage(info TypesInfo) {
	dir, err := os.ReadDir("internal/client")
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	files := []*ast.File{}

	for _, file := range dir {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			files = append(
				files, readAndParseFile(
					fset, "internal/client/"+file.Name()))
		}
	}

	p, err := doc.NewFromFiles(
		fset, files, "github.com/edgedb/edgedb-go/internal/client")
	if err != nil {
		panic(err)
	}

	rst := `
API
===`

	rst += renderTypes(fset, p, info)

	writeFile("rstdocs/api.rst", rst)
}

func renderDatatypesPage(info TypesInfo) {
	dir, err := os.ReadDir("internal/edgedbtypes")
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	files := []*ast.File{}

	for _, file := range dir {
		if !file.IsDir() {
			files = append(
				files, readAndParseFile(
					fset, "internal/edgedbtypes/"+file.Name()))
		}
	}

	p, err := doc.NewFromFiles(
		fset, files, "github.com/edgedb/edgedb-go/internal/edgedbtypes")
	if err != nil {
		panic(err)
	}

	rst := `
Datatypes
=========`

	rst += renderTypes(fset, p, info)

	writeFile("rstdocs/types.rst", rst)
}

func renderTypes(
	fset *token.FileSet,
	p *doc.Package,
	info TypesInfo,
) string {
	out := ""

	for _, t := range p.Types {
		if !slices.Contains(info.TypeNames, t.Name) {
			continue
		}

		for _, f := range t.Funcs {
			if !slices.Contains(info.FuncNames, f.Name) {
				continue
			}

			out += fmt.Sprintf(`


.. go:function:: %s

    `,
				strings.ReplaceAll(
					renderDecl(f.Decl, fset),
					"\n", "\\\n    ",
				))

			out += strings.ReplaceAll(renderDoc(p, f.Doc), "\n", "\n    ")
		}

		out += fmt.Sprintf(`


.. go:type:: %s

    `, strings.ReplaceAll(
			renderDecl(t.Decl, fset),
			"\n", "\\\n    ",
		))

		out += strings.ReplaceAll(renderDoc(p, t.Doc), "\n", "\n    ")

		for _, m := range t.Methods {
			out += fmt.Sprintf(`


.. go:method:: %s

    `, strings.ReplaceAll(
				renderDecl(m.Decl, fset),
				"\n", "\\\n    ",
			))

			out += strings.ReplaceAll(renderDoc(p, m.Doc), "\n", "\n    ")
		}
	}

	for _, f := range p.Funcs {
		if !slices.Contains(info.FuncNames, f.Name) {
			continue
		}

		out += fmt.Sprintf(`


.. go:function:: %s

    `,
			strings.ReplaceAll(
				renderDecl(f.Decl, fset),
				"\n", "\\\n    ",
			))

		out += strings.ReplaceAll(renderDoc(p, f.Doc), "\n", "\n    ")

	}

	return strings.ReplaceAll(out, "\t", "    ")
}

func renderIndexPage() {
	fset := token.NewFileSet()
	files := []*ast.File{
		readAndParseFile(fset, "doc.go"),
		readAndParseFile(fset, "doc_test.go"),
	}

	p, err := doc.NewFromFiles(fset, files, "github.com/edgedb/edgedb-go")
	if err != nil {
		panic(err)
	}

	rst := `.. _edgedb-go-intro:

================
EdgeDB Go Driver
================


.. toctree::
   :maxdepth: 3
   :hidden:

   api
   types
   codegen



`

	rst += string(printRST(p.Printer(), p.Parser().Parse(p.Doc)))

	rst += `

Usage Example
-------------

.. code-block:: go
`

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, p.Examples[0].Code); err != nil {
		panic(err)
	}

	exampleLines := strings.Split(buf.String(), "\n")

	skip := true
	for _, line := range exampleLines {
		if skip && !strings.HasPrefix(line, "//") {
			skip = false
		}
		if !skip {
			rst += "    " + strings.ReplaceAll(line, "\t", "    ") + "\n"
		}
	}

	writeFile("rstdocs/index.rst", rst)
}

func renderCodegenPage() {
	fset := token.NewFileSet()
	files := []*ast.File{
		readAndParseFile(fset, "cmd/edgeql-go/doc.go"),
	}

	p, err := doc.NewFromFiles(
		fset, files, "github.com/edgedb/edgedb-go/cmd/edgeql-go")
	if err != nil {
		panic(err)
	}

	rst := `
Codegen
=======
`

	rst += string(printRST(p.Printer(), p.Parser().Parse(p.Doc)))

	writeFile("rstdocs/codegen.rst", rst)
}
