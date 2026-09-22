package formatters

import (
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatStylish(t *testing.T) {
	tree := []diff.Node{
		{Key: "common", Status: diff.Nested, Children: []diff.Node{
			{Key: "follow", Status: diff.Added, NewValue: false},
			{Key: "setting1", Status: diff.Unchanged, OldValue: "Value 1", NewValue: "Value 1"},
			{Key: "setting3", Status: diff.Changed, OldValue: true, NewValue: nil},
			{Key: "setting5", Status: diff.Added, NewValue: map[string]any{
				"key5": "value5",
				"deep": map[string]any{"id": float64(45)},
			}},
			{Key: "wow", Status: diff.Removed, OldValue: ""},
		}},
	}

	expected := `{
    common: {
      + follow: false
        setting1: Value 1
      - setting3: true
      + setting3: null
      + setting5: {
            deep: {
                id: 45
            }
            key5: value5
        }
` + "      - wow: \n" + // пробел после двоеточия у пустой строки обязателен
		`    }
}`
	assert.Equal(t, expected, FormatStylish(tree))
}

func TestFormatStylishEmpty(t *testing.T) {
	assert.Equal(t, "{\n}", FormatStylish(nil))
}

func TestFormat(t *testing.T) {
	tree := []diff.Node{{Key: "a", Status: diff.Added, NewValue: float64(1)}}
	expected := "{\n  + a: 1\n}"

	for _, format := range []string{"", Stylish} {
		got, err := Format(tree, format)
		require.NoError(t, err)
		assert.Equal(t, expected, got, "format %q", format)
	}

	got, err := Format(tree, Plain)
	require.NoError(t, err)
	assert.Equal(t, "Property 'a' was added with value: 1", got)

	_, err = Format(tree, "unknown")
	require.ErrorContains(t, err, "unsupported output format")
}
