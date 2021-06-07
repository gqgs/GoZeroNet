package peer

import "github.com/urfave/cli/v2"

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "peer",
		Usage: "Show peer commands",
		Subcommands: []*cli.Command{
			{
				Name: "ping",
				Action: func(c *cli.Context) error {
					return ping(c.Context, c.String("addr"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "addr",
						Required: true,
					},
				},
			},
			{
				Name: "handshake",
				Action: func(c *cli.Context) error {
					return handshake(c.Context, c.String("addr"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "addr",
						Required: true,
					},
				},
			},
			{
				Name: "getFile",
				Action: func(c *cli.Context) error {
					return getFile(c.Context, c.String("addr"), c.String("site"), c.String("inner_path"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "addr",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "site",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "inner_path",
						Required: true,
					},
				},
			},
		},
	}
}
