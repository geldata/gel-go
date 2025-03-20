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
	"fmt"
	"strings"

	"github.com/geldata/gel-go/internal/codecs"
	"github.com/geldata/gel-go/internal/descriptor"
)

func generateType(
	desc descriptor.Descriptor,
	required bool,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
	isField bool,
) ([]goType, []string, error) {
	var (
		err     error
		types   []goType
		imports []string
	)

	switch desc.Type {
	case descriptor.Set, descriptor.Array:
		types, imports, err = generateSlice(
			desc,
			path,
			cmdCfg,
			isResult,
			isField,
		)
	case descriptor.Object, descriptor.NamedTuple:
		types, imports, err = generateObject(
			desc,
			required,
			path,
			cmdCfg,
			isResult,
		)
	case descriptor.Tuple:
		types, imports, err = generateTuple(
			desc,
			required,
			path,
			cmdCfg,
			isResult,
		)
	case descriptor.BaseScalar, descriptor.Scalar, descriptor.Enum:
		types, imports, err = generateBaseScalar(
			desc,
			required,
			cmdCfg,
			isResult,
			isField,
		)
	case descriptor.Range:
		types, imports, err = generateRange(desc, required)
	default:
		err = fmt.Errorf(
			"generating type: unknown descriptor type %v",
			desc.Type,
		)
	}

	if err != nil {
		return nil, nil, err
	}

	return types, imports, nil
}

func generateTypeV2(
	desc *descriptor.V2,
	required bool,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
	isField bool,
) ([]goType, []string, error) {
	var (
		err     error
		types   []goType
		imports []string
	)

	switch desc.Type {
	case descriptor.Set, descriptor.Array:
		types, imports, err = generateSliceV2(
			desc,
			path,
			cmdCfg,
			isResult,
			isField,
		)
	case descriptor.Object, descriptor.NamedTuple:
		types, imports, err = generateObjectV2(
			desc,
			required,
			path,
			cmdCfg,
			isResult,
		)
	case descriptor.Tuple:
		types, imports, err = generateTupleV2(
			desc,
			required,
			path,
			cmdCfg,
			isResult,
		)
	case descriptor.BaseScalar, descriptor.Scalar, descriptor.Enum:
		types, imports, err = generateBaseScalarV2(
			desc,
			required,
			cmdCfg,
			isResult,
			isField,
		)
	case descriptor.Range:
		types, imports, err = generateRangeV2(desc, required)
	default:
		err = fmt.Errorf(
			"generating type: unknown descriptor type %v",
			desc.Type,
		)
	}

	if err != nil {
		return nil, nil, err
	}

	return types, imports, nil
}

func generateRange(
	desc descriptor.Descriptor,
	required bool,
) ([]goType, []string, error) {
	optional := ""
	if !required {
		optional = "Optional"
	}

	var typ string
	fieldDesc := desc.Fields[0].Desc
	switch fieldDesc.ID {
	case codecs.Int32ID:
		typ = "Int32"
	case codecs.Int64ID:
		typ = "Int64"
	case codecs.Float32ID:
		typ = "Float32"
	case codecs.Float64ID:
		typ = "Float64"
	case codecs.DateTimeID:
		typ = "DateTime"
	case codecs.LocalDTID:
		typ = "LocalDateTime"
	case codecs.LocalDateID:
		typ = "LocalDate"
	default:
		return nil, nil, fmt.Errorf(
			"generating range: unknown %v with id %v",
			fieldDesc.Type,
			fieldDesc.ID,
		)
	}

	types := []goType{
		&goScalar{Name: fmt.Sprintf("geltypes.%sRange%s", optional, typ)},
	}
	return types, nil, nil
}

func generateRangeV2(
	desc *descriptor.V2,
	required bool,
) ([]goType, []string, error) {
	optional := ""
	if !required {
		optional = "Optional"
	}

	var typ string
	fieldDesc := desc.Fields[0].Desc
	switch fieldDesc.ID {
	case codecs.Int32ID:
		typ = "Int32"
	case codecs.Int64ID:
		typ = "Int64"
	case codecs.Float32ID:
		typ = "Float32"
	case codecs.Float64ID:
		typ = "Float64"
	case codecs.DateTimeID:
		typ = "DateTime"
	case codecs.LocalDTID:
		typ = "LocalDateTime"
	case codecs.LocalDateID:
		typ = "LocalDate"
	default:
		return nil, nil, fmt.Errorf(
			"generating range: unknown %v with id %v",
			fieldDesc.Type,
			fieldDesc.ID,
		)
	}

	types := []goType{
		&goScalar{Name: fmt.Sprintf("geltypes.%sRange%s", optional, typ)},
	}
	return types, nil, nil
}

func generateSlice(
	desc descriptor.Descriptor,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
	isField bool,
) ([]goType, []string, error) {
	types, imports, err := generateType(
		desc.Fields[0].Desc,
		true,
		path,
		cmdCfg,
		isResult,
		isField,
	)
	if err != nil {
		return nil, nil, err
	}

	typ := []goType{&goSlice{typ: types[0]}}
	return append(typ, types...), imports, nil
}

func generateSliceV2(
	desc *descriptor.V2,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
	isField bool,
) ([]goType, []string, error) {
	types, imports, err := generateTypeV2(
		&desc.Fields[0].Desc,
		true,
		path,
		cmdCfg,
		isResult,
		isField,
	)
	if err != nil {
		return nil, nil, err
	}

	typ := []goType{&goSlice{typ: types[0]}}
	return append(typ, types...), imports, nil
}

func generateObject(
	desc descriptor.Descriptor,
	required bool,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
) ([]goType, []string, error) {
	var imports []string
	typ := goStruct{Name: nameFromPath(path), Required: required}
	types := []goType{&typ}
	if !required {
		// This is needed for geltypes.Optional in the struct definition.
		imports = append(imports, "github.com/geldata/gel-go/geltypes")
	}

	for _, field := range desc.Fields {
		t, i, err := generateType(
			field.Desc,
			field.Required,
			append(path, field.Name),
			cmdCfg,
			isResult,
			true,
		)
		if err != nil {
			return nil, nil, err
		}

		tag := fmt.Sprintf(`gel:"%s"`, field.Name)
		name := field.Name
		if cmdCfg.mixedCaps {
			name = snakeToUpperMixedCase(name)
		}

		typ.Fields = append(typ.Fields, goStructField{
			EQLName: field.Name,
			GoName:  name,
			Type:    t[0].Reference(),
			Tag:     tag,
		})
		types = append(types, t...)
		imports = append(imports, i...)
	}

	return types, imports, nil
}

func generateObjectV2(
	desc *descriptor.V2,
	required bool,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
) ([]goType, []string, error) {
	var imports []string
	typ := goStruct{Name: nameFromPath(path), Required: required}
	types := []goType{&typ}
	if !required {
		// This is needed for geltypes.Optional in the struct definition.
		imports = append(imports, "github.com/geldata/gel-go/geltypes")
	}

	for _, field := range desc.Fields {
		t, i, err := generateTypeV2(
			&field.Desc,
			field.Required,
			append(path, field.Name),
			cmdCfg,
			isResult,
			true,
		)
		if err != nil {
			return nil, nil, err
		}

		tag := fmt.Sprintf(`gel:"%s"`, field.Name)
		name := field.Name
		if cmdCfg.mixedCaps {
			name = snakeToUpperMixedCase(name)
		}

		typ.Fields = append(typ.Fields, goStructField{
			EQLName: field.Name,
			GoName:  name,
			Type:    t[0].Reference(),
			Tag:     tag,
		})
		types = append(types, t...)
		imports = append(imports, i...)
	}

	return types, imports, nil
}

func generateTuple(
	desc descriptor.Descriptor,
	required bool,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
) ([]goType, []string, error) {
	var imports []string
	typ := &goStruct{Name: nameFromPath(path), Required: required}
	types := []goType{typ}

	for _, field := range desc.Fields {
		t, i, err := generateType(
			field.Desc,
			field.Required,
			append(path, field.Name),
			cmdCfg,
			isResult,
			true,
		)
		if err != nil {
			return nil, nil, err
		}

		name := field.Name
		if name != "" && name[0] >= '0' && name[0] <= '9' {
			name = fmt.Sprintf("Element%s", name)
		} else if cmdCfg.mixedCaps {
			name = snakeToUpperMixedCase(name)
		}

		typ.Fields = append(typ.Fields, goStructField{
			EQLName: field.Name,
			GoName:  name,
			Type:    t[0].Reference(),
			Tag:     fmt.Sprintf(`gel:"%s"`, field.Name),
		})
		types = append(types, t...)
		imports = append(imports, i...)
	}

	return types, imports, nil
}

func generateTupleV2(
	desc *descriptor.V2,
	required bool,
	path []string,
	cmdCfg *cmdConfig,
	isResult bool,
) ([]goType, []string, error) {
	var imports []string
	typ := &goStruct{Name: nameFromPath(path), Required: required}
	types := []goType{typ}

	for _, field := range desc.Fields {
		t, i, err := generateTypeV2(
			&field.Desc,
			field.Required,
			append(path, field.Name),
			cmdCfg,
			isResult,
			true,
		)
		if err != nil {
			return nil, nil, err
		}

		name := field.Name
		if name != "" && name[0] >= '0' && name[0] <= '9' {
			name = fmt.Sprintf("Element%s", name)
		} else if cmdCfg.mixedCaps {
			name = snakeToUpperMixedCase(name)
		}

		typ.Fields = append(typ.Fields, goStructField{
			EQLName: field.Name,
			GoName:  name,
			Type:    t[0].Reference(),
			Tag:     fmt.Sprintf(`gel:"%s"`, field.Name),
		})
		types = append(types, t...)
		imports = append(imports, i...)
	}

	return types, imports, nil
}

func generateBaseScalar(
	desc descriptor.Descriptor,
	required bool,
	cmdCfg *cmdConfig,
	isResult bool,
	isField bool,
) ([]goType, []string, error) {
	if desc.Type == descriptor.Scalar {
		desc = codecs.GetScalarDescriptor(desc)
	}

	var name string
	var imports []string
	if desc.Type == descriptor.Enum {
		if required {
			name = "string"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalStr"
		}

		return []goType{&goScalar{Name: name}}, imports, nil
	}

	switch desc.ID {
	case codecs.UUIDID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.UUID"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalUUID"
		}
	case codecs.StrID:
		if required {
			name = "string"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalStr"
		}
	case codecs.JSONID:
		if required {
			if cmdCfg.rawmessage && isResult && isField {
				imports = append(imports, "encoding/json")
				name = "json.RawMessage"
			} else {
				name = "[]byte"
			}
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBytes"
		}
	case codecs.BytesID:
		if required {
			name = "[]byte"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBytes"
		}
	case codecs.Int16ID:
		if required {
			name = "int16"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalInt16"
		}
	case codecs.Int32ID:
		if required {
			name = "int32"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalInt32"
		}
	case codecs.Int64ID:
		if required {
			name = "int64"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalInt64"
		}
	case codecs.Float32ID:
		if required {
			name = "float32"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalFloat32"
		}
	case codecs.Float64ID:
		if required {
			name = "float64"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalFloat64"
		}
	case codecs.BoolID:
		if required {
			name = "bool"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBool"
		}
	case codecs.DateTimeID:
		if required {
			imports = append(imports, "time")
			name = "time.Time"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalDateTime"
		}
	case codecs.LocalDTID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.LocalDateTime"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalLocalDateTime"
		}
	case codecs.LocalDateID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.LocalDate"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalLocalDate"
		}
	case codecs.LocalTimeID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.LocalTime"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalLocalTime"
		}
	case codecs.DurationID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.Duration"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalDuration"
		}
	case codecs.BigIntID:
		if required {
			imports = append(imports, "math/big")
			name = "*big.Int"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBigInt"
		}
	case codecs.RelativeDurationID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.RelativeDuration"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalRelativeDuration"
		}
	case codecs.DateDurationID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.DateDuration"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalDateDuration"
		}
	case codecs.MemoryID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.Memory"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalMemory"
		}
	}

	return []goType{&goScalar{Name: name}}, imports, nil
}

func generateBaseScalarV2(
	desc *descriptor.V2,
	required bool,
	cmdCfg *cmdConfig,
	isResult bool,
	isField bool,
) ([]goType, []string, error) {
	if desc.Type == descriptor.Scalar {
		desc = codecs.GetScalarDescriptorV2(desc)
	}

	var name string
	var imports []string
	if desc.Type == descriptor.Enum {
		if required {
			name = "string"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalStr"
		}

		return []goType{&goScalar{Name: name}}, imports, nil
	}

	switch desc.ID {
	case codecs.UUIDID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.UUID"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalUUID"
		}
	case codecs.StrID:
		if required {
			name = "string"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalStr"
		}
	case codecs.JSONID:
		if required {
			if cmdCfg.rawmessage && isResult && isField {
				imports = append(imports, "encoding/json")
				name = "json.RawMessage"
			} else {
				name = "[]byte"
			}
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBytes"
		}
	case codecs.BytesID:
		if required {
			name = "[]byte"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBytes"
		}
	case codecs.Int16ID:
		if required {
			name = "int16"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalInt16"
		}
	case codecs.Int32ID:
		if required {
			name = "int32"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalInt32"
		}
	case codecs.Int64ID:
		if required {
			name = "int64"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalInt64"
		}
	case codecs.Float32ID:
		if required {
			name = "float32"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalFloat32"
		}
	case codecs.Float64ID:
		if required {
			name = "float64"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalFloat64"
		}
	case codecs.BoolID:
		if required {
			name = "bool"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBool"
		}
	case codecs.DateTimeID:
		if required {
			imports = append(imports, "time")
			name = "time.Time"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalDateTime"
		}
	case codecs.LocalDTID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.LocalDateTime"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalLocalDateTime"
		}
	case codecs.LocalDateID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.LocalDate"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalLocalDate"
		}
	case codecs.LocalTimeID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.LocalTime"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalLocalTime"
		}
	case codecs.DurationID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.Duration"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalDuration"
		}
	case codecs.BigIntID:
		if required {
			imports = append(imports, "math/big")
			name = "*big.Int"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalBigInt"
		}
	case codecs.RelativeDurationID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.RelativeDuration"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalRelativeDuration"
		}
	case codecs.DateDurationID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.DateDuration"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalDateDuration"
		}
	case codecs.MemoryID:
		if required {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.Memory"
		} else {
			imports = append(imports, "github.com/geldata/gel-go/geltypes")
			name = "geltypes.OptionalMemory"
		}
	}

	return []goType{&goScalar{Name: name}}, imports, nil
}

func nameFromPath(path []string) string {
	if len(path) == 0 {
		return ""
	}

	if len(path) == 1 {
		return path[0]
	}

	return path[0] + strings.Join(path[1:], "Item") + "Item"
}
