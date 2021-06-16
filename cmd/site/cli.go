package site

import "github.com/urfave/cli/v2"

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "site",
		Usage: "Show site commands",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "site",
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name: "download",
				Action: func(c *cli.Context) error {
					return download(c.String("site"))
				},
			},
		},
	}
}
