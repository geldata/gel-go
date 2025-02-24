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

// This file is auto generated. Do not edit!
// run 'make errors' to regenerate

package gelerr

import (
	"fmt"
	"github.com/geldata/gel-go/gelerr"
)

func NewInternalServerError(msg string, err error) error {
	return &InternalServerError{msg, err}
}

type InternalServerError struct {
	msg string
	err error
}

func (e *InternalServerError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InternalServerError: " + msg
}

func (e *InternalServerError) Unwrap() error { return e.err }

func (e *InternalServerError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InternalServerError:
		return true
	default:
		return false
	}
}

func (e *InternalServerError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnsupportedFeatureError(msg string, err error) error {
	return &UnsupportedFeatureError{msg, err}
}

type UnsupportedFeatureError struct {
	msg string
	err error
}

func (e *UnsupportedFeatureError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnsupportedFeatureError: " + msg
}

func (e *UnsupportedFeatureError) Unwrap() error { return e.err }

func (e *UnsupportedFeatureError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnsupportedFeatureError:
		return true
	default:
		return false
	}
}

func (e *UnsupportedFeatureError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewProtocolError(msg string, err error) error {
	return &ProtocolError{msg, err}
}

type ProtocolError struct {
	msg string
	err error
}

func (e *ProtocolError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ProtocolError: " + msg
}

func (e *ProtocolError) Unwrap() error { return e.err }

func (e *ProtocolError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *ProtocolError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewBinaryProtocolError(msg string, err error) error {
	return &BinaryProtocolError{msg, err}
}

type BinaryProtocolError struct {
	msg string
	err error
}

func (e *BinaryProtocolError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.BinaryProtocolError: " + msg
}

func (e *BinaryProtocolError) Unwrap() error { return e.err }

func (e *BinaryProtocolError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.BinaryProtocolError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *BinaryProtocolError) isEdgeDBProtocolError() {}

func (e *BinaryProtocolError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnsupportedProtocolVersionError(msg string, err error) error {
	return &UnsupportedProtocolVersionError{msg, err}
}

type UnsupportedProtocolVersionError struct {
	msg string
	err error
}

func (e *UnsupportedProtocolVersionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnsupportedProtocolVersionError: " + msg
}

func (e *UnsupportedProtocolVersionError) Unwrap() error { return e.err }

func (e *UnsupportedProtocolVersionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnsupportedProtocolVersionError:
		return true
	case gelerr.BinaryProtocolError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *UnsupportedProtocolVersionError) isEdgeDBBinaryProtocolError() {}

func (e *UnsupportedProtocolVersionError) isEdgeDBProtocolError() {}

func (e *UnsupportedProtocolVersionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewTypeSpecNotFoundError(msg string, err error) error {
	return &TypeSpecNotFoundError{msg, err}
}

type TypeSpecNotFoundError struct {
	msg string
	err error
}

func (e *TypeSpecNotFoundError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.TypeSpecNotFoundError: " + msg
}

func (e *TypeSpecNotFoundError) Unwrap() error { return e.err }

func (e *TypeSpecNotFoundError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.TypeSpecNotFoundError:
		return true
	case gelerr.BinaryProtocolError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *TypeSpecNotFoundError) isEdgeDBBinaryProtocolError() {}

func (e *TypeSpecNotFoundError) isEdgeDBProtocolError() {}

func (e *TypeSpecNotFoundError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnexpectedMessageError(msg string, err error) error {
	return &UnexpectedMessageError{msg, err}
}

type UnexpectedMessageError struct {
	msg string
	err error
}

func (e *UnexpectedMessageError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnexpectedMessageError: " + msg
}

func (e *UnexpectedMessageError) Unwrap() error { return e.err }

func (e *UnexpectedMessageError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnexpectedMessageError:
		return true
	case gelerr.BinaryProtocolError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *UnexpectedMessageError) isEdgeDBBinaryProtocolError() {}

func (e *UnexpectedMessageError) isEdgeDBProtocolError() {}

func (e *UnexpectedMessageError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInputDataError(msg string, err error) error {
	return &InputDataError{msg, err}
}

type InputDataError struct {
	msg string
	err error
}

func (e *InputDataError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InputDataError: " + msg
}

func (e *InputDataError) Unwrap() error { return e.err }

func (e *InputDataError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InputDataError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *InputDataError) isEdgeDBProtocolError() {}

func (e *InputDataError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewParameterTypeMismatchError(msg string, err error) error {
	return &ParameterTypeMismatchError{msg, err}
}

type ParameterTypeMismatchError struct {
	msg string
	err error
}

func (e *ParameterTypeMismatchError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ParameterTypeMismatchError: " + msg
}

func (e *ParameterTypeMismatchError) Unwrap() error { return e.err }

func (e *ParameterTypeMismatchError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ParameterTypeMismatchError:
		return true
	case gelerr.InputDataError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *ParameterTypeMismatchError) isEdgeDBInputDataError() {}

func (e *ParameterTypeMismatchError) isEdgeDBProtocolError() {}

func (e *ParameterTypeMismatchError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewStateMismatchError(msg string, err error) error {
	return &StateMismatchError{msg, err}
}

type StateMismatchError struct {
	msg string
	err error
}

func (e *StateMismatchError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.StateMismatchError: " + msg
}

func (e *StateMismatchError) Unwrap() error { return e.err }

func (e *StateMismatchError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.StateMismatchError:
		return true
	case gelerr.InputDataError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *StateMismatchError) isEdgeDBInputDataError() {}

func (e *StateMismatchError) isEdgeDBProtocolError() {}

func (e *StateMismatchError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewResultCardinalityMismatchError(msg string, err error) error {
	return &ResultCardinalityMismatchError{msg, err}
}

type ResultCardinalityMismatchError struct {
	msg string
	err error
}

func (e *ResultCardinalityMismatchError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ResultCardinalityMismatchError: " + msg
}

func (e *ResultCardinalityMismatchError) Unwrap() error { return e.err }

func (e *ResultCardinalityMismatchError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ResultCardinalityMismatchError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *ResultCardinalityMismatchError) isEdgeDBProtocolError() {}

func (e *ResultCardinalityMismatchError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewCapabilityError(msg string, err error) error {
	return &CapabilityError{msg, err}
}

type CapabilityError struct {
	msg string
	err error
}

func (e *CapabilityError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.CapabilityError: " + msg
}

func (e *CapabilityError) Unwrap() error { return e.err }

func (e *CapabilityError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.CapabilityError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *CapabilityError) isEdgeDBProtocolError() {}

func (e *CapabilityError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnsupportedCapabilityError(msg string, err error) error {
	return &UnsupportedCapabilityError{msg, err}
}

type UnsupportedCapabilityError struct {
	msg string
	err error
}

func (e *UnsupportedCapabilityError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnsupportedCapabilityError: " + msg
}

func (e *UnsupportedCapabilityError) Unwrap() error { return e.err }

func (e *UnsupportedCapabilityError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnsupportedCapabilityError:
		return true
	case gelerr.CapabilityError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *UnsupportedCapabilityError) isEdgeDBCapabilityError() {}

func (e *UnsupportedCapabilityError) isEdgeDBProtocolError() {}

func (e *UnsupportedCapabilityError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDisabledCapabilityError(msg string, err error) error {
	return &DisabledCapabilityError{msg, err}
}

type DisabledCapabilityError struct {
	msg string
	err error
}

func (e *DisabledCapabilityError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DisabledCapabilityError: " + msg
}

func (e *DisabledCapabilityError) Unwrap() error { return e.err }

func (e *DisabledCapabilityError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DisabledCapabilityError:
		return true
	case gelerr.CapabilityError:
		return true
	case gelerr.ProtocolError:
		return true
	default:
		return false
	}
}

func (e *DisabledCapabilityError) isEdgeDBCapabilityError() {}

func (e *DisabledCapabilityError) isEdgeDBProtocolError() {}

func (e *DisabledCapabilityError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewQueryError(msg string, err error) error {
	return &QueryError{msg, err}
}

type QueryError struct {
	msg string
	err error
}

func (e *QueryError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.QueryError: " + msg
}

func (e *QueryError) Unwrap() error { return e.err }

func (e *QueryError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *QueryError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidSyntaxError(msg string, err error) error {
	return &InvalidSyntaxError{msg, err}
}

type InvalidSyntaxError struct {
	msg string
	err error
}

func (e *InvalidSyntaxError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidSyntaxError: " + msg
}

func (e *InvalidSyntaxError) Unwrap() error { return e.err }

func (e *InvalidSyntaxError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidSyntaxError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidSyntaxError) isEdgeDBQueryError() {}

func (e *InvalidSyntaxError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewEdgeQLSyntaxError(msg string, err error) error {
	return &EdgeQLSyntaxError{msg, err}
}

type EdgeQLSyntaxError struct {
	msg string
	err error
}

func (e *EdgeQLSyntaxError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.EdgeQLSyntaxError: " + msg
}

func (e *EdgeQLSyntaxError) Unwrap() error { return e.err }

func (e *EdgeQLSyntaxError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.EdgeQLSyntaxError:
		return true
	case gelerr.InvalidSyntaxError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *EdgeQLSyntaxError) isEdgeDBInvalidSyntaxError() {}

func (e *EdgeQLSyntaxError) isEdgeDBQueryError() {}

func (e *EdgeQLSyntaxError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewSchemaSyntaxError(msg string, err error) error {
	return &SchemaSyntaxError{msg, err}
}

type SchemaSyntaxError struct {
	msg string
	err error
}

func (e *SchemaSyntaxError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.SchemaSyntaxError: " + msg
}

func (e *SchemaSyntaxError) Unwrap() error { return e.err }

func (e *SchemaSyntaxError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.SchemaSyntaxError:
		return true
	case gelerr.InvalidSyntaxError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *SchemaSyntaxError) isEdgeDBInvalidSyntaxError() {}

func (e *SchemaSyntaxError) isEdgeDBQueryError() {}

func (e *SchemaSyntaxError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewGraphQLSyntaxError(msg string, err error) error {
	return &GraphQLSyntaxError{msg, err}
}

type GraphQLSyntaxError struct {
	msg string
	err error
}

func (e *GraphQLSyntaxError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.GraphQLSyntaxError: " + msg
}

func (e *GraphQLSyntaxError) Unwrap() error { return e.err }

func (e *GraphQLSyntaxError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.GraphQLSyntaxError:
		return true
	case gelerr.InvalidSyntaxError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *GraphQLSyntaxError) isEdgeDBInvalidSyntaxError() {}

func (e *GraphQLSyntaxError) isEdgeDBQueryError() {}

func (e *GraphQLSyntaxError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidTypeError(msg string, err error) error {
	return &InvalidTypeError{msg, err}
}

type InvalidTypeError struct {
	msg string
	err error
}

func (e *InvalidTypeError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidTypeError: " + msg
}

func (e *InvalidTypeError) Unwrap() error { return e.err }

func (e *InvalidTypeError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidTypeError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidTypeError) isEdgeDBQueryError() {}

func (e *InvalidTypeError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidTargetError(msg string, err error) error {
	return &InvalidTargetError{msg, err}
}

type InvalidTargetError struct {
	msg string
	err error
}

func (e *InvalidTargetError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidTargetError: " + msg
}

func (e *InvalidTargetError) Unwrap() error { return e.err }

func (e *InvalidTargetError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidTargetError:
		return true
	case gelerr.InvalidTypeError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidTargetError) isEdgeDBInvalidTypeError() {}

func (e *InvalidTargetError) isEdgeDBQueryError() {}

func (e *InvalidTargetError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidLinkTargetError(msg string, err error) error {
	return &InvalidLinkTargetError{msg, err}
}

type InvalidLinkTargetError struct {
	msg string
	err error
}

func (e *InvalidLinkTargetError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidLinkTargetError: " + msg
}

func (e *InvalidLinkTargetError) Unwrap() error { return e.err }

func (e *InvalidLinkTargetError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidLinkTargetError:
		return true
	case gelerr.InvalidTargetError:
		return true
	case gelerr.InvalidTypeError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidLinkTargetError) isEdgeDBInvalidTargetError() {}

func (e *InvalidLinkTargetError) isEdgeDBInvalidTypeError() {}

func (e *InvalidLinkTargetError) isEdgeDBQueryError() {}

func (e *InvalidLinkTargetError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidPropertyTargetError(msg string, err error) error {
	return &InvalidPropertyTargetError{msg, err}
}

type InvalidPropertyTargetError struct {
	msg string
	err error
}

func (e *InvalidPropertyTargetError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidPropertyTargetError: " + msg
}

func (e *InvalidPropertyTargetError) Unwrap() error { return e.err }

func (e *InvalidPropertyTargetError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidPropertyTargetError:
		return true
	case gelerr.InvalidTargetError:
		return true
	case gelerr.InvalidTypeError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidPropertyTargetError) isEdgeDBInvalidTargetError() {}

func (e *InvalidPropertyTargetError) isEdgeDBInvalidTypeError() {}

func (e *InvalidPropertyTargetError) isEdgeDBQueryError() {}

func (e *InvalidPropertyTargetError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidReferenceError(msg string, err error) error {
	return &InvalidReferenceError{msg, err}
}

type InvalidReferenceError struct {
	msg string
	err error
}

func (e *InvalidReferenceError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidReferenceError: " + msg
}

func (e *InvalidReferenceError) Unwrap() error { return e.err }

func (e *InvalidReferenceError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidReferenceError) isEdgeDBQueryError() {}

func (e *InvalidReferenceError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnknownModuleError(msg string, err error) error {
	return &UnknownModuleError{msg, err}
}

type UnknownModuleError struct {
	msg string
	err error
}

func (e *UnknownModuleError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownModuleError: " + msg
}

func (e *UnknownModuleError) Unwrap() error { return e.err }

func (e *UnknownModuleError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownModuleError:
		return true
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *UnknownModuleError) isEdgeDBInvalidReferenceError() {}

func (e *UnknownModuleError) isEdgeDBQueryError() {}

func (e *UnknownModuleError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnknownLinkError(msg string, err error) error {
	return &UnknownLinkError{msg, err}
}

type UnknownLinkError struct {
	msg string
	err error
}

func (e *UnknownLinkError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownLinkError: " + msg
}

func (e *UnknownLinkError) Unwrap() error { return e.err }

func (e *UnknownLinkError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownLinkError:
		return true
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *UnknownLinkError) isEdgeDBInvalidReferenceError() {}

func (e *UnknownLinkError) isEdgeDBQueryError() {}

func (e *UnknownLinkError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnknownPropertyError(msg string, err error) error {
	return &UnknownPropertyError{msg, err}
}

type UnknownPropertyError struct {
	msg string
	err error
}

func (e *UnknownPropertyError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownPropertyError: " + msg
}

func (e *UnknownPropertyError) Unwrap() error { return e.err }

func (e *UnknownPropertyError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownPropertyError:
		return true
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *UnknownPropertyError) isEdgeDBInvalidReferenceError() {}

func (e *UnknownPropertyError) isEdgeDBQueryError() {}

func (e *UnknownPropertyError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnknownUserError(msg string, err error) error {
	return &UnknownUserError{msg, err}
}

type UnknownUserError struct {
	msg string
	err error
}

func (e *UnknownUserError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownUserError: " + msg
}

func (e *UnknownUserError) Unwrap() error { return e.err }

func (e *UnknownUserError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownUserError:
		return true
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *UnknownUserError) isEdgeDBInvalidReferenceError() {}

func (e *UnknownUserError) isEdgeDBQueryError() {}

func (e *UnknownUserError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnknownDatabaseError(msg string, err error) error {
	return &UnknownDatabaseError{msg, err}
}

type UnknownDatabaseError struct {
	msg string
	err error
}

func (e *UnknownDatabaseError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownDatabaseError: " + msg
}

func (e *UnknownDatabaseError) Unwrap() error { return e.err }

func (e *UnknownDatabaseError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownDatabaseError:
		return true
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *UnknownDatabaseError) isEdgeDBInvalidReferenceError() {}

func (e *UnknownDatabaseError) isEdgeDBQueryError() {}

func (e *UnknownDatabaseError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnknownParameterError(msg string, err error) error {
	return &UnknownParameterError{msg, err}
}

type UnknownParameterError struct {
	msg string
	err error
}

func (e *UnknownParameterError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownParameterError: " + msg
}

func (e *UnknownParameterError) Unwrap() error { return e.err }

func (e *UnknownParameterError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownParameterError:
		return true
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *UnknownParameterError) isEdgeDBInvalidReferenceError() {}

func (e *UnknownParameterError) isEdgeDBQueryError() {}

func (e *UnknownParameterError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDeprecatedScopingError(msg string, err error) error {
	return &DeprecatedScopingError{msg, err}
}

type DeprecatedScopingError struct {
	msg string
	err error
}

func (e *DeprecatedScopingError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DeprecatedScopingError: " + msg
}

func (e *DeprecatedScopingError) Unwrap() error { return e.err }

func (e *DeprecatedScopingError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DeprecatedScopingError:
		return true
	case gelerr.InvalidReferenceError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DeprecatedScopingError) isEdgeDBInvalidReferenceError() {}

func (e *DeprecatedScopingError) isEdgeDBQueryError() {}

func (e *DeprecatedScopingError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewSchemaError(msg string, err error) error {
	return &SchemaError{msg, err}
}

type SchemaError struct {
	msg string
	err error
}

func (e *SchemaError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.SchemaError: " + msg
}

func (e *SchemaError) Unwrap() error { return e.err }

func (e *SchemaError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.SchemaError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *SchemaError) isEdgeDBQueryError() {}

func (e *SchemaError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewSchemaDefinitionError(msg string, err error) error {
	return &SchemaDefinitionError{msg, err}
}

type SchemaDefinitionError struct {
	msg string
	err error
}

func (e *SchemaDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.SchemaDefinitionError: " + msg
}

func (e *SchemaDefinitionError) Unwrap() error { return e.err }

func (e *SchemaDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *SchemaDefinitionError) isEdgeDBQueryError() {}

func (e *SchemaDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidDefinitionError(msg string, err error) error {
	return &InvalidDefinitionError{msg, err}
}

type InvalidDefinitionError struct {
	msg string
	err error
}

func (e *InvalidDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidDefinitionError: " + msg
}

func (e *InvalidDefinitionError) Unwrap() error { return e.err }

func (e *InvalidDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidModuleDefinitionError(msg string, err error) error {
	return &InvalidModuleDefinitionError{msg, err}
}

type InvalidModuleDefinitionError struct {
	msg string
	err error
}

func (e *InvalidModuleDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidModuleDefinitionError: " + msg
}

func (e *InvalidModuleDefinitionError) Unwrap() error { return e.err }

func (e *InvalidModuleDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidModuleDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidModuleDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidModuleDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidModuleDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidModuleDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidLinkDefinitionError(msg string, err error) error {
	return &InvalidLinkDefinitionError{msg, err}
}

type InvalidLinkDefinitionError struct {
	msg string
	err error
}

func (e *InvalidLinkDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidLinkDefinitionError: " + msg
}

func (e *InvalidLinkDefinitionError) Unwrap() error { return e.err }

func (e *InvalidLinkDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidLinkDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidLinkDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidLinkDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidLinkDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidLinkDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidPropertyDefinitionError(msg string, err error) error {
	return &InvalidPropertyDefinitionError{msg, err}
}

type InvalidPropertyDefinitionError struct {
	msg string
	err error
}

func (e *InvalidPropertyDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidPropertyDefinitionError: " + msg
}

func (e *InvalidPropertyDefinitionError) Unwrap() error { return e.err }

func (e *InvalidPropertyDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidPropertyDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidPropertyDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidPropertyDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidPropertyDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidPropertyDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidUserDefinitionError(msg string, err error) error {
	return &InvalidUserDefinitionError{msg, err}
}

type InvalidUserDefinitionError struct {
	msg string
	err error
}

func (e *InvalidUserDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidUserDefinitionError: " + msg
}

func (e *InvalidUserDefinitionError) Unwrap() error { return e.err }

func (e *InvalidUserDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidUserDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidUserDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidUserDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidUserDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidUserDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidDatabaseDefinitionError(msg string, err error) error {
	return &InvalidDatabaseDefinitionError{msg, err}
}

type InvalidDatabaseDefinitionError struct {
	msg string
	err error
}

func (e *InvalidDatabaseDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidDatabaseDefinitionError: " + msg
}

func (e *InvalidDatabaseDefinitionError) Unwrap() error { return e.err }

func (e *InvalidDatabaseDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidDatabaseDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidDatabaseDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidDatabaseDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidDatabaseDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidDatabaseDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidOperatorDefinitionError(msg string, err error) error {
	return &InvalidOperatorDefinitionError{msg, err}
}

type InvalidOperatorDefinitionError struct {
	msg string
	err error
}

func (e *InvalidOperatorDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidOperatorDefinitionError: " + msg
}

func (e *InvalidOperatorDefinitionError) Unwrap() error { return e.err }

func (e *InvalidOperatorDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidOperatorDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidOperatorDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidOperatorDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidOperatorDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidOperatorDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidAliasDefinitionError(msg string, err error) error {
	return &InvalidAliasDefinitionError{msg, err}
}

type InvalidAliasDefinitionError struct {
	msg string
	err error
}

func (e *InvalidAliasDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidAliasDefinitionError: " + msg
}

func (e *InvalidAliasDefinitionError) Unwrap() error { return e.err }

func (e *InvalidAliasDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidAliasDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidAliasDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidAliasDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidAliasDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidAliasDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidFunctionDefinitionError(msg string, err error) error {
	return &InvalidFunctionDefinitionError{msg, err}
}

type InvalidFunctionDefinitionError struct {
	msg string
	err error
}

func (e *InvalidFunctionDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidFunctionDefinitionError: " + msg
}

func (e *InvalidFunctionDefinitionError) Unwrap() error { return e.err }

func (e *InvalidFunctionDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidFunctionDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidFunctionDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidFunctionDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidFunctionDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidFunctionDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidConstraintDefinitionError(msg string, err error) error {
	return &InvalidConstraintDefinitionError{msg, err}
}

type InvalidConstraintDefinitionError struct {
	msg string
	err error
}

func (e *InvalidConstraintDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidConstraintDefinitionError: " + msg
}

func (e *InvalidConstraintDefinitionError) Unwrap() error { return e.err }

func (e *InvalidConstraintDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidConstraintDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidConstraintDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidConstraintDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidConstraintDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidConstraintDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidCastDefinitionError(msg string, err error) error {
	return &InvalidCastDefinitionError{msg, err}
}

type InvalidCastDefinitionError struct {
	msg string
	err error
}

func (e *InvalidCastDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidCastDefinitionError: " + msg
}

func (e *InvalidCastDefinitionError) Unwrap() error { return e.err }

func (e *InvalidCastDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidCastDefinitionError:
		return true
	case gelerr.InvalidDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *InvalidCastDefinitionError) isEdgeDBInvalidDefinitionError() {}

func (e *InvalidCastDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *InvalidCastDefinitionError) isEdgeDBQueryError() {}

func (e *InvalidCastDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateDefinitionError(msg string, err error) error {
	return &DuplicateDefinitionError{msg, err}
}

type DuplicateDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateDefinitionError: " + msg
}

func (e *DuplicateDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateModuleDefinitionError(msg string, err error) error {
	return &DuplicateModuleDefinitionError{msg, err}
}

type DuplicateModuleDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateModuleDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateModuleDefinitionError: " + msg
}

func (e *DuplicateModuleDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateModuleDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateModuleDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateModuleDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateModuleDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateModuleDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateModuleDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateLinkDefinitionError(msg string, err error) error {
	return &DuplicateLinkDefinitionError{msg, err}
}

type DuplicateLinkDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateLinkDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateLinkDefinitionError: " + msg
}

func (e *DuplicateLinkDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateLinkDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateLinkDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateLinkDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateLinkDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateLinkDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateLinkDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicatePropertyDefinitionError(msg string, err error) error {
	return &DuplicatePropertyDefinitionError{msg, err}
}

type DuplicatePropertyDefinitionError struct {
	msg string
	err error
}

func (e *DuplicatePropertyDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicatePropertyDefinitionError: " + msg
}

func (e *DuplicatePropertyDefinitionError) Unwrap() error { return e.err }

func (e *DuplicatePropertyDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicatePropertyDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicatePropertyDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicatePropertyDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicatePropertyDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicatePropertyDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateUserDefinitionError(msg string, err error) error {
	return &DuplicateUserDefinitionError{msg, err}
}

type DuplicateUserDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateUserDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateUserDefinitionError: " + msg
}

func (e *DuplicateUserDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateUserDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateUserDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateUserDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateUserDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateUserDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateUserDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateDatabaseDefinitionError(msg string, err error) error {
	return &DuplicateDatabaseDefinitionError{msg, err}
}

type DuplicateDatabaseDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateDatabaseDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateDatabaseDefinitionError: " + msg
}

func (e *DuplicateDatabaseDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateDatabaseDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateDatabaseDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateDatabaseDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateDatabaseDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateDatabaseDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateDatabaseDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateOperatorDefinitionError(msg string, err error) error {
	return &DuplicateOperatorDefinitionError{msg, err}
}

type DuplicateOperatorDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateOperatorDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateOperatorDefinitionError: " + msg
}

func (e *DuplicateOperatorDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateOperatorDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateOperatorDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateOperatorDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateOperatorDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateOperatorDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateOperatorDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateViewDefinitionError(msg string, err error) error {
	return &DuplicateViewDefinitionError{msg, err}
}

type DuplicateViewDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateViewDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateViewDefinitionError: " + msg
}

func (e *DuplicateViewDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateViewDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateViewDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateViewDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateViewDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateViewDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateViewDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateFunctionDefinitionError(msg string, err error) error {
	return &DuplicateFunctionDefinitionError{msg, err}
}

type DuplicateFunctionDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateFunctionDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateFunctionDefinitionError: " + msg
}

func (e *DuplicateFunctionDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateFunctionDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateFunctionDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateFunctionDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateFunctionDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateFunctionDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateFunctionDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateConstraintDefinitionError(msg string, err error) error {
	return &DuplicateConstraintDefinitionError{msg, err}
}

type DuplicateConstraintDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateConstraintDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateConstraintDefinitionError: " + msg
}

func (e *DuplicateConstraintDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateConstraintDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateConstraintDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateConstraintDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateConstraintDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateConstraintDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateConstraintDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateCastDefinitionError(msg string, err error) error {
	return &DuplicateCastDefinitionError{msg, err}
}

type DuplicateCastDefinitionError struct {
	msg string
	err error
}

func (e *DuplicateCastDefinitionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateCastDefinitionError: " + msg
}

func (e *DuplicateCastDefinitionError) Unwrap() error { return e.err }

func (e *DuplicateCastDefinitionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateCastDefinitionError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateCastDefinitionError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateCastDefinitionError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateCastDefinitionError) isEdgeDBQueryError() {}

func (e *DuplicateCastDefinitionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDuplicateMigrationError(msg string, err error) error {
	return &DuplicateMigrationError{msg, err}
}

type DuplicateMigrationError struct {
	msg string
	err error
}

func (e *DuplicateMigrationError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DuplicateMigrationError: " + msg
}

func (e *DuplicateMigrationError) Unwrap() error { return e.err }

func (e *DuplicateMigrationError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DuplicateMigrationError:
		return true
	case gelerr.DuplicateDefinitionError:
		return true
	case gelerr.SchemaDefinitionError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *DuplicateMigrationError) isEdgeDBDuplicateDefinitionError() {}

func (e *DuplicateMigrationError) isEdgeDBSchemaDefinitionError() {}

func (e *DuplicateMigrationError) isEdgeDBQueryError() {}

func (e *DuplicateMigrationError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewSessionTimeoutError(msg string, err error) error {
	return &SessionTimeoutError{msg, err}
}

type SessionTimeoutError struct {
	msg string
	err error
}

func (e *SessionTimeoutError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.SessionTimeoutError: " + msg
}

func (e *SessionTimeoutError) Unwrap() error { return e.err }

func (e *SessionTimeoutError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.SessionTimeoutError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *SessionTimeoutError) isEdgeDBQueryError() {}

func (e *SessionTimeoutError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewIdleSessionTimeoutError(msg string, err error) error {
	return &IdleSessionTimeoutError{msg, err}
}

type IdleSessionTimeoutError struct {
	msg string
	err error
}

func (e *IdleSessionTimeoutError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.IdleSessionTimeoutError: " + msg
}

func (e *IdleSessionTimeoutError) Unwrap() error { return e.err }

func (e *IdleSessionTimeoutError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.IdleSessionTimeoutError:
		return true
	case gelerr.SessionTimeoutError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *IdleSessionTimeoutError) isEdgeDBSessionTimeoutError() {}

func (e *IdleSessionTimeoutError) isEdgeDBQueryError() {}

func (e *IdleSessionTimeoutError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewQueryTimeoutError(msg string, err error) error {
	return &QueryTimeoutError{msg, err}
}

type QueryTimeoutError struct {
	msg string
	err error
}

func (e *QueryTimeoutError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.QueryTimeoutError: " + msg
}

func (e *QueryTimeoutError) Unwrap() error { return e.err }

func (e *QueryTimeoutError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.QueryTimeoutError:
		return true
	case gelerr.SessionTimeoutError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *QueryTimeoutError) isEdgeDBSessionTimeoutError() {}

func (e *QueryTimeoutError) isEdgeDBQueryError() {}

func (e *QueryTimeoutError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewTransactionTimeoutError(msg string, err error) error {
	return &TransactionTimeoutError{msg, err}
}

type TransactionTimeoutError struct {
	msg string
	err error
}

func (e *TransactionTimeoutError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.TransactionTimeoutError: " + msg
}

func (e *TransactionTimeoutError) Unwrap() error { return e.err }

func (e *TransactionTimeoutError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.TransactionTimeoutError:
		return true
	case gelerr.SessionTimeoutError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *TransactionTimeoutError) isEdgeDBSessionTimeoutError() {}

func (e *TransactionTimeoutError) isEdgeDBQueryError() {}

func (e *TransactionTimeoutError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewIdleTransactionTimeoutError(msg string, err error) error {
	return &IdleTransactionTimeoutError{msg, err}
}

type IdleTransactionTimeoutError struct {
	msg string
	err error
}

func (e *IdleTransactionTimeoutError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.IdleTransactionTimeoutError: " + msg
}

func (e *IdleTransactionTimeoutError) Unwrap() error { return e.err }

func (e *IdleTransactionTimeoutError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.IdleTransactionTimeoutError:
		return true
	case gelerr.TransactionTimeoutError:
		return true
	case gelerr.SessionTimeoutError:
		return true
	case gelerr.QueryError:
		return true
	default:
		return false
	}
}

func (e *IdleTransactionTimeoutError) isEdgeDBTransactionTimeoutError() {}

func (e *IdleTransactionTimeoutError) isEdgeDBSessionTimeoutError() {}

func (e *IdleTransactionTimeoutError) isEdgeDBQueryError() {}

func (e *IdleTransactionTimeoutError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewExecutionError(msg string, err error) error {
	return &ExecutionError{msg, err}
}

type ExecutionError struct {
	msg string
	err error
}

func (e *ExecutionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ExecutionError: " + msg
}

func (e *ExecutionError) Unwrap() error { return e.err }

func (e *ExecutionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *ExecutionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidValueError(msg string, err error) error {
	return &InvalidValueError{msg, err}
}

type InvalidValueError struct {
	msg string
	err error
}

func (e *InvalidValueError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidValueError: " + msg
}

func (e *InvalidValueError) Unwrap() error { return e.err }

func (e *InvalidValueError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidValueError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *InvalidValueError) isEdgeDBExecutionError() {}

func (e *InvalidValueError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewDivisionByZeroError(msg string, err error) error {
	return &DivisionByZeroError{msg, err}
}

type DivisionByZeroError struct {
	msg string
	err error
}

func (e *DivisionByZeroError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.DivisionByZeroError: " + msg
}

func (e *DivisionByZeroError) Unwrap() error { return e.err }

func (e *DivisionByZeroError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.DivisionByZeroError:
		return true
	case gelerr.InvalidValueError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *DivisionByZeroError) isEdgeDBInvalidValueError() {}

func (e *DivisionByZeroError) isEdgeDBExecutionError() {}

func (e *DivisionByZeroError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewNumericOutOfRangeError(msg string, err error) error {
	return &NumericOutOfRangeError{msg, err}
}

type NumericOutOfRangeError struct {
	msg string
	err error
}

func (e *NumericOutOfRangeError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.NumericOutOfRangeError: " + msg
}

func (e *NumericOutOfRangeError) Unwrap() error { return e.err }

func (e *NumericOutOfRangeError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.NumericOutOfRangeError:
		return true
	case gelerr.InvalidValueError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *NumericOutOfRangeError) isEdgeDBInvalidValueError() {}

func (e *NumericOutOfRangeError) isEdgeDBExecutionError() {}

func (e *NumericOutOfRangeError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewAccessPolicyError(msg string, err error) error {
	return &AccessPolicyError{msg, err}
}

type AccessPolicyError struct {
	msg string
	err error
}

func (e *AccessPolicyError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.AccessPolicyError: " + msg
}

func (e *AccessPolicyError) Unwrap() error { return e.err }

func (e *AccessPolicyError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.AccessPolicyError:
		return true
	case gelerr.InvalidValueError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *AccessPolicyError) isEdgeDBInvalidValueError() {}

func (e *AccessPolicyError) isEdgeDBExecutionError() {}

func (e *AccessPolicyError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewQueryAssertionError(msg string, err error) error {
	return &QueryAssertionError{msg, err}
}

type QueryAssertionError struct {
	msg string
	err error
}

func (e *QueryAssertionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.QueryAssertionError: " + msg
}

func (e *QueryAssertionError) Unwrap() error { return e.err }

func (e *QueryAssertionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.QueryAssertionError:
		return true
	case gelerr.InvalidValueError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *QueryAssertionError) isEdgeDBInvalidValueError() {}

func (e *QueryAssertionError) isEdgeDBExecutionError() {}

func (e *QueryAssertionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewIntegrityError(msg string, err error) error {
	return &IntegrityError{msg, err}
}

type IntegrityError struct {
	msg string
	err error
}

func (e *IntegrityError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.IntegrityError: " + msg
}

func (e *IntegrityError) Unwrap() error { return e.err }

func (e *IntegrityError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.IntegrityError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *IntegrityError) isEdgeDBExecutionError() {}

func (e *IntegrityError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewConstraintViolationError(msg string, err error) error {
	return &ConstraintViolationError{msg, err}
}

type ConstraintViolationError struct {
	msg string
	err error
}

func (e *ConstraintViolationError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ConstraintViolationError: " + msg
}

func (e *ConstraintViolationError) Unwrap() error { return e.err }

func (e *ConstraintViolationError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ConstraintViolationError:
		return true
	case gelerr.IntegrityError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *ConstraintViolationError) isEdgeDBIntegrityError() {}

func (e *ConstraintViolationError) isEdgeDBExecutionError() {}

func (e *ConstraintViolationError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewCardinalityViolationError(msg string, err error) error {
	return &CardinalityViolationError{msg, err}
}

type CardinalityViolationError struct {
	msg string
	err error
}

func (e *CardinalityViolationError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.CardinalityViolationError: " + msg
}

func (e *CardinalityViolationError) Unwrap() error { return e.err }

func (e *CardinalityViolationError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.CardinalityViolationError:
		return true
	case gelerr.IntegrityError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *CardinalityViolationError) isEdgeDBIntegrityError() {}

func (e *CardinalityViolationError) isEdgeDBExecutionError() {}

func (e *CardinalityViolationError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewMissingRequiredError(msg string, err error) error {
	return &MissingRequiredError{msg, err}
}

type MissingRequiredError struct {
	msg string
	err error
}

func (e *MissingRequiredError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.MissingRequiredError: " + msg
}

func (e *MissingRequiredError) Unwrap() error { return e.err }

func (e *MissingRequiredError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.MissingRequiredError:
		return true
	case gelerr.IntegrityError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *MissingRequiredError) isEdgeDBIntegrityError() {}

func (e *MissingRequiredError) isEdgeDBExecutionError() {}

func (e *MissingRequiredError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewTransactionError(msg string, err error) error {
	return &TransactionError{msg, err}
}

type TransactionError struct {
	msg string
	err error
}

func (e *TransactionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.TransactionError: " + msg
}

func (e *TransactionError) Unwrap() error { return e.err }

func (e *TransactionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.TransactionError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *TransactionError) isEdgeDBExecutionError() {}

func (e *TransactionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewTransactionConflictError(msg string, err error) error {
	return &TransactionConflictError{msg, err}
}

type TransactionConflictError struct {
	msg string
	err error
}

func (e *TransactionConflictError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.TransactionConflictError: " + msg
}

func (e *TransactionConflictError) Unwrap() error { return e.err }

func (e *TransactionConflictError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.TransactionConflictError:
		return true
	case gelerr.TransactionError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *TransactionConflictError) isEdgeDBTransactionError() {}

func (e *TransactionConflictError) isEdgeDBExecutionError() {}

func (e *TransactionConflictError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewTransactionSerializationError(msg string, err error) error {
	return &TransactionSerializationError{msg, err}
}

type TransactionSerializationError struct {
	msg string
	err error
}

func (e *TransactionSerializationError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.TransactionSerializationError: " + msg
}

func (e *TransactionSerializationError) Unwrap() error { return e.err }

func (e *TransactionSerializationError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.TransactionSerializationError:
		return true
	case gelerr.TransactionConflictError:
		return true
	case gelerr.TransactionError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *TransactionSerializationError) isEdgeDBTransactionConflictError() {}

func (e *TransactionSerializationError) isEdgeDBTransactionError() {}

func (e *TransactionSerializationError) isEdgeDBExecutionError() {}

func (e *TransactionSerializationError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewTransactionDeadlockError(msg string, err error) error {
	return &TransactionDeadlockError{msg, err}
}

type TransactionDeadlockError struct {
	msg string
	err error
}

func (e *TransactionDeadlockError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.TransactionDeadlockError: " + msg
}

func (e *TransactionDeadlockError) Unwrap() error { return e.err }

func (e *TransactionDeadlockError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.TransactionDeadlockError:
		return true
	case gelerr.TransactionConflictError:
		return true
	case gelerr.TransactionError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *TransactionDeadlockError) isEdgeDBTransactionConflictError() {}

func (e *TransactionDeadlockError) isEdgeDBTransactionError() {}

func (e *TransactionDeadlockError) isEdgeDBExecutionError() {}

func (e *TransactionDeadlockError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewWatchError(msg string, err error) error {
	return &WatchError{msg, err}
}

type WatchError struct {
	msg string
	err error
}

func (e *WatchError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.WatchError: " + msg
}

func (e *WatchError) Unwrap() error { return e.err }

func (e *WatchError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.WatchError:
		return true
	case gelerr.ExecutionError:
		return true
	default:
		return false
	}
}

func (e *WatchError) isEdgeDBExecutionError() {}

func (e *WatchError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewConfigurationError(msg string, err error) error {
	return &ConfigurationError{msg, err}
}

type ConfigurationError struct {
	msg string
	err error
}

func (e *ConfigurationError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ConfigurationError: " + msg
}

func (e *ConfigurationError) Unwrap() error { return e.err }

func (e *ConfigurationError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ConfigurationError:
		return true
	default:
		return false
	}
}

func (e *ConfigurationError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewAccessError(msg string, err error) error {
	return &AccessError{msg, err}
}

type AccessError struct {
	msg string
	err error
}

func (e *AccessError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.AccessError: " + msg
}

func (e *AccessError) Unwrap() error { return e.err }

func (e *AccessError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.AccessError:
		return true
	default:
		return false
	}
}

func (e *AccessError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewAuthenticationError(msg string, err error) error {
	return &AuthenticationError{msg, err}
}

type AuthenticationError struct {
	msg string
	err error
}

func (e *AuthenticationError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.AuthenticationError: " + msg
}

func (e *AuthenticationError) Unwrap() error { return e.err }

func (e *AuthenticationError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.AuthenticationError:
		return true
	case gelerr.AccessError:
		return true
	default:
		return false
	}
}

func (e *AuthenticationError) isEdgeDBAccessError() {}

func (e *AuthenticationError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewAvailabilityError(msg string, err error) error {
	return &AvailabilityError{msg, err}
}

type AvailabilityError struct {
	msg string
	err error
}

func (e *AvailabilityError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.AvailabilityError: " + msg
}

func (e *AvailabilityError) Unwrap() error { return e.err }

func (e *AvailabilityError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.AvailabilityError:
		return true
	default:
		return false
	}
}

func (e *AvailabilityError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewBackendUnavailableError(msg string, err error) error {
	return &BackendUnavailableError{msg, err}
}

type BackendUnavailableError struct {
	msg string
	err error
}

func (e *BackendUnavailableError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.BackendUnavailableError: " + msg
}

func (e *BackendUnavailableError) Unwrap() error { return e.err }

func (e *BackendUnavailableError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.BackendUnavailableError:
		return true
	case gelerr.AvailabilityError:
		return true
	default:
		return false
	}
}

func (e *BackendUnavailableError) isEdgeDBAvailabilityError() {}

func (e *BackendUnavailableError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewServerOfflineError(msg string, err error) error {
	return &ServerOfflineError{msg, err}
}

type ServerOfflineError struct {
	msg string
	err error
}

func (e *ServerOfflineError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ServerOfflineError: " + msg
}

func (e *ServerOfflineError) Unwrap() error { return e.err }

func (e *ServerOfflineError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ServerOfflineError:
		return true
	case gelerr.AvailabilityError:
		return true
	default:
		return false
	}
}

func (e *ServerOfflineError) isEdgeDBAvailabilityError() {}

func (e *ServerOfflineError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldReconnect:
		return true
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewUnknownTenantError(msg string, err error) error {
	return &UnknownTenantError{msg, err}
}

type UnknownTenantError struct {
	msg string
	err error
}

func (e *UnknownTenantError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownTenantError: " + msg
}

func (e *UnknownTenantError) Unwrap() error { return e.err }

func (e *UnknownTenantError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownTenantError:
		return true
	case gelerr.AvailabilityError:
		return true
	default:
		return false
	}
}

func (e *UnknownTenantError) isEdgeDBAvailabilityError() {}

func (e *UnknownTenantError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldReconnect:
		return true
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewServerBlockedError(msg string, err error) error {
	return &ServerBlockedError{msg, err}
}

type ServerBlockedError struct {
	msg string
	err error
}

func (e *ServerBlockedError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ServerBlockedError: " + msg
}

func (e *ServerBlockedError) Unwrap() error { return e.err }

func (e *ServerBlockedError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ServerBlockedError:
		return true
	case gelerr.AvailabilityError:
		return true
	default:
		return false
	}
}

func (e *ServerBlockedError) isEdgeDBAvailabilityError() {}

func (e *ServerBlockedError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewBackendError(msg string, err error) error {
	return &BackendError{msg, err}
}

type BackendError struct {
	msg string
	err error
}

func (e *BackendError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.BackendError: " + msg
}

func (e *BackendError) Unwrap() error { return e.err }

func (e *BackendError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.BackendError:
		return true
	default:
		return false
	}
}

func (e *BackendError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnsupportedBackendFeatureError(msg string, err error) error {
	return &UnsupportedBackendFeatureError{msg, err}
}

type UnsupportedBackendFeatureError struct {
	msg string
	err error
}

func (e *UnsupportedBackendFeatureError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnsupportedBackendFeatureError: " + msg
}

func (e *UnsupportedBackendFeatureError) Unwrap() error { return e.err }

func (e *UnsupportedBackendFeatureError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnsupportedBackendFeatureError:
		return true
	case gelerr.BackendError:
		return true
	default:
		return false
	}
}

func (e *UnsupportedBackendFeatureError) isEdgeDBBackendError() {}

func (e *UnsupportedBackendFeatureError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewClientError(msg string, err error) error {
	return &ClientError{msg, err}
}

type ClientError struct {
	msg string
	err error
}

func (e *ClientError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ClientError: " + msg
}

func (e *ClientError) Unwrap() error { return e.err }

func (e *ClientError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *ClientError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewClientConnectionError(msg string, err error) error {
	return &ClientConnectionError{msg, err}
}

type ClientConnectionError struct {
	msg string
	err error
}

func (e *ClientConnectionError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ClientConnectionError: " + msg
}

func (e *ClientConnectionError) Unwrap() error { return e.err }

func (e *ClientConnectionError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ClientConnectionError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *ClientConnectionError) isEdgeDBClientError() {}

func (e *ClientConnectionError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewClientConnectionFailedError(msg string, err error) error {
	return &ClientConnectionFailedError{msg, err}
}

type ClientConnectionFailedError struct {
	msg string
	err error
}

func (e *ClientConnectionFailedError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ClientConnectionFailedError: " + msg
}

func (e *ClientConnectionFailedError) Unwrap() error { return e.err }

func (e *ClientConnectionFailedError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ClientConnectionFailedError:
		return true
	case gelerr.ClientConnectionError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *ClientConnectionFailedError) isEdgeDBClientConnectionError() {}

func (e *ClientConnectionFailedError) isEdgeDBClientError() {}

func (e *ClientConnectionFailedError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewClientConnectionFailedTemporarilyError(msg string, err error) error {
	return &ClientConnectionFailedTemporarilyError{msg, err}
}

type ClientConnectionFailedTemporarilyError struct {
	msg string
	err error
}

func (e *ClientConnectionFailedTemporarilyError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ClientConnectionFailedTemporarilyError: " + msg
}

func (e *ClientConnectionFailedTemporarilyError) Unwrap() error { return e.err }

func (e *ClientConnectionFailedTemporarilyError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ClientConnectionFailedTemporarilyError:
		return true
	case gelerr.ClientConnectionFailedError:
		return true
	case gelerr.ClientConnectionError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *ClientConnectionFailedTemporarilyError) isEdgeDBClientConnectionFailedError() {}

func (e *ClientConnectionFailedTemporarilyError) isEdgeDBClientConnectionError() {}

func (e *ClientConnectionFailedTemporarilyError) isEdgeDBClientError() {}

func (e *ClientConnectionFailedTemporarilyError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldReconnect:
		return true
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewClientConnectionTimeoutError(msg string, err error) error {
	return &ClientConnectionTimeoutError{msg, err}
}

type ClientConnectionTimeoutError struct {
	msg string
	err error
}

func (e *ClientConnectionTimeoutError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ClientConnectionTimeoutError: " + msg
}

func (e *ClientConnectionTimeoutError) Unwrap() error { return e.err }

func (e *ClientConnectionTimeoutError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ClientConnectionTimeoutError:
		return true
	case gelerr.ClientConnectionError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *ClientConnectionTimeoutError) isEdgeDBClientConnectionError() {}

func (e *ClientConnectionTimeoutError) isEdgeDBClientError() {}

func (e *ClientConnectionTimeoutError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldReconnect:
		return true
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewClientConnectionClosedError(msg string, err error) error {
	return &ClientConnectionClosedError{msg, err}
}

type ClientConnectionClosedError struct {
	msg string
	err error
}

func (e *ClientConnectionClosedError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.ClientConnectionClosedError: " + msg
}

func (e *ClientConnectionClosedError) Unwrap() error { return e.err }

func (e *ClientConnectionClosedError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.ClientConnectionClosedError:
		return true
	case gelerr.ClientConnectionError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *ClientConnectionClosedError) isEdgeDBClientConnectionError() {}

func (e *ClientConnectionClosedError) isEdgeDBClientError() {}

func (e *ClientConnectionClosedError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	case gelerr.ShouldReconnect:
		return true
	case gelerr.ShouldRetry:
		return true
	default:
		return false
	}
}

func NewInterfaceError(msg string, err error) error {
	return &InterfaceError{msg, err}
}

type InterfaceError struct {
	msg string
	err error
}

func (e *InterfaceError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InterfaceError: " + msg
}

func (e *InterfaceError) Unwrap() error { return e.err }

func (e *InterfaceError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InterfaceError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *InterfaceError) isEdgeDBClientError() {}

func (e *InterfaceError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewQueryArgumentError(msg string, err error) error {
	return &QueryArgumentError{msg, err}
}

type QueryArgumentError struct {
	msg string
	err error
}

func (e *QueryArgumentError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.QueryArgumentError: " + msg
}

func (e *QueryArgumentError) Unwrap() error { return e.err }

func (e *QueryArgumentError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.QueryArgumentError:
		return true
	case gelerr.InterfaceError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *QueryArgumentError) isEdgeDBInterfaceError() {}

func (e *QueryArgumentError) isEdgeDBClientError() {}

func (e *QueryArgumentError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewMissingArgumentError(msg string, err error) error {
	return &MissingArgumentError{msg, err}
}

type MissingArgumentError struct {
	msg string
	err error
}

func (e *MissingArgumentError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.MissingArgumentError: " + msg
}

func (e *MissingArgumentError) Unwrap() error { return e.err }

func (e *MissingArgumentError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.MissingArgumentError:
		return true
	case gelerr.QueryArgumentError:
		return true
	case gelerr.InterfaceError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *MissingArgumentError) isEdgeDBQueryArgumentError() {}

func (e *MissingArgumentError) isEdgeDBInterfaceError() {}

func (e *MissingArgumentError) isEdgeDBClientError() {}

func (e *MissingArgumentError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewUnknownArgumentError(msg string, err error) error {
	return &UnknownArgumentError{msg, err}
}

type UnknownArgumentError struct {
	msg string
	err error
}

func (e *UnknownArgumentError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.UnknownArgumentError: " + msg
}

func (e *UnknownArgumentError) Unwrap() error { return e.err }

func (e *UnknownArgumentError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.UnknownArgumentError:
		return true
	case gelerr.QueryArgumentError:
		return true
	case gelerr.InterfaceError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *UnknownArgumentError) isEdgeDBQueryArgumentError() {}

func (e *UnknownArgumentError) isEdgeDBInterfaceError() {}

func (e *UnknownArgumentError) isEdgeDBClientError() {}

func (e *UnknownArgumentError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInvalidArgumentError(msg string, err error) error {
	return &InvalidArgumentError{msg, err}
}

type InvalidArgumentError struct {
	msg string
	err error
}

func (e *InvalidArgumentError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InvalidArgumentError: " + msg
}

func (e *InvalidArgumentError) Unwrap() error { return e.err }

func (e *InvalidArgumentError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InvalidArgumentError:
		return true
	case gelerr.QueryArgumentError:
		return true
	case gelerr.InterfaceError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *InvalidArgumentError) isEdgeDBQueryArgumentError() {}

func (e *InvalidArgumentError) isEdgeDBInterfaceError() {}

func (e *InvalidArgumentError) isEdgeDBClientError() {}

func (e *InvalidArgumentError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewNoDataError(msg string, err error) error {
	return &NoDataError{msg, err}
}

type NoDataError struct {
	msg string
	err error
}

func (e *NoDataError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.NoDataError: " + msg
}

func (e *NoDataError) Unwrap() error { return e.err }

func (e *NoDataError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.NoDataError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *NoDataError) isEdgeDBClientError() {}

func (e *NoDataError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func NewInternalClientError(msg string, err error) error {
	return &InternalClientError{msg, err}
}

type InternalClientError struct {
	msg string
	err error
}

func (e *InternalClientError) Error() string {
	msg := e.msg
	if e.err != nil {
		msg = e.err.Error()
	}

	return "gel.InternalClientError: " + msg
}

func (e *InternalClientError) Unwrap() error { return e.err }

func (e *InternalClientError) Category(c gelerr.ErrorCategory) bool {
	switch c {
	case gelerr.InternalClientError:
		return true
	case gelerr.ClientError:
		return true
	default:
		return false
	}
}

func (e *InternalClientError) isEdgeDBClientError() {}

func (e *InternalClientError) HasTag(tag gelerr.ErrorTag) bool {
	switch tag {
	default:
		return false
	}
}

func ErrorFromCode(code uint32, msg string) error {
	switch code {
	case 0x01_00_00_00:
		return &InternalServerError{msg: msg}
	case 0x02_00_00_00:
		return &UnsupportedFeatureError{msg: msg}
	case 0x03_00_00_00:
		return &ProtocolError{msg: msg}
	case 0x03_01_00_00:
		return &BinaryProtocolError{msg: msg}
	case 0x03_01_00_01:
		return &UnsupportedProtocolVersionError{msg: msg}
	case 0x03_01_00_02:
		return &TypeSpecNotFoundError{msg: msg}
	case 0x03_01_00_03:
		return &UnexpectedMessageError{msg: msg}
	case 0x03_02_00_00:
		return &InputDataError{msg: msg}
	case 0x03_02_01_00:
		return &ParameterTypeMismatchError{msg: msg}
	case 0x03_02_02_00:
		return &StateMismatchError{msg: msg}
	case 0x03_03_00_00:
		return &ResultCardinalityMismatchError{msg: msg}
	case 0x03_04_00_00:
		return &CapabilityError{msg: msg}
	case 0x03_04_01_00:
		return &UnsupportedCapabilityError{msg: msg}
	case 0x03_04_02_00:
		return &DisabledCapabilityError{msg: msg}
	case 0x04_00_00_00:
		return &QueryError{msg: msg}
	case 0x04_01_00_00:
		return &InvalidSyntaxError{msg: msg}
	case 0x04_01_01_00:
		return &EdgeQLSyntaxError{msg: msg}
	case 0x04_01_02_00:
		return &SchemaSyntaxError{msg: msg}
	case 0x04_01_03_00:
		return &GraphQLSyntaxError{msg: msg}
	case 0x04_02_00_00:
		return &InvalidTypeError{msg: msg}
	case 0x04_02_01_00:
		return &InvalidTargetError{msg: msg}
	case 0x04_02_01_01:
		return &InvalidLinkTargetError{msg: msg}
	case 0x04_02_01_02:
		return &InvalidPropertyTargetError{msg: msg}
	case 0x04_03_00_00:
		return &InvalidReferenceError{msg: msg}
	case 0x04_03_00_01:
		return &UnknownModuleError{msg: msg}
	case 0x04_03_00_02:
		return &UnknownLinkError{msg: msg}
	case 0x04_03_00_03:
		return &UnknownPropertyError{msg: msg}
	case 0x04_03_00_04:
		return &UnknownUserError{msg: msg}
	case 0x04_03_00_05:
		return &UnknownDatabaseError{msg: msg}
	case 0x04_03_00_06:
		return &UnknownParameterError{msg: msg}
	case 0x04_03_00_07:
		return &DeprecatedScopingError{msg: msg}
	case 0x04_04_00_00:
		return &SchemaError{msg: msg}
	case 0x04_05_00_00:
		return &SchemaDefinitionError{msg: msg}
	case 0x04_05_01_00:
		return &InvalidDefinitionError{msg: msg}
	case 0x04_05_01_01:
		return &InvalidModuleDefinitionError{msg: msg}
	case 0x04_05_01_02:
		return &InvalidLinkDefinitionError{msg: msg}
	case 0x04_05_01_03:
		return &InvalidPropertyDefinitionError{msg: msg}
	case 0x04_05_01_04:
		return &InvalidUserDefinitionError{msg: msg}
	case 0x04_05_01_05:
		return &InvalidDatabaseDefinitionError{msg: msg}
	case 0x04_05_01_06:
		return &InvalidOperatorDefinitionError{msg: msg}
	case 0x04_05_01_07:
		return &InvalidAliasDefinitionError{msg: msg}
	case 0x04_05_01_08:
		return &InvalidFunctionDefinitionError{msg: msg}
	case 0x04_05_01_09:
		return &InvalidConstraintDefinitionError{msg: msg}
	case 0x04_05_01_0a:
		return &InvalidCastDefinitionError{msg: msg}
	case 0x04_05_02_00:
		return &DuplicateDefinitionError{msg: msg}
	case 0x04_05_02_01:
		return &DuplicateModuleDefinitionError{msg: msg}
	case 0x04_05_02_02:
		return &DuplicateLinkDefinitionError{msg: msg}
	case 0x04_05_02_03:
		return &DuplicatePropertyDefinitionError{msg: msg}
	case 0x04_05_02_04:
		return &DuplicateUserDefinitionError{msg: msg}
	case 0x04_05_02_05:
		return &DuplicateDatabaseDefinitionError{msg: msg}
	case 0x04_05_02_06:
		return &DuplicateOperatorDefinitionError{msg: msg}
	case 0x04_05_02_07:
		return &DuplicateViewDefinitionError{msg: msg}
	case 0x04_05_02_08:
		return &DuplicateFunctionDefinitionError{msg: msg}
	case 0x04_05_02_09:
		return &DuplicateConstraintDefinitionError{msg: msg}
	case 0x04_05_02_0a:
		return &DuplicateCastDefinitionError{msg: msg}
	case 0x04_05_02_0b:
		return &DuplicateMigrationError{msg: msg}
	case 0x04_06_00_00:
		return &SessionTimeoutError{msg: msg}
	case 0x04_06_01_00:
		return &IdleSessionTimeoutError{msg: msg}
	case 0x04_06_02_00:
		return &QueryTimeoutError{msg: msg}
	case 0x04_06_0a_00:
		return &TransactionTimeoutError{msg: msg}
	case 0x04_06_0a_01:
		return &IdleTransactionTimeoutError{msg: msg}
	case 0x05_00_00_00:
		return &ExecutionError{msg: msg}
	case 0x05_01_00_00:
		return &InvalidValueError{msg: msg}
	case 0x05_01_00_01:
		return &DivisionByZeroError{msg: msg}
	case 0x05_01_00_02:
		return &NumericOutOfRangeError{msg: msg}
	case 0x05_01_00_03:
		return &AccessPolicyError{msg: msg}
	case 0x05_01_00_04:
		return &QueryAssertionError{msg: msg}
	case 0x05_02_00_00:
		return &IntegrityError{msg: msg}
	case 0x05_02_00_01:
		return &ConstraintViolationError{msg: msg}
	case 0x05_02_00_02:
		return &CardinalityViolationError{msg: msg}
	case 0x05_02_00_03:
		return &MissingRequiredError{msg: msg}
	case 0x05_03_00_00:
		return &TransactionError{msg: msg}
	case 0x05_03_01_00:
		return &TransactionConflictError{msg: msg}
	case 0x05_03_01_01:
		return &TransactionSerializationError{msg: msg}
	case 0x05_03_01_02:
		return &TransactionDeadlockError{msg: msg}
	case 0x05_04_00_00:
		return &WatchError{msg: msg}
	case 0x06_00_00_00:
		return &ConfigurationError{msg: msg}
	case 0x07_00_00_00:
		return &AccessError{msg: msg}
	case 0x07_01_00_00:
		return &AuthenticationError{msg: msg}
	case 0x08_00_00_00:
		return &AvailabilityError{msg: msg}
	case 0x08_00_00_01:
		return &BackendUnavailableError{msg: msg}
	case 0x08_00_00_02:
		return &ServerOfflineError{msg: msg}
	case 0x08_00_00_03:
		return &UnknownTenantError{msg: msg}
	case 0x08_00_00_04:
		return &ServerBlockedError{msg: msg}
	case 0x09_00_00_00:
		return &BackendError{msg: msg}
	case 0x09_00_01_00:
		return &UnsupportedBackendFeatureError{msg: msg}
	case 0xff_00_00_00:
		return &ClientError{msg: msg}
	case 0xff_01_00_00:
		return &ClientConnectionError{msg: msg}
	case 0xff_01_01_00:
		return &ClientConnectionFailedError{msg: msg}
	case 0xff_01_01_01:
		return &ClientConnectionFailedTemporarilyError{msg: msg}
	case 0xff_01_02_00:
		return &ClientConnectionTimeoutError{msg: msg}
	case 0xff_01_03_00:
		return &ClientConnectionClosedError{msg: msg}
	case 0xff_02_00_00:
		return &InterfaceError{msg: msg}
	case 0xff_02_01_00:
		return &QueryArgumentError{msg: msg}
	case 0xff_02_01_01:
		return &MissingArgumentError{msg: msg}
	case 0xff_02_01_02:
		return &UnknownArgumentError{msg: msg}
	case 0xff_02_01_03:
		return &InvalidArgumentError{msg: msg}
	case 0xff_03_00_00:
		return &NoDataError{msg: msg}
	case 0xff_04_00_00:
		return &InternalClientError{msg: msg}
	default:
		return &UnexpectedMessageError{
			msg: fmt.Sprintf(
				"invalid error code 0x%x with message %q", code, msg,
			),
		}
	}
}
