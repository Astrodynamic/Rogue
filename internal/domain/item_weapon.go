package domain

type Weapon struct {
	Strength int
}

func (w *Weapon) Type() ItemKind {
	return ItemWeapon
}

func (w *Weapon) Stackable() bool {
	return false
}

func (w *Weapon) MaxStack() int {
	return 1
}

func (w *Weapon) Equals(other Item) bool {
	if other.Type() != ItemWeapon {
		return false
	}
	otherWeapon, ok := other.(*Weapon)
	if !ok {
		return false
	}
	return w.Strength == otherWeapon.Strength
}

func (w *Weapon) Use() ItemUseResult {
	return ItemUseResult{
		Success:  true,
		Consumed: false,
		Message:  "Weapon equipped",
		Effects:  nil,
	}
}
