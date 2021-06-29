package site

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/peer"
)

func (s *Site) Download(since time.Time) error {
	for {
		p := s.peerManager.GetConnected()
		err := s.DownloadContentJSON(p, "content.json")
		s.peerManager.PutConnected(p)
		if err != nil {
			s.log.WithField("peer", p).Error(err)
			continue
		}

		s.Settings.Downloaded = time.Now().Unix()
		s.Settings.Peers = len(s.peers)
		s.Settings.Serving = true

		if err := s.SaveSettings(); err != nil {
			return err
		}
		return s.DownloadSince(since)
	}
}

func (s *Site) DownloadSince(since time.Time) error {
	for {
		p := s.peerManager.GetConnected()
		err := s.downloadRecent(p, since)
		s.peerManager.PutConnected(p)
		if err != nil {
			s.log.WithField("peer", p).Warn(err)
			continue
		}

		return s.SaveSettings()
	}
}

func (s *Site) downloadRecent(peer peer.Peer, since time.Time) error {
	resp, err := fileserver.ListModified(peer, s.addr, since.Unix())
	if err != nil {
		return err
	}

	updated := make([]string, 0, len(resp.ModifiedFiles))
	for innerPath, modified := range resp.ModifiedFiles {
		if info, err := s.contentDB.ContentInfo(s.addr, innerPath); err == nil {
			if modified <= info.Modified {
				s.log.WithField("peer", peer).Debugf("skipping outdated or same %s: (%d <= %d)", innerPath, modified, info.Modified)
				continue
			}
		}
		updated = append(updated, innerPath)
	}

	for _, innerPath := range updated {
		if err := peer.CheckConnection(); err != nil {
			return err
		}
		if err := s.DownloadContentJSON(peer, innerPath); err != nil {
			s.log.WithField("peer", peer).Error(err)
			continue
		}
		s.log.Debugf("done: %s", innerPath)
	}
	return nil
}

func (s *Site) DownloadContentJSON(peer peer.Peer, innerPath string) error {
	resp, err := fileserver.GetFileFull(peer, s.addr, innerPath, 0)
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

	event.BroadcastContentInfoUpdate(s.addr, s.pubsubManager, &event.ContentInfo{
		InnerPath: innerPath,
		Modified:  int(content.Modified),
		Size:      len(resp.Body),
	})

	if innerPath == "content.json" {
		s.Settings.Modified = int64(content.Modified)
	}

	if err := os.MkdirAll(path.Dir(contentPath), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(contentPath, resp.Body, os.ModePerm); err != nil {
		return err
	}

	logger := s.log.WithField("peer", peer)
	for _, filename := range sortDownloads(content.Files) {
		file := content.Files[filename]
		filename = path.Join(path.Dir(innerPath), filename)
		relPath := safe.CleanPath(filename)

		info, err := s.contentDB.FileInfo(s.addr, relPath)
		switch {
		case errors.Is(err, database.ErrFileNotFound):
		case err == nil:
		default:
			logger.Error(err)
			continue
		}

		if info.IsDownloaded && info.Hash == file.Sha512 {
			continue
		}

		info.Hash = file.Sha512
		info.Size = file.Size

		if err := s.downloadFile(peer, filename, info); err != nil {
			logger.Error(err)
		}
	}

	for _, filename := range sortDownloads(content.FilesOptional) {
		file := content.FilesOptional[filename]
		filename = path.Join(path.Dir(innerPath), filename)
		relPath := safe.CleanPath(filename)

		info, err := s.contentDB.FileInfo(s.addr, relPath)
		switch {
		case errors.Is(err, database.ErrFileNotFound):
		case err == nil:
		default:
			logger.Error(err)
			continue
		}

		if len(file.Sha512) != 64 {
			logger.Errorf("invalid hash id length: %d", len(file.Sha512))
			event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
			continue
		}

		event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, &event.FileInfo{
			InnerPath:    relPath,
			Hash:         file.Sha512,
			Size:         file.Size,
			IsDownloaded: info.IsDownloaded,
			IsPinned:     info.IsPinned,
			IsOptional:   true,
		})
	}

	for includes := range content.Includes {
		if err := s.DownloadContentJSON(peer, includes); err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) downloadFile(peer peer.Peer, innerPath string, info *event.FileInfo) error {
	resp, err := fileserver.GetFileFull(peer, s.addr, innerPath, info.Size)
	if err != nil {
		return err
	}
	body := resp.Body
	if len(body) != info.Size {
		event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
		return fmt.Errorf("ignoring file (%s) with invalid size. want: (%d), got: (%d)",
			innerPath, info.Size, len(body))
	}

	digest := sha512.Sum512(body)
	hexDigest := hex.EncodeToString(digest[:32])
	if hexDigest != info.Hash {
		event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
		return fmt.Errorf("ignoring file (%s) with invalid hash. want: %s (%d), got: %s (%d)",
			innerPath, info.Hash, info.Size, hexDigest, len(body))
	}
	s.Settings.BytesRecv += info.Size

	filePath := path.Join(config.DataDir, s.addr, innerPath)
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(filePath, body, os.ModePerm); err != nil {
		return err
	}

	s.log.WithField("inner_path", innerPath).Info("downloaded file!")

	event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: 1})
	event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, &event.FileInfo{
		InnerPath:    innerPath,
		Hash:         hexDigest,
		Size:         len(body),
		IsDownloaded: true,
		IsPinned:     info.IsPinned,
		IsOptional:   info.IsOptional,
	})

	return nil
}

func sortDownloads(files map[string]File) []string {
	var i int
	filenames := make([]string, len(files))
	for filename := range files {
		filenames[i] = filename
		i++
	}

	scoreFunc := func(filename string) int {
		suffixes := []string{
			"dbschema.json",
			"index.html",
			".css",
			".js",
			".zip",
			".png",
			".json",
		}

		var priority int
		for _, suffix := range suffixes {
			if strings.HasSuffix(filename, suffix) {
				return priority
			}
			priority += 5
		}
		return priority
	}

	sort.Slice(filenames, func(i, j int) bool {
		return scoreFunc(filenames[i]) < scoreFunc(filenames[j])
	})
	return filenames
}
