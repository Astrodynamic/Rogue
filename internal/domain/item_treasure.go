package domain

type Treasure struct {
	Value int
}

func (t *Treasure) Type() ItemKind {
	return ItemTreasure
}

func (t *Treasure) Name() string {
	return "Treasure"
}

func (t *Treasure) Stackable() bool {
	return true
}

func (t *Treasure) MaxStack() int {
	return 0
}

func (t *Treasure) Equals(other Item) bool {
	if other.Type() != ItemTreasure {
		return false
	}
	otherTreasure, ok := other.(*Treasure)
	if !ok {
		return false
	}
	return t.Value == otherTreasure.Value
}

func (t *Treasure) Use() ItemUseResult {

	return ItemUseResult{
		Success:  false,
		Consumed: false,
		Message:  t.Name() + " cannot be used",
		Effects:  nil,
	}
}
