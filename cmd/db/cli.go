package db

import (
	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "Show database commands",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "site",
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name: "rebuild",
				Action: func(c *cli.Context) error {
					return rebuild(c.String("site"))
				},
			},
		},
	}
}
