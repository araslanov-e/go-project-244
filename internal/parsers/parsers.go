package parsers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
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
func Parse(data []byte, format string) (map[string]any, error) {
	result := map[string]any{}

	switch format {
	case "json":
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("parse json: %w", err)
		}
	case "yml", "yaml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %q", format)
	}

	return result, nil
}
