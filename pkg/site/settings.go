package site

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"sync"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/lib/bigfile"
)

type Settings struct {
	Added     int64  `json:"added"`
	AjaxKey   string `json:"ajax_key"`
	AuthKey   string `json:"auth_key"`
	BytesRecv int    `json:"bytes_recv"`
	BytesSent int    `json:"bytes_sent"`
	Cache     struct {
		BadFiles         map[string]int `json:"bad_files"`
		Hashfield        string         `json:"hashfield"`
		pieceFieldsMutex sync.Mutex
		Piecefields      map[string]bigfile.PieceField `json:"piecefields"`
	} `json:"cache"`
	HasBigFile                bool              `json:"has_bigfile"`
	Downloaded                int64             `json:"downloaded"`
	Modified                  int64             `json:"modified"`
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

func loadSiteSettingsFromFile() (map[string]*Settings, error) {
	sitesFile, err := os.Open(path.Join(config.DataDir, "sites.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]*Settings), nil
		}
		return nil, err
	}
	defer sitesFile.Close()

	settings := make(map[string]*Settings)
	return settings, json.NewDecoder(sitesFile).Decode(&settings)
}
