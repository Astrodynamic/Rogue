package domain

type Armor struct {
	ArmorPart ActorPart
	ItemName  string
	Stats     Stats
}

func (a *Armor) Type() ItemKind {
	return ItemArmor
}

func (a *Armor) Name() string {
	if a.ItemName != "" {
		return a.ItemName
	}
	switch a.ArmorPart {
	case ActorPartHead:
		return "Helmet"
	case ActorPartBody:
		return "Armor"
	case ActorPartLegs:
		return "Boots"
	default:
		return "Armor"
	}
}

func (a *Armor) Stackable() bool {
	return false
}

func (a *Armor) MaxStack() int {
	return 1
}

func (a *Armor) Equals(other Item) bool {
	if other.Type() != ItemArmor {
		return false
	}
	otherArmor, ok := other.(*Armor)
	if !ok {
		return false
	}
	return a.ArmorPart == otherArmor.ArmorPart && a.Stats == otherArmor.Stats
}

func (a *Armor) Use() ItemUseResult {
	return ItemUseResult{
		Success:  true,
		Consumed: false,
		Message:  a.Name() + " equipped",
		Effects:  nil,
	}
}

func (a *Armor) GetStats() Stats {
	return a.Stats
}
