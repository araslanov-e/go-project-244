package parsers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"

	"code/internal/canonical"
)

// Контракт входных данных, одинаковый для всех форматов: файл содержит
// один документ, и этот документ — объект.
var (
	// ErrEmptyInput — в файле нет данных: он пуст, состоит из одних
	// комментариев или содержит только null. Пустую конфигурацию нужно
	// записать явно как {} — это валидно и в JSON, и в YAML.
	ErrEmptyInput = errors.New("file contains no data")

	// ErrNotAnObject — на верхнем уровне файла не объект.
	ErrNotAnObject = errors.New("expected an object at the top level")
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
// Данные должны содержать объект: пустой ввод — это ErrEmptyInput,
// документ другого вида — ErrNotAnObject (правило одно для всех форматов).
// Результат приведён к каноническому виду, поэтому одинаковые данные
// в разных форматах дают одинаковые map — их можно сравнивать напрямую.
func Parse(data []byte, format string) (map[string]any, error) {
	// Разбираем в any, а не сразу в map: так «документа нет» (пустой ввод
	// или null) выглядит одинаково в обоих форматах — как raw == nil.
	var raw any

	switch format {
	case "json":
		decoder := json.NewDecoder(bytes.NewReader(data))
		// UseNumber оставляет число текстом вместо float64: так его
		// исходная запись доходит до canonical без потери точности.
		decoder.UseNumber()

		if err := decoder.Decode(&raw); err != nil {
			// io.EOF — данных нет вовсе; обрыв в середине документа даёт
			// io.ErrUnexpectedEOF и остаётся ошибкой разбора.
			if errors.Is(err, io.EOF) {
				return nil, ErrEmptyInput
			}

			return nil, fmt.Errorf("parse json: %w", err)
		}

		// Decode читает только первое значение, поэтому всё, что идёт
		// после документа, проверяем отдельно.
		if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
			return nil, errors.New("parse json: unexpected data after the top-level value")
		}
	case "yml", "yaml":
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %q", format)
	}

	if raw == nil {
		return nil, ErrEmptyInput
	}

	result, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w, got %T", ErrNotAnObject, raw)
	}

	// Числа из JSON и YAML приходят разными типами Go, поэтому перед
	// сравнением приводим их к единому представлению.
	return canonical.Map(result), nil
}
