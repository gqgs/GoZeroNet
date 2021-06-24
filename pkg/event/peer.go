package event

// General peer info
type PeerInfo struct {
	Address         string
	ReputationDelta int
}

func (e *PeerInfo) String() string {
	return "peerInfo"
}

func BroadcastPeerInfoUpdate(site string, broadcaster EventBroadcaster, peerInfo *PeerInfo) {
	broadcaster.Broadcast(site, peerInfo)
}

// Peers that might be connected to
type PeerCandidate struct {
	Address string
}

func (e *PeerCandidate) String() string {
	return "peerCandidate"
}

func BroadcastPeerCandidate(site string, broadcaster EventBroadcaster, peerCandidate *PeerCandidate) {
	broadcaster.Broadcast(site, peerCandidate)
}

// Signals that the site needs more peers
type PeersNeed struct{}

func (e *PeersNeed) String() string {
	return "PeersNeed"
}

func BroadcastPeersNeed(site string, broadcaster EventBroadcaster, peersNeed *PeersNeed) {
	broadcaster.Broadcast(site, peersNeed)
}
