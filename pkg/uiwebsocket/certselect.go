package uiwebsocket

import (
	"encoding/json"
	"errors"
	"fmt"
)

type (
	certSelectRequest struct {
		required
		Params certSelectParams `json:"params"`
	}
	certSelectParams struct {
		AcceptedDomains []string `json:"accepted_domains"`
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

func (w *uiWebsocket) certSelect(rawMessage []byte) error {
	payload := new(certSelectRequest)
	if err := json.Unmarshal(rawMessage, payload); err != nil {
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
		account := cert.AuthuserName + "@" + domain
		title := fmt.Sprintf("<b>%s</b>", account)
		if _, ok := accepted[domain]; !ok {
			css += "disabled "
		}
		if domain == currentSizeDomain {
			css += "active "
			title += " <small>({_[currently selected]})</small>"
		}
		body += fmt.Sprintf("<a href='#Select+account' class='select select-close cert %s' title='%s'>%s</a>", css, domain, title)
	}

	id := w.ID()
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
