package main

import (
	"context"
	"os"

	"code/internal/app"
)

func main() {
	// Сообщение об ошибке печатает app.Run, main только отдаёт код
	// завершения: 2 — неправильный вызов команды, 1 — сбой выполнения.
	err := app.Run(context.Background(), os.Args, os.Stdout, os.Stderr)
	os.Exit(app.ExitCode(err))
}
