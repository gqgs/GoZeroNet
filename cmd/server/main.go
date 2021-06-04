package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gqgs/go-zeronet/pkg/file"
	"github.com/gqgs/go-zeronet/pkg/ui"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Server struct {
	FileServer file.Server
	UIServer   ui.Server
}

// The execution spawns two servers:
// FileServer serving TCP at 0.0.0.0 and random port.
// This follows the protocol at:
// https://zeronet.io/docs/help_zeronet/network_protocol/
//
// UIServer serving WSGI at 127.0.0.1:43110.
// This follows the protocol at:
// https://zeronet.io/docs/site_development/zeroframe_api_reference/

func main() {
	const port = 15441
	const ip = "127.0.0.1"

	// "Setting the Peer ID to "UT3530" tells trackers that you're using uTorrent v3.5.3"
	// https://github.com/jaruba/PowderWeb/wiki/Guide#private-torrent-trackers

	// peer_id := "-UT3530-" + random.Base62String(12)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fileServer := file.Server{}
	uiServer := ui.Server{}

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

	println("done!")
}
