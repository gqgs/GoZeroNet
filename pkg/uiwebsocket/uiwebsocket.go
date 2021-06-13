package uiwebsocket

import (
	"context"
	"sync"

	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/websocket"
	"github.com/gqgs/go-zeronet/pkg/site"
)

type uiWebsocket struct {
	reqID         int64
	conn          websocket.Conn
	log           log.Logger
	siteManager   site.SiteManager
	fileServer    fileserver.Server
	site          *site.Site
	pubsubManager pubsub.Manager
	channelsMutex sync.RWMutex
	channels      map[string]struct{}
	allChannels   bool
}

func NewUIWebsocket(conn websocket.Conn, siteManager site.SiteManager,
	fileServer fileserver.Server, site *site.Site, pubsubManager pubsub.Manager) *uiWebsocket {
	return &uiWebsocket{
		conn:          conn,
		siteManager:   siteManager,
		fileServer:    fileServer,
		log:           log.New("uiwebsocket"),
		site:          site,
		pubsubManager: pubsubManager,
		channels:      make(map[string]struct{}),
		allChannels:   false,
	}
}

func (w *uiWebsocket) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go w.handleSubsub(ctx)

	for {
		_, rawMessage, err := w.conn.ReadMessage()
		if err != nil {
			w.log.Error(err)
			return
		}
		go w.handleMessage(rawMessage)
	}
}

func (w *uiWebsocket) handleSubsub(ctx context.Context) {
	messageCh := w.pubsubManager.Register()
	defer w.pubsubManager.Unregister(messageCh)

	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messageCh:
			subscribed := w.allChannels || w.site.Address() == message.Source()
			if !subscribed {
				continue
			}

			w.channelsMutex.RLock()
			_, joinedChannel := w.channels[message.Queue()]
			w.channelsMutex.RUnlock()

			if !joinedChannel {
				continue
			}

			if err := w.conn.Write(message.Body()); err != nil {
				w.log.Error(err)
			}
		}
	}
}
