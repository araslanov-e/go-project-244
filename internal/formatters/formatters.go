package formatters

import (
	"fmt"
	"slices"

	"code/internal/diff"
)

// Имена поддерживаемых форматов вывода. Stylish — формат по умолчанию.
const (
	Stylish = "stylish"
	Plain   = "plain"
	JSON    = "json"
)

// Names возвращает имена поддерживаемых форматов вывода — для подсказок
// пользователю (usage флага, текст ошибки).
func Names() []string {
	return []string{Stylish, Plain, JSON}
}

// Supported сообщает, поддерживается ли формат вывода. Пустое имя
// означает формат по умолчанию, поэтому тоже поддерживается.
// Позволяет проверить формат до чтения файлов.
func Supported(format string) bool {
	return format == "" || slices.Contains(Names(), format)
}

// Format выводит дифф в указанном формате. Пустой format означает формат
// по умолчанию — stylish.
func Format(tree []diff.Node, format string) (string, error) {
	switch format {
	case "", Stylish:
		return FormatStylish(tree), nil
	case Plain:
		return FormatPlain(tree), nil
	case JSON:
		return FormatJSON(tree)
	default:
		return "", fmt.Errorf("unsupported output format: %q", format)
	}
}
