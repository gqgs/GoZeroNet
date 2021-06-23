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
		},
	}
}
