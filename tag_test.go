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
	assert.EqualError(t, err, "gel.InvalidArgumentError: reserved tag: edgedb/*")

	_, err = client.WithQueryTag(strings.Repeat("a", 129))
	assert.EqualError(t, err, "gel.InvalidArgumentError: tag too long (> 128 characters)")

	_, err = client.WithQueryTag(strings.Repeat("a", 128))
	assert.NoError(t, err)
}

func assertQueryTag(
	t *testing.T,
	client *Client,
	tag geltypes.OptionalStr,
) {
	// sys::QueryStats entries are not unique per tag.  Make sure that each
	// assertion uses its own query so that its sys::QueryStats entry correctly
	// reflects the tag value.
	name := randomName()
	queryWithName := fmt.Sprintf("SELECT (%s := <int64>1)", name)

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

	err = client.Execute(ctx, queryWithName)
	require.NoError(t, err)

	err = client.Query(ctx, queryStats, &results, args...)
	require.NoError(t, err)
	require.Equal(t, 1, len(results))
	assert.Equal(t, tag, results[0].Tag)
}

func TestWithQueryTag(t *testing.T) {
	assertQueryTag(t, client, geltypes.OptionalStr{})

	tag := randomName()
	taggedClient, err := client.WithQueryTag(tag)
	require.NoError(t, err)
	assertQueryTag(t, taggedClient, geltypes.NewOptionalStr(tag))

	assertQueryTag(t, client, geltypes.OptionalStr{})
}
