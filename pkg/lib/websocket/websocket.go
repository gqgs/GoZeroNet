//go:generate go run github.com/golang/mock/mockgen -package websocket -source=$GOFILE -destination=./mock.go
package websocket

import (
	"net/http"
	"sync"

	"github.com/fasthttp/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Conn interface {
	ReadMessage() (messageType int, message []byte, err error)
	WriteJSON(v interface{}) error
}

type conn struct {
	mu           sync.Mutex
	internalConn *websocket.Conn
}

func (c *conn) WriteJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.internalConn.WriteJSON(v)
}

func (c *conn) Write(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// https://www.rfc-editor.org/rfc/rfc6455.html#section-11.8
	return c.internalConn.WriteMessage(1, data)
}

func (c *conn) ReadMessage() (messageType int, message []byte, err error) {
	return c.internalConn.ReadMessage()
}

func Upgrade(w http.ResponseWriter, r *http.Request) (Conn, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &conn{
		internalConn: c,
	}, nil
}
