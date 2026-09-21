package diff

import (
	"maps"
	"reflect"
	"slices"
)

// Status описывает, что произошло с ключом при сравнении двух структур.
type Status string

const (
	Added     Status = "added"     // ключ есть только во второй структуре
	Removed   Status = "removed"   // ключ есть только в первой структуре
	Unchanged Status = "unchanged" // значения по ключу равны
	Changed   Status = "changed"   // значения по ключу различаются
	Nested    Status = "nested"    // в обеих структурах по ключу объекты, разница — в Children
)

// Node — узел внутреннего представления диффа, описывает один ключ.
// OldValue — значение из первой структуры (Removed, Unchanged, Changed),
// NewValue — из второй (Added, Unchanged, Changed). Children заполняется
// только для Nested.
type Node struct {
	Key      string
	Status   Status
	OldValue any
	NewValue any
	Children []Node
}

// Build сравнивает две структуры и возвращает дифф в виде дерева узлов,
// отсортированных по ключу. Вглубь сравниваются только значения, которые
// в обеих структурах являются объектами, остальные попадают в дифф как есть.
func Build(data1, data2 map[string]any) []Node {
	keys := slices.Collect(maps.Keys(data1))
	for key := range data2 {
		if _, ok := data1[key]; !ok {
			keys = append(keys, key)
		}
	}
	// Порядок обхода map в Go не определён, поэтому ключи сортируем явно.
	slices.Sort(keys)

	nodes := make([]Node, 0, len(keys))

	for _, key := range keys {
		value1, in1 := data1[key]
		value2, in2 := data2[key]

		switch {
		case !in1:
			nodes = append(nodes, Node{Key: key, Status: Added, NewValue: value2})
		case !in2:
			nodes = append(nodes, Node{Key: key, Status: Removed, OldValue: value1})
		default:
			nodes = append(nodes, compare(key, value1, value2))
		}
	}

	return nodes
}

func compare(key string, value1, value2 any) Node {
	map1, ok1 := value1.(map[string]any)
	map2, ok2 := value2.(map[string]any)

	switch {
	case ok1 && ok2:
		return Node{Key: key, Status: Nested, Children: Build(map1, map2)}
	case reflect.DeepEqual(value1, value2):
		return Node{Key: key, Status: Unchanged, OldValue: value1, NewValue: value2}
	default:
		return Node{Key: key, Status: Changed, OldValue: value1, NewValue: value2}
	}
}
