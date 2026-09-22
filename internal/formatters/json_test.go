package formatters

import (
	"encoding/json"
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatJSON(t *testing.T) {
	tree := []diff.Node{
		{Key: "common", Status: diff.Nested, Children: []diff.Node{
			{Key: "follow", Status: diff.Added, NewValue: false},
			{Key: "setting1", Status: diff.Unchanged, OldValue: "Value 1", NewValue: "Value 1"},
			{Key: "setting2", Status: diff.Removed, OldValue: float64(200)},
			{Key: "setting3", Status: diff.Changed, OldValue: true, NewValue: nil},
			{Key: "setting5", Status: diff.Added, NewValue: map[string]any{"key5": "value5"}},
		}},
		{Key: "list", Status: diff.Changed, OldValue: []any{float64(1)}, NewValue: float64(1.5)},
	}

	expected := `[
  {
    "children": [
      {
        "key": "follow",
        "status": "added",
        "value": false
      },
      {
        "key": "setting1",
        "status": "unchanged",
        "value": "Value 1"
      },
      {
        "key": "setting2",
        "status": "removed",
        "value": 200
      },
      {
        "key": "setting3",
        "newValue": null,
        "oldValue": true,
        "status": "changed"
      },
      {
        "key": "setting5",
        "status": "added",
        "value": {
          "key5": "value5"
        }
      }
    ],
    "key": "common",
    "status": "nested"
  },
  {
    "key": "list",
    "newValue": 1.5,
    "oldValue": [
      1
    ],
    "status": "changed"
  }
]`

	got, err := FormatJSON(tree)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
	assert.True(t, json.Valid([]byte(got)))
}

func TestFormatJSONEmpty(t *testing.T) {
	got, err := FormatJSON(nil)
	require.NoError(t, err)
	assert.Equal(t, "[]", got)

	got, err = FormatJSON([]diff.Node{{Key: "a", Status: diff.Nested}})
	require.NoError(t, err)
	assert.JSONEq(t, `[{"key": "a", "status": "nested", "children": []}]`, got)
}

func TestFormatJSONUnsupportedValue(t *testing.T) {
	tree := []diff.Node{{Key: "a", Status: diff.Added, NewValue: make(chan int)}}

	_, err := FormatJSON(tree)
	require.ErrorContains(t, err, "marshal json")
}
