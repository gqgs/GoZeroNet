package site

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
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

		s.Settings.Downloaded = time.Now().Unix()
		s.Settings.Peers = len(s.peers)
		s.Settings.Serving = true
		s.Settings.AjaxKey = random.HexString(64)
		s.Settings.AuthKey = random.HexString(64)
		s.Settings.WrapperKey = random.HexString(64)

		if err := s.SaveSettings(); err != nil {
			return err
		}
		return s.downloadRecent(p, time.Now().AddDate(0, 0, -7))
	}

	return errors.New("could not download site")
}

func (s *Site) DownloadSince(since time.Time) error {
	for _, p := range s.peers {
		if err := p.Connect(); err != nil {
			s.log.WithField("peer", p).Warn(err)
			continue
		}
		defer p.Close()
		if err := s.downloadRecent(p, since); err != nil {
			s.log.WithField("peer", p).Error(err)
			continue
		}
		return nil
	}

	return errors.New("could not files")
}

func (s *Site) downloadRecent(peer peer.Peer, since time.Time) error {
	downloaded := make(map[string]struct{})
	resp, err := fileserver.ListModified(peer, s.addr, int(since.Unix()))
	if err != nil {
		return err
	}

	for content := range resp.ModifiedFiles {
		if err := s.DownloadContentJSON(peer, content); err != nil {
			s.log.WithField("peer", peer).Error(err)
			continue
		}
		s.log.Debugf("downloaded: %s", content)
		downloaded[content] = struct{}{}
	}

	s.log.Infof("downloaded %d files", len(downloaded))
	return nil
}

func (s *Site) DownloadContentJSON(peer peer.Peer, innerPath string) error {
	resp, err := fileserver.GetFileFull(peer, s.addr, innerPath)
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

	if innerPath == "content.json" {
		s.Settings.Modified = int64(content.Modified)
	}

	contentPath := path.Join(config.DataDir, s.addr, safe.CleanPath(content.InnerPath))
	if err := os.MkdirAll(path.Dir(contentPath), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(contentPath, resp.Body, os.ModePerm); err != nil {
		return err
	}

	for filename, file := range content.Files {
		filename = path.Join(path.Dir(innerPath), filename)
		resp, err := fileserver.GetFileFull(peer, s.addr, filename)
		if err != nil {
			s.log.WithField("peer", peer).Warn(err)
			continue
		}
		body := resp.Body
		digest := sha512.Sum512(body)
		hexDigest := hex.EncodeToString(digest[:32])
		if hexDigest != file.Sha512 {
			s.log.Warnf("ignoring file with invalid hash. want: %s (%d), got: %s (%d)",
				file.Sha512, file.Size, hexDigest, len(body))
			continue
		}
		s.Settings.BytesRecv += file.Size

		filePath := path.Join(config.DataDir, s.addr, safe.CleanPath(filename))
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
