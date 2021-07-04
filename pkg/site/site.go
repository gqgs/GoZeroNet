package site

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/parser"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/peer"
	"github.com/gqgs/go-zeronet/pkg/user"
)

var errFileNotFound = fmt.Errorf("site: %w", os.ErrNotExist)

type Site struct {
	ctx                      context.Context
	addr                     string
	trackersMutex            sync.RWMutex
	trackers                 map[string]*AnnouncerStats
	peersMutex               sync.RWMutex
	peers                    map[string]peer.Peer
	pubsubManager            pubsub.Manager
	Settings                 *Settings
	user                     *user.User
	wrapperNonceMutex        sync.RWMutex
	wrapperNonce             map[string]int64
	log                      log.Logger
	db                       database.SiteDatabase
	contentDB                database.ContentDatabase
	peerManager              peer.Manager
	workerManager            Worker
	lastAnnounce             time.Time
	loading                  bool
	postmessageNonceSecurity bool
}

func (s *Site) Peers() map[string]peer.Peer {
	return s.peers
}

func (s *Site) Loading(loading bool) {
	s.loading = loading
}

func (s *Site) IsLoading() bool {
	return s.loading
}

func (s *Site) PostmessageNonceSecurity() bool {
	return s.postmessageNonceSecurity
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

	encoded, err := json.MarshalIndent(settings, "", " ")
	if err != nil {
		return err
	}

	path := path.Join(config.DataDir, "sites.json")
	return ioutil.WriteFile(path, encoded, fs.ModePerm)
}

func (s *Site) IsAdmin() bool {
	return s.addr == config.HomeSite
}

func (s *Site) BroadcastSiteChange(events ...interface{}) {
	info, err := s.Info()
	if err != nil {
		s.log.Error(err)
		return
	}

	info.Event = events

	event.BroadcastSiteChanged(s.addr, s.pubsubManager, &event.SiteChanged{
		Cmd:    "setSiteInfo",
		Params: info,
	})
}

func (s *Site) SetSiteLimit(sizeLimit int) error {
	s.Settings.SizeLimit = sizeLimit
	if err := s.SaveSettings(); err != nil {
		return err
	}

	s.BroadcastSiteChange()

	return s.DownloadSince(time.Now().AddDate(0, -1, 0))
}

func (s *Site) Address() string {
	return s.addr
}

func (s *Site) User() *user.User {
	return s.user
}

func (s *Site) DecodeJSON(filename string, v interface{}) error {
	innerPath := path.Join(config.DataDir, s.addr, safe.CleanPath(filename))
	file, err := os.Open(innerPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func (s *Site) FileNeed(innerPath string) error {
	path := path.Join(config.DataDir, s.addr, safe.CleanPath(innerPath))
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			event.BroadcastFileNeed(s.addr, s.pubsubManager, &event.FileNeed{InnerPath: innerPath})
			return nil
		}
		return err
	}
	return nil
}

func (s *Site) FileWrite(innerPath string, reader io.Reader) error {
	innerPath = safe.CleanPath(innerPath)
	writePath := path.Join(config.DataDir, s.addr, innerPath)
	if err := os.MkdirAll(path.Dir(writePath), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(writePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func (s *Site) ReadFile(ctx context.Context, innerPath string, dst io.Writer) error {
	innerPath = safe.CleanPath(innerPath)
	path := path.Join(config.DataDir, s.addr, innerPath)
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s.broadcastFileNeed(ctx, innerPath, func() error {
				return s.ReadFile(ctx, innerPath, dst)
			})
		}
		return err
	}
	defer file.Close()

	_, err = io.Copy(dst, file)
	return err
}

func (s *Site) broadcastFileNeed(ctx context.Context, innerPath string, downloadCallback func() error) error {
	msgCh := s.pubsubManager.Register("file_need("+innerPath+")", config.DefaultChannelSize)
	defer s.pubsubManager.Unregister(msgCh)
	event.BroadcastFileNeed(s.addr, s.pubsubManager, &event.FileNeed{InnerPath: innerPath})
	for {
		select {
		case <-ctx.Done():
			return errFileNotFound
		case msg := <-msgCh:
			if msg.Site() != s.addr {
				continue
			}
			if updated, ok := msg.Event().(*event.FileInfo); ok {
				if updated.InnerPath == innerPath && updated.IsDownloaded {
					return downloadCallback()
				}
			}
		}
	}
}

func (s *Site) ListFiles(innerPath string) ([]string, error) {
	files := make([]string, 0)
	root := path.Join(config.DataDir, s.addr, safe.CleanPath(innerPath))

	if _, err := os.Stat(root); err != nil {
		if errors.Is(err, os.ErrNotExist) {
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

func (s *Site) FileDelete(innerPath string) error {
	filePath := path.Join(config.DataDir, s.addr, safe.CleanPath(innerPath))
	return os.Remove(filePath)
}

func (s *Site) Update(daysAgo int) error {
	now := time.Now().UTC()
	s.BroadcastSiteChange("updating", true)
	defer s.BroadcastSiteChange("updated", true)

	if err := s.DownloadSince(now.AddDate(0, 0, -daysAgo)); err != nil {
		return err
	}
	if err := s.OpenDB(); err != nil {
		return err
	}
	return s.UpdateDB(now)
}

// Verify verifies a content.json file.
// It returns nil if the file is valid.
func (s *Site) Verify(innerPath string) error {
	s.log.WithField("inner_path", innerPath).Debug("verifying file")

	if !strings.HasSuffix(innerPath, "content.json") {
		return errors.New("can only verifiy content.json files")
	}
	path := path.Join(config.DataDir, s.addr, safe.CleanPath(innerPath))
	contentFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer contentFile.Close()

	c := new(Content)
	if err := json.NewDecoder(contentFile).Decode(c); err != nil {
		return err
	}

	if c.isValid() {
		return nil
	}

	return errors.New("content file is invalid")
}

// Sign signs a content.json file.
func (s *Site) Sign(innerPath, privateKey string, user *user.User) error {
	s.log.WithField("inner_path", innerPath).Debug("signing file")

	if !strings.HasSuffix(innerPath, "content.json") {
		return errors.New("can only verifiy content.json files")
	}
	filePath := path.Join(config.DataDir, s.addr, safe.CleanPath(innerPath))
	contentFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer contentFile.Close()

	c := new(Content)
	if err := json.NewDecoder(contentFile).Decode(c); err != nil {
		return err
	}

	if innerPath == "content.json" {
		signers := make([]string, len(c.Signs))
		var i int
		for addr := range c.Signs {
			signers[i] = addr
			i++
		}
		signerdMsg := fmt.Sprintf("%d:%s", c.SignsRequired, strings.Join(signers, ","))
		c.SignersSign, err = crypto.Sign([]byte(signerdMsg), privateKey)
		if err != nil {
			return err
		}
	}

	files := make(map[string]File)
	filesOptional := make(map[string]File)

	ignoreRegex, err := regexp.Compile(c.Ignore)
	if err != nil {
		return err
	}

	optionalRegex, err := regexp.Compile(c.Optional)
	if err != nil {
		return err
	}

	root := path.Join(config.DataDir, s.addr)

	err = filepath.WalkDir(filepath.Dir(filePath), func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relativePath := strings.TrimPrefix(path, root)
		relativePath = strings.TrimLeft(relativePath, "/")

		if strings.HasPrefix(relativePath, ".") || ignoreRegex.MatchString(relativePath) {
			return nil
		}

		fileData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hasher := sha512.New()
		if _, err := hasher.Write(fileData); err != nil {
			return err
		}

		// TODO: handle bigfiles
		file := File{
			Size:   len(fileData),
			Sha512: hex.EncodeToString(hasher.Sum(nil))[:64],
		}

		if optionalRegex.MatchString(relativePath) {
			filesOptional[relativePath] = file
		} else {
			files[relativePath] = file
		}

		return nil
	})
	if err != nil {
		return err
	}

	signs := c.Signs

	c.Signs = nil
	c.Modified = float64(time.Now().Unix())
	c.Files = files
	c.FilesOptional = filesOptional
	c.Address = s.addr
	c.InnerPath = innerPath

	if user != nil {
		c.CertAuthType = user.Certs[s.addr].AuthType
		c.CertSign = user.Certs[s.addr].CertSign
		c.CertUserID = user.CertUserID(s.addr)
	}

	contentJSON, err := json.Marshal(c)
	if err != nil {
		return err
	}

	contentJSON, err = parser.FixJSONSpacing(bytes.NewReader(contentJSON))
	if err != nil {
		return err
	}

	sign, err := crypto.Sign(contentJSON, privateKey)
	if err != nil {
		return err
	}
	address, err := crypto.PrivateKeyToAddress(privateKey)
	if err != nil {
		return err
	}

	if signs == nil {
		signs = make(map[string]string)
	}

	signs[address] = sign

	c.Signs = signs

	content, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, content, os.ModePerm)
}

func (s *Site) Publish(innerPath string) error {
	published := make(map[string]struct{})

	if !strings.HasSuffix(innerPath, "content.json") {
		return errors.New("can only update content.json files")
	}

	for i := 0; i < 10; i++ {
		connected, err := s.peerManager.GetConnected(s.ctx)
		if err != nil {
			return err
		}

		if _, alreadyPublished := published[connected.ID()]; alreadyPublished {
			continue
		}

		logger := s.log.WithField("peer", connected)
		logger.Info("publishing update")
		resp, err := fileserver.Update(connected, s.addr, innerPath)
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Info("updated published ", resp.Ok)
		published[connected.ID()] = struct{}{}
	}
	return nil
}
