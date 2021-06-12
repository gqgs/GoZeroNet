package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/websocket"
	"github.com/gqgs/go-zeronet/pkg/site"
)

type uiWebsocket struct {
	conn        websocket.Conn
	log         log.Logger
	siteManager site.SiteManager
	reqID       int64
	fileServer  fileserver.Server
	site        *site.Site
}

func NewUIWebsocket(conn websocket.Conn, siteManager site.SiteManager, fileServer fileserver.Server, site *site.Site) *uiWebsocket {
	return &uiWebsocket{
		conn:        conn,
		siteManager: siteManager,
		fileServer:  fileServer,
		log:         log.New("uiwebsocket"),
		site:        site,
	}
}

func (w *uiWebsocket) Serve() {
	for {
		_, rawMessage, err := w.conn.ReadMessage()
		if err != nil {
			w.log.Error(err)
			return
		}
		go w.handleMessage(rawMessage)
	}
}
