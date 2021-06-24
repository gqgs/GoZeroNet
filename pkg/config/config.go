package config

import (
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type config struct {
	Language                 string   `toml:"language"`
	DataDir                  string   `toml:"data_dir"`
	SiteSizeLimit            int      `toml:"site_size_limit"`
	Trackers                 []string `toml:"trackers"`
	ValidateDatabaseQueries  bool     `toml:"validate_database_queries"`
	ContentBufferSize        int      `toml:"content_buffer_size"`
	WebsocketBufferSize      int      `toml:"websocket_buffer_size"`
	PeerCandidatesBufferSize int      `toml:"peer_candidates_buffer_size"`
	MaxConnectedPeers        int      `toml:"max_connected_peers"`
	ConnectionDeadline       duration `toml:"connection_deadline"`
	FileServerDeadline       duration `toml:"fileserver_deadline"`
	FileNeedDeadline         duration `toml:"file_need_deadline"`
}

func init() {
	//nolint:dogsled
	_, configFilename, _, _ := runtime.Caller(0)
	root := strings.TrimSuffix(configFilename, "pkg/config/config.go")
	for _, configFile := range []string{"zeronet.toml", "zeronet.toml.example"} {
		filename := path.Join(root, configFile)
		if _, err := os.Stat(filename); err != nil {
			continue
		}

		c := new(config)
		if _, err := toml.DecodeFile(filename, c); err != nil {
			panic(err)
		}

		Language = c.Language
		DataDir = c.DataDir
		SiteSizeLimit = c.SiteSizeLimit
		Trackers = c.Trackers
		ValidateDatabaseQueries = c.ValidateDatabaseQueries
		ContentBufferSize = c.ContentBufferSize
		WebsocketBufferSize = c.WebsocketBufferSize
		PeerCandidatesBufferSize = c.PeerCandidatesBufferSize
		MaxConnectedPeers = c.MaxConnectedPeers
		ConnectionDeadline = c.ConnectionDeadline.Duration
		FileServerDeadline = c.FileServerDeadline.Duration
		FileNeedDeadline = c.FileNeedDeadline.Duration

		return
	}

	panic("configuration file `zeronet.toml` not found")
}

const (
	Rev        = 4555
	Version    = "0.7.2"
	Protocol   = "v2"
	PortOpened = true
	// UseBinType tells msgpack to use the bin type
	// instead of the deprecated raw type.
	UseBinType = true

	DefaultFileServerAddr = "127.0.0.1:0"
	DefaultUIServerAddr   = "127.0.0.1:43111"

	RandomIPv4Addr = "127.0.0.1:"
	UpdateSite     = "1uPDaT3uSyWAPdCv1WkMb5hBQjWSNNACf" // TODO: ZN updater. We would need a new zite for this.
)

var (
	SiteSizeLimit            int
	DataDir                  string
	Language                 string
	ContentBufferSize        int
	WebsocketBufferSize      int
	PeerCandidatesBufferSize int
	MaxConnectedPeers        int
	ValidateDatabaseQueries  bool // Validate database queries for correctness
	Trackers                 []string
	ConnectionDeadline       time.Duration
	FileServerDeadline       time.Duration
	FileNeedDeadline         time.Duration

	FileServerHost = "127.0.0.1"
	FileServerPort = 0
	UIServerHost   = "127.0.0.1"
	UIServerPort   = 43111

	Debug = strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") ||
		strings.EqualFold(os.Getenv("LOG_LEVEL"), "trace")
)

// needed for decoding time from TOML files
type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
