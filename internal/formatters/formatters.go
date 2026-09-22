package formatters

import (
	"fmt"

	"code/internal/diff"
)

// Имена поддерживаемых форматов вывода. Stylish — формат по умолчанию.
const (
	Stylish = "stylish"
	Plain   = "plain"
)

// Format выводит дифф в указанном формате. Пустой format означает формат
// по умолчанию — stylish.
func Format(tree []diff.Node, format string) (string, error) {
	switch format {
	case "", Stylish:
		return FormatStylish(tree), nil
	case Plain:
		return FormatPlain(tree), nil
	default:
		return "", fmt.Errorf("unsupported output format: %q", format)
	}
}
