package main

import (
	"context"
	"log"
	"os"

	"code/internal/app"
)

func main() {
	if err := app.Run(context.Background(), os.Args, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
