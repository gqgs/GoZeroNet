package plugin

type pluginWriter interface {
	WriteJSON(v interface{}) error
	Write(data []byte) error
}

type Plugin interface {
	Name() string
	Description() string
	// Handles returns true if the plugin handles this command
	Handles(cmd string) bool
	// Handle parses and handles the message writing the result to w
	Handle(w pluginWriter, cmd string, to, id int64, message []byte) error
}
