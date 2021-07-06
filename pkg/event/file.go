package event

type FileInfo struct {
	InnerPath         string  `json:"inner_path"`
	Hash              string  `json:"hash"`
	Size              int     `json:"size"`
	IsDownloaded      bool    `json:"is_downloaded"`
	IsPinned          bool    `json:"is_pinned"`
	IsOptional        bool    `json:"is_optional"`
	Uploaded          int     `json:"uploaded"`
	PieceSize         int     `json:"piece_size"`
	Piecemap          string  `json:"piecemap"`
	Downloaded        int     `json:"downloaded"`
	DownloadedPercent float64 `json:"downloaded_percent"`
}

func (e *FileInfo) IsBigFile() bool {
	return e.PieceSize > 0
}

func (e *FileInfo) String() string {
	return "fileInfo"
}

func BroadcastFileInfoUpdate(site string, broadcaster EventBroadcaster, fileInfo *FileInfo) {
	go broadcaster.Broadcast(site, fileInfo)
}

type FileNeed struct {
	InnerPath string `json:"inner_path"`
	Tries     int    `json:"tries"`
}

func (e *FileNeed) String() string {
	return "fileNeed"
}

func BroadcastFileNeed(site string, broadcaster EventBroadcaster, fileNeed *FileNeed) {
	go broadcaster.Broadcast(site, fileNeed)
}
