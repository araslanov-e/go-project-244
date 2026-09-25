package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"

	"code/internal/canonical"
)

// ParseFile читает файл по пути (относительному или абсолютному) и
// разбирает его содержимое в map. Формат определяется по расширению файла.
func ParseFile(path string) (map[string]any, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path %q: %w", path, err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	format := strings.ToLower(strings.TrimPrefix(filepath.Ext(absPath), "."))

	return Parse(data, format)
}

// Parse разбирает данные в указанном формате: "json", "yml" или "yaml".
// Результат приведён к каноническому виду, поэтому одинаковые данные
// в разных форматах дают одинаковые map — их можно сравнивать напрямую.
func Parse(data []byte, format string) (map[string]any, error) {
	result := map[string]any{}

	switch format {
	case "json":
		decoder := json.NewDecoder(bytes.NewReader(data))
		// UseNumber оставляет число текстом вместо float64: так его
		// исходная запись доходит до canonical без потери точности.
		decoder.UseNumber()

		if err := decoder.Decode(&result); err != nil {
			return nil, fmt.Errorf("parse json: %w", err)
		}
	case "yml", "yaml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %q", format)
	}

	// Числа из JSON и YAML приходят разными типами Go, поэтому перед
	// сравнением приводим их к единому представлению.
	return canonical.Map(result), nil
}
