package site

import (
	"encoding/binary"
	"errors"
	"sync"

	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
	"github.com/gqgs/go-zeronet/pkg/lib/ip"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
	"github.com/gqgs/go-zeronet/pkg/lib/safe"
	"github.com/gqgs/go-zeronet/pkg/peer"
)

type Worker interface {
	Close()
}

type worker struct {
	log         log.Logger
	queue       <-chan pubsub.Message
	closeCh     chan struct{}
	site        *Site
	peerManager peer.Manager
	mu          sync.RWMutex
	downloading map[string]struct{}
}

// Creates a new site worker.
// Caller is responsible for calling close when the worker is no longer needed.
func (s *Site) NewWorker(peerManager peer.Manager) *worker {
	w := &worker{
		log:         log.New("site_worker"),
		queue:       s.pubsubManager.Register(50),
		closeCh:     make(chan struct{}),
		site:        s,
		peerManager: peerManager,
		downloading: make(map[string]struct{}),
	}
	go w.run()
	return w
}

func (w *worker) run() {
	for msg := range w.queue {
		switch payload := msg.Event().(type) {
		case *event.PeersNeed:
			w.log.Debug("peer need event")
			go w.site.Announce()
		case *event.FileNeed:
			w.log.Debug("file need event")

			w.mu.RLock()
			_, downloading := w.downloading[payload.InnerPath]
			w.mu.RUnlock()

			if downloading {
				w.log.Debug("already downlading file. skipping ", payload.InnerPath)
				continue
			}
			w.mu.Lock()
			w.downloading[payload.InnerPath] = struct{}{}
			w.mu.Unlock()

			go func() {
				if err := w.downloadFile(payload); err != nil {
					w.log.Error(err)
				}

				w.mu.Lock()
				delete(w.downloading, payload.InnerPath)
				w.mu.Unlock()
			}()

			// TODO: if file is already downloaded on other site (check via content hash)
			// create a hard link for this one instead of downloading again
		}
	}
	close(w.closeCh)
}

func (w *worker) downloadFile(fileNeed *event.FileNeed) error {
	info, err := w.site.contentDB.FileInfo(w.site.addr, fileNeed.InnerPath)
	if err != nil {
		return err
	}

	hashID, err := crypto.HashID(info.Hash)
	if err != nil {
		return err
	}

	p, err := w.peerManager.GetConnected()
	if err != nil {
		return err
	}
	defer w.peerManager.PutConnected(p)

	resp, err := fileserver.FindHashIDs(p, w.site.addr, hashID)
	if err != nil {
		return err
	}

	w.log.Debugf("found %d results for %d", len(resp.Peers), hashID)
	for id, addresses := range resp.Peers {
		w.log.Debug(id, len(addresses), id, hashID)
		if id == hashID {
			for _, addr := range addresses {
				parsed := ip.ParseIPv4(addr, binary.LittleEndian)
				w.log.Debug("connection to new new peer ", parsed)
				newPeer := peer.NewPeer(parsed)
				if err := newPeer.Connect(); err != nil {
					w.log.Warn(err)
					continue
				}
				defer newPeer.Close()

				w.log.Debugf("downloading file %s from %s", fileNeed.InnerPath, newPeer)
				if err := w.site.downloadFile(newPeer, safe.CleanPath(fileNeed.InnerPath), info); err != nil {
					w.log.WithField("peer", newPeer).Warn(err)
					continue
				}
				return nil
			}
			break
		}
	}

	return errors.New("could not download file")
}

func (w *worker) Close() {
	w.log.Debug("closing")
	w.site.pubsubManager.Unregister(w.queue)
	<-w.closeCh
}
