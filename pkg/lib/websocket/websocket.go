//go:generate go run github.com/golang/mock/mockgen -package websocket -source=$GOFILE -destination=./mock.go
package websocket

import (
	"net/http"
	"sync"

	"github.com/bytedance/sonic"
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
	mu           *sync.Mutex
	internalConn *websocket.Conn
}

func (c *conn) WriteJSON(v interface{}) error {
	data, err := sonic.Marshal(v)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.internalConn.WriteMessage(websocket.TextMessage, data)
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
		mu:           new(sync.Mutex),
		internalConn: c,
	}, nil
}
