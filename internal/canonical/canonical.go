// Package canonical приводит разобранные конфигурации к единому
// («каноническому») представлению. Парсеры разных форматов возвращают
// одно и то же число разными типами Go: 50 из JSON — это float64,
// а из YAML — int. Без приведения к общему виду такая пара попала бы
// в дифф как изменение, хотя смысл конфигурации одинаков.
package canonical

import (
	"encoding/json"
	"math"
	"strconv"
)

// maxExactInt — граница, до которой float64 представляет целые числа
// точно (2^53). Выше неё целое значение не восстановить без потерь,
// поэтому такие числа записываем как дробные.
const maxExactInt = 1 << 53

// Map возвращает копию m, в которой все значения приведены
// к каноническому виду (см. Value).
func Map(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))

	for key, value := range m {
		result[key] = Value(value)
	}

	return result
}

// Value приводит значение к каноническому виду: любое число (int, float64,
// json.Number) становится json.Number с единой текстовой записью,
// объекты и массивы обходятся рекурсивно, остальные значения (строки,
// булевы, nil) возвращаются как есть.
//
// json.Number выбран как канонический тип для чисел потому, что это
// строка: значения сравниваются посимвольно, а форматтеры выводят её
// без изменений — и как число в JSON-выводе.
func Value(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return Map(v)
	case []any:
		items := make([]any, len(v))
		for i, item := range v {
			items[i] = Value(item)
		}

		return items
	case json.Number:
		return numberFromText(v.String())
	case int:
		return json.Number(strconv.Itoa(v))
	case int64:
		return json.Number(strconv.FormatInt(v, 10))
	case uint64:
		return json.Number(strconv.FormatUint(v, 10))
	case float32:
		return numberFromFloat(float64(v))
	case float64:
		return numberFromFloat(v)
	default:
		return value
	}
}

// numberFromText приводит к канонической записи число, полученное как
// текст (json.Number при разборе JSON с UseNumber).
func numberFromText(text string) json.Number {
	if i, err := strconv.ParseInt(text, 10, 64); err == nil {
		return json.Number(strconv.FormatInt(i, 10))
	}

	if f, err := strconv.ParseFloat(text, 64); err == nil {
		return numberFromFloat(f)
	}

	// Разобрать не удалось — оставляем исходную запись.
	return json.Number(text)
}

// numberFromFloat приводит число с плавающей точкой к канонической
// записи. Целое значение записывается как целое (50.0 → "50"), чтобы
// совпасть с таким же числом из другого формата.
func numberFromFloat(f float64) json.Number {
	if f == math.Trunc(f) && math.Abs(f) < maxExactInt {
		return json.Number(strconv.FormatInt(int64(f), 10))
	}

	return json.Number(strconv.FormatFloat(f, 'g', -1, 64))
}
