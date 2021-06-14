package site

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/content"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/template"
	"github.com/gqgs/go-zeronet/pkg/user"
)

type Site struct {
	addr          string
	trackers      map[string]*AnnouncerStats
	peers         map[string]struct{}
	pubsubManager pubsub.Manager
	Settings      *Settings
}

type Settings struct {
	Added     int    `json:"added"`
	AjaxKey   string `json:"ajax_key"`
	AuthKey   string `json:"auth_key"`
	BytesRecv int    `json:"bytes_recv"`
	BytesSent int    `json:"bytes_sent"`
	Cache     struct {
		BadFiles  map[string]int `json:"bad_files"`
		Hashfield string         `json:"hashfield"`
	} `json:"cache"`
	PieceFields               map[string]string `json:"piecefields"`
	HasBigFile                bool              `json:"has_bigfile"`
	Downloaded                int               `json:"downloaded"`
	Modified                  int               `json:"modified"`
	ModifiedFilesModification bool              `json:"modified_files_notification"`
	OptionalDownloaded        int               `json:"optional_downloaded"`
	OptionalHelp              map[string]string `json:"optional_help"`
	Own                       bool              `json:"own"`
	Peers                     int               `json:"peers"`
	Permissions               []string          `json:"permissions"`
	Serving                   bool              `json:"serving"`
	Size                      int               `json:"size"`
	SizeFilesOptional         int               `json:"size_files_optional"`
	SizeLimit                 int               `json:"size_limit"`
	SizeOptional              int               `json:"size_optional"`
	WrapperKey                string            `json:"wrapper_key"`
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

func (s *Site) Download() error {
	// TODO: implement me
	return nil
}

func (s *Site) SetSiteLimit(sizeLimit int, user user.User) error {
	s.Settings.SizeLimit = sizeLimit
	if err := s.SaveSettings(); err != nil {
		return err
	}

	info, err := s.Info(user)
	if err != nil {
		return err
	}

	event := SiteChangedEvent{
		Cmd:    "setSiteInfo",
		Params: info,
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	s.pubsubManager.Broadcast(s.addr, "siteChanged", encoded)

	return s.Download()
}

func (s *Site) Address() string {
	return s.addr
}

func (s Site) DecodeJSON(filename string, v interface{}) error {
	innerPath := path.Join(config.DataDir, s.addr, filename)
	file, err := os.Open(innerPath)
	if err != nil {
		// TODO: download file
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

func (s Site) ReadFile(innerPath string, dst io.Writer) error {
	path := path.Join(config.DataDir, s.addr, innerPath)
	file, err := os.Open(path)
	if err != nil {
		// TODO: download file
		return err
	}
	defer file.Close()

	_, err = io.Copy(dst, file)
	return err
}

type SiteManager interface {
	Site(addr string) *Site
	RenderIndex(site, indexFilename string, dst io.Writer) error
	ReadFile(site, innerPath string, dst io.Writer) error
	SiteByWrapperKey(wrapperKey string) *Site
}

type siteManager struct {
	// Address -> Site
	sites map[string]*Site
	// WrapperKey -> Site
	wrapperKeyMap map[string]*Site
}

func (m *siteManager) Site(addr string) *Site {
	return m.sites[addr]
}

func (m *siteManager) SiteByWrapperKey(wrapperKey string) *Site {
	return m.wrapperKeyMap[wrapperKey]
}

func (m *siteManager) RenderIndex(site, indexFilename string, dst io.Writer) error {
	s, ok := m.sites[site]
	if !ok {
		// TODO: download site
		return errors.New("site not found")
	}

	var siteContent content.Content
	if err := s.DecodeJSON("content.json", &siteContent); err != nil {
		return err
	}

	vars := struct {
		ServerURL                string
		InnerPath                string
		FileURL                  string // TODO: escape?
		FileInnerPath            string // TODO: escape?
		Address                  string
		Title                    string // TODO: escape?
		BodyStyle                string
		MetaTags                 string
		QueryString              string // TODO: escape?
		WrapperKey               string
		AjaxKey                  string
		WrapperNonce             string
		PostMessageNonceSecurity bool
		Permissions              []string
		ShowLoadingScreen        bool
		SandboxPermissions       string
		Rev                      int
		Lang                     string
		HomePage                 string
		ThemeClass               string
		ScriptNonce              string
	}{
		Address:       site,
		Title:         siteContent.Title,
		Rev:           config.Rev,
		Lang:          config.Language,
		FileURL:       "1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D/index.html",
		FileInnerPath: "index.html",
		Permissions:   []string{"ADMIN"},
		WrapperNonce:  "f9b6fc1fc24bd5e6ae7c3cd5761520466000d36e2e1f0f46d3d5c308a126bb56",
		WrapperKey:    "e02c32aa7bf2625c81808ff55d98b58f93b6fba8cbda0702033cdd8cd5463d27",
		// AjaxKey:                  "bcf959ce5ac90fa70e1ac2499b19de92e031aa9cd87c6ade6ca4a7ed91b7b002",
		// ScriptNonce:              "iiz9PAl7yqImqqntjJ67TuyWvdk8GMUJ3rHc2mOSc0OkddjqaOHxhOpKjJ9xIIUJ",
		QueryString:              "?wrapper_nonce=f9b6fc1fc24bd5e6ae7c3cd5761520466000d36e2e1f0f46d3d5c308a126bb56",
		PostMessageNonceSecurity: false,
		ShowLoadingScreen:        false,
		ThemeClass:               "theme-light",
		BodyStyle:                "background-color: #F2F4F6",
	}

	return template.Wrapper.ExecuteHTML(dst, vars)
}

func (m *siteManager) ReadFile(site, innerPath string, dst io.Writer) error {
	s, ok := m.sites[site]
	if !ok {
		// TODO: download site
		return errors.New("site not found")
	}

	return s.ReadFile(innerPath, dst)
}

func NewSiteManager(pubsubManager pubsub.Manager) (*siteManager, error) {
	settings, err := loadSiteSettingsFromFile()
	if err != nil {
		return nil, err
	}

	sites := make(map[string]*Site)
	wrapperKeyMap := make(map[string]*Site)
	for addr, siteSettings := range settings {
		site := new(Site)
		site.Settings = siteSettings
		site.addr = addr
		site.trackers = make(map[string]*AnnouncerStats)
		site.peers = make(map[string]struct{})
		wrapperKeyMap[siteSettings.WrapperKey] = site
		site.pubsubManager = pubsubManager
		sites[addr] = site
	}

	return &siteManager{
		sites:         sites,
		wrapperKeyMap: wrapperKeyMap,
	}, nil
}

func loadSiteSettingsFromFile() (map[string]*Settings, error) {
	sitesFile, err := os.Open(path.Join(config.DataDir, "sites.json"))
	if err != nil {
		// TODO: ignore error if file not found
		return nil, err
	}
	defer sitesFile.Close()

	settings := make(map[string]*Settings)
	return settings, json.NewDecoder(sitesFile).Decode(&settings)
}
