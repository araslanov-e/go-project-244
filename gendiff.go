package code

import (
	"fmt"

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

	// Само сравнение появится на следующих шагах, пока выводим разобранные данные.
	return fmt.Sprintf("format: %s\nfile1: %v\nfile2: %v", format, data1, data2), nil
}
