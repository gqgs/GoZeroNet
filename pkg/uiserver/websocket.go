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

		cmd, err := decodeCmd(message)
		if err != nil {
			// TODO: handle error
			w.log.WithField("err", err).Warn("cmd decode error")
			continue
		}

		switch cmd.CMD {
		case "userGetGlobalSettings":
			w.UserGetGlobalSettings(message, cmd.ID)
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
		}

		w.log.Info(cmd, err)
		w.log.Info(string(message))
	}
}

func decodeCmd(data []byte) (Cmd, error) {
	var cmd Cmd
	err := easyjson.Unmarshal(data, &cmd)
	return cmd, err
}
