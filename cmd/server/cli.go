package server

import "github.com/urfave/cli/v2"

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start file and UI servers",
		Action: func(c *cli.Context) error {
			return serve(c.Context)
		},
	}
}
