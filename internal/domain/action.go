package domain

type ActionKind uint8

const (
	ActionNone ActionKind = iota
	ActionMove
	ActionAttack
	ActionUseFood
	ActionUseElixir
	ActionUseScroll
	ActionEquipWeapon
	ActionUnequipWeapon
)

type Action struct {
	Kind ActionKind
	Dx   int
	Dy   int
	Idx  int // menu selection (0-based)
}
