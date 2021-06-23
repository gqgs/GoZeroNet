package site

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
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

func (s *Site) Download(peerManager peer.Manager, since time.Time) error {
	for {
		p, err := peerManager.GetConnected()
		if err != nil {
			return fmt.Errorf("could not download files: %s", err)
		}

		err = s.DownloadContentJSON(p, "content.json")
		peerManager.PutConnected(p)
		if err != nil {
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
		return s.DownloadSince(peerManager, since)
	}
}

func (s *Site) DownloadSince(peerManager peer.Manager, since time.Time) error {
	for {
		p, err := peerManager.GetConnected()
		if err != nil {
			return fmt.Errorf("could not download files: %s", err)
		}

		err = s.downloadRecent(p, since)
		peerManager.PutConnected(p)
		if err != nil {
			s.log.WithField("peer", p).Warn(err)
			continue
		}

		return s.SaveSettings()
	}
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
		s.log.Debugf("done: %s", innerPath)
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
		event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
		return fmt.Errorf("invalid content.json inner path: %s", content.InnerPath)
	}

	if !content.isValid() {
		event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
		return fmt.Errorf("invalid content.json: %s", content.InnerPath)
	}

	event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: 1})

	contentPath := path.Join(config.DataDir, s.addr, safe.CleanPath(content.InnerPath))
	file, err := os.Open(contentPath)
	if err == nil {
		defer file.Close()
		currentContent := new(Content)
		if err := json.NewDecoder(file).Decode(currentContent); err == nil {
			if content.Modified <= currentContent.Modified {
				s.log.Debugf("outdated %s, skipping...", contentPath)
				return nil
			}
		}
	}

	if innerPath == "content.json" {
		s.Settings.Modified = int64(content.Modified)
	}

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
			event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
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

		event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: 1})
		event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, &event.FileInfo{
			InnerPath:    relPath,
			Hash:         hexDigest,
			Size:         len(body),
			IsDownloaded: true,
		})
	}

	for filename, file := range content.FilesOptional {
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

		if len(file.Sha512) != 64 {
			s.log.WithField("peer", peer).Errorf("invalid hash id length: %d", len(file.Sha512))
			event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
			continue
		}

		event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, &event.FileInfo{
			InnerPath:    relPath,
			Hash:         file.Sha512,
			Size:         file.Size,
			IsDownloaded: info.IsDownloaded,
		})
	}

	for includes := range content.Includes {
		if err := s.DownloadContentJSON(peer, includes); err != nil {
			return err
		}
	}

	return nil
}
