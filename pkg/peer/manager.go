package peer

import (
	"context"
	"errors"
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
	GetConnected(ctx context.Context) (*peer, error)
	PutConnected(p *peer)
	Close()
}

type manager struct {
	log           log.Logger
	pubsubManager pubsub.Manager
	site          string
	connectedSem  chan struct{}
	connectedCh   chan *peer
	msgCh         <-chan pubsub.Message
	closeCh       chan struct{}
}

func NewManager(pubsubManager pubsub.Manager, site string) *manager {
	m := &manager{
		log:           log.New("peer_manager"),
		site:          site,
		pubsubManager: pubsubManager,
		connectedSem:  make(chan struct{}, config.MaxConnectedPeers),
		connectedCh:   make(chan *peer, config.MaxConnectedPeers),
		msgCh:         pubsubManager.Register("peer_manager", config.MaxConnectedPeers),
		closeCh:       make(chan struct{}),
	}
	go m.processPeerCandidates()
	return m
}

func (m *manager) GetConnected(ctx context.Context) (*peer, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("context canceled")
	case connected := <-m.connectedCh:
		if err := connected.CheckConnection(); err == nil {
			return connected, nil
		}
		connected.Close()
		<-m.connectedSem
	case <-time.After(waitForconnectedTimeout):
		event.BroadcastPeersNeed(m.site, m.pubsubManager, &event.PeersNeed{})
		time.Sleep(time.Second * 5)
	}
	return m.GetConnected(ctx)
}

func (m *manager) PutConnected(p *peer) {
	select {
	case m.connectedCh <- p:
	default:
		p.Close()
		<-m.connectedSem
	}
}

func (m *manager) Close() {
	m.log.Debug("closing")
	close(m.closeCh)
	m.pubsubManager.Unregister(m.msgCh)
}

func (m *manager) processPeerCandidates() {
	ctx, cancel := context.WithCancel(context.Background())
	for {
		select {
		case <-m.closeCh:
			cancel()
			return
		case msg := <-m.msgCh:
			if msg.Site() != m.site {
				continue
			}
			switch candidate := msg.Event().(type) {
			case *event.PeerCandidate:
				m.log.WithField("queue", len(m.msgCh)).Debug("new peer candidate event")
				if candidate.IsOnion && !config.TorEnabled {
					continue
				}

				go func() {
					select {
					case <-ctx.Done():
					case <-time.After(time.Minute):
					case m.connectedSem <- struct{}{}:
						peer := NewPeer(candidate.Address)
						if err := peer.Connect(); err != nil {
							peer.Warn(err)
							<-m.connectedSem
							return
						}

						peer.Info("new connected peer: ", peer)
						m.connectedCh <- peer
					}
				}()
			}
		}
	}
}
