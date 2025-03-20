package main

import (
	"testing"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestJSONScalar(t *testing.T) {
	client, err := gel.CreateClient(gelcfg.Options{})
	require.NoError(t, err)

	actual, err := jsonScalar(ctx, client, []byte(`"json"`))
	require.NoError(t, err)

	expected := []byte(`"json"`)
	assert.Equal(t, expected, actual)
}
