package formatters

import (
	"testing"

	"code/internal/diff"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tree := []diff.Node{{Key: "a", Status: diff.Added, NewValue: float64(1)}}
	expected := "{\n  + a: 1\n}"

	// Пустое имя формата — тот же форматтер, что и stylish.
	for _, format := range []string{"", Stylish} {
		formatter, err := New(format)
		require.NoError(t, err, "format %q", format)

		got, err := formatter.Format(tree)
		require.NoError(t, err)
		assert.Equal(t, expected, got, "format %q", format)
	}

	formatter, err := New(Plain)
	require.NoError(t, err)

	got, err := formatter.Format(tree)
	require.NoError(t, err)
	assert.Equal(t, "Property 'a' was added with value: 1", got)

	formatter, err = New(JSON)
	require.NoError(t, err)

	got, err = formatter.Format(tree)
	require.NoError(t, err)
	assert.JSONEq(t, `{"a": {"status": "added", "value": 1}}`, got)

	_, err = New("unknown")
	require.ErrorContains(t, err, "unsupported output format")
}

// Names — источник подсказки пользователю, поэтому порядок должен быть
// постоянным, а формат по умолчанию — первым.
func TestNames(t *testing.T) {
	assert.Equal(t, []string{Stylish, Plain, JSON}, Names())
}

// Supported должен согласоваться с New: то, что Supported признаёт,
// New обязан создать, и наоборот.
func TestSupported(t *testing.T) {
	for _, format := range append(Names(), "") {
		assert.True(t, Supported(format), "format %q", format)

		_, err := New(format)
		assert.NoError(t, err, "format %q", format)
	}

	for _, format := range []string{"unknown", "STYLISH", "yaml"} {
		assert.False(t, Supported(format), "format %q", format)

		_, err := New(format)
		assert.Error(t, err, "format %q", format)
	}
}
