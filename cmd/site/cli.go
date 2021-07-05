package site

import (
	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "site",
		Usage: "Show site commands",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "site",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "since",
				Usage: "Download content.json files since this many days ago",
				Value: 7,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name: "download",
				Action: func(c *cli.Context) error {
					return download(c.String("site"), c.Int("since"))
				},
			},
			{
				Name: "download-recent",
				Action: func(c *cli.Context) error {
					return downloadRecent(c.String("site"), c.Int("since"))
				},
			},
			{
				Name: "verify",
				Action: func(c *cli.Context) error {
					return verify(c.String("site"), c.String("inner_path"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "inner_path",
						Required: true,
					},
				},
			},
			{
				Name: "sign",
				Action: func(c *cli.Context) error {
					return sign(c.String("site"), c.String("inner_path"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "inner_path",
						Required: true,
					},
				},
			},
			{
				Name: "publish",
				Action: func(c *cli.Context) error {
					return publish(c.String("site"), c.String("inner_path"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "inner_path",
						Required: true,
					},
				},
			},
		},
	}
}
