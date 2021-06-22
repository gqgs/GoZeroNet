package pubsub

import (
	"bytes"
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
		case <-time.After(time.Second):
			require.Fail(t, "message timeout")
		}
	}()

	manager.Broadcast("test site", bytes.NewBuffer([]byte("test message")))
	manager.Unregister(messageCh)
	require.Len(t, manager.queue, 0)
}
