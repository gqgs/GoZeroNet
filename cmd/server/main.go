package server

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gqgs/go-zeronet/pkg/file"
	"github.com/gqgs/go-zeronet/pkg/ui"
	"github.com/urfave/cli/v2"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start file and UI servers",
		Action: func(c *cli.Context) error {
			return serve(c.Context)
		},
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// The Python execution spawns two servers:
// FileServer serving TCP at 0.0.0.0 and random port.
// This follows the protocol at:
// https://zeronet.io/docs/help_zeronet/network_protocol/
//
// UIServer serving WSGI at 127.0.0.1:43110.
// This follows the protocol at:
// https://zeronet.io/docs/site_development/zeroframe_api_reference/

func serve(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileServer := file.NewServer()
	uiServer := ui.NewServer()

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := fileServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		if err := uiServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		idleConnsClosed <- struct{}{}
		idleConnsClosed <- struct{}{}
	}()

	go fileServer.Listen()
	go uiServer.Listen()

	<-idleConnsClosed
	<-idleConnsClosed

	return nil
}
