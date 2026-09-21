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

func TestParseFileYAML(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected map[string]any
	}{
		{
			name: "yml extension",
			file: "file1.yml",
			expected: map[string]any{
				"host":    "hexlet.io",
				"timeout": 50,
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
		{
			name: "yaml extension",
			file: "file2.yaml",
			expected: map[string]any{
				"timeout": 20,
				"verbose": true,
				"host":    "hexlet.io",
			},
		},
		{
			name:     "empty mapping",
			file:     "empty.yml",
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFile(filepath.Join("..", "..", "testdata", "fixture", tt.file))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseFormatCaseInsensitive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.YML")
	require.NoError(t, os.WriteFile(path, []byte("key: value\n"), 0o600))

	got, err := ParseFile(path)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"key": "value"}, got)
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

	invalidYAML := filepath.Join(t.TempDir(), "broken.yml")
	require.NoError(t, os.WriteFile(invalidYAML, []byte("key: [unclosed"), 0o600))
	_, err = ParseFile(invalidYAML)
	require.ErrorContains(t, err, "parse yaml")
}
