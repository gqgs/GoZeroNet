package site

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/gqgs/go-zeronet/pkg/user"
)

type Info struct {
	Address        string    `json:"address"`
	AddressHash    string    `json:"address_hash"`
	AddressShort   string    `json:"address_short"`
	AuthAddress    string    `json:"auth_address"`
	BadFiles       int       `json:"bad_files"`
	CertUserID     string    `json:"cert_user_id"`
	Peers          int       `json:"peers"`
	NextSizeLimit  int       `json:"next_size_limit"`
	SizeLimit      int       `json:"size_limit"`
	Workers        int       `json:"workers"`
	ContentUpdated float64   `json:"content_updated"`
	Content        *Content  `json:"content"`
	Settings       *Settings `json:"settings"`
}

type File struct {
	Sha512 string `json:"sha512"`
	Size   int    `json:"size"`
}

type Include struct {
	Signers         []string `json:"signers"`
	SignersRequired int      `json:"signers_required"`
}

type Content struct {
	Address                  string             `json:"address"`
	AddressIndex             int                `json:"address_index"`
	BackgroundColor          string             `json:"background-color"`
	BackgroundColorDark      string             `json:"background-color-dark"`
	CloneRoot                string             `json:"clone_root"`
	Cloneable                bool               `json:"cloneable"`
	ClonedFrom               string             `json:"cloned_from"`
	Description              string             `json:"description"`
	Favicon                  string             `json:"favicon"`
	Files                    map[string]File    `json:"files"`
	FilesOptional            map[string]File    `json:"files_optional"`
	Ignore                   string             `json:"ignore"`
	Includes                 map[string]Include `json:"includes"`
	InnerPath                string             `json:"inner_path"`
	Modified                 int                `json:"modified"`
	Optional                 string             `json:"optional"`
	PostmessageNonceSecurity bool               `json:"postmessage_nonce_security"`
	SignersSign              string             `json:"signers_sign"`
	Signs                    map[string]string  `json:"signs"`
	SignsRequired            int                `json:"signs_required"`
	Title                    string             `json:"title"`
	Translate                []string           `json:"translate"`
	Viewport                 string             `json:"viewport"`
	ZeronetVersion           string             `json:"zeronet_version"`
}

func (s *Site) Info(user user.User) (*Info, error) {
	content := new(Content)
	if err := s.DecodeJSON("content.json", content); err != nil {
		return nil, err
	}

	return &Info{
		Address:        s.addr,
		AddressHash:    addressHash(s.addr),
		AddressShort:   addressShort(s.addr),
		AuthAddress:    user.AuthAddress(s.addr),
		CertUserID:     user.CertUserID(s.addr),
		Peers:          len(s.peers),
		SizeLimit:      s.Settings.SizeLimit,
		NextSizeLimit:  s.Settings.SizeLimit * 2,
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
