package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

// jsonIndent — отступ одного уровня вложенности в выводе JSON.
const jsonIndent = "  "

// jsonNode — представление узла диффа в формате JSON. Набор полей
// зависит от статуса, поэтому узел собирается в map, а не в структуру
// с omitempty: иначе значение null (nil) пропадало бы из вывода.
type jsonNode map[string]any

// FormatJSON выводит дерево диффа как JSON-массив узлов. У каждого узла
// есть key и status, остальные поля зависят от статуса:
//   - added — value (новое значение);
//   - removed, unchanged — value (значение из первой структуры);
//   - changed — oldValue и newValue;
//   - nested — children (массив вложенных узлов).
func FormatJSON(tree []diff.Node) (string, error) {
	data, err := json.MarshalIndent(jsonNodes(tree), "", jsonIndent)
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}

	return string(data), nil
}

func jsonNodes(nodes []diff.Node) []jsonNode {
	// Пустой, но не nil слайс — чтобы в выводе был [], а не null.
	result := make([]jsonNode, 0, len(nodes))

	for _, node := range nodes {
		item := jsonNode{"key": node.Key, "status": node.Status}

		switch node.Status {
		case diff.Added:
			item["value"] = node.NewValue
		case diff.Removed, diff.Unchanged:
			item["value"] = node.OldValue
		case diff.Changed:
			item["oldValue"] = node.OldValue
			item["newValue"] = node.NewValue
		case diff.Nested:
			item["children"] = jsonNodes(node.Children)
		}

		result = append(result, item)
	}

	return result
}
