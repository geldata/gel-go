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

// Package marshal documents marshaling interfaces.
//
// User defined marshaler/unmarshalers can be defined for any scalar Gel
// type except arrays. They must implement the interface for their type.
// For example a custom int64 unmarshaler should implement Int64Unmarshaler.
//
// # Optional Fields
//
// When shape fields in a query result are optional (not required) the client
// requires the out value's optional fields to implement OptionalUnmarshaler.
// For scalar types, this means that the field value will need to implement a
// custom marshaler interface i.e. Int64Unmarshaler AND OptionalUnmarshaler.
// For shapes, only OptionalUnmarshaler needs to be implemented.
package marshal

// OptionalUnmarshaler is used for optional (not required) shape field values.
type OptionalUnmarshaler interface {
	// SetMissing is call with true when the value is missing and false when
	// the value is present.
	SetMissing(bool)
}

// OptionalScalarUnmarshaler is implemented by optional scalar types.
type OptionalScalarUnmarshaler interface {
	Unset()
}

// OptionalMarshaler is used for optional (not required) shape field values.
type OptionalMarshaler interface {
	// Missing returns true when the value is missing.
	Missing() bool
}

// StrMarshaler is the interface implemented by an object
// that can marshal itself into the [std::str] wire format.
//
// [std::str]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-str
type StrMarshaler interface {
	// MarshalEdgeDBStr encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBStr() ([]byte, error)
}

// StrUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::str] wire format representation of itself.
//
// [std::str]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-str
type StrUnmarshaler interface {
	// UnmarshalEdgeDBStr must be able to decode the str wire format.
	// UnmarshalEdgeDBStr must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBStr(data []byte) error
}

// BoolMarshaler is the interface implemented by an object
// that can marshal itself into the [std::bool] wire format.
//
// [std::bool]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-bool
type BoolMarshaler interface {
	// MarshalEdgeDBBool encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBBool() ([]byte, error)
}

// BoolUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::bool] wire format representation of itself.
//
// [std::bool]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-bool
type BoolUnmarshaler interface {
	// UnmarshalEdgeDBBool must be able to decode the bool wire format.
	// UnmarshalEdgeDBBool must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBBool(data []byte) error
}

// JSONMarshaler is the interface implemented by an object
// that can marshal itself into the [std::json] wire format.
//
// [std::json]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-json
type JSONMarshaler interface {
	// MarshalEdgeDBJSON encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBJSON() ([]byte, error)
}

// JSONUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::json] wire format representation of itself.
//
// [std::json]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-json
type JSONUnmarshaler interface {
	// UnmarshalEdgeDBJSON must be able to decode the json wire format.
	// UnmarshalEdgeDBJSON must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBJSON(data []byte) error
}

// UUIDMarshaler is the interface implemented by an object
// that can marshal itself into the [std::uuid] wire format.
//
// [std::uuid]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-uuid
type UUIDMarshaler interface {
	// MarshalEdgeDBUUID encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBUUID() ([]byte, error)
}

// UUIDUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::uuid] wire format representation of itself.
//
// [std::uuid]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-uuid
type UUIDUnmarshaler interface {
	// UnmarshalEdgeDBUUID must be able to decode the uuid wire format.
	// UnmarshalEdgeDBUUID must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBUUID(data []byte) error
}

// BytesMarshaler is the interface implemented by an object
// that can marshal itself into the [std::bytes] wire format.
//
// [std::bytes]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-bytes
type BytesMarshaler interface {
	// MarshalEdgeDBBytes encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBBytes() ([]byte, error)
}

// BytesUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::bytes] wire format representation of itself.
//
// [std::bytes]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-bytes
type BytesUnmarshaler interface {
	// UnmarshalEdgeDBBytes must be able to decode the bytes wire format.
	// UnmarshalEdgeDBBytes must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBBytes(data []byte) error
}

// BigIntMarshaler is the interface implemented by an object
// that can marshal itself into the [std::bigint] wire format.
//
// [std::bigint]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-bigint
type BigIntMarshaler interface {
	// MarshalEdgeDBBigInt encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBBigInt() ([]byte, error)
}

// BigIntUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::bigint] wire format representation of itself.
//
// [std::bigint]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-bigint
type BigIntUnmarshaler interface {
	// UnmarshalEdgeDBBigInt must be able to decode the bigint wire format.
	// UnmarshalEdgeDBBigInt must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBBigInt(data []byte) error
}

// DecimalMarshaler is the interface implemented by an object
// that can marshal itself into the [std::decimal] wire format.
//
// [std::decimal]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-decimal
type DecimalMarshaler interface {
	// MarshalEdgeDBDecimal encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBDecimal() ([]byte, error)
}

// DecimalUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::decimal] wire format representation of itself.
//
// [std::decimal]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-decimal
type DecimalUnmarshaler interface {
	// UnmarshalEdgeDBDecimal must be able to decode the decimal wire format.
	// UnmarshalEdgeDBDecimal must copy the data if it wishes to retain the
	// data after returning.
	UnmarshalEdgeDBDecimal(data []byte) error
}

// DateTimeMarshaler is the interface implemented by an object
// that can marshal itself into the [std::datetime] wire format.
//
// [std::datetime]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-datetime
type DateTimeMarshaler interface {
	// MarshalEdgeDBDateTime encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBDateTime() ([]byte, error)
}

// DateTimeUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::datetime] wire format representation of itself.
//
// [std::datetime]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-datetime
type DateTimeUnmarshaler interface {
	// UnmarshalEdgeDBDateTime must be able to decode the datetime wire format.
	// UnmarshalEdgeDBDateTime must copy the data if it wishes to retain the
	// data after returning.
	UnmarshalEdgeDBDateTime(data []byte) error
}

// LocalDateTimeMarshaler is the interface implemented by an object
// that can marshal itself into the [cal::local_datetime] wire format.
//
// [cal::local_datetime]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-local-datetime
type LocalDateTimeMarshaler interface {
	// MarshalEdgeDBLocalDateTime encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBLocalDateTime() ([]byte, error)
}

// LocalDateTimeUnmarshaler is the interface implemented by an object
// that can unmarshal the [cal::local_datetime] wire format representation of
// itself.
//
// [cal::local_datetime]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-local-datetime
type LocalDateTimeUnmarshaler interface {
	// UnmarshalEdgeDBLocalDateTime must be able to decode the local_datetime
	// wire format. UnmarshalEdgeDBLocalDateTime must copy the data if it
	// wishes to retain the data after returning.
	UnmarshalEdgeDBLocalDateTime(data []byte) error
}

// LocalDateMarshaler is the interface implemented by an object
// that can marshal itself into the [cal::locl_date] wire format.
//
// [cal::locl_date]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-local-date
type LocalDateMarshaler interface {
	// MarshalEdgeDBLocalDate encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBLocalDate() ([]byte, error)
}

// LocalDateUnmarshaler is the interface implemented by an object
// that can unmarshal the [cal::locl_date] wire format representation of
// itself.
//
// [cal::locl_date]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-local-date
type LocalDateUnmarshaler interface {
	// UnmarshalEdgeDBLocalDate must be able to decode the local_date wire
	// format.  UnmarshalEdgeDBLocalDate must copy the data if it wishes to
	// retain the data after returning.
	UnmarshalEdgeDBLocalDate(data []byte) error
}

// LocalTimeMarshaler is the interface implemented by an object
// that can marshal itself into the [cal::local_time] wire format.
//
// [cal::local_time]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-local-time
type LocalTimeMarshaler interface {
	// MarshalEdgeDBLocalTime encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBLocalTime() ([]byte, error)
}

// LocalTimeUnmarshaler is the interface implemented by an object
// that can unmarshal the [cal::local_time] wire format representation of
// itself.
//
// [cal::local_time]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-local-time
type LocalTimeUnmarshaler interface {
	// UnmarshalEdgeDBLocalTime must be able to decode the local_time wire
	// format.  UnmarshalEdgeDBLocalTime must copy the data if it wishes to
	// retain the data after returning.
	UnmarshalEdgeDBLocalTime(data []byte) error
}

// DurationMarshaler is the interface implemented by an object
// that can marshal itself into the [std::duration] wire format.
//
// [std::duration]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-duration
type DurationMarshaler interface {
	// MarshalEdgeDBDuration encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBDuration() ([]byte, error)
}

// DurationUnmarshaler is the interface implemented by an object
// that can unmarshal the [std::duration] wire format representation of itself.
//
// [std::duration]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-duration
type DurationUnmarshaler interface {
	// UnmarshalEdgeDBDuration must be able to decode the duration wire format.
	// UnmarshalEdgeDBDuration must copy the data if it wishes to retain the
	// data after returning.
	UnmarshalEdgeDBDuration(data []byte) error
}

// RelativeDurationMarshaler is the interface implemented by an object that can
// marshal itself into the [cal::relative_duration] wire format.
//
// [cal::relative_duration]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-relative-duration
type RelativeDurationMarshaler interface {
	// MarshalEdgeDBRelativeDuration encodes the receiver into a binary form
	// and returns the result.
	MarshalEdgeDBRelativeDuration() ([]byte, error)
}

// RelativeDurationUnmarshaler is the interface implemented by an object that
// can unmarshal the [cal::relative_duration] wire format representation of
// itself.
//
// [cal::relative_duration]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-relative-duration
type RelativeDurationUnmarshaler interface {
	// UnmarshalEdgeDBRelativeDuration must be able to decode the
	// cal::relative_duration wire format.  UnmarshalEdgeDBRelativeDuration
	// must copy data if it wishes to retain the bytes after returning.
	UnmarshalEdgeDBRelativeDuration(data []byte) error
}

// DateDurationMarshaler is the interface implemented by an object that can
// marshal itself into the [cal::date_duration] wire format.
//
// [cal::date_duration]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-date-duration
type DateDurationMarshaler interface {
	// MarshalEdgeDBDateDuration encodes the receiver into a binary form and
	// returns the result.
	MarshalEdgeDBDateDuration() ([]byte, error)
}

// DateDurationUnmarshaler is the interface implemented by an object that
// can unmarshal the [cal::date_duration] wire format representation of
// itself.
//
// [cal::date_duration]: https://docs.geldata.com/reference/reference/protocol/dataformats#cal-date-duration
type DateDurationUnmarshaler interface {
	// UnmarshalEdgeDBDateDuration must be able to decode the
	// cal::relative_duration wire format.  UnmarshalEdgeDBDateDuration must
	// copy the data if it wishes to retain the data after returning.
	UnmarshalEdgeDBDateDuration(data []byte) error
}

// Int16Marshaler is the interface implemented by an object
// that can marshal itself into the [std::int16] wire format.
//
// [std::int16]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-int16
type Int16Marshaler interface {
	// MarshalEdgeDBInt16 encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBInt16() ([]byte, error)
}

// Int16Unmarshaler is the interface implemented by an object
// that can unmarshal the [std::int16] wire format representation of itself.
//
// [std::int16]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-int16
type Int16Unmarshaler interface {
	// UnmarshalEdgeDBInt16 must be able to decode the int16 wire format.
	// UnmarshalEdgeDBInt16 must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBInt16(data []byte) error
}

// Int32Marshaler is the interface implemented by an object
// that can marshal itself into the [std::int32] wire format.
//
// [std::int32]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-int32
type Int32Marshaler interface {
	// MarshalEdgeDBInt32 encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBInt32() ([]byte, error)
}

// Int32Unmarshaler is the interface implemented by an object
// that can unmarshal the [std::int32] wire format representation of itself.
//
// [std::int32]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-int32
type Int32Unmarshaler interface {
	// UnmarshalEdgeDBInt32 must be able to decode the int32 wire format.
	// UnmarshalEdgeDBInt32 must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBInt32(data []byte) error
}

// Int64Marshaler is the interface implemented by an object
// that can marshal itself into the [std::int64] wire format.
//
// [std::int64]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-int64
type Int64Marshaler interface {
	// MarshalEdgeDBInt64 encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBInt64() ([]byte, error)
}

// Int64Unmarshaler is the interface implemented by an object
// that can unmarshal the [std::int64] wire format representation of itself.
//
// [std::int64]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-int64
type Int64Unmarshaler interface {
	// UnmarshalEdgeDBInt64 must be able to decode the int64 wire format.
	// UnmarshalEdgeDBInt64 must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBInt64(data []byte) error
}

// Float32Marshaler is the interface implemented by an object
// that can marshal itself into the [std::float32] wire format.
//
// [std::float32]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-float32
type Float32Marshaler interface {
	// MarshalEdgeDBFloat32 encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBFloat32() ([]byte, error)
}

// Float32Unmarshaler is the interface implemented by an object
// that can unmarshal the [std::float32] wire format representation of itself.
//
// [std::float32]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-float32
type Float32Unmarshaler interface {
	// UnmarshalEdgeDBFloat32 must be able to decode the float32 wire format.
	// UnmarshalEdgeDBFloat32 must copy the data if it wishes to retain the
	// data after returning.
	UnmarshalEdgeDBFloat32(data []byte) error
}

// Float64Marshaler is the interface implemented by an object
// that can marshal itself into the [std::float64] wire format.
//
// [std::float64]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-float64
type Float64Marshaler interface {
	// MarshalEdgeDBFloat64 encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBFloat64() ([]byte, error)
}

// Float64Unmarshaler is the interface implemented by an object
// that can unmarshal the [std::float64] wire format representation of itself.
//
// [std::float64]: https://docs.geldata.com/reference/reference/protocol/dataformats#std-float64
type Float64Unmarshaler interface {
	// UnmarshalEdgeDBFloat64 must be able to decode the float64 wire format.
	// UnmarshalEdgeDBFloat64 must copy the data if it wishes to retain the
	// data after returning.
	UnmarshalEdgeDBFloat64(data []byte) error
}

// MemoryMarshaler is the interface implemented by an object
// that can marshal itself into the [cfg::memory] wire format.
//
// [cfg::memory]: https://docs.geldata.com/reference/reference/protocol/dataformats#cfg-memory
type MemoryMarshaler interface {
	// MarshalEdgeDBMemory encodes the receiver
	// into a binary form and returns the result.
	MarshalEdgeDBMemory() ([]byte, error)
}

// MemoryUnmarshaler is the interface implemented by an object
// that can unmarshal the [cfg::memory] wire format representation of itself.
//
// [cfg::memory]: https://docs.geldata.com/reference/reference/protocol/dataformats#cfg-memory
type MemoryUnmarshaler interface {
	// UnmarshalEdgeDBMemory must be able to decode the memory wire format.
	// UnmarshalEdgeDBMemory must copy the data if it wishes to retain the data
	// after returning.
	UnmarshalEdgeDBMemory(data []byte) error
}
