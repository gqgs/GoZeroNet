package event

type FileInfo struct {
	InnerPath    string `json:"inner_path"`
	Hash         string `json:"hash"`
	Size         int    `json:"size"`
	IsDownloaded bool   `json:"is_downloaded"`
}

func (e *FileInfo) String() string {
	return "fileInfo"
}

func BroadcastFileDone(site string, broadcaster EventBroadcaster, fileInfo *FileInfo) {
	broadcaster.Broadcast(site, fileInfo)
}
