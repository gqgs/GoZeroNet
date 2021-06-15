package site

import (
	"errors"
	"io"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/content"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/template"
	"github.com/gqgs/go-zeronet/pkg/user"
)

type Manager interface {
	Site(addr string) *Site
	RenderIndex(site, indexFilename string, dst io.Writer) error
	ReadFile(site, innerPath string, dst io.Writer) error
	SiteByWrapperKey(wrapperKey string) *Site
	SiteList() ([]*Info, error)
	SetUser(user user.User)
}

type manager struct {
	// Address -> Site
	sites map[string]*Site
	// WrapperKey -> Site
	wrapperKeyMap map[string]*Site
}

func NewSiteManager(pubsubManager pubsub.Manager) (*manager, error) {
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

func (m *manager) SetUser(user user.User) {
	for _, site := range m.sites {
		site.user = user
	}
}

func (m *manager) RenderIndex(site, indexFilename string, dst io.Writer) error {
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
