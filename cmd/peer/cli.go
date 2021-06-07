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
					return getFile(c.Context)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "addr",
						Required: true,
					},
				},
			},
			{
				Name: "cmd",
				Action: func(c *cli.Context) error {
					return cmd(c.Context)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "addr",
						Required: true,
					},
				},
			},
		},
	}
}
