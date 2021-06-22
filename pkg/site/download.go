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
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
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
		if err := s.downloadRecent(p, time.Now().AddDate(0, 0, -7)); err != nil {
			return err
		}
		return s.SaveSettings()
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
		return s.SaveSettings()
	}

	return errors.New("could not files")
}

func (s *Site) downloadRecent(peer peer.Peer, since time.Time) error {
	resp, err := fileserver.ListModified(peer, s.addr, int(since.Unix()))
	if err != nil {
		return err
	}

	for innerPath := range resp.ModifiedFiles {
		if err := s.DownloadContentJSON(peer, innerPath); err != nil {
			s.log.WithField("peer", peer).Error(err)
			continue
		}
		s.log.Debugf("downloaded: %s", innerPath)
	}
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
		relPath := safe.CleanPath(filename)

		info, err := s.contentDB.FileInfo(s.addr, relPath)
		switch err {
		case database.ErrFileNotFound:
		case nil:
		default:
			s.log.WithField("peer", peer).Error(err)
			continue
		}

		if info.IsDownloaded {
			continue
		}
		resp, err := fileserver.GetFileFull(peer, s.addr, relPath)
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

		filePath := path.Join(config.DataDir, s.addr, relPath)
		if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(filePath, body, os.ModePerm); err != nil {
			return err
		}

		fileDone, _ := json.Marshal(&event.FileInfo{
			InnerPath:    relPath,
			Hash:         hexDigest,
			Size:         len(body),
			IsDownloaded: true,
		})
		s.pubsubManager.Broadcast(s.addr, "file-done", fileDone)
	}

	for includes := range content.Includes {
		if err := s.DownloadContentJSON(peer, includes); err != nil {
			return err
		}
	}

	return nil
}
