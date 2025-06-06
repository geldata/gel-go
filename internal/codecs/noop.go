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

package codecs

import (
	"unsafe"

	types "github.com/geldata/gel-go/geltypes"
	"github.com/geldata/gel-go/internal/buff"
	"github.com/geldata/gel-go/internal/descriptor"
)

var (
	// NoOpDecoder is a noOpDecoder
	NoOpDecoder = noOpDecoder{}

	// NoOpEncoder is a noOpEncoder
	NoOpEncoder = noOpEncoder{}
)

// noOpDecoder decodes [empty blocks] i.e. does nothing.
//
//	There is one special type with type id of zero:
//	00000000-0000-0000-0000-000000000000.
//	The describe result of this type contains zero blocks.
//	It’s used when a statement returns no meaningful results,
//	e.g. the CREATE DATABASE example statement.
//
// [empty blocks]: https://docs.geldata.com/reference/reference/protocol/typedesc#type-descriptors
type noOpDecoder struct{}

func (c noOpDecoder) DescriptorID() types.UUID { return descriptor.IDZero }

func (c noOpDecoder) Decode(_ *buff.Reader, _ unsafe.Pointer) error {
	return nil
}

type noOpEncoder struct{}

func (c noOpEncoder) DescriptorID() types.UUID { return descriptor.IDZero }

func (c noOpEncoder) Encode(
	w *buff.Writer,
	_ interface{},
	_ Path,
	_ bool,
) error {
	w.PushUint32(0)
	return nil
}
