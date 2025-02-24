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
	"encoding/json"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/internal/buff"
	"github.com/geldata/gel-go/internal/header"
)

func ignoreHeaders(r *buff.Reader) {
	n := int(r.PopUint16())

	for i := 0; i < n; i++ {
		r.PopUint16()
		r.PopBytes()
	}
}

func discardHeaders1pX(r *buff.Reader) {
	n := int(r.PopUint16())

	for i := 0; i < n; i++ {
		r.PopUint16()
		r.PopBytes()
	}
}

var discardHeaders2pX = discardHeaders1pX

func decodeHeaders1pX(
	r *buff.Reader,
	query string,
	warningHandler gelcfg.WarningHandler,
) (header.Header1pX, error) {
	n := int(r.PopUint16())

	headers := make(header.Header1pX, n)
	for i := 0; i < n; i++ {
		headers[r.PopString()] = r.PopString()
	}

	if data, ok := headers["warnings"]; ok {
		var warnings []*Warning
		err := json.Unmarshal([]byte(data), &warnings)
		if err != nil {
			return nil, err
		}

		errors := make([]error, len(warnings))
		for i := range warnings {
			errors[i] = warnings[i].Err(query)
		}

		err = warningHandler(errors)
		if err != nil {
			return nil, err
		}
	}

	return headers, nil
}

var decodeHeaders2pX = decodeHeaders1pX
