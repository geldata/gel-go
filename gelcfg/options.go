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

package gelcfg

import (
	"errors"
	"log"
	"time"

	types "github.com/geldata/gel-go/geltypes"
)

// WarningHandler takes an error slice that represent warnings and optionally
// returns an error. This can be used to log warnings, increment metrics,
// promote warnings to errors by returning them etc.
type WarningHandler = func([]error) error

// LogWarnings is a gelcfg.WarningHandler that logs warnings.
func LogWarnings(errors []error) error {
	for _, err := range errors {
		log.Println("Gel warning:", err.Error())
	}

	return nil
}

// WarningsAsErrors is a gelcfg.WarningHandler that returns warnings as
// errors.
func WarningsAsErrors(warnings []error) error {
	return errors.Join(warnings...)
}

// Options for connecting to a Gel server
type Options struct {
	// Host is a Gel server host address, given as either an IP address or
	// domain name. (Unix-domain socket paths are not supported)
	//
	// Host cannot be specified alongside the 'dsn' argument, or
	// CredentialsFile option. Host will override all other credentials
	// resolved from any environment variables, or project credentials with
	// their defaults.
	Host string

	// Port is a port number to connect to at the server host.
	//
	// Port cannot be specified alongside the 'dsn' argument, or
	// CredentialsFile option. Port will override all other credentials
	// resolved from any environment variables, or project credentials with
	// their defaults.
	Port int

	// Credentials is a JSON string containing connection credentials.
	//
	// Credentials cannot be specified alongside the 'dsn' argument, Host,
	// Port, or CredentialsFile.  Credentials will override all other
	// credentials not present in the credentials string with their defaults.
	Credentials []byte

	// CredentialsFile is a path to a file containing connection credentials.
	//
	// CredentialsFile cannot be specified alongside the 'dsn' argument, Host,
	// Port, or Credentials.  CredentialsFile will override all other
	// credentials not present in the credentials file with their defaults.
	CredentialsFile string

	// User is the name of the database role used for authentication.
	//
	// If not specified, the value is resolved from any compound
	// argument/option, then from GEL_USER, then any compound environment
	// variable, then project credentials.
	User string

	// Database is the name of the database to connect to.
	//
	// If not specified, the value is resolved from any compound
	// argument/option, then from EDGEDB_DATABASE, then any compound
	// environment variable, then project credentials.
	//
	// Deprecated: Database has been replaced by Branch
	Database string

	// Branch is the name of the branch to use.
	//
	// If not specified, the value is resolved from any compound
	// argument/option, then from GEL_BRANCH, then any compound environment
	// variable, then project credentials.
	Branch string

	// Password to be used for authentication, if the server requires one.
	//
	// If not specified, the value is resolved from any compound
	// argument/option, then from GEL_PASSWORD, then any compound
	// environment variable, then project credentials.
	// Note that the use of the environment variable is discouraged
	// as other users and applications may be able to read it
	// without needing specific privileges.
	Password types.OptionalStr

	// ConnectTimeout is used when establishing connections in the background.
	ConnectTimeout time.Duration

	// WaitUntilAvailable determines how long to wait
	// to reestablish a connection.
	WaitUntilAvailable time.Duration

	// Concurrency determines the maximum number of connections.
	// If Concurrency is zero, max(4, runtime.NumCPU()) will be used.
	// Has no effect for single connections.
	Concurrency uint

	// Parameters used to configure TLS connections to Gel server.
	TLSOptions TLSOptions

	// Read the TLS certificate from this file.
	// DEPRECATED, use TLSOptions.CAFile instead.
	TLSCAFile string

	// Specifies how strict TLS validation is.
	// DEPRECATED, use TLSOptions.SecurityMode instead.
	TLSSecurity string

	// ServerSettings is currently unused.
	ServerSettings map[string][]byte

	// SecretKey is used to connect to cloud instances.
	SecretKey string

	// WarningHandler is invoked when Gel returns warnings. Defaults to
	// gelcfg.LogWarnings.
	WarningHandler WarningHandler
}

// TLSOptions contains the parameters needed to configure TLS on Gel
// server connections.
type TLSOptions struct {
	// PEM-encoded CA certificate
	CA []byte
	// Path to a PEM-encoded CA certificate file
	CAFile string
	// Determines how strict we are with TLS checks
	SecurityMode TLSSecurityMode
	// Used to verify the hostname on the returned certificates
	ServerName string
}

// TLSSecurityMode specifies how strict TLS validation is.
type TLSSecurityMode string

const (
	// TLSModeDefault makes security mode inferred from other options
	TLSModeDefault TLSSecurityMode = "default"
	// TLSModeInsecure results in no certificate verification whatsoever
	TLSModeInsecure TLSSecurityMode = "insecure"
	// TLSModeNoHostVerification enables certificate verification
	// against CAs, but hostname matching is not performed.
	TLSModeNoHostVerification TLSSecurityMode = "no_host_verification"
	// TLSModeStrict enables full certificate and hostname verification.
	TLSModeStrict TLSSecurityMode = "strict"
)

// ModuleAlias is an alias name and module name pair.
type ModuleAlias struct {
	Alias  string
	Module string
}
