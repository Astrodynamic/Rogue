package game

type ItemType uint8

const (
	ItemTreasure ItemType = iota
	ItemFood
	ItemElixir
	ItemScroll
	ItemWeapon
)

type ItemSubtype uint8

const (
	SubtypeNone ItemSubtype = iota
	SubtypeDexterity
	SubtypeStrength
	SubtypeMaxHealth
)

type Item struct {
	ID          int64
	Type        ItemType
	Subtype     ItemSubtype
	Health      int
	MaxHealth   int
	Dexterity   int
	Strength    int
	Value       int
	Temporary   bool // for elixirs; scrolls are permanent
	Duration    int  // turns, if Temporary
	Description string
}

func (it Item) DisplayName() string {
	if it.Description != "" {
		return it.Description
	}
	switch it.Type {
	case ItemFood:
		return "Food"
	case ItemElixir:
		return "Elixir"
	case ItemScroll:
		return "Scroll"
	case ItemWeapon:
		return "Weapon"
	case ItemTreasure:
		return "Treasure"
	default:
		return "Item"
	}
}
