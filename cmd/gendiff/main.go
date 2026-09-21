package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format",
				Value:   "stylish",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			if cmd.NArg() != 2 {
				return fmt.Errorf("expected 2 arguments (file paths), got %d", cmd.NArg())
			}

			filepath1 := cmd.Args().Get(0)
			filepath2 := cmd.Args().Get(1)
			format := cmd.String("format")

			// Само сравнение появится на следующих шагах.
			fmt.Printf("file1: %s, file2: %s, format: %s\n", filepath1, filepath2, format)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
