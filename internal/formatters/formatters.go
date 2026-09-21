package formatters

import (
	"fmt"

	"code/internal/diff"
)

// Stylish — имя формата вывода по умолчанию.
const Stylish = "stylish"

// Format выводит дифф в указанном формате. Пустой format означает формат
// по умолчанию — stylish.
func Format(tree []diff.Node, format string) (string, error) {
	switch format {
	case "", Stylish:
		return FormatStylish(tree), nil
	default:
		return "", fmt.Errorf("unsupported output format: %q", format)
	}
}
