package main

import (
	"log"
	"os"

	"github.com/mtsmfm/sqlcodegen/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "sqlcodegen",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "init",
				Action: func(c *cli.Context) error {
					return cmd.RunInit()
				},
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate codes",
				Action: func(c *cli.Context) error {
					return cmd.RunGenerate()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
