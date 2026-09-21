package parsers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFileJSON(t *testing.T) {
	got, err := ParseFile(filepath.Join("..", "..", "testdata", "fixture", "file1.json"))
	require.NoError(t, err)

	expected := map[string]any{
		"host":    "hexlet.io",
		"timeout": float64(50),
		"proxy":   "123.234.53.22",
		"follow":  false,
	}
	assert.Equal(t, expected, got)
}

func TestParseFileAbsolutePath(t *testing.T) {
	absPath, err := filepath.Abs(filepath.Join("..", "..", "testdata", "fixture", "file2.json"))
	require.NoError(t, err)

	got, err := ParseFile(absPath)
	require.NoError(t, err)

	expected := map[string]any{
		"timeout": float64(20),
		"verbose": true,
		"host":    "hexlet.io",
	}
	assert.Equal(t, expected, got)
}

func TestParseFileErrors(t *testing.T) {
	_, err := ParseFile(filepath.Join("..", "..", "testdata", "fixture", "missing.json"))
	require.Error(t, err)

	unsupported := filepath.Join(t.TempDir(), "config.txt")
	require.NoError(t, os.WriteFile(unsupported, []byte("key=value"), 0o600))
	_, err = ParseFile(unsupported)
	require.ErrorContains(t, err, "unsupported file format")

	invalid := filepath.Join(t.TempDir(), "broken.json")
	require.NoError(t, os.WriteFile(invalid, []byte("{not json"), 0o600))
	_, err = ParseFile(invalid)
	require.ErrorContains(t, err, "parse json")
}
