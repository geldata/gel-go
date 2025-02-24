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
	"testing"

	"github.com/geldata/gel-go/gelerr"
	gel "github.com/geldata/gel-go/internal/client"
	gelerrint "github.com/geldata/gel-go/internal/gelerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyntaxError(t *testing.T) {
	samples := []struct {
		protocolVersion uint16
		query           string
		err             string
	}{
		{
			1,
			"SELECT 1 2 3",
			`gel.EdgeQLSyntaxError: Unexpected '2'
query:1:10

SELECT 1 2 3
         ^ error`,
		},
		{
			2,
			"SELECT 1 2 3",
			`gel.EdgeQLSyntaxError: Unexpected '2'
query:1:10

SELECT 1 2 3
         ^ error`,
		},
		{
			1,
			"SELECT (foo (((1 2) 3)) 4)",
			`gel.EdgeQLSyntaxError: Unexpected token: <Token ICONST "2">
query:1:18

SELECT (foo (((1 2) 3)) 4)
                 ^ It appears that a ',' is missing in a tuple before '2'`,
		},
		{
			2,
			"SELECT (foo (((1 2) 3)) 4)",
			`gel.EdgeQLSyntaxError: Missing ','
query:1:17

SELECT (foo (((1 2) 3)) 4)
                ^ error`,
		},
		{
			1,
			`SELECT (Foo {
				foo
				bar
			} 2);`,
			`gel.EdgeQLSyntaxError: Unexpected token: <Token IDENT "bar">
query:3:5

    bar
    ^ It appears that a ',' is missing in a shape before 'bar'`,
		},
		{
			2,
			`SELECT (Foo {
				foo
				bar
			} 2);`,
			`gel.EdgeQLSyntaxError: Missing ','
query:2:1

    foo
^ error`,
		},
	}

	for _, s := range samples {
		t.Run(s.query, func(t *testing.T) {
			var result int64
			ctx := context.Background()
			pv, err := gel.ProtocolVersion(ctx, client.pool)
			assert.NoError(t, err)
			if pv.Major == s.protocolVersion {
				err := client.QuerySingle(ctx, s.query, &result)
				assert.EqualError(t, err, s.err)
			}
		})
	}
}

func TestNewErrorFromCodeAs(t *testing.T) {
	msg := "example error message"
	err := gelerrint.NewDuplicateCastDefinitionError(msg, nil)
	require.NotNil(t, err)

	assert.EqualError(t, err, "gel.DuplicateCastDefinitionError: "+msg)

	var edbErr gelerr.Error
	require.True(t, errors.As(err, &edbErr))

	assert.True(t, edbErr.Category(gelerr.DuplicateCastDefinitionError))
	assert.True(t, edbErr.Category(gelerr.DuplicateDefinitionError))
	assert.True(t, edbErr.Category(gelerr.SchemaDefinitionError))
	assert.True(t, edbErr.Category(gelerr.QueryError))
}
