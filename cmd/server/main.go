package server

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/uiserver"
	"github.com/gqgs/go-zeronet/pkg/user"
)

// The execution spawns two servers:
// FileServer serving TCP at 127.0.0.1 and a random port, by default.
// This follows the protocol at:
// https://zeronet.io/docs/help_zeronet/network_protocol/
//
// UIServer serving at 127.0.0.1:43111, by default.
// This follows the protocol at:
// https://zeronet.io/docs/site_development/zeroframe_api_reference/

func serve(ctx context.Context, fileServerAddr, uiServerAddr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pubsubManager := pubsub.NewManager()

	siteManager, err := site.NewSiteManager(pubsubManager)
	if err != nil {
		return err
	}

	userManager, err := user.NewUserManager()
	if err != nil {
		return err
	}

	fileServer, err := fileserver.NewServer(fileServerAddr)
	if err != nil {
		return err
	}

	uiServer, err := uiserver.NewServer(uiServerAddr, siteManager, fileServer, pubsubManager, userManager)
	if err != nil {
		return err
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := fileServer.Shutdown(); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		if err := uiServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		cancel()

		idleConnsClosed <- struct{}{}
		idleConnsClosed <- struct{}{}
	}()

	go fileServer.Listen()
	go uiServer.Listen(ctx)

	<-idleConnsClosed
	<-idleConnsClosed

	return nil
}
