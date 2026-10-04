package formatters

import (
	"fmt"
	"strings"

	"code/internal/diff"
)

// plainFormatter выводит дифф построчно, без состояния.
type plainFormatter struct{}

// Format — реализация Formatter. Ошибок при выводе не бывает, но
// сигнатура общая для всех форматтеров.
func (plainFormatter) Format(tree []diff.Node) (string, error) {
	return formatPlain(tree), nil
}

// formatPlain выводит дифф построчно: по строке на каждое добавленное,
// удалённое или изменённое свойство с полным путём от корня.
// Неизменённые свойства не выводятся.
func formatPlain(tree []diff.Node) string {
	return strings.Join(buildPlainLines(tree, ""), "\n")
}

// buildPlainLines строит строки вывода для узлов одного уровня; parent —
// путь до этого уровня от корня (пустой для верхнего).
func buildPlainLines(nodes []diff.Node, parent string) []string {
	var lines []string

	for _, node := range nodes {
		path := node.Key
		if parent != "" {
			path = parent + "." + node.Key
		}

		switch node.Status {
		case diff.Added:
			lines = append(lines, fmt.Sprintf("Property '%s' was added with value: %s",
				path, formatPlainValue(node.NewValue)))
		case diff.Removed:
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", path))
		case diff.Changed:
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s",
				path, formatPlainValue(node.OldValue), formatPlainValue(node.NewValue)))
		case diff.Nested:
			lines = append(lines, buildPlainLines(node.Children, path)...)
		case diff.Unchanged:
		}
	}

	return lines
}

// formatPlainValue приводит значение к строке: составные значения (объекты
// и массивы) заменяются на [complex value], строки берутся в одинарные
// кавычки, остальное выводится как есть.
func formatPlainValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case map[string]any, []any:
		return "[complex value]"
	case string:
		return "'" + v + "'"
	default:
		return fmt.Sprint(v)
	}
}
