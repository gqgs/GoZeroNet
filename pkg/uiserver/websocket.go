package uiserver

import (
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

type Cmd struct {
	ID  int    `json:"id"`
	CMD string `json:"cmd"`
}

func (w *uiWebsocket) Serve() {
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			w.log.Error(err)
			return
		}
		w.reqID++

		cmd, err := decodeCmd(message)
		if err != nil {
			// TODO: handle error
			w.log.WithField("err", err).Warn("cmd decode error")
			continue
		}

		switch cmd.CMD {
		case "userGetGlobalSettings":
			w.userGetGlobalSettings(message, cmd.ID)
			continue
		case "channelJoin":
			w.channelInfo(message, cmd.ID)
			continue
		case "siteInfo":
			w.siteInfo(message, cmd.ID)
			continue
		case "siteSetLimit":
			w.siteLimit(message, cmd.ID)
			continue
		case "optionalLimitStats":
			w.optionalLimitStats(message, cmd.ID)
			continue
		case "userGetSettings":
			w.userGetSettings(message, cmd.ID)
			continue
		case "serverInfo":
			w.serverInfo(message, cmd.ID)
			continue
		case "serverErrors":
			w.serverErrors(message, cmd.ID)
			continue
		case "announcerStats":
			w.announcerStats(message, cmd.ID)
			continue
		case "siteList":
			w.siteList(message, cmd.ID)
			continue
		case "channelJoinAllsite":
			w.channelJoinAllsite(message, cmd.ID)
			continue
		case "feedQuery":
			w.feedQuery(message, cmd.ID)
			continue
		case "filterIncludeList":
			w.filterIncludeList(message, cmd.ID)
			continue
		}

		w.log.WithField("cmd", cmd).Warn("unknown cmd")
		w.log.Warn(string(message))
	}
}

func decodeCmd(data []byte) (Cmd, error) {
	var cmd Cmd
	err := easyjson.Unmarshal(data, &cmd)
	return cmd, err
}
