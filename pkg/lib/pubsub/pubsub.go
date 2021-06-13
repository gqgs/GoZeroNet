package pubsub

import (
	"sync"
	"time"

	"github.com/gqgs/go-zeronet/pkg/lib/log"
)

const sendTimeout = time.Millisecond * 100

func NewManager() *manager {
	return &manager{
		queue: make(map[<-chan Message]chan Message),
		log:   log.New("pubsub"),
	}
}

type (
	Message interface {
		Body() []byte
		Queue() string
		Source() string
	}

	manager struct {
		mu    sync.RWMutex
		queue map[<-chan Message]chan Message
		log   log.Logger
	}

	message struct {
		body   []byte
		queue  string
		source string
	}

	Manager interface {
		// Register creates a new subscriber to all queues
		// The client MUST unregister the channel after using it
		Register() <-chan Message
		Unregister(messageCh <-chan Message)
		Broadcast(source, queue string, body []byte)
	}
)

func (m *message) Body() []byte {
	return m.body
}

func (m *message) Queue() string {
	return m.queue
}

func (m *message) Source() string {
	return m.source
}

func (m *manager) Register() <-chan Message {
	messageCh := make(chan Message)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queue[messageCh] = messageCh
	return messageCh
}

func (m *manager) Unregister(messageCh <-chan Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.queue, messageCh)
}

func (m *manager) Broadcast(source, queue string, body []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var wg sync.WaitGroup
	wg.Add(len(m.queue))
	for _, channel := range m.queue {
		channel := channel
		go func() {
			defer wg.Done()
			select {
			case channel <- &message{
				source: source,
				queue:  queue,
				body:   body,
			}:
			case <-time.After(sendTimeout):
				m.log.Warn("dropped message")
			}
		}()
	}
	wg.Wait()
}
