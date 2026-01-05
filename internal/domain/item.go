package domain

type ItemKind int

const (
	ItemFood ItemKind = iota
	ItemElixir
	ItemScroll
	ItemWeapon
	ItemTreasure
)

type Item interface {
	Type() ItemKind
	Name() string
	Stackable() bool
	MaxStack() int
	Equals(other Item) bool
	Use() ItemUseResult
}

type ItemUseResult struct {
	Success  bool
	Consumed bool
	Message  string
	Error    error
	Effects  []Effect
}
