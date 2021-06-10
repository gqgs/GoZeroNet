package uiserver

import "github.com/gqgs/go-zeronet/pkg/site"

type (
	SiteInfoRequest struct {
		CMD    string         `json:"cmd"`
		ID     int            `json:"id"`
		Params SiteInfoParams `json:"params"`
	}
	SiteInfoParams map[struct{}]struct{}

	SiteInfoResponse struct {
		CMD    string         `json:"cmd"`
		ID     int            `json:"id"`
		To     int            `json:"to"`
		Result SiteInfoResult `json:"result"`
	}

	SiteInfoResult struct {
		Address        string    `json:"address"`
		AddressHash    string    `json:"address_hash"`
		AddressShort   string    `json:"address_short"`
		AuthAddress    string    `json:"auth_address"`
		BadFiles       int       `json:"bad_files"`
		Peers          int       `json:"peers"`
		NextSizeLimit  int       `json:"next_size_limit"`
		SizeLimit      int       `json:"size_limit"`
		Workers        int       `json:"workers"`
		ContentUpdated float64   `json:"content_updated"`
		Content        Content   `json:"content"`
		Settings       site.Site `json:"settings"`
	}

	Content struct {
		Address                  string   `json:"address"`
		BackgroundColor          string   `json:"background-color"`
		BackgroundColorDark      string   `json:"background-color-dark"`
		Description              string   `json:"description"`
		Files                    int      `json:"files"`
		FilesOptional            int      `json:"files_optional"`
		Ignore                   string   `json:"ignore"`
		Includes                 int      `json:"includes"`
		InnerPath                string   `json:"inner_path"`
		Modified                 int      `json:"modified"`
		PostmessageNonceSecurity bool     `json:"postmessage_nonce_security"`
		SignsRequired            int      `json:"signs_required"`
		Title                    string   `json:"title"`
		Translate                []string `json:"translate"`
		ViewPort                 string   `json:"viewport"`
		ZeronetVersion           string   `json:"zeronet_version"`
	}
)

func (w *uiWebsocket) siteInfo(message []byte, id int) {
	err := w.conn.WriteJSON(SiteInfoResponse{
		CMD: "response",
		To:  id,
		Result: SiteInfoResult{
			Address:        "1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D",
			AddressHash:    "f69941233e191d9e00f0cd16c5da10b0124d1c0a498b5ecfa1448b21a3eb0094",
			AddressShort:   "1HeLLo..Tf3D",
			AuthAddress:    "15yYcwrSCKLg3KsokUzwgNMxvapMmbz6iW",
			Peers:          200,
			SizeLimit:      100,
			NextSizeLimit:  200,
			ContentUpdated: 1623347716.8319452,
			Content: Content{
				Address:                  "1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D",
				BackgroundColor:          "#F2F4F6",
				BackgroundColorDark:      "#383F46",
				Description:              "New ZeroHello",
				Files:                    32,
				Ignore:                   "(js|css)/(?!all.(js|css))",
				InnerPath:                "content.json",
				Modified:                 1604369175,
				PostmessageNonceSecurity: true,
				Title:                    "ZeroHello",
				ViewPort:                 "width=device-width, initial-scale=0.8",
				ZeronetVersion:           "0.7.2",
			},
			Settings: w.siteManager.Site("1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D"),
		},
	})
	w.log.IfError(err)
}
