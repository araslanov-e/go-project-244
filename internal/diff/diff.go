package diff

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

// Build сравнивает две плоские структуры и возвращает дифф: ключи в
// алфавитном порядке, "-" — значение из первой структуры, "+" — из второй,
// без знака — ключ есть в обеих с одинаковым значением.
func Build(data1, data2 map[string]any) string {
	keys := slices.Collect(maps.Keys(data1))
	for key := range data2 {
		if _, ok := data1[key]; !ok {
			keys = append(keys, key)
		}
	}
	slices.Sort(keys)

	var b strings.Builder
	b.WriteString("{\n")

	for _, key := range keys {
		value1, in1 := data1[key]
		value2, in2 := data2[key]

		switch {
		case in1 && in2 && reflect.DeepEqual(value1, value2):
			writeLine(&b, " ", key, value1)
		case in1 && in2:
			writeLine(&b, "-", key, value1)
			writeLine(&b, "+", key, value2)
		case in1:
			writeLine(&b, "-", key, value1)
		default:
			writeLine(&b, "+", key, value2)
		}
	}

	b.WriteString("}")

	return b.String()
}

func writeLine(b *strings.Builder, sign, key string, value any) {
	fmt.Fprintf(b, "  %s %s: %s\n", sign, key, stringify(value))
}

func stringify(value any) string {
	if value == nil {
		return "null"
	}

	return fmt.Sprint(value)
}
