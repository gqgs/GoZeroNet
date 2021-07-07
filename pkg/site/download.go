package site

import (
	"bytes"
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
	"github.com/gqgs/go-zeronet/pkg/lib/bigfile"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/peer"
)

func (s *Site) Download(since time.Time) error {
	for {
		p, err := s.peerManager.GetConnected(s.ctx)
		if err != nil {
			return err
		}
		err = s.DownloadContentJSON(p, "content.json")
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
		p, err := s.peerManager.GetConnected(s.ctx)
		if err != nil {
			return err
		}
		err = s.downloadRecent(p, since)
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

	pieceMap := func(innerPath, pieceMap string) string {
		if pieceMap == "" {
			return ""
		}
		return safe.CleanPath(path.Join(path.Dir(innerPath), pieceMap))
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

		info.InnerPath = relPath
		info.Hash = file.Sha512
		info.Size = file.Size
		info.PieceSize = file.PieceSize
		info.Piecemap = pieceMap(innerPath, file.Piecemap)

		event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, &event.FileInfo{
			InnerPath:  info.InnerPath,
			Hash:       info.Hash,
			Size:       info.Size,
			IsPinned:   info.IsPinned,
			IsOptional: info.IsOptional,
			PieceSize:  info.PieceSize,
			Piecemap:   info.Piecemap,
			Downloaded: info.Downloaded,
		})

		if err := s.downloadFile(peer, info); err != nil {
			return err
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
			InnerPath:  relPath,
			Hash:       file.Sha512,
			Size:       file.Size,
			IsPinned:   info.IsPinned,
			IsOptional: true,
			PieceSize:  file.PieceSize,
			Piecemap:   pieceMap(innerPath, file.Piecemap),
			Downloaded: info.Downloaded,
		})
	}

	for includes := range content.Includes {
		if err := s.DownloadContentJSON(peer, includes); err != nil {
			return err
		}
	}

	return nil
}

func (s *Site) downloadChunks(peer peer.Peer, info *event.FileInfo) error {
	if !info.IsBigFile() {
		resp, err := fileserver.StreamFileFull(peer, s.addr, info.InnerPath, info.Size)
		if err != nil {
			return err
		}
		if err := s.verifyDownload(resp.Body, info.Size, info.Hash); err != nil {
			return err
		}

		filePath := path.Join(config.DataDir, s.addr, info.InnerPath)
		if err := os.WriteFile(filePath, resp.Body, os.ModePerm); err != nil {
			return err
		}

		info.Downloaded = len(resp.Body)
		event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, info)
		return nil
	}

	// Parse and verify piece info

	pieceInfo, err := s.contentDB.FileInfo(s.addr, info.Piecemap)
	if err != nil {
		return err
	}

	var pieceBody []byte
	if pieceInfo.IsDownloaded {
		pieceBody, _ = os.ReadFile(path.Join(config.DataDir, s.addr, pieceInfo.InnerPath))
	}

	if len(pieceBody) == 0 {
		pieceResp, err := fileserver.StreamFileFull(peer, s.addr, pieceInfo.InnerPath, pieceInfo.Size)
		if err != nil {
			return err
		}
		pieceBody = pieceResp.Body

		if err := s.verifyDownload(pieceBody, pieceInfo.Size, pieceInfo.Hash); err != nil {
			return err
		}

		pieceInfo.Downloaded = len(pieceBody)
		event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, pieceInfo)

		filePath := path.Join(config.DataDir, s.addr, pieceInfo.InnerPath)
		if err := os.WriteFile(filePath, pieceBody, os.ModePerm); err != nil {
			return err
		}
	}

	pieceMap, err := bigfile.ParsePieceMap(bytes.NewReader(pieceBody))
	if err != nil {
		return err
	}

	hashes, err := pieceMap.Hashes(path.Base(info.InnerPath))
	if err != nil {
		return err
	}

	// Create disk file

	filePath := path.Join(config.DataDir, s.addr, info.InnerPath)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() != int64(info.Size) {
		if err := file.Truncate(int64(info.Size)); err != nil {
			return err
		}
	}

	// Download and verify file pieces

	piecemap := bigfile.UnpackPieceField(s.Settings.Cache.Piecefields[info.Hash])
	if len(piecemap) != len(hashes) {
		piecemap = strings.Repeat("0", len(hashes))
	}

	s.Settings.Cache.pieceFieldsMutex.Lock()
	if s.Settings.Cache.Piecefields == nil {
		s.Settings.Cache.Piecefields = make(map[string]bigfile.PieceField)
	}
	s.Settings.Cache.pieceFieldsMutex.Unlock()

	for i, hash := range hashes {
		if piecemap[i] == '1' {
			continue
		}

		resp, err := fileserver.StreamAtMost(peer, s.addr, info.InnerPath, i*info.PieceSize, info.PieceSize, info.Size)
		if err != nil {
			return err
		}
		if err := s.verifyDownload(resp.Body, len(resp.Body), hash); err != nil {
			return err
		}

		if _, err := file.WriteAt(resp.Body, int64(i*info.PieceSize)); err != nil {
			return err
		}

		piecemap = piecemap[:i] + "1" + piecemap[i+1:]
		s.Settings.Cache.pieceFieldsMutex.Lock()
		s.Settings.Cache.Piecefields[info.Hash] = bigfile.PackPieceField(piecemap)
		s.Settings.Cache.pieceFieldsMutex.Unlock()

		info.Downloaded += len(resp.Body)
		event.BroadcastFileInfoUpdate(s.addr, s.pubsubManager, info)
	}
	return nil
}

func (s *Site) verifyDownload(body []byte, size int, hash string) error {
	if len(body) != size {
		return fmt.Errorf("file with invalid size. want: (%d), got: (%d)", size, len(body))
	}

	digest := sha512.Sum512(body)
	hexDigest := hex.EncodeToString(digest[:32])
	if hexDigest != hash {
		return fmt.Errorf("file with invalid hash. want: %s (%d), got: %s (%d)", hash, size, hexDigest, len(body))
	}

	return nil
}

func (s *Site) downloadFile(peer peer.Peer, info *event.FileInfo) error {
	filePath := path.Join(config.DataDir, s.addr, info.InnerPath)
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	if err := s.downloadChunks(peer, info); err != nil {
		event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: -1})
		return err
	}

	s.log.WithField("inner_path", info.InnerPath).Info("downloaded file!")

	event.BroadcastPeerInfoUpdate(s.addr, s.pubsubManager, &event.PeerInfo{Address: peer.String(), ReputationDelta: 1})
	s.BroadcastSiteChange("file_done", info.InnerPath)

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
