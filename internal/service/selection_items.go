package service

import "rogue/internal/domain"

type ItemSelectionItem struct {
	Item       domain.Item
	ItemKind   domain.ItemKind
	StackIndex int
}

func (i *ItemSelectionItem) DisplayText() string {
	return i.Item.Name()
}

func (i *ItemSelectionItem) Data() interface{} {
	return i
}

type EquipmentSelectionItem struct {
	Part domain.ActorPart
	Item domain.Item
}

func (e *EquipmentSelectionItem) DisplayText() string {
	partNames := map[domain.ActorPart]string{
		domain.ActorPartHead: "Head",
		domain.ActorPartBody: "Body",
		domain.ActorPartHand: "Hand",
		domain.ActorPartLegs: "Legs",
	}
	partName := partNames[e.Part]
	if e.Item == nil {
		return partName + ": -"
	}
	return partName + ": " + e.Item.Name()
}

func (e *EquipmentSelectionItem) Data() interface{} {
	return e
}
