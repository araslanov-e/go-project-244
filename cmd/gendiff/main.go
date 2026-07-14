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
		Name:        "gendiff",
		Usage:       "Compares two configuration files and shows a difference.",
		Description: "Description.",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("gendiff!")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
