package site

import (
	"errors"
	"io"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/random"
	"github.com/gqgs/go-zeronet/pkg/template"
	"github.com/gqgs/go-zeronet/pkg/user"
)

type Manager interface {
	Site(addr string) *Site
	RenderIndex(site, indexFilename string, dst io.Writer) error
	ReadFile(site, innerPath string, dst io.Writer) error
	SiteByWrapperKey(wrapperKey string) *Site
	SiteList() ([]*Info, error)
}

type manager struct {
	// Address -> Site
	sites map[string]*Site
	// WrapperKey -> Site
	wrapperKeyMap map[string]*Site
}

func NewSiteManager(pubsubManager pubsub.Manager, userManager user.Manager) (*manager, error) {
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
		site.peers = make(map[string]struct{})
		wrapperKeyMap[siteSettings.WrapperKey] = site
		site.pubsubManager = pubsubManager
		site.user = user
		sites[addr] = site
	}

	return &manager{
		sites:         sites,
		wrapperKeyMap: wrapperKeyMap,
	}, nil
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
		// TODO: download site
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
		Favicon:                  info.Content.Favicon,
		FileInnerPath:            indexFilename,
		FileURL:                  info.Address + "/" + indexFilename,
		HomePage:                 info.Address,
		InnerPath:                indexFilename,
		Lang:                     config.Language,
		Permissions:              info.Settings.Permissions,
		PostMessageNonceSecurity: info.Content.PostmessageNonceSecurity,
		QueryString:              "?wrapper_nonce=" + wrapperNonce,
		Rev:                      config.Rev,
		ScriptNonce:              scriptNonce,
		ServerURL:                "", // only used for proxy requests
		ShowLoadingScreen:        false,
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
		// TODO: download site
		return errors.New("site not found")
	}

	return s.ReadFile(innerPath, dst)
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
