package config

import (
	"os"
	"strings"
	"time"
)

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

	ConnectionDeadline = time.Second * 10
	FileServerDeadline = time.Second * 5
	RandomIPv4Addr     = "127.0.0.1:"
	DataDir            = "./data/"
	Language           = "en"
	UpdateSite         = "1uPDaT3uSyWAPdCv1WkMb5hBQjWSNNACf" // TODO: ZN updater. We would need a new zite for this.
	SiteSizeLimit      = 10

	// Validate database queries for correctness
	ValidateDatabaseQueries = true

	ContentBufferSize        = 50
	WebsocketBufferSize      = 50
	PeerCandidatesBufferSize = 25
	MaxConnectedPeers        = 10
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
