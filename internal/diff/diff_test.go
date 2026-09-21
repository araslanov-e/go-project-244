package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	data1 := map[string]any{
		"host":    "hexlet.io",
		"timeout": float64(50),
		"proxy":   "123.234.53.22",
		"follow":  false,
	}
	data2 := map[string]any{
		"timeout": float64(20),
		"verbose": true,
		"host":    "hexlet.io",
	}

	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`
	assert.Equal(t, expected, Build(data1, data2))
}

func TestBuildEqual(t *testing.T) {
	data := map[string]any{"a": float64(1), "b": nil}

	assert.Equal(t, "{\n    a: 1\n    b: null\n}", Build(data, data))
}

func TestBuildEmpty(t *testing.T) {
	assert.Equal(t, "{\n}", Build(map[string]any{}, map[string]any{}))
}
