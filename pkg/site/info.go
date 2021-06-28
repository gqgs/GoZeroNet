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
	Event          []interface{} `json:"event,omitempty"`
	CertUserID     string        `json:"cert_user_id"`
	Peers          int           `json:"peers"`
	NextSizeLimit  int           `json:"next_size_limit"`
	SizeLimit      int           `json:"size_limit"`
	Workers        int           `json:"workers"`
	ContentUpdated float64       `json:"content_updated"`
	Content        *Content      `json:"content"`
	Settings       *Settings     `json:"settings"`
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
		ContentUpdated: content.Modified,
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
		return config.SiteSizeLimit
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
