package formatters

import (
	"fmt"
	"strings"

	"code/internal/diff"
)

// FormatPlain выводит дифф построчно: по строке на каждое добавленное,
// удалённое или изменённое свойство с полным путём от корня.
// Неизменённые свойства не выводятся.
func FormatPlain(tree []diff.Node) string {
	return strings.Join(plainLines(tree, ""), "\n")
}

func plainLines(nodes []diff.Node, parent string) []string {
	var lines []string

	for _, node := range nodes {
		path := node.Key
		if parent != "" {
			path = parent + "." + node.Key
		}

		switch node.Status {
		case diff.Added:
			lines = append(lines, fmt.Sprintf("Property '%s' was added with value: %s",
				path, plainValue(node.NewValue)))
		case diff.Removed:
			lines = append(lines, fmt.Sprintf("Property '%s' was removed", path))
		case diff.Changed:
			lines = append(lines, fmt.Sprintf("Property '%s' was updated. From %s to %s",
				path, plainValue(node.OldValue), plainValue(node.NewValue)))
		case diff.Nested:
			lines = append(lines, plainLines(node.Children, path)...)
		case diff.Unchanged:
		}
	}

	return lines
}

// plainValue приводит значение к строке: составные значения (объекты
// и массивы) заменяются на [complex value], строки берутся в одинарные
// кавычки, остальное выводится как есть.
func plainValue(value any) string {
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
