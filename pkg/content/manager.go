package content

import (
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
)

type Manager interface {
	Close()
}

type manager struct {
	log           log.Logger
	pubsubManager pubsub.Manager
	queue         <-chan pubsub.Message
	closeCh       chan struct{}
}

// Createa a new content manager.
// Caller is responsible for calling close when the manager is no longer needed.
func NewManager() *manager {
	pubsubManager := pubsub.NewManager()
	queue := pubsubManager.Register()
	return &manager{
		log:           log.New("content_manager"),
		pubsubManager: pubsubManager,
		queue:         queue,
		closeCh:       make(chan struct{}),
	}
}

func (m *manager) Run() {
	for msg := range m.queue {
		m.log.Info(msg)
	}
	close(m.closeCh)
}

func (m *manager) Close() {
	m.pubsubManager.Unregister(m.queue)
	<-m.closeCh
}
