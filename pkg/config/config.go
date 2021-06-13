package config

import (
	"os"
	"strings"
	"time"
)

const (
	Rev        = 4556
	Version    = "0.8.0"
	Protocol   = "v2"
	PortOpened = true
	// UseBinType tells msgpack to use the bin type
	// instead of the deprecated raw type.
	UseBinType = true

	DefaultFileServerAddr = "127.0.0.1:0"
	DefaultUIServerAddr   = "127.0.0.1:43111"

	ConnectionDeadline = time.Second * 5
	FileServerDeadline = time.Second
	RandomIPv4Addr     = "127.0.0.1:"
	DataDir            = "./data/"
	Language           = "en"
	UpdateSite         = "1uPDaT3uSyWAPdCv1WkMb5hBQjWSNNACf"
)

var (
	Trackers = []string{
		"http://h4.trakx.nibba.trade:80/announce",  // US/VA
		"http://open.acgnxtracker.com:80/announce", // DE
		"http://tracker.bt4g.com:2095/announce",    // Cloudflare
	}

	Debug = strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") ||
		strings.EqualFold(os.Getenv("LOG_LEVEL"), "trace")

	FileServerHost = "127.0.0.1"
	FileServerPort = 0
	UIServerHost   = "127.0.0.1"
	UIServerPort   = 43111
)
