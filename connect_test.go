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

package gel

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	types "github.com/geldata/gel-go/geltypes"
	gel "github.com/geldata/gel-go/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	ctx := context.Background()
	p, err := CreateClient(gelcfg.Options{
		Host:       opts.Host,
		Port:       opts.Port,
		User:       "user_with_password",
		Password:   types.NewOptionalStr("secret"),
		Database:   opts.Database, //nolint:staticcheck // SA1019
		TLSOptions: opts.TLSOptions,
	})
	require.NoError(t, err)

	var result string
	err = p.QuerySingle(ctx, "SELECT 'It worked!';", &result)
	assert.NoError(t, err)
	assert.Equal(t, "It worked!", result)

	clientCopy := p.WithTxOptions(gelcfg.NewTxOptions())

	err = p.Close()
	assert.NoError(t, err)

	// A connection should not be closeable more than once.
	err = p.Close()
	msg := "gel.InterfaceError: client closed"
	assert.EqualError(t, err, msg)

	// Copied connections should not be closeable after another copy is closed.
	err = clientCopy.Close()
	assert.EqualError(t, err, msg)
}

func TestCloudClientHandshakeMessage(t *testing.T) {
	params := map[string]string{
		"database":   "mydb",
		"secret_key": "mysecret",
		"user":       "myuser",
	}
	got, err := gel.ClientHandshakeMessage(params, []byte{})
	assert.NoError(t, err)
	majorUpper, majorLower := convertUint16ToUint8(
		gel.ProtocolVersionMax.Major,
	)
	minorUpper, minorLower := convertUint16ToUint8(
		gel.ProtocolVersionMax.Minor,
	)

	want := []byte{
		uint8(gel.ClientHandshake), // mtype (uint8)
		0, 0, 0, 76,                // message_length (uint32)
		majorLower, majorUpper, // major_ver (uint16)
		minorLower, minorUpper, // minor_ver (uint16)
		0, 3, // num_params (uint16)

		// Parameter 1: database
		0, 0, 0, 8, // param1 name length (uint32)
		'd', 'a', 't', 'a', 'b', 'a', 's', 'e', // param1 name ("database")
		0, 0, 0, 4, // param1 value length (uint32)
		'm', 'y', 'd', 'b', // param1 value ("mydb")

		// Parameter 2: secret_key
		0, 0, 0, 10, // param3 name length (uint32)
		's', 'e', 'c', 'r', 'e', 't', '_', 'k', 'e', 'y', // p3 ("secret_key")
		0, 0, 0, 8, // param3 value length (uint32)
		'm', 'y', 's', 'e', 'c', 'r', 'e', 't', // param3 value ("mysecret")

		// Parameter 3: user
		0, 0, 0, 4, // param2 name length (uint32)
		'u', 's', 'e', 'r', // param2 name ("user")
		0, 0, 0, 6, // param2 value length (uint32)
		'm', 'y', 'u', 's', 'e', 'r', // param2 value ("myuser")

		0, 0, // num_extensions (uint16)
	}

	assert.EqualValues(t, got.Unwrap(), want)
}

func convertUint16ToUint8(value uint16) (uint8, uint8) {
	lowerByte := uint8(value & 0xFF)
	upperByte := uint8((value >> 8) & 0xFF)

	return lowerByte, upperByte
}

func TestConnectTimeout(t *testing.T) {
	ctx := context.Background()
	p, err := CreateClient(gelcfg.Options{
		Host:               opts.Host,
		Port:               opts.Port,
		User:               opts.User,
		Password:           opts.Password,
		Database:           opts.Database, //nolint:staticcheck // SA1019
		ConnectTimeout:     2 * time.Nanosecond,
		WaitUntilAvailable: 1 * time.Nanosecond,
	})

	if p != nil {
		err = p.EnsureConnected(ctx)
		_ = p.Close()
	}

	require.NotNil(t, err, "connection didn't timeout")

	var edbErr gelerr.Error

	require.True(t, errors.As(err, &edbErr), "wrong error: %v", err)
	assert.True(
		t,
		edbErr.Category(gelerr.ClientConnectionTimeoutError),
		"wrong error: %v",
		err,
	)
}

func TestConnectRefused(t *testing.T) {
	ctx := context.Background()
	p, err := CreateClient(gelcfg.Options{
		Host:               "localhost",
		Port:               23456,
		WaitUntilAvailable: 1 * time.Nanosecond,
	})

	if p != nil {
		err = p.EnsureConnected(ctx)
		_ = p.Close()
	}

	require.NotNil(t, err, "connection wasn't refused")

	msg := "wrong error: " + err.Error()
	var edbErr gelerr.Error
	require.True(t, errors.As(err, &edbErr), msg)
	assert.True(
		t,
		edbErr.Category(gelerr.ClientConnectionFailedError),
		msg,
	)
}

func TestConnectInvalidName(t *testing.T) {
	ctx := context.Background()
	p, err := CreateClient(gelcfg.Options{
		Host:               "invalid.example.org",
		Port:               23456,
		WaitUntilAvailable: 1 * time.Nanosecond,
	})

	if p != nil {
		err = p.EnsureConnected(ctx)
		_ = p.Close()
	}

	require.NotNil(t, err, "name was resolved")

	var edbErr gelerr.Error
	require.True(t, errors.As(err, &edbErr), "wrong error: %v", err)
	assert.True(
		t,
		edbErr.Category(gelerr.ClientConnectionFailedTemporarilyError),
		"wrong error: %v",
		err,
	)
	// Match lookup error agnostic to OS. Examples:
	// dial tcp: lookup invalid.example.org: no such host
	// dial tcp: lookup invalid.example.org on 127.0.0.1:53: no such host
	assert.Contains(t, err.Error(),
		"gel.ClientConnectionFailedTemporarilyError: "+
			"dial tcp: lookup invalid.example.org")
	assert.Contains(t, err.Error(), "no such host")

	var errNotFound *net.DNSError
	assert.True(t, errors.As(err, &errNotFound))
}

func TestConnectRefusedUnixSocket(t *testing.T) {
	ctx := context.Background()
	p, err := CreateClient(gelcfg.Options{
		Host:               "/tmp/non-existent",
		WaitUntilAvailable: 1 * time.Nanosecond,
	})

	if p != nil {
		err = p.EnsureConnected(ctx)
		_ = p.Close()
	}

	require.NotNil(t, err, "connection wasn't refused")

	var edbErr gelerr.Error
	require.True(t, errors.As(err, &edbErr), "wrong error: %v", err)
	assert.True(
		t,
		edbErr.Category(gelerr.ConfigurationError),
		"wrong error: %v",
		err,
	)
}
