package domain

type Room struct {
	Rect
}

func NewRoom(x, y, w, h int) Room {
	return Room{Rect: NewRect(x, y, w, h)}
}
