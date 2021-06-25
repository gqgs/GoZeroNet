package content

import (
	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/database"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
)

type Worker interface {
	Close()
}

type worker struct {
	log           log.Logger
	pubsubManager pubsub.Manager
	queue         <-chan pubsub.Message
	closeCh       chan struct{}
	db            database.ContentDatabase
}

// Creates a new content worker.
// Caller is responsible for calling close when the worker is no longer needed.
func NewWorker(contentDB database.ContentDatabase, pubsubManager pubsub.Manager) *worker {
	w := &worker{
		log:           log.New("content_worker"),
		pubsubManager: pubsubManager,
		queue:         pubsubManager.Register(config.ContentBufferSize),
		closeCh:       make(chan struct{}),
		db:            contentDB,
	}
	go w.run()
	return w
}

func (w *worker) run() {
	for msg := range w.queue {
		switch payload := msg.Event().(type) {
		case *event.FileInfo:
			w.log.Debug("file update event")
			if err := w.db.UpdateFile(msg.Site(), payload); err != nil {
				w.log.Error(err)
			}
		case *event.PeerInfo:
			w.log.Debug("peer update event")
			if err := w.db.UpdatePeer(msg.Site(), payload); err != nil {
				w.log.Error(err)
			}
		case *event.ContentInfo:
			w.log.Debug("content update event")
			if err := w.db.UpdateContent(msg.Site(), payload); err != nil {
				w.log.Error(err)
			}
		}
	}
	close(w.closeCh)
}

func (w *worker) Close() {
	w.log.Debug("closing")
	w.pubsubManager.Unregister(w.queue)
	<-w.closeCh
}
