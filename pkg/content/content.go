package content

type File struct {
	Sha512 string `json:"sha512"`
	Size   int    `json:"size"`
}

type Content struct {
	Address                  string            `json:"address"`
	BackgroundColor          string            `json:"background-color"`
	BackgroundColorDark      string            `json:"background-color-dark"`
	Description              string            `json:"description"`
	Files                    map[string]File   `json:"files"`
	Ignore                   string            `json:"ignore"`
	InnerPath                string            `json:"inner_path"`
	Modified                 int               `json:"modified"`
	PostmessageNonceSecurity bool              `json:"postmessage_nonce_security"`
	SignersSign              string            `json:"signers_sign"`
	Signs                    map[string]string `json:"signs"`
	SignsRequired            int               `json:"signs_required"`
	Title                    string            `json:"title"`
	Translate                []string          `json:"translate"`
	Viewport                 string            `json:"viewport"`
	ZeronetVersion           string            `json:"zeronet_version"`
}
