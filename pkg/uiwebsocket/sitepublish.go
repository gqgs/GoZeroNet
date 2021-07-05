package uiwebsocket

import "errors"

type (
	sitePublishRequest struct {
		required
		Params sitePublishParams `json:"params"`
	}
	sitePublishParams struct {
		InnerPath  string `json:"inner_path"`
		PrivateKey string `json:"private_key"`
		Sign       *bool  `json:"sign"`
	}

	sitePublishResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) sitePublish(rawMessage []byte, message Message) error {
	payload := new(sitePublishRequest)
	if err := jsonUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	user := w.site.User()
	if user == nil {
		return errors.New("user not found")
	}

	privateKey := payload.Params.PrivateKey
	if privateKey == "" {
		site := user.Sites[w.site.Address()]
		if site == nil {
			return errors.New("site not found")
		}
		privateKey = site.AuthPrivatekey
	}

	sign := true
	if payload.Params.Sign != nil {
		sign = *payload.Params.Sign
	}

	if sign {
		if err := w.site.Sign(payload.Params.InnerPath, privateKey, user); err != nil {
			return err
		}
	}

	go w.site.Publish(payload.Params.InnerPath)

	return w.conn.WriteJSON(sitePublishResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
