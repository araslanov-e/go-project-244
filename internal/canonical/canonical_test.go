package canonical

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueNumbers(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected any
	}{
		{name: "int from yaml", value: 50, expected: json.Number("50")},
		{name: "int64", value: int64(50), expected: json.Number("50")},
		{name: "uint64", value: uint64(50), expected: json.Number("50")},
		{name: "float from json", value: float64(50), expected: json.Number("50")},
		{name: "float32", value: float32(1.5), expected: json.Number("1.5")},
		{name: "json.Number int", value: json.Number("50"), expected: json.Number("50")},
		{name: "json.Number float", value: json.Number("1.50"), expected: json.Number("1.5")},
		{name: "json.Number integral float", value: json.Number("50.0"), expected: json.Number("50")},
		{name: "json.Number exponent", value: json.Number("5e1"), expected: json.Number("50")},
		{name: "negative float", value: -0.5, expected: json.Number("-0.5")},
		{name: "big int stays exact", value: int64(9007199254740993), expected: json.Number("9007199254740993")},
		{name: "big float", value: 1e30, expected: json.Number("1e+30")},
		{name: "out of range stays as is", value: json.Number("1e400"), expected: json.Number("1e400")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Value(tt.value))
		})
	}
}

func TestValueOther(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "string", value: "hexlet.io"},
		{name: "numeric string stays string", value: "50"},
		{name: "bool", value: false},
		{name: "nil", value: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, Value(tt.value))
		})
	}
}

func TestValueNested(t *testing.T) {
	value := map[string]any{
		"limits": map[string]any{"port": 8080, "ratio": float64(1.5)},
		"ports":  []any{80, float64(443), "http"},
	}

	expected := map[string]any{
		"limits": map[string]any{"port": json.Number("8080"), "ratio": json.Number("1.5")},
		"ports":  []any{json.Number("80"), json.Number("443"), "http"},
	}

	assert.Equal(t, expected, Value(value))
}

// Ключевое свойство: одно и то же число, пришедшее из JSON (float64)
// и из YAML (int), после приведения становится одинаковым значением.
func TestMapSameNumberFromDifferentFormats(t *testing.T) {
	fromJSON := Map(map[string]any{"timeout": json.Number("50"), "ratio": float64(1.5)})
	fromYAML := Map(map[string]any{"timeout": 50, "ratio": float64(1.5)})

	assert.Equal(t, fromJSON, fromYAML)
}

// Map не меняет исходную map, а возвращает новую.
func TestMapDoesNotMutateSource(t *testing.T) {
	source := map[string]any{"timeout": 50}

	result := Map(source)

	assert.Equal(t, map[string]any{"timeout": 50}, source)
	assert.Equal(t, map[string]any{"timeout": json.Number("50")}, result)
}

// Канонические числа остаются числами в JSON-выводе, а не строками.
func TestCanonicalNumberMarshalsAsNumber(t *testing.T) {
	data, err := json.Marshal(Map(map[string]any{"timeout": 50, "ratio": 1.5}))

	assert.NoError(t, err)
	assert.JSONEq(t, `{"timeout": 50, "ratio": 1.5}`, string(data))
}
