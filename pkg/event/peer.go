package event

// General peer info
type PeerInfo struct {
	Address         string
	ReputationDelta int
}

func (e *PeerInfo) String() string {
	return "peerInfo"
}

func BroadcastPeerInfoUpdate(site string, broadcaster Broadcaster, peerInfo *PeerInfo) {
	go broadcaster.Broadcast(site, peerInfo)
}

// Peers that might be connected to
type PeerCandidate struct {
	Address string
	IsOnion bool
}

func (e *PeerCandidate) String() string {
	return "peerCandidate"
}

func BroadcastPeerCandidate(site string, broadcaster Broadcaster, peerCandidate *PeerCandidate) {
	go broadcaster.Broadcast(site, peerCandidate)
}

// Signals that the site needs more peers
type PeersNeed struct{}

func (e *PeersNeed) String() string {
	return "peersNeed"
}

func BroadcastPeersNeed(site string, broadcaster Broadcaster, peersNeed *PeersNeed) {
	go broadcaster.Broadcast(site, peersNeed)
}
