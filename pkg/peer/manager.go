package peer

import (
	"sync/atomic"
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
)

const waitForconnectedTimeout = time.Second

type Manager interface {
	// Returns a connected peer.
	// The caller MUST return the peer after doing using
	// by calling the PutConnected method.
	GetConnected() *peer
	PutConnected(p *peer)
	Close()
}

type manager struct {
	log           log.Logger
	pubsubManager pubsub.Manager
	site          string
	connected     int64
	connectedCh   chan *peer
	msgCh         <-chan pubsub.Message
	closeCh       chan struct{}
}

func NewManager(pubsubManager pubsub.Manager, site string) *manager {
	m := &manager{
		log:           log.New("peer_manager"),
		site:          site,
		pubsubManager: pubsubManager,
		connectedCh:   make(chan *peer, config.MaxConnectedPeers),
		msgCh:         pubsubManager.Register("peer_manager", config.MaxConnectedPeers),
		closeCh:       make(chan struct{}),
	}
	go m.processPeerCandidates()
	return m
}

func (m *manager) GetConnected() *peer {
	select {
	case connected := <-m.connectedCh:
		if err := connected.CheckConnection(); err == nil {
			return connected
		}
		connected.Close()
		atomic.AddInt64(&m.connected, -1)
	case <-time.After(waitForconnectedTimeout):
		event.BroadcastPeersNeed(m.site, m.pubsubManager, &event.PeersNeed{})
		time.Sleep(time.Second * 5)
	}
	return m.GetConnected()
}

func (m *manager) PutConnected(p *peer) {
	select {
	case m.connectedCh <- p:
	default:
		p.Close()
		atomic.AddInt64(&m.connected, -1)
	}
}

func (m *manager) Close() {
	m.log.Debug("closing")
	close(m.closeCh)
	m.pubsubManager.Unregister(m.msgCh)
}

func (m *manager) processPeerCandidates() {
	for {
		select {
		case <-m.closeCh:
			return
		case msg := <-m.msgCh:
			switch candidate := msg.Event().(type) {
			case *event.PeerCandidate:
				if msg.Site() != m.site {
					continue
				}

				m.log.WithField("queue", len(m.msgCh)).Debug("new peer candidate event")

				if int(atomic.LoadInt64(&m.connected)) >= config.MaxConnectedPeers {
					m.log.Debug("already have as many connected pers as we want")
					continue
				}
				atomic.AddInt64(&m.connected, 1)

				go func() {
					peer := NewPeer(candidate.Address)
					logger := m.log.WithField("peer", peer)
					if err := peer.Connect(); err != nil {
						logger.Warn(err)
						atomic.AddInt64(&m.connected, -1)
						return
					}

					logger.Info("new connected peer: ", peer)
					m.connectedCh <- peer
				}()
			}
		}
	}
}
