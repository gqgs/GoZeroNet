package peer

import (
	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "peer",
		Usage: "Show peer commands",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "addr",
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name: "ping",
				Action: func(c *cli.Context) error {
					return ping(c.String("addr"))
				},
			},
			{
				Name: "handshake",
				Action: func(c *cli.Context) error {
					return handshake(c.String("addr"))
				},
			},
			{
				Name: "getFile",
				Action: func(c *cli.Context) error {
					return getFile(
						c.String("addr"),
						c.String("site"),
						c.String("inner_path"),
						c.Int("location"),
						c.Int("size"),
					)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "site",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "inner_path",
						Required: true,
					},
					&cli.IntFlag{
						Name: "location",
					},
					&cli.IntFlag{
						Name: "size",
					},
				},
			},
			{
				Name: "streamFile",
				Action: func(c *cli.Context) error {
					return streamFile(
						c.String("addr"),
						c.String("site"),
						c.String("inner_path"),
						c.Int("location"),
						c.Int("size"),
					)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "site",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "inner_path",
						Required: true,
					},
					&cli.IntFlag{
						Name: "location",
					},
					&cli.IntFlag{
						Name: "size",
					},
				},
			},
			{
				Name: "checkPort",
				Action: func(c *cli.Context) error {
					return checkPort(c.String("addr"), c.Int("port"))
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "port",
						Required: true,
					},
				},
			},
			{
				Name: "pex",
				Action: func(c *cli.Context) error {
					return pex(c.String("addr"), c.String("site"), c.Int("need"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "site",
						Required: true,
					},
					&cli.IntFlag{
						Name:  "need",
						Value: 5,
					},
				},
			},
			{
				Name: "listModified",
				Action: func(c *cli.Context) error {
					return listModified(c.String("addr"), c.String("site"), c.Int("since"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "site",
						Required: true,
					},
					&cli.IntFlag{
						Name: "since",
					},
				},
			},
			{
				Name: "update",
				Action: func(c *cli.Context) error {
					return update(c.String("addr"), c.String("site"), c.String("inner_path"))
				},
				Flags: []cli.Flag{
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
