package parsers

import (
	"fmt"

	"go.yaml.in/yaml/v3"
)

// yamlParser разбирает YAML. Состояния у него нет, поэтому фабрика
// создаёт его на каждый вызов, не заботясь о переиспользовании.
type yamlParser struct{}

// Parse — реализация Parser. Возвращает значение как есть: проверки
// общего контракта делает Parse пакета, они одинаковы для всех форматов.
func (yamlParser) Parse(data []byte) (any, error) {
	// Разбираем в any, а не сразу в map: так «документа нет» (пустой ввод,
	// одни комментарии или null) приходит как raw == nil.
	var raw any

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	return raw, nil
}
