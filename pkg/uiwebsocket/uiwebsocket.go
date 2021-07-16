package uiwebsocket

import (
	"context"
	"sync"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/lib/websocket"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/gqgs/go-zeronet/pkg/uiwebsocket/plugin"
)

type uiWebsocket struct {
	ctx              context.Context
	conn             websocket.Conn
	log              log.Logger
	siteManager      site.Manager
	fileServer       fileserver.Server
	site             *site.Site
	pubsubManager    pubsub.Manager
	channelsMutex    sync.RWMutex
	channels         map[string]struct{}
	allChannels      bool
	plugins          []plugin.Plugin
	ID               func() int64
	waitingMutex     sync.Mutex
	waitingResponses map[int64]func([]byte) error
}

func NewUIWebsocket(ctx context.Context, conn websocket.Conn, siteManager site.Manager, fileServer fileserver.Server,
	site *site.Site, pubsubManager pubsub.Manager) *uiWebsocket {
	counter := safe.Counter()
	return &uiWebsocket{
		ctx:           ctx,
		conn:          conn,
		siteManager:   siteManager,
		fileServer:    fileServer,
		log:           log.New("uiwebsocket"),
		site:          site,
		pubsubManager: pubsubManager,
		channels:      make(map[string]struct{}),
		allChannels:   false,
		plugins: []plugin.Plugin{
			plugin.NewNewsFeed(counter),
			plugin.NewOptionalManager(counter),
			plugin.NewContentFilter(counter),
			plugin.NewBigFile(counter),
		},
		ID:               counter,
		waitingResponses: make(map[int64]func([]byte) error),
	}
}
func (w *uiWebsocket) Serve() {
	ctx, cancel := context.WithCancel(w.ctx)
	defer cancel()

	if err := w.site.OpenDB(); err != nil {
		w.log.Fatal(err)
	}

	go w.handlePubsub(ctx)

	for {
		_, rawMessage, err := w.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err) {
				w.log.Error(err)
			}
			return
		}
		go w.handleMessage(rawMessage)
	}
}

func (w *uiWebsocket) handlePubsub(ctx context.Context) {
	messageCh := w.pubsubManager.Register("uiwebsocket", config.WebsocketBufferSize)
	defer w.pubsubManager.Unregister(messageCh)

	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messageCh:
			subscribed := w.allChannels || w.site.Address() == message.Site()
			if !subscribed {
				continue
			}

			event := message.Event()

			w.channelsMutex.RLock()
			_, joinedChannel := w.channels[event.String()]
			w.channelsMutex.RUnlock()

			if !joinedChannel {
				continue
			}

			if err := w.conn.WriteJSON(event); err != nil {
				w.log.Error(err)
			}
		}
	}
}
