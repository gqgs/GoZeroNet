//go:generate go run github.com/golang/mock/mockgen -package websocket -source=$GOFILE -destination=./mock.go
package websocket

import (
	"net/http"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/gqgs/go-zeronet/pkg/lib/serialize"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Conn interface {
	ReadMessage() (messageType int, message []byte, err error)
	WriteJSON(v interface{}) error
	Write([]byte) error
}

type conn struct {
	mu           *sync.Mutex
	internalConn *websocket.Conn
}

func IsCloseError(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure)
}

func (c *conn) WriteJSON(v interface{}) error {
	data, err := serialize.JSONMarshal(v)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.internalConn.WriteMessage(websocket.TextMessage, data)
}

func (c *conn) Write(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.internalConn.WriteMessage(websocket.TextMessage, msg)
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
