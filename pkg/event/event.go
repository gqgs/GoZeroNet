package event

import "encoding/json"

type EventBroadcaster interface {
	Broadcast(site, event string, body []byte)
}

type FileInfo struct {
	InnerPath    string
	Hash         string
	Size         int
	IsDownloaded bool
}

func BroadcastFileDone(site string, broadcaster EventBroadcaster, fileInfo *FileInfo) {
	body, _ := json.Marshal(fileInfo)
	broadcaster.Broadcast(site, "file-done", body)
}

type PeerInfo struct {
	Address         string
	ReputationDelta int
}

func BroadcastPeerInfoUpdate(site string, broadcaster EventBroadcaster, peerInfo *PeerInfo) {
	body, _ := json.Marshal(peerInfo)
	broadcaster.Broadcast(site, "peer-info", body)
}
