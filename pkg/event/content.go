package event

type ContentInfo struct {
	InnerPath string
	Modified  int
	Size      int
}

func (e *ContentInfo) AddUploaded(uploaded int) {

}

func (e *ContentInfo) GetIsDownloaded() bool {
	return true
}

func (e *ContentInfo) Update(site string, broadcaster Broadcaster) {
	go broadcaster.Broadcast(site, e)
}

func (e *ContentInfo) GetSize() int {
	return e.Size
}

func (e *ContentInfo) String() string {
	return "contentInfo"
}

func BroadcastContentInfoUpdate(site string, broadcaster Broadcaster, contentInfo *ContentInfo) {
	go broadcaster.Broadcast(site, contentInfo)
}
