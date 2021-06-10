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

	ConnectionDeadline = time.Second * 5
	FileServerDeadline = time.Second
	FileServerAddr     = "127.0.0.1:"
	UIServerAddr       = "127.0.0.1:43111"
	RandomIPv4Addr     = "127.0.0.1:"
	DataDir            = "./data/"
	Language           = "en"
)

var (
	Debug = strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") ||
		strings.EqualFold(os.Getenv("LOG_LEVEL"), "trace")
)
