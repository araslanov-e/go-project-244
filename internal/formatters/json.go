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

// FormatJSON выводит дерево диффа как JSON-объект, в котором ключи —
// имена свойств, а значения — описания изменений. У каждого описания
// есть status, остальные поля зависят от статуса:
//   - added — value (новое значение);
//   - removed, unchanged — value (значение из первой структуры);
//   - changed — oldValue и newValue;
//   - nested — children (объект того же вида для вложенных свойств).
//
// Ключи объектов encoding/json выводит в отсортированном порядке.
func FormatJSON(tree []diff.Node) (string, error) {
	data, err := json.MarshalIndent(jsonNodes(tree), "", jsonIndent)
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}

	return string(data), nil
}

func jsonNodes(nodes []diff.Node) map[string]jsonNode {
	// Пустая, но не nil map — чтобы в выводе был {}, а не null.
	result := make(map[string]jsonNode, len(nodes))

	for _, node := range nodes {
		item := jsonNode{"status": node.Status}

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

		result[node.Key] = item
	}

	return result
}
