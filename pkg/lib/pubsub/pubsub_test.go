package pubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Pubsub(t *testing.T) {
	manager := NewManager()
	messageCh := manager.Register()
	require.Len(t, manager.queue, 1)

	go func() {
		select {
		case message := <-messageCh:
			require.Equal(t, "test site", message.Site())
			require.Equal(t, "test queue", message.Event())
			require.Equal(t, []byte("test message"), message.Body())
		case <-time.After(time.Second):
			require.Fail(t, "message timeout")
		}
	}()

	manager.Broadcast("test site", "test queue", []byte("test message"))
	manager.Unregister(messageCh)
	require.Len(t, manager.queue, 0)
}
