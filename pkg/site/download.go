package site

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
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

	if content.InnerPath != innerPath {
		return fmt.Errorf("invalid content.json inner path: %s", content.InnerPath)
	}

	if !content.isValid() {
		return fmt.Errorf("invalid content.json: %s", content.InnerPath)
	}

	contentPath := path.Join(config.DataDir, s.addr, content.InnerPath)
	if err := os.MkdirAll(path.Dir(contentPath), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(contentPath, resp.Body, os.ModePerm); err != nil {
		return err
	}

	for filename, file := range content.Files {
		var body []byte
		var location int
		for {
			resp, err := fileserver.GetFile(peer, s.addr, filename, location, file.Size)
			if err != nil {
				return err
			}
			body = append(body, resp.Body...)
			if len(body) >= file.Size {
				break
			}
			location = resp.Location
		}

		digest := sha512.Sum512(body)
		hexDigest := hex.EncodeToString(digest[:32])
		if hexDigest != file.Sha512 {
			s.log.Warnf("ignoring file with invalid hash. want: %s (%d), got: %s (%d)",
				file.Sha512, file.Size, hexDigest, len(body))
			continue
		}

		filePath := path.Join(config.DataDir, s.addr, filename)
		if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(filePath, body, os.ModePerm); err != nil {
			return err
		}
	}

	for includes := range content.Includes {
		if err := s.DownloadContentJSON(peer, includes); err != nil {
			return err
		}
	}

	return nil
}
