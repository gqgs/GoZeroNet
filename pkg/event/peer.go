package event

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
