package formatters

import (
	"encoding/json"
	"fmt"

	"code/internal/diff"
)

// jsonIndent — отступ одного уровня вложенности в выводе JSON.
const jsonIndent = "  "

// jsonNode — представление узла диффа в формате JSON. Имена полей формата
// заданы тегами; какие поля заполнены, зависит от статуса, остальные
// в вывод не попадают.
//
// Значения хранятся по указателю, чтобы отличать «поля нет» от null:
// nil-указатель omitempty пропускает, а указатель на nil выводится как null.
// У Children стоит omitzero, а не omitempty: он пропускает только nil,
// так что вложенный объект без свойств остаётся в выводе как {}.
type jsonNode struct {
	Status   diff.Status         `json:"status"`
	Value    *any                `json:"value,omitempty"`
	OldValue *any                `json:"oldValue,omitempty"`
	NewValue *any                `json:"newValue,omitempty"`
	Children map[string]jsonNode `json:"children,omitzero"`
}

// jsonFormatter выводит дифф как JSON-объект, без состояния.
type jsonFormatter struct{}

// Format — реализация Formatter.
func (jsonFormatter) Format(tree []diff.Node) (string, error) {
	return formatJSON(tree)
}

// formatJSON выводит дерево диффа как JSON-объект, в котором ключи —
// имена свойств, а значения — описания изменений. У каждого описания
// есть status, остальные поля зависят от статуса:
//   - added — value (новое значение);
//   - removed, unchanged — value (значение из первой структуры);
//   - changed — oldValue и newValue;
//   - nested — children (объект того же вида для вложенных свойств).
//
// Свойства encoding/json выводит в отсортированном порядке, поля
// описания — в порядке их объявления в jsonNode.
func formatJSON(tree []diff.Node) (string, error) {
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
		item := jsonNode{Status: node.Status}

		switch node.Status {
		case diff.Added:
			item.Value = &node.NewValue
		case diff.Removed, diff.Unchanged:
			item.Value = &node.OldValue
		case diff.Changed:
			item.OldValue = &node.OldValue
			item.NewValue = &node.NewValue
		case diff.Nested:
			item.Children = jsonNodes(node.Children)
		}

		result[node.Key] = item
	}

	return result
}
