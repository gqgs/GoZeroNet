package config

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Rev        = 4556
	Version    = "0.8.0"
	Protocol   = "v2"
	PortOpened = false
	// UseBinType tells msgpack to use the bin type
	// instead of the deprecated raw type.
	UseBinType = true
)

func init() {
	rand.Seed(time.Now().UnixNano())
	// User ports [1024-49151)
	FileServer.Port = rand.Intn(48127) + 1024
}

var (
	FileServer = Server{
		Port: 43112,
	}
	UIServer = Server{
		Port: 43111,
	}
)

type Server struct {
	Port int
}

func (s Server) Addr() string {
	return fmt.Sprintf("0.0.0.0:%d", s.Port)
}
