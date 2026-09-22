package code

import (
	"encoding/json"
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

func TestGenDiffNested(t *testing.T) {
	tests := []struct {
		name   string
		file1  string
		file2  string
		format string
	}{
		{name: "json", file1: "nested1.json", file2: "nested2.json", format: "stylish"},
		{name: "yaml", file1: "nested1.yml", file2: "nested2.yaml", format: "stylish"},
		{name: "json and yaml", file1: "nested1.json", file2: "nested2.yaml", format: "stylish"},
		{name: "default format", file1: "nested1.json", file2: "nested2.json", format: ""},
	}

	expected := readFixture(t, "nested_result.txt")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tt.file1), fixturePath(tt.file2), tt.format)
			require.NoError(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestGenDiffPlain(t *testing.T) {
	tests := []struct {
		name  string
		file1 string
		file2 string
	}{
		{name: "json", file1: "nested1.json", file2: "nested2.json"},
		{name: "yaml", file1: "nested1.yml", file2: "nested2.yaml"},
	}

	expected := readFixture(t, "plain_result.txt")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tt.file1), fixturePath(tt.file2), "plain")
			require.NoError(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestGenDiffJSON(t *testing.T) {
	tests := []struct {
		name  string
		file1 string
		file2 string
	}{
		{name: "json", file1: "nested1.json", file2: "nested2.json"},
		{name: "yaml", file1: "nested1.yml", file2: "nested2.yaml"},
	}

	expected := readFixture(t, "json_result.json")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenDiff(fixturePath(tt.file1), fixturePath(tt.file2), "json")
			require.NoError(t, err)
			assert.JSONEq(t, expected, got)

			var parsed map[string]any
			require.NoError(t, json.Unmarshal([]byte(got), &parsed))
		})
	}
}

func TestGenDiffUnsupportedFormat(t *testing.T) {
	_, err := GenDiff(fixturePath("nested1.json"), fixturePath("nested2.json"), "unknown")
	require.ErrorContains(t, err, "unsupported output format")
}

func TestGenDiffMissingFile(t *testing.T) {
	_, err := GenDiff(fixturePath("missing.json"), fixturePath("nested2.json"), "stylish")
	require.ErrorContains(t, err, "missing.json")

	_, err = GenDiff(fixturePath("nested1.json"), fixturePath("missing.json"), "stylish")
	require.ErrorContains(t, err, "missing.json")
}
