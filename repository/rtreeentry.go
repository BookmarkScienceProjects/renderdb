package repository

import "github.com/dhconnelly/rtreego"

type rtreeEntry struct {
	id     int64
	bounds *rtreego.Rect
}

func (e *rtreeEntry) Bounds() *rtreego.Rect {
	return e.bounds
}
