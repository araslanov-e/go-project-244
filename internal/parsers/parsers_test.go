package parsers

import (
	"encoding/json"
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
		"timeout": json.Number("50"),
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
		"timeout": json.Number("20"),
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
				"timeout": json.Number("50"),
				"proxy":   "123.234.53.22",
				"follow":  false,
			},
		},
		{
			name: "yaml extension",
			file: "file2.yaml",
			expected: map[string]any{
				"timeout": json.Number("20"),
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

// Одинаковые данные в JSON и YAML должны разбираться в одинаковые map:
// без этого одно и то же число попадало бы в дифф как изменение.
func TestParseSameDataDifferentFormats(t *testing.T) {
	tests := []struct {
		name string
		json string
		yaml string
	}{
		{name: "int", json: `{"timeout": 50}`, yaml: "timeout: 50"},
		{name: "float", json: `{"ratio": 1.5}`, yaml: "ratio: 1.5"},
		{name: "integral float", json: `{"timeout": 50.0}`, yaml: "timeout: 50"},
		{name: "big int", json: `{"id": 9007199254740993}`, yaml: "id: 9007199254740993"},
		{name: "nested", json: `{"limits": {"port": 8080}}`, yaml: "limits:\n  port: 8080"},
		{name: "array", json: `{"ports": [80, 443]}`, yaml: "ports:\n  - 80\n  - 443"},
		{name: "mixed types", json: `{"a": 1, "b": "1", "c": true, "d": null}`,
			yaml: "a: 1\nb: \"1\"\nc: true\nd: null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fromJSON, err := Parse([]byte(tt.json), "json")
			require.NoError(t, err)

			fromYAML, err := Parse([]byte(tt.yaml), "yaml")
			require.NoError(t, err)

			assert.Equal(t, fromJSON, fromYAML)
		})
	}
}

// Правило для пустого ввода одно для всех форматов: это ошибка.
// Иначе случайно очищенный файл выглядел бы как конфигурация без ключей.
func TestParseEmptyInput(t *testing.T) {
	tests := []struct {
		name   string
		format string
		data   string
	}{
		{name: "empty json", format: "json", data: ""},
		{name: "blank json", format: "json", data: "   \n\t"},
		{name: "null json", format: "json", data: "null"},
		{name: "empty yaml", format: "yaml", data: ""},
		{name: "blank yaml", format: "yaml", data: "   \n"},
		{name: "comments only yaml", format: "yaml", data: "# только комментарий\n"},
		{name: "null yaml", format: "yaml", data: "null"},
		{name: "tilde yaml", format: "yaml", data: "~"},
		{name: "empty document yaml", format: "yaml", data: "---\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data), tt.format)
			require.ErrorIs(t, err, ErrEmptyInput)
		})
	}
}

// Пустую конфигурацию нужно записывать явно как {} — и это работает
// в обоих форматах (фикстуры empty.json и empty.yml как раз такие).
func TestParseEmptyObject(t *testing.T) {
	for _, format := range []string{"json", "yaml"} {
		t.Run(format, func(t *testing.T) {
			got, err := Parse([]byte("{}"), format)
			require.NoError(t, err)
			assert.Equal(t, map[string]any{}, got)
		})
	}
}

// На верхнем уровне должен быть объект, иначе сравнивать нечего.
func TestParseNotAnObject(t *testing.T) {
	tests := []struct {
		name   string
		format string
		data   string
	}{
		{name: "json array", format: "json", data: "[1, 2]"},
		{name: "json number", format: "json", data: "42"},
		{name: "json string", format: "json", data: `"text"`},
		{name: "json bool", format: "json", data: "true"},
		{name: "yaml array", format: "yaml", data: "- 1\n- 2\n"},
		{name: "yaml number", format: "yaml", data: "42"},
		{name: "yaml string", format: "yaml", data: "text"},
		{name: "yaml bool", format: "yaml", data: "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse([]byte(tt.data), tt.format)
			require.ErrorIs(t, err, ErrNotAnObject)
		})
	}
}

func TestParseJSONTrailingData(t *testing.T) {
	_, err := Parse([]byte(`{"key": "value"} trailing`), "json")
	require.ErrorContains(t, err, "unexpected data after the top-level value")
}

func TestParseFileEmpty(t *testing.T) {
	for _, name := range []string{"config.json", "config.yml"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			require.NoError(t, os.WriteFile(path, nil, 0o600))

			_, err := ParseFile(path)
			require.ErrorIs(t, err, ErrEmptyInput)
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
	// Обрыв в середине документа — это ошибка разбора, а не пустой файл.
	require.NotErrorIs(t, err, ErrEmptyInput)

	invalidYAML := filepath.Join(t.TempDir(), "broken.yml")
	require.NoError(t, os.WriteFile(invalidYAML, []byte("key: [unclosed"), 0o600))
	_, err = ParseFile(invalidYAML)
	require.ErrorContains(t, err, "parse yaml")
}
