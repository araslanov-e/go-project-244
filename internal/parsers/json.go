package parsers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// jsonParser разбирает JSON. Состояния у него нет, поэтому фабрика
// создаёт его на каждый вызов, не заботясь о переиспользовании.
type jsonParser struct{}

// Parse — реализация Parser. Возвращает значение как есть: проверки
// общего контракта делает Parse пакета, они одинаковы для всех форматов.
func (jsonParser) Parse(data []byte) (any, error) {
	// Разбираем в any, а не сразу в map: так «документа нет» (пустой ввод
	// или null) выглядит одинаково в обоих форматах — как raw == nil.
	var raw any

	decoder := json.NewDecoder(bytes.NewReader(data))
	// UseNumber оставляет число текстом вместо float64: так его
	// исходная запись доходит до canonical без потери точности.
	decoder.UseNumber()

	if err := decoder.Decode(&raw); err != nil {
		// io.EOF — данных нет вовсе; обрыв в середине документа даёт
		// io.ErrUnexpectedEOF и остаётся ошибкой разбора.
		if errors.Is(err, io.EOF) {
			return nil, ErrEmptyInput
		}

		return nil, fmt.Errorf("parse json: %w", err)
	}

	// Decode читает только первое значение, поэтому всё, что идёт
	// после документа, проверяем отдельно.
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return nil, errors.New("parse json: unexpected data after the top-level value")
	}

	return raw, nil
}
