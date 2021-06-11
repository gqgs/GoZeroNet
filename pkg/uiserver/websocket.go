package uiserver

import (
	"errors"

	"github.com/fasthttp/websocket"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/site"
	"github.com/mailru/easyjson"
)

//go:generate go run github.com/mailru/easyjson/easyjson -all

//easyjson:skip
type uiWebsocket struct {
	// mu   sync.Mutex
	conn        *websocket.Conn
	log         log.Logger
	siteManager site.SiteManager
	reqID       int
}

func newUIWebsocket(conn *websocket.Conn, siteManager site.SiteManager) *uiWebsocket {
	return &uiWebsocket{
		conn:        conn,
		siteManager: siteManager,
		log:         log.New("uiwebsocket"),
	}
}

type Message struct {
	ID           int    `json:"id"`
	CMD          string `json:"cmd"`
	WrapperNonce string `json:"wrapper_nonce"`
}

func (w *uiWebsocket) Serve() {
	for {
		_, rawMessage, err := w.conn.ReadMessage()
		if err != nil {
			w.log.Error(err)
			return
		}
		w.reqID++

		message, err := decode(rawMessage)
		if err != nil {
			w.log.WithField("err", err).Warn("cmd decode error")
			continue
		}

		if err := w.route(rawMessage, message); err != nil {
			w.log.WithField("rawMessage", rawMessage).Error(err)
		}
	}
}

func (w *uiWebsocket) route(rawMessage []byte, message Message) error {
	switch message.CMD {
	case "userGetGlobalSettings":
		return w.userGetGlobalSettings(rawMessage, message)
	case "channelJoin":
		return w.channelInfo(rawMessage, message)
	case "siteInfo":
		return w.siteInfo(rawMessage, message)
	case "siteSetLimit":
		return w.siteLimit(rawMessage, message)
	case "optionalLimitStats":
		return w.optionalLimitStats(rawMessage, message)
	case "userGetSettings":
		return w.userGetSettings(rawMessage, message)
	case "serverInfo":
		return w.serverInfo(rawMessage, message)
	case "serverErrors":
		return w.serverErrors(rawMessage, message)
	case "announcerStats":
		return w.announcerStats(rawMessage, message)
	case "siteList":
		return w.siteList(rawMessage, message)
	case "channelJoinAllsite":
		return w.channelJoinAllsite(rawMessage, message)
	case "feedQuery":
		return w.feedQuery(rawMessage, message)
	case "filterIncludeList":
		return w.filterIncludeList(rawMessage, message)
	default:
		return errors.New("unknown cmd")
	}
}

func decode(payload []byte) (Message, error) {
	var message Message
	err := easyjson.Unmarshal(payload, &message)
	return message, err
}
