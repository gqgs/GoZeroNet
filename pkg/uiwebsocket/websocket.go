package uiwebsocket

import (
	"errors"

	"github.com/bytedance/sonic"
)

type Message struct {
	ID           int64  `json:"id"`
	CMD          string `json:"cmd"`
	WrapperNonce string `json:"wrapper_nonce"`
	To           int64  `json:"to"`
}

// fields required for every message
type required struct {
	CMD string `json:"cmd"`
	ID  int64  `json:"id"`
	To  int64  `json:"to"`
}

var jsonUnmarshal = sonic.Unmarshal

type wsHandlerFunc func(rawMessage []byte, message Message) error

func (w *uiWebsocket) handleMessage(rawMessage []byte) {
	message, err := decode(rawMessage)
	if err != nil {
		w.log.WithField("err", err).Warn("cmd decode error")
		return
	}

	if w.site.PostmessageNonceSecurity() && isInnerFrameCmd(message) {
		if !w.site.HasValidWrapperNonce(message.WrapperNonce) {
			w.log.WithField("wrapper_nonce", message.WrapperNonce).
				WithField("cmd", message.CMD).
				WithField("id", message.ID).
				Warn("unknown wrapper nonce")
			return
		}
	}

	if err := w.route(rawMessage, message); err != nil {
		w.log.WithField("rawMessage", string(rawMessage)).Error(err)
	}
}

func (w *uiWebsocket) route(rawMessage []byte, message Message) error {
	switch message.CMD {
	case "certAdd":
		return w.certAdd(rawMessage, message)
	case "certSelect":
		return w.certSelect(rawMessage, message)
	case "channelJoin":
		return w.channelJoin(rawMessage, message)
	case "channelJoinAllsite":
		return w.adminOnly(w.channelJoinAllsite)(rawMessage, message)
	case "siteSetLimit":
		return w.adminOnly(w.siteSetLimit)(rawMessage, message)
	case "userGetSettings":
		return w.userGetSettings(message)
	case "userGetGlobalSettings":
		return w.userGetGlobalSettings(message)
	case "userSetSettings":
		return w.userSetSettings(rawMessage, message)
	case "siteInfo":
		return w.siteInfo(rawMessage, message)
	case "serverInfo":
		return w.serverInfo(message)
	case "siteUpdate":
		return w.adminOnly(w.siteUpdate)(rawMessage, message)
	case "serverErrors":
		return w.adminOnly(w.serverErrors)(rawMessage, message)
	case "announcerInfo":
		return w.announcerInfo(message)
	case "announcerStats":
		return w.adminOnly(w.announcerStats)(rawMessage, message)
	case "siteList":
		return w.adminOnly(w.siteList)(rawMessage, message)
	case "fileDelete":
		return w.fileDelete(rawMessage, message)
	case "fileGet":
		return w.fileGet(rawMessage, message)
	case "fileList":
		return w.fileList(rawMessage, message)
	case "fileNeed":
		return w.fileNeed(rawMessage, message)
	case "fileWrite":
		return w.fileWrite(rawMessage, message)
	case "dbQuery":
		return w.dbQuery(rawMessage, message)
	case "serverShutdown":
		return w.adminOnly(w.serverShutdown)(rawMessage, message)
	case "ping":
		return w.ping(message)
	case "response":
		return w.response(rawMessage, message)
	}

	for _, plugin := range w.plugins {
		if handler, ok := plugin.Handler(message.CMD); ok {
			return handler(w.conn, w.site, rawMessage)
		}
	}
	return errors.New("unknown cmd")
}

func decode(payload []byte) (Message, error) {
	var message Message
	err := jsonUnmarshal(payload, &message)
	return message, err
}

func (w *uiWebsocket) adminOnly(handler wsHandlerFunc) wsHandlerFunc {
	if w.site.IsAdmin() {
		return handler
	}

	return func(rawMessage []byte, message Message) error {
		return w.conn.WriteJSON(serverErrorRsponse{
			required{
				CMD: "response",
				ID:  w.ID(),
				To:  message.ID,
			},
			"Forbidden",
		})
	}
}

type serverErrorRsponse struct {
	required
	Error string `json:"error"`
}

func isInnerFrameCmd(message Message) bool {
	return message.ID < 1000000 && message.ID >= 0
}
