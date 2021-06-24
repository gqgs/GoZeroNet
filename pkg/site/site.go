package site

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/peer"
	"github.com/gqgs/go-zeronet/pkg/user"
)

type Site struct {
	addr              string
	trackersMutex     sync.RWMutex
	trackers          map[string]*AnnouncerStats
	peersMutex        sync.RWMutex
	peers             map[string]peer.Peer
	pubsubManager     pubsub.Manager
	Settings          *Settings
	user              user.User
	isAdmin           bool
	wrapperNonceMutex sync.RWMutex
	wrapperNonce      map[string]int64
	log               log.Logger
	db                database.SiteDatabase
	contentDB         database.ContentDatabase
}

func (s *Site) Peers() map[string]peer.Peer {
	return s.peers
}

func (s *Site) HasValidWrapperNonce(wrapperNonce string) bool {
	s.wrapperNonceMutex.RLock()
	defer s.wrapperNonceMutex.RUnlock()
	created, exists := s.wrapperNonce[wrapperNonce]
	// Nonces are valid for 24 hours
	if exists && time.Since(time.Unix(created, 0)).Hours() > 24 {
		s.unRegisterWrapperNonce(wrapperNonce)
		return false
	}
	return exists
}

func (s *Site) unRegisterWrapperNonce(wrapperNonce string) {
	s.wrapperNonceMutex.Lock()
	delete(s.wrapperNonce, wrapperNonce)
	s.wrapperNonceMutex.Unlock()
}

func (s *Site) registerWrapperNonce(wrapperNonce string) {
	s.wrapperNonceMutex.Lock()
	s.wrapperNonce[wrapperNonce] = time.Now().Unix()
	s.wrapperNonceMutex.Unlock()
}

func (s *Site) SaveSettings() error {
	settings, err := loadSiteSettingsFromFile()
	if err != nil {
		return err
	}

	settings[s.addr] = s.Settings

	encoded, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	path := path.Join(config.DataDir, "sites.json")
	return ioutil.WriteFile(path, encoded, fs.ModePerm)
}

func (s *Site) IsAdmin() bool {
	return s.isAdmin
}

func (s *Site) broadcastSiteChange(events ...interface{}) error {
	info, err := s.Info()
	if err != nil {
		return err
	}

	info.Events = events

	event.BroadcastSiteChanged(s.addr, s.pubsubManager, &event.SiteChanged{
		Cmd:    "setSiteInfo",
		Params: info,
	})

	return nil
}

func (s *Site) SetSiteLimit(sizeLimit int) error {
	s.Settings.SizeLimit = sizeLimit
	if err := s.SaveSettings(); err != nil {
		return err
	}

	if err := s.broadcastSiteChange(); err != nil {
		return err
	}

	// TODO: return s.Download()
	return nil
}

func (s *Site) Address() string {
	return s.addr
}

func (s *Site) User() user.User {
	return s.user
}

func (s *Site) DecodeJSON(filename string, v interface{}) error {
	innerPath := path.Join(config.DataDir, s.addr, safe.CleanPath(filename))
	file, err := os.Open(innerPath)
	if err != nil {
		// TODO: download file
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func (s *Site) ReadFile(ctx context.Context, innerPath string, dst io.Writer) error {
	path := path.Join(config.DataDir, s.addr, safe.CleanPath(innerPath))
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			msgCh := s.pubsubManager.Register(10)
			defer s.pubsubManager.Unregister(msgCh)
			event.BroadcastFileNeed(s.addr, s.pubsubManager, &event.FileNeed{InnerPath: innerPath})
			for {
				select {
				case <-ctx.Done():
					return err
				case msg := <-msgCh:
					if msg.Site() != s.addr {
						continue
					}
					if updated, ok := msg.Event().(*event.FileInfo); ok {
						if updated.InnerPath == innerPath && updated.IsDownloaded {
							s.log.WithField("inner_path", innerPath).Info("file downloaded!")
							return s.ReadFile(ctx, innerPath, dst)
						}
					}
				}
			}
		}
		return err
	}
	defer file.Close()

	_, err = io.Copy(dst, file)
	return err
}

func (s *Site) ListFiles(innerPath string) ([]string, error) {
	files := make([]string, 0)
	root := path.Join(config.DataDir, s.addr, safe.CleanPath(innerPath))

	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return files, nil
		}
		return nil, err
	}

	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.Type().IsRegular() {
			files = append(files, info.Name())
		}

		return nil
	})
	return files, err
}
