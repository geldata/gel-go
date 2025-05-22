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

// This file is auto generated. Do not edit!
// run 'make errors' to regenerate

package gelerr

const (
	ShouldRetry     ErrorTag = "SHOULD_RETRY"
	ShouldReconnect ErrorTag = "SHOULD_RECONNECT"
)

const (
	InternalServerError                    ErrorCategory = "errors::InternalServerError"
	UnsupportedFeatureError                ErrorCategory = "errors::UnsupportedFeatureError"
	ProtocolError                          ErrorCategory = "errors::ProtocolError"
	BinaryProtocolError                    ErrorCategory = "errors::BinaryProtocolError"
	UnsupportedProtocolVersionError        ErrorCategory = "errors::UnsupportedProtocolVersionError"
	TypeSpecNotFoundError                  ErrorCategory = "errors::TypeSpecNotFoundError"
	UnexpectedMessageError                 ErrorCategory = "errors::UnexpectedMessageError"
	InputDataError                         ErrorCategory = "errors::InputDataError"
	ParameterTypeMismatchError             ErrorCategory = "errors::ParameterTypeMismatchError"
	StateMismatchError                     ErrorCategory = "errors::StateMismatchError"
	ResultCardinalityMismatchError         ErrorCategory = "errors::ResultCardinalityMismatchError"
	CapabilityError                        ErrorCategory = "errors::CapabilityError"
	UnsupportedCapabilityError             ErrorCategory = "errors::UnsupportedCapabilityError"
	DisabledCapabilityError                ErrorCategory = "errors::DisabledCapabilityError"
	UnsafeIsolationLevelError              ErrorCategory = "errors::UnsafeIsolationLevelError"
	QueryError                             ErrorCategory = "errors::QueryError"
	InvalidSyntaxError                     ErrorCategory = "errors::InvalidSyntaxError"
	EdgeQLSyntaxError                      ErrorCategory = "errors::EdgeQLSyntaxError"
	SchemaSyntaxError                      ErrorCategory = "errors::SchemaSyntaxError"
	GraphQLSyntaxError                     ErrorCategory = "errors::GraphQLSyntaxError"
	InvalidTypeError                       ErrorCategory = "errors::InvalidTypeError"
	InvalidTargetError                     ErrorCategory = "errors::InvalidTargetError"
	InvalidLinkTargetError                 ErrorCategory = "errors::InvalidLinkTargetError"
	InvalidPropertyTargetError             ErrorCategory = "errors::InvalidPropertyTargetError"
	InvalidReferenceError                  ErrorCategory = "errors::InvalidReferenceError"
	UnknownModuleError                     ErrorCategory = "errors::UnknownModuleError"
	UnknownLinkError                       ErrorCategory = "errors::UnknownLinkError"
	UnknownPropertyError                   ErrorCategory = "errors::UnknownPropertyError"
	UnknownUserError                       ErrorCategory = "errors::UnknownUserError"
	UnknownDatabaseError                   ErrorCategory = "errors::UnknownDatabaseError"
	UnknownParameterError                  ErrorCategory = "errors::UnknownParameterError"
	DeprecatedScopingError                 ErrorCategory = "errors::DeprecatedScopingError"
	SchemaError                            ErrorCategory = "errors::SchemaError"
	SchemaDefinitionError                  ErrorCategory = "errors::SchemaDefinitionError"
	InvalidDefinitionError                 ErrorCategory = "errors::InvalidDefinitionError"
	InvalidModuleDefinitionError           ErrorCategory = "errors::InvalidModuleDefinitionError"
	InvalidLinkDefinitionError             ErrorCategory = "errors::InvalidLinkDefinitionError"
	InvalidPropertyDefinitionError         ErrorCategory = "errors::InvalidPropertyDefinitionError"
	InvalidUserDefinitionError             ErrorCategory = "errors::InvalidUserDefinitionError"
	InvalidDatabaseDefinitionError         ErrorCategory = "errors::InvalidDatabaseDefinitionError"
	InvalidOperatorDefinitionError         ErrorCategory = "errors::InvalidOperatorDefinitionError"
	InvalidAliasDefinitionError            ErrorCategory = "errors::InvalidAliasDefinitionError"
	InvalidFunctionDefinitionError         ErrorCategory = "errors::InvalidFunctionDefinitionError"
	InvalidConstraintDefinitionError       ErrorCategory = "errors::InvalidConstraintDefinitionError"
	InvalidCastDefinitionError             ErrorCategory = "errors::InvalidCastDefinitionError"
	DuplicateDefinitionError               ErrorCategory = "errors::DuplicateDefinitionError"
	DuplicateModuleDefinitionError         ErrorCategory = "errors::DuplicateModuleDefinitionError"
	DuplicateLinkDefinitionError           ErrorCategory = "errors::DuplicateLinkDefinitionError"
	DuplicatePropertyDefinitionError       ErrorCategory = "errors::DuplicatePropertyDefinitionError"
	DuplicateUserDefinitionError           ErrorCategory = "errors::DuplicateUserDefinitionError"
	DuplicateDatabaseDefinitionError       ErrorCategory = "errors::DuplicateDatabaseDefinitionError"
	DuplicateOperatorDefinitionError       ErrorCategory = "errors::DuplicateOperatorDefinitionError"
	DuplicateViewDefinitionError           ErrorCategory = "errors::DuplicateViewDefinitionError"
	DuplicateFunctionDefinitionError       ErrorCategory = "errors::DuplicateFunctionDefinitionError"
	DuplicateConstraintDefinitionError     ErrorCategory = "errors::DuplicateConstraintDefinitionError"
	DuplicateCastDefinitionError           ErrorCategory = "errors::DuplicateCastDefinitionError"
	DuplicateMigrationError                ErrorCategory = "errors::DuplicateMigrationError"
	SessionTimeoutError                    ErrorCategory = "errors::SessionTimeoutError"
	IdleSessionTimeoutError                ErrorCategory = "errors::IdleSessionTimeoutError"
	QueryTimeoutError                      ErrorCategory = "errors::QueryTimeoutError"
	TransactionTimeoutError                ErrorCategory = "errors::TransactionTimeoutError"
	IdleTransactionTimeoutError            ErrorCategory = "errors::IdleTransactionTimeoutError"
	ExecutionError                         ErrorCategory = "errors::ExecutionError"
	InvalidValueError                      ErrorCategory = "errors::InvalidValueError"
	DivisionByZeroError                    ErrorCategory = "errors::DivisionByZeroError"
	NumericOutOfRangeError                 ErrorCategory = "errors::NumericOutOfRangeError"
	AccessPolicyError                      ErrorCategory = "errors::AccessPolicyError"
	QueryAssertionError                    ErrorCategory = "errors::QueryAssertionError"
	IntegrityError                         ErrorCategory = "errors::IntegrityError"
	ConstraintViolationError               ErrorCategory = "errors::ConstraintViolationError"
	CardinalityViolationError              ErrorCategory = "errors::CardinalityViolationError"
	MissingRequiredError                   ErrorCategory = "errors::MissingRequiredError"
	TransactionError                       ErrorCategory = "errors::TransactionError"
	TransactionConflictError               ErrorCategory = "errors::TransactionConflictError"
	TransactionSerializationError          ErrorCategory = "errors::TransactionSerializationError"
	TransactionDeadlockError               ErrorCategory = "errors::TransactionDeadlockError"
	QueryCacheInvalidationError            ErrorCategory = "errors::QueryCacheInvalidationError"
	WatchError                             ErrorCategory = "errors::WatchError"
	ConfigurationError                     ErrorCategory = "errors::ConfigurationError"
	AccessError                            ErrorCategory = "errors::AccessError"
	AuthenticationError                    ErrorCategory = "errors::AuthenticationError"
	AvailabilityError                      ErrorCategory = "errors::AvailabilityError"
	BackendUnavailableError                ErrorCategory = "errors::BackendUnavailableError"
	ServerOfflineError                     ErrorCategory = "errors::ServerOfflineError"
	UnknownTenantError                     ErrorCategory = "errors::UnknownTenantError"
	ServerBlockedError                     ErrorCategory = "errors::ServerBlockedError"
	BackendError                           ErrorCategory = "errors::BackendError"
	UnsupportedBackendFeatureError         ErrorCategory = "errors::UnsupportedBackendFeatureError"
	ClientError                            ErrorCategory = "errors::ClientError"
	ClientConnectionError                  ErrorCategory = "errors::ClientConnectionError"
	ClientConnectionFailedError            ErrorCategory = "errors::ClientConnectionFailedError"
	ClientConnectionFailedTemporarilyError ErrorCategory = "errors::ClientConnectionFailedTemporarilyError"
	ClientConnectionTimeoutError           ErrorCategory = "errors::ClientConnectionTimeoutError"
	ClientConnectionClosedError            ErrorCategory = "errors::ClientConnectionClosedError"
	InterfaceError                         ErrorCategory = "errors::InterfaceError"
	QueryArgumentError                     ErrorCategory = "errors::QueryArgumentError"
	MissingArgumentError                   ErrorCategory = "errors::MissingArgumentError"
	UnknownArgumentError                   ErrorCategory = "errors::UnknownArgumentError"
	InvalidArgumentError                   ErrorCategory = "errors::InvalidArgumentError"
	NoDataError                            ErrorCategory = "errors::NoDataError"
	InternalClientError                    ErrorCategory = "errors::InternalClientError"
)
