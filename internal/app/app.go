package app

import (
	"context"
	"fmt"
	"io"

	"code"

	"github.com/urfave/cli/v3"
)

// Run собирает CLI-команду gendiff и запускает её с аргументами args,
// результат пишется в output.
func Run(ctx context.Context, args []string, output io.Writer) error {
	return newCommand(output).Run(ctx, args)
}

func newCommand(output io.Writer) *cli.Command {
	return &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<filepath1> <filepath2>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format",
				Value:   "stylish",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			return run(output, cmd)
		},
	}
}

func run(output io.Writer, cmd *cli.Command) error {
	if cmd.NArg() != 2 {
		return fmt.Errorf("expected 2 arguments (file paths), got %d", cmd.NArg())
	}

	diff, err := code.GenDiff(cmd.Args().Get(0), cmd.Args().Get(1), cmd.String("format"))
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(output, diff); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
