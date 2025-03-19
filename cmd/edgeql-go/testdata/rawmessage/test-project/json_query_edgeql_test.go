package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

var (
	ctx = context.Background()
)

func TestJSONQuery(t *testing.T) {
	client, err := gel.CreateClient(gelcfg.Options{})
	require.NoError(t, err)

	actual, err := jsonQuery(
		ctx,
		client,
		[]byte("bytes"),
		geltypes.NewOptionalBytes([]byte("optional bytes")),
		[]byte(`"json"`),
		geltypes.NewOptionalBytes([]byte(`"optional json"`)),
	)
	require.NoError(t, err)

	expected := jsonQueryResult{
		Bytes:         []byte("bytes"),
		OptionalBytes: geltypes.NewOptionalBytes([]byte("optional bytes")),
		Json:          json.RawMessage(`"json"`),
		OptionalJson:  geltypes.NewOptionalBytes([]byte(`"optional json"`)),
	}
	assert.Equal(t, expected, actual)
}

func TestJSONQueryMissingOptional(t *testing.T) {
	client, err := gel.CreateClient(gelcfg.Options{})
	require.NoError(t, err)

	actual, err := jsonQuery(
		ctx,
		client,
		[]byte("bytes"),
		geltypes.OptionalBytes{},
		[]byte(`"json"`),
		geltypes.OptionalBytes{},
	)
	require.NoError(t, err)

	expected := jsonQueryResult{
		Bytes:         []byte("bytes"),
		OptionalBytes: geltypes.OptionalBytes{},
		Json:          json.RawMessage(`"json"`),
		OptionalJson:  geltypes.OptionalBytes{},
	}
	assert.Equal(t, expected, actual)
}
