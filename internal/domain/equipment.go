package domain

type Equipment struct {
	Parts map[ActorPart]*Item
}

func NewEquipment() *Equipment {
	return &Equipment{
		Parts: make(map[ActorPart]*Item),
	}
}

func (e *Equipment) Get(part ActorPart) *Item {
	return e.Parts[part]
}

func (e *Equipment) Set(part ActorPart, item *Item) {
	e.Parts[part] = item
}
