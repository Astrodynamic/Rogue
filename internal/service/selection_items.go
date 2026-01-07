package service

import "rogue/internal/domain"

type ItemSelItem struct {
	Item       domain.Item
	ItemKind   domain.ItemKind
	StackIndex int
}

func (i *ItemSelItem) DisplayText() string {
	return i.Item.Name()
}

func (i *ItemSelItem) Data() interface{} {
	return i
}

type EquipSelItem struct {
	Part domain.ActorPart
	Item domain.Item
}

func (e *EquipSelItem) DisplayText() string {
	partName := domain.GetPartName(e.Part)
	if e.Item == nil {
		return partName + ": -"
	}
	return partName + ": " + e.Item.Name()
}

func (e *EquipSelItem) Data() interface{} {
	return e
}
