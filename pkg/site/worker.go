package site

import (
	"encoding/binary"
	"sync"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/crypto"
	"github.com/gqgs/go-zeronet/pkg/lib/ip"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
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
	mu          *sync.Mutex
	downloading map[string]struct{}
}

// Creates a new site worker.
// Caller is responsible for calling close when the worker is no longer needed.
func (s *Site) NewWorker() *worker {
	w := &worker{
		log:         log.New("site_worker").WithField("site", s.addr),
		queue:       s.pubsubManager.Register("site_worker", 50),
		closeCh:     make(chan struct{}),
		site:        s,
		peerManager: s.peerManager,
		downloading: make(map[string]struct{}),
		mu:          new(sync.Mutex),
	}
	go w.run()
	return w
}

func (w *worker) run() {
	var wg sync.WaitGroup
	for msg := range w.queue {
		if msg.Site() != w.site.addr {
			continue
		}

		switch payload := msg.Event().(type) {
		case *event.PeersNeed:
			w.log.WithField("queue", len(w.queue)).Debug("peer need event")
			go w.site.Announce()

		case *event.FileInfo:
			w.log.WithField("queue", len(w.queue)).Debug("file info event")
		case *event.ContentInfo:
			w.log.WithField("queue", len(w.queue)).Debug("content info event")
			if payload.Modified > int(w.site.Settings.Modified) {
				w.site.Settings.Modified = int64(payload.Modified)
			}
			go w.site.BroadcastSiteChange("file_done", payload.InnerPath)
		case *event.SiteUpdate:
			w.log.WithField("queue", len(w.queue)).WithField("inner_path", payload.InnerPath).Debug("site update event")
			go func() {
				if err := w.site.Update(7); err != nil {
					w.log.Error(err)
				}
			}()
		case *event.FileNeed:
			w.log.WithField("queue", len(w.queue)).Debug("file need event")
			if payload.Tries >= config.MaxDownloadTries {
				w.log.WithField("inner_path", payload.InnerPath).Error("failed to download file")
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.mu.Lock()
				if _, downloading := w.downloading[payload.InnerPath]; downloading {
					w.log.Debug("already downlading file. skipping ", payload.InnerPath)
					w.mu.Unlock()
					return
				}
				w.downloading[payload.InnerPath] = struct{}{}
				w.mu.Unlock()

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
	w.log.Debug("waiting for wg")
	wg.Wait()
	w.log.Debug("wg done")
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

	p, err := w.peerManager.GetConnected(w.site.ctx)
	if err != nil {
		return err
	}
	defer w.peerManager.PutConnected(p)

	resp, err := fileserver.FindHashIDs(p, w.site.addr, int64(hashID))
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

				w.log.Debugf("downloading file %s from %s", fileNeed.InnerPath, newPeer)
				if err := w.site.downloadFile(newPeer, info); err != nil {
					w.log.WithField("peer", newPeer).Warn(err)
					continue
				}
				return nil
			}
			break
		}
	}

	fileNeed.Tries++
	event.BroadcastFileNeed(w.site.addr, w.site.pubsubManager, fileNeed)

	return nil
}

func (w *worker) Close() {
	w.log.Debug("closing")
	w.site.pubsubManager.Unregister(w.queue)
	<-w.closeCh
}
