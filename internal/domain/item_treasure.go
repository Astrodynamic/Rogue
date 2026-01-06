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
	return other.Type() == ItemTreasure
}

func (t *Treasure) Use() ItemUseResult {

	return ItemUseResult{
		Success:  false,
		Consumed: false,
		Message:  t.Name() + " cannot be used",
		Effects:  nil,
	}
}

func (t *Treasure) GetStats() Stats {
	return Stats{}
}
