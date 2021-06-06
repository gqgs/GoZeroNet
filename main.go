package main

import (
	"log"
	"os"

	"github.com/gqgs/go-zeronet/cmd/peer"
	"github.com/gqgs/go-zeronet/cmd/server"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			server.NewCommand(),
			peer.NewCommand(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Print(err)
	}
}
