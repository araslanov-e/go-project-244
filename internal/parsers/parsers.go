package parsers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code/internal/canonical"
)

// Имена поддерживаемых форматов входных данных. YAML и YML — одно и то же,
// различаются только расширением файла.
const (
	JSON = "json"
	YAML = "yaml"
	YML  = "yml"
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

// Parser разбирает содержимое одного формата в произвольное значение.
// Проверки общего контракта (документ есть, на верхнем уровне объект)
// в реализации не входят: они одинаковы для всех форматов и живут в Parse.
type Parser interface {
	Parse(data []byte) (any, error)
}

// registry — единственное место, где перечислены форматы входных данных:
// отсюда берут данные и NewParser, и Names, поэтому новый формат
// добавляется одной записью. Слайс, а не map, чтобы порядок в сообщениях
// пользователю был постоянным.
var registry = []struct {
	name string
	new  func() Parser
}{
	{JSON, func() Parser { return jsonParser{} }},
	{YAML, func() Parser { return yamlParser{} }},
	{YML, func() Parser { return yamlParser{} }},
}

// NewParser возвращает парсер для формата: "json", "yaml" или "yml".
func NewParser(format string) (Parser, error) {
	for _, entry := range registry {
		if entry.name == format {
			return entry.new(), nil
		}
	}

	return nil, fmt.Errorf("unsupported file format: %q (supported: %s)",
		format, strings.Join(Names(), ", "))
}

// Names возвращает имена поддерживаемых форматов входных данных —
// для сообщений пользователю.
func Names() []string {
	names := make([]string, 0, len(registry))
	for _, entry := range registry {
		names = append(names, entry.name)
	}

	return names
}

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
// документ другого вида — ErrNotAnObject (правило одно для всех форматов,
// поэтому проверки живут здесь, а не в реализациях Parser).
// Результат приведён к каноническому виду, поэтому одинаковые данные
// в разных форматах дают одинаковые map — их можно сравнивать напрямую.
func Parse(data []byte, format string) (map[string]any, error) {
	parser, err := NewParser(format)
	if err != nil {
		return nil, err
	}

	raw, err := parser.Parse(data)
	if err != nil {
		return nil, err
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
