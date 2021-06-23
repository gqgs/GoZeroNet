package peer

import (
	"errors"
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
	GetConnected() (*peer, error)
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
		msgCh:         pubsubManager.Register(config.PeerCandidatesBufferSize),
		closeCh:       make(chan struct{}),
	}
	go m.processPeerCandidates()
	return m
}

func (m *manager) GetConnected() (*peer, error) {
	select {
	case connected := <-m.connectedCh:
		return connected, nil
	case <-time.After(waitForconnectedTimeout):
		return nil, errors.New("could not find any connected peers")
	}
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

				// already have as many connected pers as we want
				if len(m.connectedCh) == config.MaxConnectedPeers {
					continue
				}

				m.log.Debug("new peer candidate event")
				go func() {
					peer := NewPeer(candidate.Address)
					if err := peer.Connect(); err != nil {
						m.log.WithField("peer", candidate).Warn(err)
						return
					}
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