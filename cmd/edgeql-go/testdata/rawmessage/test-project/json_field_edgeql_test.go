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

func TestJSONField(t *testing.T) {
	client, err := gel.CreateClient(gelcfg.Options{})
	require.NoError(t, err)

	actual, err := jsonField(
		ctx,
		client,
		[]byte("bytes"),
		geltypes.NewOptionalBytes([]byte("optional bytes")),
		json.RawMessage(`"json"`),
		geltypes.NewOptionalBytes([]byte(`"optional json"`)),
	)
	require.NoError(t, err)

	expected := jsonFieldResult{
		Bytes:         []byte("bytes"),
		OptionalBytes: geltypes.NewOptionalBytes([]byte("optional bytes")),
		Json:          json.RawMessage(`"json"`),
		OptionalJson:  geltypes.NewOptionalBytes([]byte(`"optional json"`)),
	}
	assert.Equal(t, expected, actual)
}

func TestJSONFieldMissingOptional(t *testing.T) {
	client, err := gel.CreateClient(gelcfg.Options{})
	require.NoError(t, err)

	actual, err := jsonField(
		ctx,
		client,
		[]byte("bytes"),
		geltypes.OptionalBytes{},
		json.RawMessage(`"json"`),
		geltypes.OptionalBytes{},
	)
	require.NoError(t, err)

	expected := jsonFieldResult{
		Bytes:         []byte("bytes"),
		OptionalBytes: geltypes.OptionalBytes{},
		Json:          json.RawMessage(`"json"`),
		OptionalJson:  geltypes.OptionalBytes{},
	}
	assert.Equal(t, expected, actual)
}
