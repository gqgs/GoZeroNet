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
			&cli.IntFlag{
				Name:  "file_server_port",
				Value: config.FileServer.Port,
			},
			&cli.IntFlag{
				Name:  "ui_server_port",
				Value: config.UIServer.Port,
			},
		},
		Action: func(c *cli.Context) error {
			config.FileServer.Port = c.Int("file_server_port")
			config.UIServer.Port = c.Int("ui_server_port")
			return serve(c.Context)
		},
	}
}
