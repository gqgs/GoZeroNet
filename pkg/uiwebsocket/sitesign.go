package uiwebsocket

import "errors"

type (
	siteSignRequest struct {
		required
		Params siteSignParams `json:"params"`
	}
	siteSignParams struct {
		InnerPath  string `json:"inner_path"`
		PrivateKey string `json:"private_key"`
	}

	siteSignResponse struct {
		required
		Result string `json:"result"`
	}
)

func (w *uiWebsocket) siteSign(rawMessage []byte, message Message) error {
	payload := new(siteSignRequest)
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

	if err := w.site.Sign(payload.Params.InnerPath, privateKey, user); err != nil {
		return err
	}

	return w.conn.WriteJSON(siteSignResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		"ok",
	})
}
