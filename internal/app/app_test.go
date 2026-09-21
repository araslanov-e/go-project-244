package app

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	var out bytes.Buffer

	err := Run(context.Background(), []string{
		"gendiff",
		filepath.Join("..", "..", "testdata", "file1.json"),
		filepath.Join("..", "..", "testdata", "file2.json"),
	}, &out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "format: stylish")
}

func TestRunWrongArgsCount(t *testing.T) {
	var out bytes.Buffer

	err := Run(context.Background(), []string{"gendiff", "only-one.json"}, &out)
	require.ErrorContains(t, err, "expected 2 arguments")
	assert.Empty(t, out.String())
}
