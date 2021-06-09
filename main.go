package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gqgs/go-zeronet/cmd/db"
	"github.com/gqgs/go-zeronet/cmd/peer"
	"github.com/gqgs/go-zeronet/cmd/server"
	"github.com/urfave/cli/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := &cli.App{
		Commands: []*cli.Command{
			server.NewCommand(),
			peer.NewCommand(),
			db.NewCommand(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Print(err)
	}
}
