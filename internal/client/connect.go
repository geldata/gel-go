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

package gel

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/geldata/gel-go/internal"
	"github.com/geldata/gel-go/internal/buff"
	"github.com/geldata/gel-go/internal/gelerr"
	"github.com/xdg/scram"
)

// ClientHandshakeMessage writes a client handshake message.
func ClientHandshakeMessage(
	params map[string]string,
	alocatedMemory []byte,
) (*buff.Writer, error) {
	if len(params) > math.MaxUint16 {
		return nil, errors.New("too many connection parameters")
	}

	numParams := uint16(len(params))
	paramKeys := make([]string, 0, len(params))
	for k := range params {
		paramKeys = append(paramKeys, k)
	}
	slices.Sort(paramKeys)
	w := buff.NewWriter(alocatedMemory)
	w.BeginMessage(uint8(ClientHandshake))
	w.PushUint16(ProtocolVersionMax.Major)
	w.PushUint16(ProtocolVersionMax.Minor)
	w.PushUint16(numParams)
	for _, pk := range paramKeys {
		w.PushString(pk)
		w.PushString(params[pk])
	}
	w.PushUint16(0) // no extensions
	w.EndMessage()

	return w, nil
}

func (c *protocolConnection) connect(r *buff.Reader, cfg *connConfig) error {
	var err error

	params := map[string]string{
		"branch":     cfg.branch,
		"database":   cfg.database,
		"user":       cfg.user,
		"secret_key": cfg.secretKey,
	}

	w, err := ClientHandshakeMessage(params, c.writeMemory[:0])
	if err != nil {
		return err
	}

	c.protocolVersion = ProtocolVersionMax

	if err = c.soc.WriteAll(w.Unwrap()); err != nil {
		return err
	}

	done := buff.NewSignal()

	for r.Next(done.Chan) {
		switch Message(r.MsgType) {
		case ServerHandshake:
			protocolVersion := internal.ProtocolVersion{
				Major: r.PopUint16(),
				Minor: r.PopUint16(),
			}

			// The client _MUST_ close the connection
			// if the protocol version can't be supported.
			// https://docs.geldata.com/reference/reference/protocol#connection-phase
			if protocolVersion.LT(protocolVersionMin) ||
				protocolVersion.GT(ProtocolVersionMax) {
				_ = c.soc.Close()
				msg := fmt.Sprintf(
					"unsupported protocol version: %v.%v",
					protocolVersion.Major,
					protocolVersion.Minor,
				)
				return gelerr.NewUnsupportedProtocolVersionError(msg, nil)
			}

			c.protocolVersion = protocolVersion

			n := r.PopUint16()
			for i := uint16(0); i < n; i++ {
				r.PopBytes() // extension name
				ignoreHeaders(r)
			}
		case ServerKeyData:
			r.DiscardMessage() // key data
		case ReadyForCommand:
			ignoreHeaders(r)
			r.Discard(1) // transaction state
			done.Signal()
		case Authentication:
			if r.PopUint32() == 0 { // auth status
				continue
			}

			// skip supported SASL methods
			n := int(r.PopUint32()) // method count
			for i := 0; i < n; i++ {
				r.PopBytes()
			}

			if e := c.authenticate(r, cfg); e != nil {
				return e
			}

			done.Signal()
		case StateDataDescription:
			if e := c.decodeStateDataDescription(r); e != nil {
				err = wrapAll(err, e)
			}
		case ErrorResponse:
			err = wrapAll(err, decodeErrorResponseMsg(r, "", ""))
			done.Signal()
		default:
			if e := c.fallThrough(r); e != nil {
				// the connection will not be usable after this x_x
				return e
			}
		}
	}

	return wrapAll(err, r.Err)
}

func (c *protocolConnection) authenticate(
	r *buff.Reader,
	cfg *connConfig,
) error {
	client, err := scram.SHA256.NewClient(cfg.user, cfg.password, "")
	if err != nil {
		return gelerr.NewAuthenticationError(err.Error(), nil)
	}

	conv := client.NewConversation()
	scramMsg, err := conv.Step("")
	if err != nil {
		return gelerr.NewAuthenticationError(err.Error(), nil)
	}

	w := buff.NewWriter(c.writeMemory[:0])
	w.BeginMessage(uint8(AuthenticationSASLInitialResponse))
	w.PushString("SCRAM-SHA-256")
	w.PushString(scramMsg)
	w.EndMessage()

	if e := c.soc.WriteAll(w.Unwrap()); e != nil {
		return e
	}

	done := buff.NewSignal()

	for r.Next(done.Chan) {
		switch Message(r.MsgType) {
		case Authentication:
			authStatus := r.PopUint32()
			if authStatus != 0xb {
				// the connection will not be usable after this x_x
				return gelerr.NewAuthenticationError(fmt.Sprintf(
					"unexpected authentication status: 0x%x", authStatus,
				), nil)
			}

			scramRcv := r.PopString()
			scramMsg, err = conv.Step(scramRcv)
			if err != nil {
				// the connection will not be usable after this x_x
				return gelerr.NewAuthenticationError(err.Error(), nil)
			}

			done.Signal()
		case ErrorResponse:
			err = decodeErrorResponseMsg(r, "", "")
		default:
			if e := c.fallThrough(r); e != nil {
				// the connection will not be usable after this x_x
				return e
			}
		}
	}

	if err != nil || r.Err != nil {
		return wrapAll(err, r.Err)
	}

	w = buff.NewWriter(c.writeMemory[:0])
	w.BeginMessage(uint8(AuthenticationSASLResponse))
	w.PushString(scramMsg)
	w.EndMessage()

	if e := c.soc.WriteAll(w.Unwrap()); e != nil {
		return e
	}

	done = buff.NewSignal()

	for r.Next(done.Chan) {
		switch Message(r.MsgType) {
		case Authentication:
			authStatus := r.PopUint32()
			switch authStatus {
			case 0:
			case 0xc:
				scramRcv := r.PopString()
				_, e := conv.Step(scramRcv)
				if e != nil {
					// the connection will not be usable after this x_x
					return gelerr.NewAuthenticationError(e.Error(), nil)
				}
			default:
				// the connection will not be usable after this x_x
				return gelerr.NewAuthenticationError(fmt.Sprintf(
					"unexpected authentication status: 0x%x", authStatus,
				), nil)
			}
		case ServerKeyData:
			r.DiscardMessage() // key data
		case ReadyForCommand:
			ignoreHeaders(r)
			r.Discard(1) // transaction state
			done.Signal()
		case StateDataDescription:
			if e := c.decodeStateDataDescription(r); e != nil {
				err = wrapAll(err, e)
			}
		case ErrorResponse:
			err = wrapAll(decodeErrorResponseMsg(r, "", ""))
		default:
			if e := c.fallThrough(r); e != nil {
				// the connection will not be usable after this x_x
				return e
			}
		}
	}

	return wrapAll(err, r.Err)
}

func (c *protocolConnection) terminate() error {
	w := buff.NewWriter(c.writeMemory[:0])
	w.BeginMessage(uint8(Terminate))
	w.EndMessage()
	return c.soc.WriteAll(w.Unwrap())
}
