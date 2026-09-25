package code

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

// Одинаковые данные в JSON и YAML не должны давать различий: до
// приведения к единому представлению число 50 из JSON (float64) и из
// YAML (int) считались разными значениями.
func TestGenDiffSameDataDifferentFormats(t *testing.T) {
	files := []string{"numbers1.json", "numbers1.yml"}

	for _, file1 := range files {
		for _, file2 := range files {
			t.Run(file1+" vs "+file2, func(t *testing.T) {
				plain, err := GenDiff(fixturePath(file1), fixturePath(file2), "plain")
				require.NoError(t, err)
				assert.Empty(t, plain, "изменений быть не должно")

				stylish, err := GenDiff(fixturePath(file1), fixturePath(file2), "stylish")
				require.NoError(t, err)

				for _, line := range strings.Split(stylish, "\n") {
					trimmed := strings.TrimSpace(line)
					assert.False(t, strings.HasPrefix(trimmed, "+ ") || strings.HasPrefix(trimmed, "- "),
						"строка %q помечена как изменение", line)
				}
			})
		}
	}
}

// Дифф не должен зависеть от того, в каких форматах записаны файлы.
func TestGenDiffIndependentOfFormatPair(t *testing.T) {
	expected, err := GenDiff(fixturePath("numbers1.json"), fixturePath("numbers2.json"), "stylish")
	require.NoError(t, err)
	require.Contains(t, expected, "- retries: 3")

	pairs := [][2]string{
		{"numbers1.json", "numbers2.yml"},
		{"numbers1.yml", "numbers2.json"},
		{"numbers1.yml", "numbers2.yml"},
	}

	for _, pair := range pairs {
		t.Run(pair[0]+" vs "+pair[1], func(t *testing.T) {
			got, err := GenDiff(fixturePath(pair[0]), fixturePath(pair[1]), "stylish")
			require.NoError(t, err)
			assert.Equal(t, expected, got)
		})
	}
}

func TestGenDiffUnsupportedFormat(t *testing.T) {
	_, err := GenDiff(fixturePath("nested1.json"), fixturePath("nested2.json"), "unknown")
	require.ErrorContains(t, err, "unsupported output format")
}

// Пустой файл отклоняется в любом формате и в любом из аргументов,
// а в сообщении видно, какой именно файл пуст.
func TestGenDiffEmptyFile(t *testing.T) {
	for _, name := range []string{"cleared.json", "cleared.yml"} {
		t.Run(name, func(t *testing.T) {
			empty := filepath.Join(t.TempDir(), name)
			require.NoError(t, os.WriteFile(empty, nil, 0o600))

			_, err := GenDiff(empty, fixturePath("file2.json"), "stylish")
			require.ErrorContains(t, err, name)
			require.ErrorContains(t, err, "file contains no data")

			_, err = GenDiff(fixturePath("file1.json"), empty, "stylish")
			require.ErrorContains(t, err, name)
			require.ErrorContains(t, err, "file contains no data")
		})
	}
}

func TestGenDiffMissingFile(t *testing.T) {
	_, err := GenDiff(fixturePath("missing.json"), fixturePath("nested2.json"), "stylish")
	require.ErrorContains(t, err, "missing.json")

	_, err = GenDiff(fixturePath("nested1.json"), fixturePath("missing.json"), "stylish")
	require.ErrorContains(t, err, "missing.json")
}
