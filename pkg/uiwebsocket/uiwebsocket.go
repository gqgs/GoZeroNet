package uiwebsocket

import (
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/websocket"
	"github.com/gqgs/go-zeronet/pkg/site"
)

type websocketWriter interface {
	WriteJSON(v interface{}) error
}

//easyjson:skip
type uiWebsocket struct {
	conn        websocket.Conn
	log         log.Logger
	siteManager site.SiteManager
	reqID       int64
}

func NewUIWebsocket(conn websocket.Conn, siteManager site.SiteManager) *uiWebsocket {
	return &uiWebsocket{
		conn:        conn,
		siteManager: siteManager,
		log:         log.New("uiwebsocket"),
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
