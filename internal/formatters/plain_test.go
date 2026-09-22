package formatters

import (
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/assert"
)

func TestFormatPlain(t *testing.T) {
	tree := []diff.Node{
		{Key: "common", Status: diff.Nested, Children: []diff.Node{
			{Key: "follow", Status: diff.Added, NewValue: false},
			{Key: "setting1", Status: diff.Unchanged, OldValue: "Value 1", NewValue: "Value 1"},
			{Key: "setting2", Status: diff.Removed, OldValue: float64(200)},
			{Key: "setting3", Status: diff.Changed, OldValue: true, NewValue: nil},
			{Key: "setting5", Status: diff.Added, NewValue: map[string]any{"key5": "value5"}},
			{Key: "setting6", Status: diff.Nested, Children: []diff.Node{
				{Key: "wow", Status: diff.Changed, OldValue: "", NewValue: "so much"},
			}},
		}},
		{Key: "list", Status: diff.Changed, OldValue: []any{float64(1)}, NewValue: float64(1.5)},
	}

	expected := `Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'common.setting5' was added with value: [complex value]
Property 'common.setting6.wow' was updated. From '' to 'so much'
Property 'list' was updated. From [complex value] to 1.5`
	assert.Equal(t, expected, FormatPlain(tree))
}

func TestFormatPlainNoChanges(t *testing.T) {
	tree := []diff.Node{{Key: "a", Status: diff.Unchanged, OldValue: "x", NewValue: "x"}}
	assert.Empty(t, FormatPlain(tree))
	assert.Empty(t, FormatPlain(nil))
}
