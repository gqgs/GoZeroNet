package event

type FileInfo struct {
	InnerPath    string `json:"inner_path"`
	Hash         string `json:"hash"`
	Size         int    `json:"size"`
	IsDownloaded bool   `json:"is_downloaded"`
	IsPinned     bool   `json:"is_pinned"`
	IsOptional   bool   `json:"is_optional"`
	Uploaded     int    `json:"uploaded"`
}

func (e *FileInfo) String() string {
	return "fileInfo"
}

func BroadcastFileInfoUpdate(site string, broadcaster EventBroadcaster, fileInfo *FileInfo) {
	go broadcaster.Broadcast(site, fileInfo)
}

type FileNeed struct {
	InnerPath string `json:"inner_path"`
}

func (e *FileNeed) String() string {
	return "fileNeed"
}

func BroadcastFileNeed(site string, broadcaster EventBroadcaster, fileNeed *FileNeed) {
	go broadcaster.Broadcast(site, fileNeed)
}
