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
	LogLevel                string   `toml:"log_level"`
	Language                string   `toml:"language"`
	DataDir                 string   `toml:"data_dir"`
	SiteSizeLimit           int      `toml:"site_size_limit"`
	Trackers                []string `toml:"trackers"`
	ValidateDatabaseQueries bool     `toml:"validate_database_queries"`
	ContentBufferSize       int      `toml:"content_buffer_size"`
	WebsocketBufferSize     int      `toml:"websocket_buffer_size"`
	MaxConnectedPeers       int      `toml:"max_connected_peers"`
	DefaultChannelSize      int      `toml:"default_channel_size"`
	ConnectionDeadline      duration `toml:"connection_deadline"`
	FileServerDeadline      duration `toml:"fileserver_deadline"`
	FileNeedDeadline        duration `toml:"file_need_deadline"`
	FileServerAddress       string   `toml:"fileserver_address"`
	UIServerAddress         string   `toml:"uiserver_address"`
	MaxDownloadTries        int      `toml:"max_download_tries"`
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
		MaxConnectedPeers = c.MaxConnectedPeers
		ConnectionDeadline = c.ConnectionDeadline.Duration
		FileServerDeadline = c.FileServerDeadline.Duration
		FileNeedDeadline = c.FileNeedDeadline.Duration
		DefaultChannelSize = c.DefaultChannelSize
		LogLevel = c.LogLevel
		FileServerAddress = c.FileServerAddress
		UIServerAddress = c.UIServerAddress
		MaxDownloadTries = c.MaxDownloadTries

		if err := os.MkdirAll(path.Dir(DataDir), os.ModePerm); err != nil {
			panic(err)
		}

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
	UseBinType       = true
	FileGetSizeLimit = 512 * 1024

	RandomIPv4Addr = "127.0.0.1:"
	UpdateSite     = "1uPDaT3uSyWAPdCv1WkMb5hBQjWSNNACf" // TODO: ZN updater. We would need a new zite for this.
	HomeSite       = "1HeLLo4uzjaLetFx6NH3PMwFP3qbRbTf3D"
)

var (
	LogLevel                string
	SiteSizeLimit           int
	DataDir                 string
	Language                string
	ContentBufferSize       int
	WebsocketBufferSize     int
	MaxConnectedPeers       int
	DefaultChannelSize      int
	ValidateDatabaseQueries bool // Validate database queries for correctness
	Trackers                []string
	ConnectionDeadline      time.Duration
	FileServerDeadline      time.Duration
	FileNeedDeadline        time.Duration
	FileServerAddress       string
	UIServerAddress         string
	MaxDownloadTries        int

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
