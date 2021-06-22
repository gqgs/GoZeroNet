package content

import (
	"encoding/json"

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
		queue:         pubsubManager.Register(),
		closeCh:       make(chan struct{}),
		db:            contentDB,
	}
	go m.listen()
	return m
}

func (m *manager) listen() {
	for msg := range m.queue {
		switch msg.Event() {
		case "file-done":
			m.log.Debug("file done event")
			payload := new(event.FileInfo)
			if err := json.Unmarshal(msg.Body(), payload); err != nil {
				m.log.Error(err)
				continue
			}
			if err := m.db.UpdateFile(msg.Site(), payload.InnerPath, payload.Hash, payload.Size); err != nil {
				m.log.Error(err)
			}
		}
	}
	close(m.closeCh)
}

func (m *manager) Close() {
	m.pubsubManager.Unregister(m.queue)
	<-m.closeCh
}
