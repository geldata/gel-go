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
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/geldata/gel-go/internal/gelerr"
)

// Warning is used to decode warnings in the protocol.
type Warning struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
	Line    *int   `json:"line,omitempty"`
	Start   *int   `json:"start,omitempty"`
}

// Err returns a formatted error for a query
func (w *Warning) Err(query, filename string) error {
	if w.Line == nil || w.Start == nil {
		return gelerr.ErrorFromCode(w.Code, w.Message)
	}

	lineNo := *w.Line - 1
	byteNo := *w.Start
	lines := strings.Split(query, "\n")
	if lineNo >= len(lines) {
		return gelerr.ErrorFromCode(w.Code, w.Message)
	}

	// replace tabs with a single space
	// because we don't know how they will be printed.
	line := strings.ReplaceAll(lines[lineNo], "\t", " ")

	for i := 0; i < lineNo; i++ {
		byteNo -= 1 + len(lines[i])
	}

	if byteNo >= len(line) {
		byteNo = 0
	}

	hint := w.Hint
	if hint == "" {
		hint = "error"
	}

	runeCount := utf8.RuneCountInString(line[:byteNo])
	padding := strings.Repeat(" ", runeCount)
	msg := w.Message + fmt.Sprintf(
		"\n%s:%v:%v\n\n%v\n%v^ %v",
		filename,
		1+lineNo,
		1+runeCount,
		line,
		padding,
		hint,
	)

	return gelerr.ErrorFromCode(w.Code, msg)
}
