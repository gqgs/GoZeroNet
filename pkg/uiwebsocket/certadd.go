package uiwebsocket

import (
	"errors"
	"html/template"
	"strings"

	"github.com/gqgs/go-zeronet/pkg/user"
)

type (
	certAddRequest struct {
		required
		Params []string `json:"params"`
	}
	certAddParams struct {
		Domain       string `json:"domain"`
		AuthType     string `json:"auth_type"`
		AuthUserName string `json:"auth_user_name"`
		Cert         string `json:"cert"`
	}

	certAddResponse struct {
		required
		Result string `json:"result"`
	}

	confirmResponse struct {
		required
		Params []string `json:"params"`
	}
)

var certTemplate *template.Template

func init() {
	certTemplate = template.Must(template.New("cert").Parse("{{.AuthType}}/{{.AuthUserName}}@{{.Domain}}"))
}

func (w *uiWebsocket) certAdd(rawMessage []byte, message Message) error {
	payload := new(certAddRequest)
	if err := jsonUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	if len(payload.Params) < 4 {
		return errors.New("missing arguments for request")
	}

	params := certAddParams{
		Domain:       payload.Params[0],
		AuthType:     payload.Params[1],
		AuthUserName: payload.Params[2],
		Cert:         payload.Params[3],
	}

	u := w.site.User()
	if u == nil {
		return errors.New("user not found")
	}

	var response string
	err := u.AddCert(w.site.Address(), params.Domain,
		params.AuthType, params.AuthUserName, params.Cert, false)
	switch err {
	case nil:
		builder := new(strings.Builder)
		if err := certTemplate.Execute(builder, params); err != nil {
			return err
		}
		if err := u.UpdateCert(w.site.Address(), params.Domain); err != nil {
			return err
		}
		w.site.BroadcastSiteChange("cert_changed", params.Domain)
		if err := w.conn.WriteJSON(notificationResponse{
			required{
				CMD: "notification",
				ID:  w.ID(),
			},
			[]string{"done", "New certificate added: " + builder.String()},
		}); err != nil {
			return err
		}
		response = "ok"
	case user.ErrCertAlreadyExists:
		builder := new(strings.Builder)
		currentCert := certAddParams{
			Domain:       params.Domain,
			AuthType:     u.Certs[params.Domain].AuthType,
			AuthUserName: u.Certs[params.Domain].AuthUserName,
		}
		if err := certTemplate.Execute(builder, currentCert); err != nil {
			return err
		}
		current := "Your current certificate: " + builder.String()
		builder.Reset()
		if err := certTemplate.Execute(builder, params); err != nil {
			return err
		}
		updated := "Change it to " + builder.String()

		id := w.ID()
		w.waitingMutex.Lock()
		w.waitingResponses[id] = func(rawMessage []byte) error {
			type request struct {
				Choice int `json:"result"`
			}
			payload := new(request)
			if err := jsonUnmarshal(rawMessage, payload); err != nil {
				return err
			}
			if payload.Choice == 1 {
				if err := u.AddCert(w.site.Address(), params.Domain,
					params.AuthType, params.AuthUserName, params.Cert, true); err != nil {
					return err
				}
				w.site.BroadcastSiteChange("cert_changed", params.Domain)
			}
			return w.conn.WriteJSON(certSelectResponse{
				required{
					CMD: "response",
					ID:  w.ID(),
					To:  message.ID,
				},
				"ok",
			})
		}
		w.waitingMutex.Unlock()
		return w.conn.WriteJSON(confirmResponse{
			required{
				CMD: "confirm",
				ID:  id,
			},
			[]string{current, updated},
		})
	case user.ErrCertNotChanged:
		response = "Not changed"
	default:
		return err
	}

	return w.conn.WriteJSON(certAddResponse{
		required{
			CMD: "response",
			ID:  w.ID(),
			To:  message.ID,
		},
		response,
	})
}
