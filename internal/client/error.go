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
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"syscall"

	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/internal/buff"
	gelerrint "github.com/geldata/gel-go/internal/gelerr"
)

var (
	errNoTOMLFound = errors.New("no gel.toml found")

	// ErrZeroResults is returned when a query didn't return any data.
	ErrZeroResults error = gelerrint.NewNoDataError("zero results", nil)
)

// FirstError returns the first non nil error or nil.
func FirstError(a, b error) error {
	if a != nil {
		return a
	}

	return b
}

const (
	hint          = 0x0001
	positionStart = 0xfff1
	lineStart     = 0xfff3
)

func positionFromHeaders(headers map[uint16]string) (*int, *int, error) {
	lineNoRaw, ok := headers[lineStart]
	if !ok {
		return nil, nil, nil
	}

	byteNoRaw, ok := headers[positionStart]
	if !ok {
		return nil, nil, nil
	}

	lineNo, err := strconv.Atoi(lineNoRaw)
	if err != nil {
		return nil, nil, gelerrint.NewBinaryProtocolError(
			"", fmt.Errorf("decode lineNo: %q: %w", lineNoRaw, err))
	}
	byteNo, err := strconv.Atoi(byteNoRaw)
	if err != nil {
		return nil, nil, gelerrint.NewBinaryProtocolError(
			"", fmt.Errorf("decode byteNo: %q: %w", byteNoRaw, err))
	}

	return &lineNo, &byteNo, nil
}

// decodeErrorResponseMsg decodes an error response
// https://docs.geldata.com/reference/reference/protocol/messages#errorresponse
func decodeErrorResponseMsg(r *buff.Reader, query, filename string) error {
	r.Discard(1) // severity
	w := Warning{
		Code:    r.PopUint32(),
		Message: r.PopString(),
	}

	n := int(r.PopUint16())
	headers := make(map[uint16]string, n)
	for i := 0; i < n; i++ {
		headers[r.PopUint16()] = r.PopString()
	}

	var err error
	w.Line, w.Start, err = positionFromHeaders(headers)
	if err != nil {
		return errors.Join(w.Err(query, filename), err)
	}

	w.Hint = headers[hint]
	return w.Err(query, filename)
}

type wrappedManyError struct {
	msg  string
	errs []error
}

func (e *wrappedManyError) Error() string {
	return e.msg
}

func (e *wrappedManyError) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (e *wrappedManyError) As(target interface{}) bool {
	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

func wrapAll(errs ...error) error {
	err := &wrappedManyError{}
	for _, e := range errs {
		if e != nil {
			err.errs = append(err.errs, e)
		}
	}

	if len(err.errs) == 0 {
		return nil
	}

	if len(err.errs) == 1 {
		return err.errs[0]
	}

	err.msg = err.errs[0].Error()
	for _, e := range err.errs[1:] {
		err.msg += "; " + e.Error()
	}

	return err
}

func isClientConnectionError(err error) bool {
	var edbErr gelerr.Error
	return errors.As(err, &edbErr) &&
		edbErr.Category(gelerr.ClientConnectionError)
}

func wrapNetError(err error) error {
	var errEDB gelerr.Error
	var errNetOp *net.OpError
	var errDSN *net.DNSError

	switch {
	case err == nil:
		return err
	case errors.As(err, &errEDB):
		return err

	case errors.Is(err, context.Canceled):
		fallthrough
	case errors.Is(err, context.DeadlineExceeded):
		fallthrough
	case errors.As(err, &errNetOp) && errNetOp.Timeout():
		return gelerrint.NewClientConnectionTimeoutError("", err)

	case errors.Is(err, io.EOF):
		fallthrough
	case errors.Is(err, syscall.ECONNREFUSED):
		fallthrough
	case errors.Is(err, syscall.ECONNABORTED):
		fallthrough
	case errors.Is(err, syscall.ECONNRESET):
		fallthrough
	case errors.Is(err, syscall.EADDRINUSE):
		fallthrough
	case errors.As(err, &errDSN):
		fallthrough
	case errors.Is(err, syscall.ENOENT):
		return gelerrint.NewClientConnectionFailedTemporarilyError("", err)

	case errors.Is(err, net.ErrClosed):
		return gelerrint.NewClientConnectionClosedError("", err)

	default:
		return gelerrint.NewClientConnectionFailedError("", err)
	}
}

func invalidTLSSecurity(val string) error {
	return fmt.Errorf(
		"invalid TLSSecurity value: expected one of %v, got: %q",
		englishList(
			[]string{"insecure", "no_host_verification", "strict"},
			"or"),
		val,
	)
}
