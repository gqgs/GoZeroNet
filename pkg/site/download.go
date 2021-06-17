package site

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
	"github.com/gqgs/go-zeronet/pkg/lib/parser"
	"github.com/gqgs/go-zeronet/pkg/peer"
)

func (s *Site) Download() error {
	for _, p := range s.peers {
		if err := p.Connect(); err != nil {
			s.log.WithField("peer", p).Warn(err)
			continue
		}
		defer p.Close()
		if err := s.DownloadContentJSON(p, "content.json"); err != nil {
			s.log.WithField("peer", p).Error(err)
			continue
		}
		return nil
	}

	return errors.New("could not download site")
}

func (s *Site) DownloadContentJSON(peer peer.Peer, innerPath string) error {
	resp, err := fileserver.GetFile(peer, s.addr, innerPath, 0, 0)
	if err != nil {
		return err
	}
	content := new(Content)
	if err := json.Unmarshal(resp.Body, content); err != nil {
		return err
	}

	contentPath := path.Join(config.DataDir, s.addr, content.InnerPath)
	if err := os.MkdirAll(path.Dir(contentPath), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(contentPath, resp.Body, os.ModePerm); err != nil {
		return err
	}

	for filename := range content.Files {
		resp, err := fileserver.GetFile(peer, s.addr, filename, 0, 0)
		if err != nil {
			return err
		}

		filePath := path.Join(config.DataDir, s.addr, filename)
		if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(filePath, resp.Body, os.ModePerm); err != nil {
			return err
		}
	}

	for includes := range content.Includes {
		if err := s.DownloadContentJSON(peer, includes); err != nil {
			return err
		}
	}

	// TODO:
	// Validate downloaded file
	// Bitcoin signature + SHA512_64 of downloaded files

	return nil
}

func (c *Content) isValid() bool {
	if c == nil {
		return false
	}

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
