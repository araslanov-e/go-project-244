package formatters

import (
	"fmt"

	"code/internal/diff"
)

// Имена поддерживаемых форматов вывода. Stylish — формат по умолчанию.
const (
	Stylish = "stylish"
	Plain   = "plain"
	JSON    = "json"
)

// Formatter выводит дерево диффа в одном из форматов. Выбор реализации
// отделён от её запуска: реализацию отдаёт New, данные ей передаёт
// вызывающий код.
type Formatter interface {
	Format(tree []diff.Node) (string, error)
}

// registry — единственное место, где перечислены форматы вывода: отсюда
// берут данные и New, и Names, и Supported, поэтому новый формат
// добавляется одной записью. Слайс, а не map, чтобы порядок в подсказке
// пользователю был постоянным, а формат по умолчанию — первым.
var registry = []struct {
	name string
	new  func() Formatter
}{
	{Stylish, func() Formatter { return stylishFormatter{} }},
	{Plain, func() Formatter { return plainFormatter{} }},
	{JSON, func() Formatter { return jsonFormatter{} }},
}

// New возвращает форматтер по имени формата. Пустое имя означает формат
// по умолчанию — stylish.
func New(format string) (Formatter, error) {
	if format == "" {
		format = Stylish
	}

	for _, entry := range registry {
		if entry.name == format {
			return entry.new(), nil
		}
	}

	return nil, fmt.Errorf("unsupported output format: %q", format)
}

// Names возвращает имена поддерживаемых форматов вывода — для подсказок
// пользователю (usage флага, текст ошибки).
func Names() []string {
	names := make([]string, 0, len(registry))
	for _, entry := range registry {
		names = append(names, entry.name)
	}

	return names
}

// Supported сообщает, поддерживается ли формат вывода. Пустое имя
// означает формат по умолчанию, поэтому тоже поддерживается.
// Позволяет проверить формат до чтения файлов.
func Supported(format string) bool {
	_, err := New(format)

	return err == nil
}
