package code

import (
	"fmt"

	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/parsers"
)

// GenDiff сравнивает два конфигурационных файла и возвращает разницу
// в указанном формате вывода. Пустой format означает формат stylish.
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parsers.ParseFile(filepath1)
	if err != nil {
		return "", fmt.Errorf("parse %q: %w", filepath1, err)
	}

	data2, err := parsers.ParseFile(filepath2)
	if err != nil {
		return "", fmt.Errorf("parse %q: %w", filepath2, err)
	}

	formatter, err := formatters.New(format)
	if err != nil {
		return "", err
	}

	return formatter.Format(diff.Build(data1, data2))
}
