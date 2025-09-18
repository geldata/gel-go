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

package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/geldata/gel-go/geltypes"
	"github.com/geldata/gel-go/internal"
	gelint "github.com/geldata/gel-go/internal/client"
	"github.com/geldata/gel-go/internal/descriptor"
)

var (
	errQueryNum = errors.New(
		"numbered query arguments detected, use named arguments instead",
	)
)

type queryConfig struct {
	name    string
	method  string
	file    string
	imports []string
	structs []*goStruct
	sTypes  *goStruct
	rTypes  []goType
}

type queryConfigV1 struct{}

type queryConfigV2 struct{}

type querySetup interface {
	setup(
		ctx context.Context,
		cmd,
		qryFile,
		outFile string,
		cfg *cmdConfig,
		c *gelint.Pool,
	) (*queryConfig, error)
}

func newQuery(
	ctx context.Context,
	p *gelint.Pool,
	qryFile,
	outFile string,
	cfg *cmdConfig,
) (*Query, error) {
	var err error
	qryFile, err = filepath.Abs(qryFile)
	if err != nil {
		return nil, err
	}

	queryBytes, err := os.ReadFile(qryFile)
	if err != nil {
		log.Fatalf("error reading %q: %s", qryFile, err)
	}

	var qs querySetup

	if cfg.protocolVersion.GTE(internal.ProtocolVersion{Major: 2, Minor: 0}) {
		qs = &queryConfigV2{}
	} else {
		qs = &queryConfigV1{}
	}

	q, err := qs.setup(ctx, string(queryBytes), qryFile, outFile, cfg, p)
	if err != nil {
		log.Fatalf("failed to setup query: %s", err)
	}

	return &Query{
		imports:             q.imports,
		QueryFile:           q.file,
		QueryName:           q.name,
		CMDVarName:          cmdVarName(qryFile),
		ResultTypes:         q.structs,
		SignatureReturnType: q.rTypes[0].Reference(),
		SignatureArgs:       q.sTypes.Fields,
		Method:              q.method,
	}, nil
}

func (r *queryConfigV1) setup(
	ctx context.Context,
	cmd,
	qryFile,
	outFile string,
	cmdCfg *cmdConfig,
	p *gelint.Pool,
) (*queryConfig, error) {
	description, err := gelint.Describe(ctx, p, cmd, qryFile)

	if err != nil {
		return nil, fmt.Errorf("error introspecting query %q: %s", qryFile,
			err)
	}

	if isNumberedArgs(description.In) {
		return nil, errQueryNum
	}

	qryName := queryName(qryFile, cmdCfg)
	rTypes, imports, err := resultTypes(qryFile, description, cmdCfg)
	if err != nil {
		log.Fatal(err)
	}
	var rStructs []*goStruct
	for _, typ := range rTypes {
		if t, ok := typ.(*goStruct); ok {
			t.QueryFuncName = qryName
			rStructs = append(rStructs, t)
		}
	}

	sTypes, i, err := signatureTypes(description, cmdCfg)
	if err != nil {
		log.Fatal(err)
	}
	imports = append(imports, i...)

	qryFile, err = queryFile(outFile, qryFile)
	if err != nil {
		log.Fatal(err)
	}

	m, err := method(description)
	if err != nil {
		log.Fatal(err)
	}

	return &queryConfig{
		name:    qryName,
		method:  m,
		file:    qryFile,
		imports: imports,
		structs: rStructs,
		sTypes:  sTypes,
		rTypes:  rTypes,
	}, nil
}

func (r *queryConfigV2) setup(
	ctx context.Context,
	cmd,
	qryFile,
	outFile string,
	cmdCfg *cmdConfig,
	p *gelint.Pool,
) (*queryConfig, error) {
	description, err := gelint.DescribeV2(ctx, p, cmd, qryFile)

	if err != nil {
		return nil, fmt.Errorf("error introspecting query %q: %s", qryFile,
			err)
	}

	if isNumberedArgsV2(&description.In) {
		return nil, errQueryNum
	}

	qryName := queryName(qryFile, cmdCfg)
	rTypes, imports, err := resultTypesV2(qryFile, description, cmdCfg)
	if err != nil {
		log.Fatal(err)
	}
	var rStructs []*goStruct
	for _, typ := range rTypes {
		if t, ok := typ.(*goStruct); ok {
			t.QueryFuncName = qryName
			rStructs = append(rStructs, t)
		}
	}

	sTypes, i, err := signatureTypesV2(description, cmdCfg)
	if err != nil {
		log.Fatal(err)
	}
	imports = append(imports, i...)

	qryFile, err = queryFile(outFile, qryFile)
	if err != nil {
		log.Fatal(err)
	}

	m, err := methodV2(description)
	if err != nil {
		log.Fatal(err)
	}

	return &queryConfig{
		name:    qryName,
		method:  m,
		file:    qryFile,
		imports: imports,
		structs: rStructs,
		sTypes:  sTypes,
		rTypes:  rTypes,
	}, nil
}

func queryFile(outFile, queryFile string) (string, error) {
	return filepath.Rel(filepath.Dir(outFile), queryFile)
}

func cmdVarName(qryFile string) string {
	name := filepath.Base(qryFile)
	name = strings.TrimSuffix(name, ".edgeql")
	name = fmt.Sprintf("%s_cmd", name)
	return snakeToLowerMixedCase(name)
}

func queryName(qryFile string, cmdCfg *cmdConfig) string {
	name := filepath.Base(qryFile)
	name = strings.TrimSuffix(name, ".edgeql")
	if cmdCfg.pubfuncs {
		return snakeToUpperMixedCase(name)
	}
	return snakeToLowerMixedCase(name)
}

func typeName(qryFile string, cmdCfg *cmdConfig) string {
	name := filepath.Base(qryFile)
	name = strings.TrimSuffix(name, ".edgeql")
	if cmdCfg.pubtypes {
		return snakeToUpperMixedCase(name)
	}
	return snakeToLowerMixedCase(name)
}

func signatureTypes(
	description *gelint.CommandDescription,
	cmdCfg *cmdConfig,
) (*goStruct, []string, error) {
	types, imports, err := generateType(
		description.In,
		true,
		nil,
		cmdCfg,
		false,
		false,
	)
	if err != nil {
		return &goStruct{}, nil, err
	}

	return types[0].(*goStruct), imports, nil
}

func signatureTypesV2(
	description *gelint.CommandDescriptionV2,
	cmdCfg *cmdConfig,
) (*goStruct, []string, error) {
	types, imports, err := generateTypeV2(
		&description.In,
		true,
		nil,
		cmdCfg,
		false,
		false,
	)
	if err != nil {
		return &goStruct{}, nil, err
	}

	return types[0].(*goStruct), imports, nil
}

func resultTypes(
	qryFile string,
	description *gelint.CommandDescription,
	cmdCfg *cmdConfig,
) ([]goType, []string, error) {
	outDesc := description.Out
	var required bool
	switch description.Card {
	case gelint.Many, gelint.AtLeastOne:
		id, err := randomID()
		if err != nil {
			return nil, nil, err
		}

		required = true
		outDesc = descriptor.Descriptor{
			Type: descriptor.Set,
			ID:   id,
			Fields: []*descriptor.Field{{
				Desc: description.Out,
			}},
		}
	case gelint.One:
		required = true
	}

	name := typeName(qryFile, cmdCfg)
	return generateType(
		outDesc,
		required,
		[]string{name + "Result"},
		cmdCfg,
		true,
		false,
	)
}

func resultTypesV2(
	qryFile string,
	description *gelint.CommandDescriptionV2,
	cmdCfg *cmdConfig,
) ([]goType, []string, error) {
	outDesc := description.Out
	var required bool
	switch description.Card {
	case gelint.Many, gelint.AtLeastOne:
		id, err := randomID()
		if err != nil {
			return nil, nil, err
		}

		required = true
		outDesc = descriptor.V2{
			Type: descriptor.Set,
			ID:   id,
			Fields: []*descriptor.FieldV2{{
				Desc: description.Out,
			}},
		}
	case gelint.One:
		required = true
	}

	name := typeName(qryFile, cmdCfg)
	return generateTypeV2(
		&outDesc,
		required,
		[]string{name + "Result"},
		cmdCfg,
		true,
		false,
	)
}

func randomID() (geltypes.UUID, error) {
	var id geltypes.UUID
	_, err := rand.Read(id[:])
	return id, err
}

func method(description *gelint.CommandDescription) (string, error) {
	switch description.Card {
	case gelint.AtMostOne, gelint.One:
		return "QuerySingle", nil
	case gelint.NoResult, gelint.Many, gelint.AtLeastOne:
		return "Query", nil
	default:
		return "", errors.New("unreachable 20135")
	}
}

func methodV2(description *gelint.CommandDescriptionV2) (string, error) {
	switch description.Card {
	case gelint.AtMostOne, gelint.One:
		return "QuerySingle", nil
	case gelint.NoResult, gelint.Many, gelint.AtLeastOne:
		return "Query", nil
	default:
		return "", errors.New("unreachable 20135")
	}
}
