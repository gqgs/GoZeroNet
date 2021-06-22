package plugin

type pluginWriter interface {
	WriteJSON(v interface{}) error
}

type errorMsg struct {
	Msg string `json:"error"`
	To  int64  `json:"to"`
	ID  int64  `json:"id"`
}

type Plugin interface {
	Name() string
	Description() string
	// Handles returns true if the plugin handles this command
	Handles(cmd string) bool
	// Handle parses and handles the message writing the result to w
	Handle(w pluginWriter, cmd string, to, id int64, message []byte) error
}
