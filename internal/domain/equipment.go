package domain

type Equipment struct {
	Parts map[ActorPart]Item
}

func NewEquipment() *Equipment {
	return &Equipment{
		Parts: make(map[ActorPart]Item),
	}
}

func (e *Equipment) Get(part ActorPart) Item {
	return e.Parts[part]
}

func (e *Equipment) Set(part ActorPart, item Item) {
	e.Parts[part] = item
}

func GetEquipPart(item Item) (ActorPart, bool) {
	equippable, ok := item.(Equippable)
	if !ok {
		return 0, false
	}
	return equippable.EquipPart(), true
}

func (e *Equipment) Equip(part ActorPart, item Item) Item {
	oldItem := e.Parts[part]
	e.Parts[part] = item
	return oldItem
}

func (e *Equipment) Unequip(part ActorPart) Item {
	item := e.Parts[part]
	if item != nil {
		delete(e.Parts, part)
	}
	return item
}

func (e *Equipment) GetTotalStats() Stats {
	total := Stats{}
	for _, item := range e.Parts {
		if item != nil {
			stats := item.GetStats()
			total.AddMaxHealth(stats.MaxHealth)
			total.AddHealth(stats.Health)
			total.AddDexterity(stats.Dexterity)
			total.AddStrength(stats.Strength)
		}
	}
	return total
}

func GetAllParts() []ActorPart {
	return []ActorPart{
		ActorPartHead,
		ActorPartBody,
		ActorPartHand,
		ActorPartLegs,
	}
}

func GetPartName(part ActorPart) string {
	partNames := map[ActorPart]string{
		ActorPartHead: "Head",
		ActorPartBody: "Body",
		ActorPartHand: "Hand",
		ActorPartLegs: "Legs",
	}
	return partNames[part]
}
