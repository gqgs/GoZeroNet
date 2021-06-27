package peer

import (
	"time"

	"github.com/gqgs/go-zeronet/pkg/config"
	"github.com/gqgs/go-zeronet/pkg/event"
	"github.com/gqgs/go-zeronet/pkg/fileserver"
	"github.com/gqgs/go-zeronet/pkg/lib/log"
	"github.com/gqgs/go-zeronet/pkg/lib/pubsub"
)

const waitForconnectedTimeout = time.Second * 30

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
	connectedCh   chan *peer
	doneCh        chan *peer
	msgCh         <-chan pubsub.Message
	closeCh       chan struct{}
}

func NewManager(pubsubManager pubsub.Manager, site string) *manager {
	m := &manager{
		log:           log.New("peer_manager"),
		site:          site,
		pubsubManager: pubsubManager,
		connectedCh:   make(chan *peer, config.MaxConnectedPeers),
		doneCh:        make(chan *peer, config.MaxConnectedPeers),
		msgCh:         pubsubManager.Register("peer_manager", config.PeerCandidatesBufferSize),
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
	case <-time.After(waitForconnectedTimeout):
		event.BroadcastPeersNeed(m.site, m.pubsubManager, &event.PeersNeed{})
	}
	return m.GetConnected()
}

func (m *manager) PutConnected(p *peer) {
	m.doneCh <- p
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

				m.log.WithField("queue", m.msgCh).Debug("new peer candidate event")

				// already have as many connected pers as we want
				if len(m.connectedCh) == config.MaxConnectedPeers {
					continue
				}

				go func() {
					peer := NewPeer(candidate.Address)
					if err := peer.Connect(); err != nil {
						m.log.WithField("peer", peer).Warn(err)
						return
					}
					if err := peer.CheckConnection(); err != nil {
						m.log.WithField("peer", peer).Warn(err)
						peer.Close()
						return
					}

					m.log.Info("new connected peer")

					m.connectedCh <- peer
				}()
			}
		case done := <-m.doneCh:
			if _, err := fileserver.Ping(done); err != nil {
				m.log.WithField("peer", done).Warnf("closing connection: %s", err)
				if err := done.Close(); err != nil {
					m.log.WithField("peer", done).Error(err)
					continue
				}
			}
			m.connectedCh <- done
		}
	}
}
