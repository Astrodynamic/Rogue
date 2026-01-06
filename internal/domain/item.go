package domain

type ItemKind int

const (
	ItemFood ItemKind = iota
	ItemElixir
	ItemScroll
	ItemWeapon
	ItemArmor
	ItemTreasure
)

type Item interface {
	Type() ItemKind
	Name() string
	Stackable() bool
	MaxStack() int
	Equals(other Item) bool
	Use() ItemUseResult
	GetStats() Stats
}

type ItemUseResult struct {
	Success  bool
	Consumed bool
	Message  string
	Error    error
	Effects  []Effect
}
