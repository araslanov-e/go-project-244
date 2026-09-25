package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixturePath(name string) string {
	return filepath.Join("..", "..", "testdata", "fixture", name)
}

// runCommand создаёт команду и запускает её — так же, как это делает
// точка входа, — и возвращает её вывод, поток ошибок и ошибку.
func runCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	var out, errOut bytes.Buffer

	cmd := New(&out, &errOut)
	err := cmd.Run(context.Background(), append([]string{"gendiff"}, args...))

	return out.String(), errOut.String(), err
}

// New только собирает команду: пока её не запустили, ничего не происходит,
// а запускает её вызывающий код — тогда, когда решит сам.
func TestNew(t *testing.T) {
	var out, errOut bytes.Buffer

	cmd := New(&out, &errOut)

	require.NotNil(t, cmd)
	assert.Equal(t, "gendiff", cmd.Name)
	assert.Empty(t, out.String(), "создание команды ничего не выводит")
	assert.Empty(t, errOut.String())

	err := cmd.Run(context.Background(), []string{
		"gendiff", fixturePath("nested1.json"), fixturePath("nested2.json"),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, out.String())
}

func TestRun(t *testing.T) {
	out, errOut, err := runCommand(t, fixturePath("nested1.json"), fixturePath("nested2.json"))
	require.NoError(t, err)
	assert.Equal(t, ExitSuccess, ExitCode(err))

	expected, readErr := os.ReadFile(fixturePath("nested_result.txt"))
	require.NoError(t, readErr)
	assert.Equal(t, string(expected)+"\n", out)
	assert.Empty(t, errOut)
}

// Неправильный вызов: сообщение об ошибке, подсказка по вызову
// и отдельный код завершения.
func TestRunUsageErrors(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "no arguments",
			args:     nil,
			contains: "expected 2 arguments",
		},
		{
			name:     "one argument",
			args:     []string{"only-one.json"},
			contains: "expected 2 arguments",
		},
		{
			name:     "unsupported format",
			args:     []string{"--format", "unknown", fixturePath("nested1.json"), fixturePath("nested2.json")},
			contains: `unsupported output format: "unknown" (supported: stylish, plain, json)`,
		},
		{
			name:     "unknown flag",
			args:     []string{"--bogus", fixturePath("nested1.json"), fixturePath("nested2.json")},
			contains: "flag provided but not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, errOut, err := runCommand(t, tt.args...)

			require.ErrorContains(t, err, tt.contains)
			assert.Equal(t, ExitUsage, ExitCode(err), "неправильный вызов — свой код завершения")
			assert.Empty(t, out, "результат в stdout не пишем")
			assert.Contains(t, errOut, tt.contains)
			assert.Contains(t, errOut, "USAGE", "при ошибке вызова показываем usage")
			assert.Contains(t, errOut, "<filepath1> <filepath2>")
			assert.Equal(t, 1, strings.Count(errOut, "gendiff: "), "сообщение печатается один раз")
		})
	}
}

// Сбой выполнения — другая категория: код завершения другой,
// usage не показывается (вызов-то правильный).
func TestRunFailureErrors(t *testing.T) {
	empty := filepath.Join(t.TempDir(), "cleared.json")
	require.NoError(t, os.WriteFile(empty, nil, 0o600))

	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "missing file",
			args:     []string{fixturePath("missing.json"), fixturePath("nested2.json")},
			contains: "missing.json",
		},
		{
			name:     "empty file",
			args:     []string{empty, fixturePath("nested2.json")},
			contains: "file contains no data",
		},
		{
			name:     "unsupported file format",
			args:     []string{fixturePath("nested_result.txt"), fixturePath("nested2.json")},
			contains: "unsupported file format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, errOut, err := runCommand(t, tt.args...)

			require.ErrorContains(t, err, tt.contains)
			assert.Equal(t, ExitFailure, ExitCode(err), "сбой выполнения — не ошибка вызова")
			assert.Empty(t, out)
			assert.Contains(t, errOut, tt.contains)
			assert.NotContains(t, errOut, "USAGE", "usage тут не при чём: вызов правильный")
			assert.Equal(t, 1, strings.Count(errOut, "gendiff: "), "сообщение печатается один раз")
		})
	}
}

// --help — это не ошибка: помощь идёт в обычный вывод и код завершения нулевой.
func TestRunHelp(t *testing.T) {
	out, errOut, err := runCommand(t, "--help")

	require.NoError(t, err)
	assert.Equal(t, ExitSuccess, ExitCode(err))
	assert.Contains(t, out, "<filepath1> <filepath2>")
	assert.Contains(t, out, "stylish, plain, json")
	assert.Empty(t, errOut)
}

func TestExitCode(t *testing.T) {
	assert.Equal(t, ExitSuccess, ExitCode(nil))
	assert.Equal(t, ExitFailure, ExitCode(os.ErrNotExist))
	assert.Equal(t, ExitUsage, ExitCode(usageErrorf("bad call")))
	// Ошибка вызова остаётся такой же, даже если её обернули по пути наверх.
	assert.Equal(t, ExitUsage, ExitCode(fmt.Errorf("run: %w", usageErrorf("bad call"))))
}
