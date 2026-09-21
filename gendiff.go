package code

import (
	"fmt"

	"code/internal/diff"
	"code/internal/parsers"
)

// GenDiff сравнивает два конфигурационных файла и возвращает разницу
// в указанном формате вывода.
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parsers.ParseFile(filepath1)
	if err != nil {
		return "", fmt.Errorf("parse %q: %w", filepath1, err)
	}

	data2, err := parsers.ParseFile(filepath2)
	if err != nil {
		return "", fmt.Errorf("parse %q: %w", filepath2, err)
	}

	// Пока поддерживается только формат stylish, выбор форматтера появится позже.
	_ = format

	return diff.Build(data1, data2), nil
}
