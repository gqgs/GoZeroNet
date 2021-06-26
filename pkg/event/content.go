package event

type ContentInfo struct {
	InnerPath string
	Modified  int
	Size      int
}

func (e *ContentInfo) String() string {
	return "contentInfo"
}

func BroadcastContentInfoUpdate(site string, broadcaster EventBroadcaster, contentInfo *ContentInfo) {
	go broadcaster.Broadcast(site, contentInfo)
}
