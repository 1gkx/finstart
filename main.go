package main

import (
	"log"
	"os"

	"github.com/1gkx/finstar/internal/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "Finstar",
		Action: cmd.Start,
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}
