package domain

const (
	MaxUniqueItems = 20
)

type Backpack struct {
	stacks map[ItemKind][]*ItemStack
}

func NewBackpack() *Backpack {
	return &Backpack{
		stacks: make(map[ItemKind][]*ItemStack),
	}
}

func (b *Backpack) CountUniqueItems() int {
	total := 0
	for _, stacks := range b.stacks {
		total += len(stacks)
	}
	return total
}

func (b *Backpack) Add(item Item) bool {
	itemType := item.Type()

	if stacks, exists := b.stacks[itemType]; exists {
		for _, stack := range stacks {
			if stack.CanAdd(item) {
				stack.Add(item)
				return true
			}
		}
	}

	if b.CountUniqueItems() >= MaxUniqueItems {
		return false
	}

	if _, exists := b.stacks[itemType]; !exists {
		b.stacks[itemType] = make([]*ItemStack, 0)
	}

	newStack := NewItemStack(item)
	b.stacks[itemType] = append(b.stacks[itemType], newStack)
	return true
}

func (b *Backpack) Remove(itemKind ItemKind, stackIndex int) Item {
	if stacks, exists := b.stacks[itemKind]; exists {
		if stackIndex < 0 || stackIndex >= len(stacks) {
			return nil
		}

		stack := stacks[stackIndex]
		removedItem := stack.Remove()

		if stack.IsEmpty() {
			b.stacks[itemKind] = append(stacks[:stackIndex], stacks[stackIndex+1:]...)
			if len(b.stacks[itemKind]) == 0 {
				delete(b.stacks, itemKind)
			}
		}

		return removedItem
	}
	return nil
}

func (b *Backpack) GetStack(itemKind ItemKind, stackIndex int) *ItemStack {
	if stacks, exists := b.stacks[itemKind]; exists {
		if stackIndex >= 0 && stackIndex < len(stacks) {
			return stacks[stackIndex]
		}
	}
	return nil
}

func (b *Backpack) GetStacks(itemKind ItemKind) []*ItemStack {
	if stacks, exists := b.stacks[itemKind]; exists {
		return stacks
	}
	return []*ItemStack{}
}

func (b *Backpack) Count(itemKind ItemKind) int {
	total := 0
	if stacks, exists := b.stacks[itemKind]; exists {
		for _, stack := range stacks {
			total += stack.Count
		}
	}
	return total
}

func (b *Backpack) HasSpace(item Item) bool {
	itemKind := item.Type()

	if stacks, exists := b.stacks[itemKind]; exists {
		for _, stack := range stacks {
			if stack.CanAdd(item) {
				return true
			}
		}

		if b.CountUniqueItems() >= MaxUniqueItems {
			return false
		}
		return true
	}

	if b.CountUniqueItems() >= MaxUniqueItems {
		return false
	}
	return true
}

func (b *Backpack) Use(itemKind ItemKind, stackIndex int, actor *Actor) ItemUseResult {
	stack := b.GetStack(itemKind, stackIndex)
	if stack == nil || stack.IsEmpty() {
		return ItemUseResult{
			Success: false,
			Message: "Item not found",
		}
	}

	item := stack.Peek()
	result := actor.UseItem(item)

	if result.Consumed {
		b.Remove(itemKind, stackIndex)
	}

	return result
}

func (b *Backpack) GetTotalTreasure() int {
	if stacks, exists := b.stacks[ItemTreasure]; exists && len(stacks) > 0 {
		if treasure, ok := stacks[0].Item.(*Treasure); ok {
			return treasure.Value
		}
	}
	return 0
}

func (b *Backpack) GetAllStacks() map[ItemKind][]*ItemStack {
	result := make(map[ItemKind][]*ItemStack)
	for kind, stacks := range b.stacks {
		result[kind] = stacks
	}
	return result
}

func (b *Backpack) AddStack(kind ItemKind, stack *ItemStack) {
	if _, exists := b.stacks[kind]; !exists {
		b.stacks[kind] = make([]*ItemStack, 0)
	}
	b.stacks[kind] = append(b.stacks[kind], stack)
}
