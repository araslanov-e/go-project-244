package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"code/internal/diff"
)

// Каждый уровень вложенности сдвигает строку на indentSize пробелов.
// Знак (+, - или пробел) вместе с пробелом после него занимают signWidth
// символов и «съедают» часть отступа: отступ до знака равен
// depth*indentSize - signWidth.
const (
	indentSize = 4
	signWidth  = 2
)

// FormatStylish выводит дифф в виде вложенных блоков в фигурных скобках:
// "+" — добавленное значение, "-" — удалённое, пробел — без изменений.
func FormatStylish(tree []diff.Node) string {
	var b strings.Builder

	b.WriteString("{\n")
	writeNodes(&b, tree, 1)
	b.WriteString("}")

	return b.String()
}

func writeNodes(b *strings.Builder, nodes []diff.Node, depth int) {
	for _, node := range nodes {
		switch node.Status {
		case diff.Added:
			writeLine(b, depth, "+", node.Key, stringify(node.NewValue, depth))
		case diff.Removed:
			writeLine(b, depth, "-", node.Key, stringify(node.OldValue, depth))
		case diff.Unchanged:
			writeLine(b, depth, " ", node.Key, stringify(node.OldValue, depth))
		case diff.Changed:
			writeLine(b, depth, "-", node.Key, stringify(node.OldValue, depth))
			writeLine(b, depth, "+", node.Key, stringify(node.NewValue, depth))
		case diff.Nested:
			fmt.Fprintf(b, "%s  %s: {\n", signIndent(depth), node.Key)
			writeNodes(b, node.Children, depth+1)
			fmt.Fprintf(b, "%s}\n", indent(depth))
		}
	}
}

func writeLine(b *strings.Builder, depth int, sign, key, value string) {
	fmt.Fprintf(b, "%s%s %s: %s\n", signIndent(depth), sign, key, value)
}

// stringify приводит значение к строке. Объекты выводятся блоком без знаков,
// depth — глубина строки, в которой находится значение.
func stringify(value any, depth int) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case map[string]any:
		var b strings.Builder

		b.WriteString("{\n")

		for _, key := range slices.Sorted(maps.Keys(v)) {
			fmt.Fprintf(&b, "%s%s: %s\n", indent(depth+1), key, stringify(v[key], depth+1))
		}

		b.WriteString(indent(depth) + "}")

		return b.String()
	default:
		return fmt.Sprint(v)
	}
}

func indent(depth int) string {
	return strings.Repeat(" ", depth*indentSize)
}

func signIndent(depth int) string {
	return strings.Repeat(" ", depth*indentSize-signWidth)
}
