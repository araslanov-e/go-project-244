package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	data1 := map[string]any{
		"common": map[string]any{
			"setting1": "Value 1",
			"setting2": float64(200),
			"deep":     map[string]any{"id": float64(45)},
		},
		"group": map[string]any{"key": "value"},
		"same":  nil,
	}
	data2 := map[string]any{
		"common": map[string]any{
			"setting1": "Value 1",
			"setting3": true,
			"deep":     map[string]any{"id": float64(46)},
		},
		"group": "str",
		"same":  nil,
		"added": map[string]any{"key": "value"},
	}

	expected := []Node{
		{Key: "added", Status: Added, NewValue: map[string]any{"key": "value"}},
		{Key: "common", Status: Nested, Children: []Node{
			{Key: "deep", Status: Nested, Children: []Node{
				{Key: "id", Status: Changed, OldValue: float64(45), NewValue: float64(46)},
			}},
			{Key: "setting1", Status: Unchanged, OldValue: "Value 1", NewValue: "Value 1"},
			{Key: "setting2", Status: Removed, OldValue: float64(200)},
			{Key: "setting3", Status: Added, NewValue: true},
		}},
		{Key: "group", Status: Changed, OldValue: map[string]any{"key": "value"}, NewValue: "str"},
		{Key: "same", Status: Unchanged},
	}

	assert.Equal(t, expected, Build(data1, data2))
}

func TestBuildEmpty(t *testing.T) {
	assert.Empty(t, Build(map[string]any{}, map[string]any{}))
}
