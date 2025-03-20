package main

import (
	"testing"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestJSONScalarMulti(t *testing.T) {
	client, err := gel.CreateClient(gelcfg.Options{})
	require.NoError(t, err)

	actual, err := jsonScalarMulti(
		ctx,
		client,
		[]byte(`"json"`),
		[]byte(`"json"`),
	)
	require.NoError(t, err)

	expected := [][]byte{[]byte(`"json"`), []byte(`"json"`)}
	assert.Equal(t, expected, actual)
}
