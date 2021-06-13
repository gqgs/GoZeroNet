package uiwebsocket

import (
	"errors"

	"github.com/mailru/easyjson"
)

//go:generate go run github.com/mailru/easyjson/easyjson -all

type Message struct {
	ID           int64  `json:"id"`
	CMD          string `json:"cmd"`
	WrapperNonce string `json:"wrapper_nonce"`
}

func (w *uiWebsocket) handleMessage(rawMessage []byte) {
	message, err := decode(rawMessage)
	if err != nil {
		w.log.WithField("err", err).Warn("cmd decode error")
	}

	if err := w.route(rawMessage, message); err != nil {
		w.log.WithField("rawMessage", string(rawMessage)).Error(err)
	}
}

func (w *uiWebsocket) route(rawMessage []byte, message Message) error {
	switch message.CMD {
	case "channelJoin":
		return w.channelJoin(rawMessage, message)
	case "channelJoinAllsite":
		return w.channelJoinAllsite(rawMessage, message)
	case "siteSetLimit":
		return w.siteLimit(rawMessage, message)
	case "userGetSettings":
		return w.userGetSettings(rawMessage, message)
	case "userGetGlobalSettings":
		return w.userGetGlobalSettings(rawMessage, message)
	case "siteInfo":
		return w.siteInfo(rawMessage, message)
	case "serverInfo":
		return w.serverInfo(rawMessage, message)
	case "serverErrors":
		return w.serverErrors(rawMessage, message)
	case "announcerStats":
		return w.announcerStats(rawMessage, message)
	case "siteList":
		return w.siteList(rawMessage, message)
	case "filterIncludeList":
		return w.filterIncludeList(rawMessage, message)
	}

	for _, plugin := range w.plugins {
		if plugin.Handles(message.CMD) {
			return plugin.Handle(w.conn, message.CMD, message.ID, w.ID(), rawMessage)
		}
	}
	return errors.New("unknown cmd")
}

func decode(payload []byte) (Message, error) {
	var message Message
	err := easyjson.Unmarshal(payload, &message)
	return message, err
}