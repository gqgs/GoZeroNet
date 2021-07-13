package info

import (
	"runtime"

	"github.com/gqgs/go-zeronet/pkg/config"
)

type Server struct {
	IPExternal bool `json:"ip_external"`
	PortOpened struct {
		Ipv4 bool `json:"ipv4"`
		Ipv6 bool `json:"ipv6"`
	} `json:"port_opened"`
	Platform          string   `json:"platform"`
	FileserverIP      string   `json:"fileserver_ip"`
	FileserverPort    int      `json:"fileserver_port"`
	TorEnabled        bool     `json:"tor_enabled"`
	TorStatus         string   `json:"tor_status"`
	TorHasMeekBridges bool     `json:"tor_has_meek_bridges"`
	TorUseBridges     bool     `json:"tor_use_bridges"`
	UIIP              string   `json:"ui_ip"`
	UIPort            int      `json:"ui_port"`
	Version           string   `json:"version"`
	Rev               int      `json:"rev"`
	Timecorrection    int      `json:"timecorrection"`
	Language          string   `json:"language"`
	Debug             bool     `json:"debug"`
	Offline           bool     `json:"offline"`
	Plugins           []string `json:"plugins"`
	PluginsRev        struct {
	} `json:"plugins_rev"`
	UserSettings struct {
	} `json:"user_settings"`
	Updatesite string `json:"updatesite"`
	DistType   string `json:"dist_type"`
}

func ServerInfo(admin bool) Server {
	s := Server{
		UIIP:              config.UIServerHost,
		UIPort:            config.UIServerPort,
		TorEnabled:        config.TorEnabled,
		TorStatus:         "Disabled", // TODO: check Tor status
		TorHasMeekBridges: false,
		FileserverIP:      "*",
		FileserverPort:    config.FileServerPort,
		Platform:          runtime.GOOS,
		PortOpened: struct {
			Ipv4 bool `json:"ipv4"`
			Ipv6 bool `json:"ipv6"`
		}{
			Ipv4: true,
			Ipv6: false,
		},
		Timecorrection: 0,
		IPExternal:     true,
		Version:        config.Version,
		Rev:            config.Rev,
		Debug:          config.Debug,
		Offline:        false,
		Language:       config.Language,
		Plugins:        make([]string, 0), // TODO: return enabled plugins
	}

	if admin {
		s.Updatesite = config.UpdateSite
		s.DistType = "source"
	}

	return s
}
