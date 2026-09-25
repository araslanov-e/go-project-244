package main

import (
	"context"
	"os"

	"code/internal/app"
)

func main() {
	cmd := app.New(os.Stdout, os.Stderr)

	// Сообщение об ошибке печатает сама команда, точка входа только
	// решает, когда её запустить, и отдаёт код завершения:
	// 2 — неправильный вызов команды, 1 — сбой выполнения.
	err := cmd.Run(context.Background(), os.Args)

	os.Exit(app.ExitCode(err))
}
