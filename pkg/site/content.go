package site

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
	"github.com/gqgs/go-zeronet/pkg/lib/parser"
)

type Content struct {
	Address                  string                 `json:"address,omitempty"`
	AddressIndex             int                    `json:"address_index,omitempty"`
	BackgroundColor          string                 `json:"background-color,omitempty"`
	BackgroundColorDark      string                 `json:"background-color-dark,omitempty"`
	CertAuthType             string                 `json:"cert_auth_type,omitempty"`
	CertSign                 string                 `json:"cert_sign,omitempty"`
	CertUserID               string                 `json:"cert_user_id,omitempty"`
	CloneRoot                string                 `json:"clone_root,omitempty"`
	Cloneable                *bool                  `json:"cloneable,omitempty"`
	ClonedFrom               string                 `json:"cloned_from,omitempty"`
	DefaultPage              string                 `json:"default_page,omitempty"`
	Description              json.RawMessage        `json:"description,omitempty"`
	Domain                   string                 `json:"domain,omitempty"`
	Favicon                  string                 `json:"favicon,omitempty"`
	Files                    map[string]File        `json:"files"`
	FilesOptional            map[string]File        `json:"files_optional,omitempty"`
	Ignore                   string                 `json:"ignore,omitempty"`
	Includes                 map[string]Include     `json:"includes,omitempty"`
	InnerPath                string                 `json:"inner_path,omitempty"`
	Modified                 float64                `json:"modified,omitempty"`
	Optional                 string                 `json:"optional,omitempty"`
	PostmessageNonceSecurity bool                   `json:"postmessage_nonce_security,omitempty"`
	Settings                 map[string]interface{} `json:"settings,omitempty"`
	SignersSign              string                 `json:"signers_sign,omitempty"`
	Signs                    map[string]string      `json:"signs,omitempty"`
	SignsRequired            int                    `json:"signs_required,omitempty"`
	Title                    string                 `json:"title,omitempty"`
	Translate                []string               `json:"translate,omitempty"`
	Viewport                 string                 `json:"viewport,omitempty"`
	UserContents             *UserContents          `json:"user_contents,omitempty"`
	ZeronetVersion           string                 `json:"zeronet_version,omitempty"`
}

type File struct {
	ContentInnerPath string `json:"content_inner_path,omitempty"`
	Optional         bool   `json:"optional,omitempty"`
	PieceSize        int    `json:"piece_size,omitempty"`
	Piecemap         string `json:"piecemap,omitempty"`
	RelativePath     string `json:"relative_path,omitempty"`
	Sha512           string `json:"sha512"`
	Size             int    `json:"size"`
}

type Include struct {
	Signers         []string `json:"signers"`
	SignersRequired int      `json:"signers_required"`
}

type UserContents struct {
	Archived         map[string]int         `json:"archived,omitempty"`
	ArchivedBefore   int                    `json:"archived_before,omitempty"`
	CertSigners      map[string][]string    `json:"cert_signers,omitempty"`
	ContentInnerPath string                 `json:"content_inner_path,omitempty"`
	Optional         json.RawMessage        `json:"optional,omitempty"`
	PermissionRules  map[string]interface{} `json:"permission_rules,omitempty"`
	Permissions      map[string]interface{} `json:"permissions"`
	RelativePath     string                 `json:"relative_path,omitempty"`
}

func (c *Content) isValid() bool {
	if c == nil || !strings.HasSuffix(c.InnerPath, "content.json") {
		return false
	}

	if c.InnerPath == "content.json" {
		signers := make([]string, len(c.Signs))
		var i int
		for addr := range c.Signs {
			signers[i] = addr
			i++
		}
		signerdMsg := fmt.Sprintf("%d:%s", c.SignsRequired, strings.Join(signers, ","))
		if !crypto.IsValidSignature([]byte(signerdMsg), c.SignersSign, c.Address) {
			return false
		}
	}

	signs := make(map[string]string)
	for key, value := range c.Signs {
		signs[key] = value
	}

	// file was signed without signs
	c.Signs = nil
	contentJSON, err := json.Marshal(c)

	// restore signs
	c.Signs = signs

	if err != nil {
		return false
	}

	contentJSON, err = parser.FixJSONSpacing(bytes.NewReader(contentJSON))
	if err != nil {
		return false
	}

	var validSigns int
	for addr, sign := range c.Signs {
		if crypto.IsValidSignature(contentJSON, sign, addr) {
			validSigns++
			if validSigns >= c.SignsRequired {
				return true
			}
		}
	}

	return false
}
