package db

import (
	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "Show database commands",
		Subcommands: []*cli.Command{
			{
				Name: "rebuild",
				Action: func(c *cli.Context) error {
					return rebuild(c.Context)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "site",
						Required: true,
					},
				},
			},
			{
				Name: "query",
				Action: func(c *cli.Context) error {
					return query(c.Context)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "site",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "query",
						Required: true,
					},
				},
			},
		},
	}
}
