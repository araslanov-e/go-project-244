package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	var out bytes.Buffer

	err := Run(context.Background(), []string{
		"gendiff",
		filepath.Join("..", "..", "testdata", "fixture", "nested1.json"),
		filepath.Join("..", "..", "testdata", "fixture", "nested2.json"),
	}, &out)
	require.NoError(t, err)

	expected, err := os.ReadFile(filepath.Join("..", "..", "testdata", "fixture", "nested_result.txt"))
	require.NoError(t, err)
	assert.Equal(t, string(expected)+"\n", out.String())
}

func TestRunUnsupportedFormat(t *testing.T) {
	var out bytes.Buffer

	err := Run(context.Background(), []string{
		"gendiff", "--format", "unknown",
		filepath.Join("..", "..", "testdata", "fixture", "nested1.json"),
		filepath.Join("..", "..", "testdata", "fixture", "nested2.json"),
	}, &out)
	require.ErrorContains(t, err, "unsupported output format")
	assert.Empty(t, out.String())
}

func TestRunWrongArgsCount(t *testing.T) {
	var out bytes.Buffer

	err := Run(context.Background(), []string{"gendiff", "only-one.json"}, &out)
	require.ErrorContains(t, err, "expected 2 arguments")
	assert.Empty(t, out.String())
}
