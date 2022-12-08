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
	"fmt"
	"strings"

	"github.com/edgedb/edgedb-go/internal/codecs"
	"github.com/edgedb/edgedb-go/internal/descriptor"
)

func generateType(
	desc descriptor.Descriptor,
	required bool,
	path []string,
) ([]goType, []string, error) {
	var (
		err     error
		types   []goType
		imports []string
	)

	switch desc.Type {
	case descriptor.Set, descriptor.Array:
		types, imports, err = generateSlice(desc, path)
	case descriptor.Object, descriptor.NamedTuple:
		types, imports, err = generateObject(desc, required, path)
	case descriptor.Tuple:
		types, imports, err = generateTuple(desc, required, path)
	case descriptor.BaseScalar, descriptor.Scalar, descriptor.Enum:
		types, imports, err = generateBaseScalar(desc, required)
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
		&goScalar{Name: fmt.Sprintf("edgedb.%sRange%s", optional, typ)},
	}
	return types, nil, nil
}

func generateSlice(
	desc descriptor.Descriptor,
	path []string,
) ([]goType, []string, error) {
	types, imports, err := generateType(
		desc.Fields[0].Desc,
		desc.Fields[0].Required,
		path,
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
) ([]goType, []string, error) {
	var imports []string
	typ := goStruct{Name: nameFromPath(path)}
	types := []goType{&typ}

	for _, field := range desc.Fields {
		t, i, err := generateType(
			field.Desc,
			field.Required,
			append(path, field.Name),
		)
		if err != nil {
			return nil, nil, err
		}

		typ.Fields = append(typ.Fields, goStructField{
			Name: field.Name,
			Type: t[0].Reference(),
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
) ([]goType, []string, error) {
	var imports []string
	typ := &goStruct{Name: nameFromPath(path)}
	types := []goType{typ}

	for _, field := range desc.Fields {
		t, i, err := generateType(
			field.Desc,
			field.Required,
			append(path, field.Name),
		)
		if err != nil {
			return nil, nil, err
		}

		typ.Fields = append(typ.Fields, goStructField{
			Name: field.Name,
			Type: t[0].Reference(),
		})
		types = append(types, t...)
		imports = append(imports, i...)
	}

	return types, imports, nil
}

func generateBaseScalar(
	desc descriptor.Descriptor,
	required bool,
) ([]goType, []string, error) {
	if desc.Type == descriptor.Scalar {
		desc = codecs.GetScalarDescriptor(desc)
	}

	var name string
	if desc.Type == descriptor.Enum {
		if required {
			name = "string"
		} else {
			name = "edgedb.OptionalStr"
		}

		return []goType{&goScalar{Name: name}}, nil, nil
	}

	var imports []string
	switch desc.ID {
	case codecs.UUIDID:
		if required {
			name = "edgedb.UUID"
		} else {
			name = "edgedb.OptionalUUID"
		}
	case codecs.StrID:
		if required {
			name = "string"
		} else {
			name = "edgedb.OptionalStr"
		}
	case codecs.BytesID, codecs.JSONID:
		if required {
			name = "[]byte"
		} else {
			name = "edgedb.OptionalBytes"
		}
	case codecs.Int16ID:
		if required {
			name = "int16"
		} else {
			name = "edgedb.OptionalInt16"
		}
	case codecs.Int32ID:
		if required {
			name = "int32"
		} else {
			name = "edgedb.OptionalInt32"
		}
	case codecs.Int64ID:
		if required {
			name = "int64"
		} else {
			name = "edgedb.OptionalInt64"
		}
	case codecs.Float32ID:
		if required {
			name = "float32"
		} else {
			name = "edgedb.OptionalFloat32"
		}
	case codecs.Float64ID:
		if required {
			name = "float64"
		} else {
			name = "edgedb.OptionalFloat64"
		}
	case codecs.BoolID:
		if required {
			name = "bool"
		} else {
			name = "edgedb.OptionalBool"
		}
	case codecs.DateTimeID:
		if required {
			imports = append(imports, "time")
			name = "time.Time"
		} else {
			name = "edgedb.OptionalDateTime"
		}
	case codecs.LocalDTID:
		if required {
			name = "edgedb.LocalDateTime"
		} else {
			name = "edgedb.OptionalLocalDateTime"
		}
	case codecs.LocalDateID:
		if required {
			name = "edgedb.LocalDate"
		} else {
			name = "edgedb.OptionalLocalDate"
		}
	case codecs.LocalTimeID:
		if required {
			name = "edgedb.LocalTime"
		} else {
			name = "edgedb.OptionalLocalTime"
		}
	case codecs.DurationID:
		if required {
			name = "edgedb.Duration"
		} else {
			name = "edgedb.OptionalDuration"
		}
	case codecs.BigIntID:
		if required {
			imports = append(imports, "math/big")
			name = "*big.Int"
		} else {
			name = "edgedb.OptionalBigInt"
		}
	case codecs.RelativeDurationID:
		if required {
			name = "edgedb.RelativeDuration"
		} else {
			name = "edgedb.OptionalRelativeDuration"
		}
	case codecs.DateDurationID:
		if required {
			name = "edgedb.DateDuration"
		} else {
			name = "edgedb.OptionalDateDuration"
		}
	case codecs.MemoryID:
		if required {
			name = "edgedb.Memory"
		} else {
			name = "edgedb.OptionalMemory"
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
