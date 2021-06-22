package plugin

import "github.com/gqgs/go-zeronet/pkg/site"

type pluginWriter interface {
	WriteJSON(v interface{}) error
}

// HandlerFunc parses and handles the message writing the result to w.
type HandlerFunc func(w pluginWriter, site *site.Site, rawMesage []byte) error

// Function that generates a new ID response.
type IDFunc = func() int64

type Plugin interface {
	// Name of the plugin.
	Name() string
	// Short description of the plugin's functionality.
	Description() string
	// Handler returns true and the handler if plugin handles this command.
	// It returns nil and false otherwise.
	Handler(cmd string) (HandlerFunc, bool)
}
