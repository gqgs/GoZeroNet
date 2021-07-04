package server

import (
	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start file and UI servers",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file_server_addr",
				Value: config.FileServerAddress,
			},
			&cli.StringFlag{
				Name:  "ui_server_addr",
				Value: config.UIServerAddress,
			},
		},
		Action: func(c *cli.Context) error {
			return serve(c.Context, c.String("file_server_addr"), c.String("ui_server_addr"))
		},
	}
}
