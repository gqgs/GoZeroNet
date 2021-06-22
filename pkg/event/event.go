package event

type Event interface {
	String() string
}

type EventBroadcaster interface {
	Broadcast(site string, event Event)
}

type FileInfo struct {
	InnerPath    string
	Hash         string
	Size         int
	IsDownloaded bool
}

func (e *FileInfo) String() string {
	return "fileInfo"
}

func BroadcastFileDone(site string, broadcaster EventBroadcaster, fileInfo *FileInfo) {
	broadcaster.Broadcast(site, fileInfo)
}

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

type SiteChanged struct {
	Params interface{} `json:"params"`
	Cmd    string      `json:"cmd"`
}

func (e *SiteChanged) String() string {
	return "siteChanged"
}

func BroadcastSiteChanged(site string, broadcaster EventBroadcaster, siteChanged *SiteChanged) {
	broadcaster.Broadcast(site, siteChanged)
}
