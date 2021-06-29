package site

import (
	"context"
	"errors"
	"io"
	"path"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/gqgs/go-zeronet/pkg/peer"
	"github.com/gqgs/go-zeronet/pkg/template"
	"github.com/gqgs/go-zeronet/pkg/user"
)

type Manager interface {
	Site(addr string) *Site
	RenderIndex(site, indexFilename string, dst io.Writer) error
	ReadFile(site, innerPath string, dst io.Writer) error
	SiteByWrapperKey(wrapperKey string) *Site
	SiteList() ([]*Info, error)
	NewSite(addr string) (*Site, error)
	Close()
}

type manager struct {
	ctx context.Context
	// Address -> Site
	sites map[string]*Site
	// WrapperKey -> Site
	wrapperKeyMap map[string]*Site

	pubsubManager pubsub.Manager
	userManager   user.Manager
	contentDB     database.ContentDatabase
}

func NewManager(ctx context.Context, pubsubManager pubsub.Manager, userManager user.Manager,
	contentDB database.ContentDatabase) (*manager, error) {
	settings, err := loadSiteSettingsFromFile()
	if err != nil {
		return nil, err
	}

	user := userManager.User()
	sites := make(map[string]*Site)
	wrapperKeyMap := make(map[string]*Site)
	for addr, siteSettings := range settings {
		site := new(Site)
		site.Settings = siteSettings
		site.addr = addr
		site.trackers = make(map[string]*AnnouncerStats)
		site.peers = make(map[string]peer.Peer)
		site.wrapperNonce = make(map[string]int64)
		site.pubsubManager = pubsubManager
		site.user = user
		site.log = log.New("site").WithField("site", addr)
		site.contentDB = contentDB
		site.peerManager = peer.NewManager(pubsubManager, addr)
		site.workerManager = site.NewWorker()
		site.ctx = ctx

		sites[addr] = site
		wrapperKeyMap[siteSettings.WrapperKey] = site
	}

	return &manager{
		ctx:           ctx,
		sites:         sites,
		wrapperKeyMap: wrapperKeyMap,
		pubsubManager: pubsubManager,
		userManager:   userManager,
		contentDB:     contentDB,
	}, nil
}

func (m *manager) NewSite(addr string) (*Site, error) {
	if site, alreadyExists := m.sites[addr]; alreadyExists {
		return site, errors.New("site already exists")
	}
	site := new(Site)
	site.addr = addr
	site.Settings = new(Settings)
	site.trackers = make(map[string]*AnnouncerStats)
	site.peers = make(map[string]peer.Peer)
	site.wrapperNonce = make(map[string]int64)
	site.user = m.userManager.User()
	site.pubsubManager = m.pubsubManager
	site.log = log.New(addr)
	site.contentDB = m.contentDB
	site.peerManager = peer.NewManager(m.pubsubManager, addr)
	site.workerManager = site.NewWorker()
	site.Settings.Added = time.Now().Unix()

	site.Settings.AjaxKey = random.HexString(64)
	site.Settings.AuthKey = random.HexString(64)
	site.Settings.WrapperKey = random.HexString(64)

	m.wrapperKeyMap[site.Settings.WrapperKey] = site
	m.sites[addr] = site
	return site, nil
}

func (m *manager) Site(addr string) *Site {
	return m.sites[addr]
}

func (m *manager) SiteByWrapperKey(wrapperKey string) *Site {
	return m.wrapperKeyMap[wrapperKey]
}

func (m *manager) RenderIndex(siteAddress, indexFilename string, dst io.Writer) error {
	site, ok := m.sites[siteAddress]
	if !ok {
		return errors.New("site not found")
	}

	info, err := site.Info()
	if err != nil {
		return err
	}

	userSettings := site.user.GlobalSettings()
	theme := "theme-light"
	backgroundColor := info.Content.BackgroundColor
	if userSettings.Theme == "dark" {
		theme = "theme-dark"
		backgroundColor = info.Content.BackgroundColorDark
	}
	bodyStyle := "background-color: " + backgroundColor

	wrapperNonce := random.HexString(64)
	scriptNonce := random.Base62String(64)

	site.registerWrapperNonce(wrapperNonce)

	favicon := "uimedia/img/favicon.ico"
	if info.Content.Favicon != "" {
		favicon = path.Join(info.Address, info.Content.Favicon)
	}

	permissions := info.Settings.Permissions
	if site.IsAdmin() {
		permissions = []string{"ADMIN"}
	}

	vars := struct {
		Address                  string
		AjaxKey                  string
		BodyStyle                string
		Favicon                  string
		FileInnerPath            string
		FileURL                  string
		HomePage                 string
		InnerPath                string
		Lang                     string
		Permissions              []string
		PostMessageNonceSecurity bool
		QueryString              string
		Rev                      int
		SandboxPermissions       string
		ScriptNonce              string
		ServerURL                string
		ShowLoadingScreen        bool
		ThemeClass               string
		Title                    string
		WrapperKey               string
		WrapperNonce             string
		ViewPort                 string
	}{
		Address:                  info.Address,
		AjaxKey:                  info.Settings.AjaxKey,
		BodyStyle:                bodyStyle,
		Favicon:                  favicon,
		FileInnerPath:            indexFilename,
		FileURL:                  path.Join(info.Address, indexFilename),
		HomePage:                 "/" + config.HomeSite,
		InnerPath:                indexFilename,
		Lang:                     config.Language,
		Permissions:              permissions,
		PostMessageNonceSecurity: info.Content.PostmessageNonceSecurity,
		QueryString:              "?wrapper_nonce=" + wrapperNonce,
		Rev:                      config.Rev,
		ScriptNonce:              scriptNonce,
		ServerURL:                "", // only used for proxy requests
		ShowLoadingScreen:        site.IsLoading(),
		ThemeClass:               theme,
		Title:                    info.Content.Title,
		WrapperKey:               info.Settings.WrapperKey,
		WrapperNonce:             wrapperNonce,
		ViewPort:                 info.Content.Viewport,
	}

	return template.Wrapper.ExecuteHTML(dst, vars)
}

func (m *manager) ReadFile(site, innerPath string, dst io.Writer) error {
	s, ok := m.sites[site]
	if !ok {
		return errors.New("site not found")
	}

	ctx, cancel := context.WithTimeout(m.ctx, config.FileNeedDeadline)
	defer cancel()
	return s.ReadFile(ctx, innerPath, dst)
}

func (m *manager) SiteList() ([]*Info, error) {
	list := make([]*Info, len(m.sites))
	var err error
	var i int
	for _, site := range m.sites {
		list[i], err = site.Info()
		if err != nil {
			return nil, err
		}
		i++
	}
	return list, nil
}

func (m *manager) Close() {
	for _, site := range m.sites {
		site.SaveSettings()
		site.peerManager.Close()
		site.workerManager.Close()
	}
}
