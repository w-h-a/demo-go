package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/demo-go/cmd"
)

func main() {
	app := &cli.App{
		Name: "demo",
		Commands: []*cli.Command{
			{
				Name: "demo",
				Action: func(ctx *cli.Context) error {
					return cmd.Run(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
