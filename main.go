package main

import (
	"XChange/cmd"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "XChange",
		Usage: "xchange [amount] [initial currency code] [first final currency code] ... [N final currency code]",
		Action: func(cCtx *cli.Context) error {
			return cmd.Exchange(cCtx)
		},
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"conf"},
				Usage:   "Setup config file",
				Action: func(c *cli.Context) error {
					loadConfig, err := cmd.LoadConfig()
					if err != nil {
						return err
					}
					fmt.Println("token: " + loadConfig.Token)
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "create",
						Usage: "Create config file",
						Action: func(c *cli.Context) error {
							return cmd.CreateConfigFile()
						},
					},
					{
						Name:  "token",
						Usage: "Set token",
						Action: func(c *cli.Context) error {
							return cmd.SetToken(c.Args().Get(0))
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
