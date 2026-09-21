package code

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenDiff(t *testing.T) {
	got, err := GenDiff(
		filepath.Join("testdata", "file1.json"),
		filepath.Join("testdata", "file2.json"),
		"stylish",
	)
	require.NoError(t, err)

	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`
	assert.Equal(t, expected, got)
}

func TestGenDiffMissingFile(t *testing.T) {
	_, err := GenDiff(
		filepath.Join("testdata", "missing.json"),
		filepath.Join("testdata", "file2.json"),
		"stylish",
	)
	require.ErrorContains(t, err, "missing.json")
}
