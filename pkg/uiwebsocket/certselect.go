package uiwebsocket

import (
	"errors"
	"fmt"

	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

type (
	certSelectRequest struct {
		required
		Params certSelectParams `json:"params"`
	}
	certSelectParams struct {
		AcceptedDomains []string `json:"accepted_domains"`
	}
	certSelectResponse struct {
		required
		Result string `json:"result"`
	}

	notificationResponse struct {
		required
		Params []string `json:"params"`
	}

	scriptResponse struct {
		required
		Params string `json:"params"`
	}
)

func (w *uiWebsocket) certSelect(rawMessage []byte, message Message) error {
	payload := new(certSelectRequest)
	if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
		return err
	}

	user := w.site.User()
	if user == nil {
		return errors.New("user not found")
	}

	var currentSizeDomain string
	site := user.Sites[w.site.Address()]
	if site != nil {
		currentSizeDomain = site.Cert
	}
	accepted := make(map[string]struct{})
	for _, domain := range payload.Params.AcceptedDomains {
		accepted[domain] = struct{}{}
	}

	body := "<span style='padding-bottom: 5px; display: inline-block'>Select account you want to use in this site:</span>"
	body += "<a href='#Select+account' class='select select-close cert'><b>No certificate</b></a>"
	for domain, cert := range user.Certs {
		var css string
		account := cert.AuthUserName + "@" + domain
		title := fmt.Sprintf("<b>%s</b>", account)
		if _, ok := accepted[domain]; !ok {
			css += "disabled "
		}
		if domain == currentSizeDomain {
			css += "active "
			title += " <small>(currently selected)</small>"
		}
		body += fmt.Sprintf("<a href='#Select+account' class='select select-close cert %s' title='%s'>%s</a>", css, domain, title)
	}

	id := w.ID()
	w.waitingMutex.Lock()
	w.waitingResponses[id] = func(rawMessage []byte) error {
		type request struct {
			Cert string `json:"result"`
		}
		payload := new(request)
		if err := serialize.JSONUnmarshal(rawMessage, payload); err != nil {
			return err
		}
		if err := user.UpdateCert(w.site.Address(), payload.Cert); err != nil {
			return err
		}
		w.site.BroadcastSiteChange("cert_changed", payload.Cert)
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

	if err := w.conn.WriteJSON(notificationResponse{
		required{
			CMD: "notification",
			ID:  id,
		},
		[]string{"ask", body},
	}); err != nil {
		return err
	}

	script := fmt.Sprintf(`
		$(".notification .select.cert").on("click", function() {
			$(".notification .select").removeClass('active')
			zeroframe.response(%d, this.title)
			return false
		})
	`, id)

	return w.conn.WriteJSON(scriptResponse{
		required{
			CMD: "injectScript",
			ID:  w.ID(),
		},
		script,
	})
}
