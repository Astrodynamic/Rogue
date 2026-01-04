package domain

type Food struct {
	Health int
}

func (f *Food) Type() ItemKind {
	return ItemFood
}

func (f *Food) Stackable() bool {
	return true
}

func (f *Food) MaxStack() int {
	return 9
}

func (f *Food) Equals(other Item) bool {
	if other.Type() != ItemFood {
		return false
	}
	otherFood, ok := other.(*Food)
	if !ok {
		return false
	}
	return f.Health == otherFood.Health
}

func (f *Food) Use() ItemUseResult {
	return ItemUseResult{
		Success:  true,
		Consumed: true,
		Message:  "Food consumed",
		Effects: []Effect{
			&HealthEffect{
				BaseEffect: BaseEffect{
					duration: 0,
				},
				Stats: Stats{Health: f.Health},
			},
		},
	}
}
