package content

import (
	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
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
	db            database.ContentDatabase
}

// Creates a new content manager.
// Caller is responsible for calling close when the manager is no longer needed.
func NewManager(contentDB database.ContentDatabase, pubsubManager pubsub.Manager) *manager {
	m := &manager{
		log:           log.New("content_manager"),
		pubsubManager: pubsubManager,
		queue:         pubsubManager.Register(config.ContentBufferSize),
		closeCh:       make(chan struct{}),
		db:            contentDB,
	}
	go m.listen()
	return m
}

func (m *manager) listen() {
	for msg := range m.queue {
		switch payload := msg.Event().(type) {
		case *event.FileInfo:
			m.log.Debug("file done event")
			if err := m.db.UpdateFile(msg.Site(), payload.InnerPath, payload.Hash, payload.Size); err != nil {
				m.log.Error(err)
			}
		case *event.PeerInfo:
			m.log.Debug("peer info event")
			if err := m.db.UpdatePeer(msg.Site(), payload.Address, payload.ReputationDelta); err != nil {
				m.log.Error(err)
			}
		}
	}
	close(m.closeCh)
}

func (m *manager) Close() {
	m.log.Debug("closing")
	m.pubsubManager.Unregister(m.queue)
	<-m.closeCh
}
