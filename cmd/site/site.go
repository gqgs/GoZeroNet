package site

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"

	"github.com/gqgs/go-zeronet/pkg/config"
)

type Site struct {
	addr      string
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

func (s Site) ReadFile(filename string, dst io.Writer) error {
	innerPath := path.Join(config.DataDir, s.addr, filename)
	file, err := os.Open(innerPath)
	if err != nil {
		// TODO: download file
		return err
	}
	defer file.Close()

	_, err = io.Copy(dst, file)
	return err
}

type SiteManager interface {
	ReadFile(site, filename string, dst io.Writer) error
}

type siteManager struct {
	// Address -> Site info
	sites map[string]*Site
}

func (m *siteManager) ReadFile(site, filename string, dst io.Writer) error {
	s, ok := m.sites[site]
	if !ok {
		// TODO: download site
		return errors.New("site not found")
	}

	return s.ReadFile(filename, dst)
}

func NewSiteManager() (*siteManager, error) {
	sitesFile, err := os.Open(path.Join(config.DataDir, "sites.json"))
	if err != nil {
		// TODO: ignore error if file not found
		return nil, err
	}
	defer sitesFile.Close()

	sites := make(map[string]*Site)
	if err = json.NewDecoder(sitesFile).Decode(&sites); err != nil {
		return nil, err
	}

	for addr, site := range sites {
		site.addr = addr
	}

	return &siteManager{
		sites: sites,
	}, nil
}
