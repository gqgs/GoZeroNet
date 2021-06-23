package pubsub

import (
	"sync"

	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
)

func NewManager() *manager {
	return &manager{
		queue: make(map[<-chan Message]chan Message),
		log:   log.New("pubsub"),
	}
}

type (
	Message interface {
		Event() event.Event
		Site() string
	}

	manager struct {
		mu    sync.RWMutex
		queue map[<-chan Message]chan Message
		log   log.Logger
	}

	message struct {
		event event.Event
		site  string
	}

	Manager interface {
		// Register creates a new subscriber to all events.
		// bufferSize defines the buffer size of the message channel.
		// If the buffer is full when the manager tries to send a new message
		// the message will be discarted.
		// The client MUST unregister the channel after using it.
		Register(bufferSize int) <-chan Message
		Unregister(messageCh <-chan Message)
		Broadcast(site string, event event.Event)
	}
)

func (m *message) Event() event.Event {
	return m.event
}

func (m *message) Site() string {
	return m.site
}

func (m *manager) Register(bufferSize int) <-chan Message {
	messageCh := make(chan Message, bufferSize)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queue[messageCh] = messageCh
	return messageCh
}

func (m *manager) Unregister(messageCh <-chan Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	close(m.queue[messageCh])
	delete(m.queue, messageCh)
}

func (m *manager) Broadcast(site string, event event.Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var wg sync.WaitGroup
	wg.Add(len(m.queue))
	m.log.WithField("site", site).WithField("event", event).Debug("broadcasting new event")
	for _, channel := range m.queue {
		channel := channel
		go func() {
			defer wg.Done()
			select {
			case channel <- &message{
				site:  site,
				event: event,
			}:
			default:
				m.log.Warn("dropped message")
			}
		}()
	}
	wg.Wait()
}
