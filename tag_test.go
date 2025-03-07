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
	"fmt"
	"strings"
	"testing"

	"github.com/geldata/gel-go/geltypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type QueryStats struct {
	Query geltypes.OptionalStr `gel:"query"`
	Tag   geltypes.OptionalStr `gel:"tag"`
}

func TestWithQueryTagInvalid(t *testing.T) {
	_, err := client.WithQueryTag("gel/tag")
	assert.EqualError(t, err, "gel.InvalidArgumentError: reserved tag: gel/*")

	_, err = client.WithQueryTag("edgedb/tag")
	assert.EqualError(
		t,
		err,
		"gel.InvalidArgumentError: reserved tag: edgedb/*",
	)

	_, err = client.WithQueryTag(strings.Repeat("a", 129))
	assert.EqualError(
		t,
		err,
		"gel.InvalidArgumentError: tag too long (> 128 characters)",
	)

	_, err = client.WithQueryTag(strings.Repeat("a", 128))
	assert.NoError(t, err)
}

type assertQueryTagQueryHandler func(
	context.Context,
	*testing.T,
	*Client,
	string,
)

func assertQueryTag(
	t *testing.T,
	client *Client,
	tag geltypes.OptionalStr,
	cb assertQueryTagQueryHandler,
) {
	// sys::QueryStats entries are not unique per tag.  Make sure that each
	// assertion uses its own query so that its sys::QueryStats entry correctly
	// reflects the tag value.
	name := randomName()
	queryWithName := fmt.Sprintf("SELECT (%s := <int64>1).%s", name, name)

	args := []interface{}{name}
	queryStats := `
		SELECT sys::QueryStats { query, tag }
		FILTER contains(.query, <str>$0)
	`

	tagStr, isTagged := tag.Get()
	if isTagged {
		queryStats += " AND .tag = <str>$1"
		args = append(args, tagStr)
	}

	var results []QueryStats
	ctx := context.Background()
	err := client.Query(ctx, queryStats, &results, args...)
	require.NoError(t, err)
	require.Equal(t, 0, len(results))

	cb(ctx, t, client, queryWithName)

	err = client.Query(ctx, queryStats, &results, args...)
	require.NoError(t, err)
	require.Equal(t, 1, len(results))
	assert.Equal(t, tag, results[0].Tag)
}

var TestWithQueryTagTestCases = []struct {
	name    string
	handler assertQueryTagQueryHandler
}{
	{
		"Execute",
		func(ctx context.Context, t *testing.T, c *Client, query string) {
			err := c.Execute(ctx, query)
			require.NoError(t, err)
		},
	},
	{
		"Query",
		func(ctx context.Context, t *testing.T, c *Client, query string) {
			var result []int64
			err := c.Query(ctx, query, &result)
			require.NoError(t, err)
			require.Equal(t, []int64{1}, result)
		},
	},
	{
		"QueryJSON",
		func(ctx context.Context, t *testing.T, c *Client, query string) {
			var result []byte
			err := c.QueryJSON(ctx, query, &result)
			require.NoError(t, err)
			require.Equal(t, "[1]", string(result))
		},
	},
	{
		"QuerySingle",
		func(ctx context.Context, t *testing.T, c *Client, query string) {
			var result int64
			err := c.QuerySingle(ctx, query, &result)
			require.NoError(t, err)
			require.Equal(t, int64(1), result)
		},
	},
	{
		"QuerySingleJSON",
		func(ctx context.Context, t *testing.T, c *Client, query string) {
			var result []byte
			err := c.QuerySingleJSON(ctx, query, &result)
			require.NoError(t, err)
			require.Equal(t, "1", string(result))
		},
	},
}

func TestWithQueryTag(t *testing.T) {
	for _, testCase := range TestWithQueryTagTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			expected := geltypes.OptionalStr{}
			assertQueryTag(t, client, expected, testCase.handler)

			tag := randomName()
			taggedClient, err := client.WithQueryTag(tag)
			require.NoError(t, err)

			expected = geltypes.NewOptionalStr(tag)
			assertQueryTag(t, taggedClient, expected, testCase.handler)

			assertQueryTag(t, client, geltypes.OptionalStr{}, testCase.handler)

			expected = geltypes.OptionalStr{}
			untaggedClient := taggedClient.WithoutQueryTag()
			assertQueryTag(t, untaggedClient, expected, testCase.handler)
		})
	}
}
