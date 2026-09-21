package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixturePath(name string) string {
	return filepath.Join("testdata", "fixture", name)
}

func readFixture(t *testing.T, name string) string {
	t.Helper()

	data, err := os.ReadFile(fixturePath(name))
	require.NoError(t, err)

	return string(data)
}

func TestGenDiffFlatJSON(t *testing.T) {
	tests := []struct {
		name     string
		file1    string
		file2    string
		expected string
	}{
		{
			name:     "changed, added and removed keys",
			file1:    "file1.json",
			file2:    "file2.json",
			expected: "flat_result.txt",
		},
		{
			name:     "identical files",
			file1:    "file1.json",
			file2:    "file1.json",
			expected: "flat_same_result.txt",
		},
		{
			name:     "all keys added",
			file1:    "empty.json",
			file2:    "file2.json",
			expected: "flat_added_result.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tt.file1), fixturePath(tt.file2), "stylish")
			require.NoError(t, err)
			assert.Equal(t, readFixture(t, tt.expected), got)
		})
	}
}

func TestGenDiffFlatYAML(t *testing.T) {
	tests := []struct {
		name     string
		file1    string
		file2    string
		expected string
	}{
		{
			name:     "changed, added and removed keys",
			file1:    "file1.yml",
			file2:    "file2.yaml",
			expected: "flat_result.txt",
		},
		{
			name:     "identical files",
			file1:    "file1.yml",
			file2:    "file1.yml",
			expected: "flat_same_result.txt",
		},
		{
			name:     "all keys added",
			file1:    "empty.yml",
			file2:    "file2.yaml",
			expected: "flat_added_result.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tt.file1), fixturePath(tt.file2), "stylish")
			require.NoError(t, err)
			assert.Equal(t, readFixture(t, tt.expected), got)
		})
	}
}

func TestGenDiffMissingFile(t *testing.T) {
	_, err := GenDiff(fixturePath("missing.json"), fixturePath("file2.json"), "stylish")
	require.ErrorContains(t, err, "missing.json")

	_, err = GenDiff(fixturePath("file1.json"), fixturePath("missing.json"), "stylish")
	require.ErrorContains(t, err, "missing.json")
}
