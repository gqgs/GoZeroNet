package event

type ContentInfo struct {
	InnerPath string
	Modified  int
}

func (e *ContentInfo) String() string {
	return "contentInfo"
}

func BroadcastContentInfoUpdate(site string, broadcaster EventBroadcaster, contentInfo *ContentInfo) {
	broadcaster.Broadcast(site, contentInfo)
}
