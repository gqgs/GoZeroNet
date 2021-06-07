package config

import (
	"time"
)

const (
	Rev        = 4556
	Version    = "0.8.0"
	Protocol   = "v2"
	PortOpened = true
	// UseBinType tells msgpack to use the bin type
	// instead of the deprecated raw type.
	UseBinType     = true
	Deadline       = time.Second
	FileServerAddr = "127.0.0.1:"
	UIServerAddr   = "127.0.0.1:43111"
	RandomIPv4Addr = "127.0.0.1:"
)
