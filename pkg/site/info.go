package site

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/gqgs/go-zeronet/pkg/config"
)

type Info struct {
	Address        string        `json:"address"`
	AddressHash    string        `json:"address_hash"`
	AddressShort   string        `json:"address_short"`
	AuthAddress    string        `json:"auth_address"`
	BadFiles       int           `json:"bad_files"`
	Events         []interface{} `json:"events,omitempty"`
	CertUserID     string        `json:"cert_user_id"`
	Peers          int           `json:"peers"`
	NextSizeLimit  int           `json:"next_size_limit"`
	SizeLimit      int           `json:"size_limit"`
	Workers        int           `json:"workers"`
	ContentUpdated float64       `json:"content_updated"`
	Content        *Content      `json:"content"`
	Settings       *Settings     `json:"settings"`
}

type File struct {
	Sha512 string `json:"sha512"`
	Size   int    `json:"size"`
}

type Include struct {
	Signers         []string `json:"signers,omitempty"`
	SignersRequired int      `json:"signers_required,omitempty"`
}

type Content struct {
	Address                  string             `json:"address,omitempty"`
	AddressIndex             int                `json:"address_index,omitempty"`
	BackgroundColor          string             `json:"background-color,omitempty"`
	BackgroundColorDark      string             `json:"background-color-dark,omitempty"`
	CloneRoot                string             `json:"clone_root,omitempty"`
	Cloneable                bool               `json:"cloneable,omitempty"`
	ClonedFrom               string             `json:"cloned_from,omitempty"`
	Description              string             `json:"description,omitempty"`
	Favicon                  string             `json:"favicon,omitempty"`
	Files                    map[string]File    `json:"files,omitempty"`
	FilesOptional            map[string]File    `json:"files_optional,omitempty"`
	Ignore                   string             `json:"ignore,omitempty"`
	Includes                 map[string]Include `json:"includes,omitempty"`
	InnerPath                string             `json:"inner_path,omitempty"`
	Modified                 int                `json:"modified,omitempty"`
	Optional                 string             `json:"optional,omitempty"`
	PostmessageNonceSecurity bool               `json:"postmessage_nonce_security,omitempty"`
	SignersSign              string             `json:"signers_sign,omitempty"`
	Signs                    map[string]string  `json:"signs,omitempty"`
	SignsRequired            int                `json:"signs_required,omitempty"`
	Title                    string             `json:"title,omitempty"`
	Translate                []string           `json:"translate,omitempty"`
	Viewport                 string             `json:"viewport,omitempty"`
	ZeronetVersion           string             `json:"zeronet_version,omitempty"`
}

func (s *Site) Info() (*Info, error) {
	content := new(Content)
	if err := s.DecodeJSON("content.json", content); err != nil {
		return nil, err
	}

	return &Info{
		Address:        s.addr,
		AddressHash:    addressHash(s.addr),
		AddressShort:   addressShort(s.addr),
		AuthAddress:    s.user.AuthAddress(s.addr),
		CertUserID:     s.user.CertUserID(s.addr),
		Peers:          len(s.peers),
		SizeLimit:      sizeLimit(s.Settings.SizeLimit),
		NextSizeLimit:  nextSizeLimit(s.Settings.Size),
		ContentUpdated: float64(content.Modified),
		Content:        content,
		Settings:       s.Settings,
	}, nil
}

func addressShort(addr string) string {
	return addr[:6] + ".." + addr[len(addr)-4:]
}

func addressHash(addr string) string {
	h := sha256.New()
	_, _ = io.WriteString(h, addr)
	return hex.EncodeToString(h.Sum(nil))
}

func sizeLimit(size int) int {
	if size == 0 {
		return config.SizeLimit
	}
	return size
}

// https://github.com/HelloZeroNet/ZeroNet/blob/454c0b2e7e000fda7000cba49027541fbf327b96/src/Site/Site.py#L145
func nextSizeLimit(size int) int {
	sizeLimits := []int{10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 20000, 50000, 100000}
	for _, limit := range sizeLimits {
		if float64(size)*1.2 < float64(limit)*1024*1024 {
			return limit
		}
	}

	return 1e6
}
