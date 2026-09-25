package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"code"
	"code/internal/formatters"

	"github.com/urfave/cli/v3"
)

// Коды завершения команды. По ним скрипт различает, что исправлять:
// сам вызов gendiff или причину сбоя.
const (
	ExitSuccess = 0 // успех
	ExitFailure = 1 // сбой выполнения: файл не читается, не разбирается и т.п.
	ExitUsage   = 2 // неправильный вызов: аргументы, флаги, формат вывода
)

// UsageError — команда вызвана неправильно: не то число аргументов,
// неизвестный флаг или неизвестный формат вывода. Run печатает такие
// ошибки вместе с подсказкой по вызову, а ExitCode отдаёт для них
// отдельный код завершения.
//
// Метод ExitCode у типа намеренно не объявлен: тогда cli считал бы его
// своим ExitCoder и вызывал бы os.Exit прямо внутри cmd.Run — команда
// перестала бы быть тестируемой. Код завершения даёт функция ExitCode.
type UsageError struct {
	err error
}

func (e *UsageError) Error() string { return e.err.Error() }

func (e *UsageError) Unwrap() error { return e.err }

func usageErrorf(format string, args ...any) error {
	return &UsageError{err: fmt.Errorf(format, args...)}
}

// ExitCode возвращает код завершения для ошибки, полученной от Run.
func ExitCode(err error) int {
	var usageErr *UsageError

	switch {
	case err == nil:
		return ExitSuccess
	case errors.As(err, &usageErr):
		return ExitUsage
	default:
		return ExitFailure
	}
}

// Run собирает CLI-команду gendiff и запускает её с аргументами args:
// результат пишется в output, сообщения об ошибках — в errOutput.
// Ошибка возвращается только для того, чтобы вызывающий получил по ней
// код завершения через ExitCode; печатать её повторно не нужно.
func Run(ctx context.Context, args []string, output, errOutput io.Writer) error {
	cmd := newCommand(output, errOutput)

	err := cmd.Run(ctx, args)
	if err == nil {
		return nil
	}

	// Писать в поток ошибок больше некуда, поэтому результат записи
	// здесь и ниже не проверяем.
	_, _ = fmt.Fprintf(errOutput, "%s: %v\n", cmd.Name, err)

	var usageErr *UsageError
	if errors.As(err, &usageErr) {
		// Вызов неправильный — показываем, как вызывать правильно.
		// Помощь cli печатает в cmd.Writer, поэтому на этом пути
		// направляем его в поток ошибок.
		_, _ = fmt.Fprintln(errOutput)
		cmd.Writer = errOutput
		_ = cli.ShowAppHelp(cmd)
	}

	return err
}

func newCommand(output, errOutput io.Writer) *cli.Command {
	return &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<filepath1> <filepath2>",
		Writer:    output,
		ErrWriter: errOutput,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format: " + strings.Join(formatters.Names(), ", "),
				Value:   formatters.Stylish,
			},
		},
		// Ошибки разбора флагов (неизвестный флаг, отсутствующее
		// значение) cli по умолчанию печатает сам. Перехватываем их,
		// чтобы все неправильные вызовы шли одним путём и с одним
		// кодом завершения.
		OnUsageError: func(_ context.Context, _ *cli.Command, err error, _ bool) error {
			return &UsageError{err: err}
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			return run(output, cmd)
		},
	}
}

func run(output io.Writer, cmd *cli.Command) error {
	if cmd.NArg() != 2 {
		return usageErrorf("expected 2 arguments (file paths), got %d", cmd.NArg())
	}

	// Формат проверяем до чтения файлов: неизвестный формат — ошибка
	// вызова, а не сбой обработки данных.
	format := cmd.String("format")
	if !formatters.Supported(format) {
		return usageErrorf("unsupported output format: %q (supported: %s)",
			format, strings.Join(formatters.Names(), ", "))
	}

	diff, err := code.GenDiff(cmd.Args().Get(0), cmd.Args().Get(1), format)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(output, diff); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
