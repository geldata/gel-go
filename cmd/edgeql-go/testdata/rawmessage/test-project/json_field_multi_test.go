package main

import (
	"encoding/json"
	"testing"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestJsonFieldMulti(t *testing.T) {
	client, err := gel.CreateClient(gelcfg.Options{})
	require.NoError(t, err)

	actual, err := jsonFieldMulti(
		ctx,
		client,
	)
	require.NoError(t, err)

	expected := []jsonFieldMultiResult{
		{Json: json.RawMessage(`"std::str_lower"`)},
		{Json: json.RawMessage(`"std::str_repeat"`)},
	}
	assert.Equal(t, expected, actual)
}
