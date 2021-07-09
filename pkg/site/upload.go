package site

import "time"

type Upload struct {
	Added     time.Time
	Site      string
	InnerPath string
	Size      int
	PieceSize int
	Piecemap  string
}

func (s *Site) UploadInit(upload Upload, nonce string) {
	s.uploadMutex.Lock()
	defer s.uploadMutex.Unlock()
	s.uploads[nonce] = upload
}
