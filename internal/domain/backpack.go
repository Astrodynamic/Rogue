package domain

type Backpack struct {
	Items []Item
}

func NewBackpack() *Backpack {
	return &Backpack{Items: make([]Item, 0)}
}
