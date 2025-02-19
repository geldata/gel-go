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
	"errors"
	"testing"

	gelerrint "github.com/edgedb/edgedb-go/internal/gelerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapAllAs(t *testing.T) {
	err1 := gelerrint.NewBinaryProtocolError("bad bits!", nil)
	err2 := gelerrint.NewInvalidValueError("guess again...", nil)
	err := wrapAll(err1, err2)

	require.NotNil(t, err)
	assert.Equal(
		t,
		"gel.BinaryProtocolError: bad bits!; "+
			"gel.InvalidValueError: guess again...",
		err.Error(),
	)

	var bin *gelerrint.BinaryProtocolError
	require.True(t, errors.As(err, &bin), "errors.As failed")
	assert.Equal(t, "gel.BinaryProtocolError: bad bits!", bin.Error())

	var val *gelerrint.InvalidValueError
	require.True(t, errors.As(err, &val))
	assert.Equal(t, "gel.InvalidValueError: guess again...", val.Error())
}

func TestWrapAllIs(t *testing.T) {
	errA := errors.New("error A")
	errB := errors.New("error B")
	err := wrapAll(errA, errB)

	require.NotNil(t, err)
	assert.Equal(t, "error A; error B", err.Error())
	assert.True(t, errors.Is(err, errA))
	assert.True(t, errors.Is(err, errB))
}
