package domain

type ItemContainer interface {
	Add(Item) bool
	TreasureAmount() int
	AddTreasure(value int)
	Count(ItemType) int
	List(ItemType) []Item
	RemoveAt(ItemType, idx int) (Item, bool)
	FindByID(id int64) (Item, bool)
	RemoveByID(id int64) (Item, bool)
}

type Backpack struct {
	Treasure int
	Items    []Item
}

const MaxItemsPerType = 9

func (b *Backpack) Add(it Item) bool {
	switch it.Type {
	case ItemTreasure:
		b.Treasure += it.Value
		return true
	case ItemFood, ItemElixir, ItemScroll, ItemWeapon:
		if b.Count(it.Type) >= MaxItemsPerType {
			return false
		}
		b.Items = append(b.Items, it)
		return true
	default:
		return false
	}
}

func (b *Backpack) TreasureAmount() int { return b.Treasure }

func (b *Backpack) AddTreasure(v int) {
	if v <= 0 {
		return
	}
	b.Treasure += v
}

func (b *Backpack) Count(t ItemType) int {
	if t == ItemTreasure {
		if b.Treasure > 0 {
			return 1
		}
		return 0
	}
	n := 0
	for i := range b.Items {
		if b.Items[i].Type == t {
			n++
		}
	}
	return n
}

func (b *Backpack) List(t ItemType) []Item {
	if t == ItemTreasure {
		return nil
	}
	out := make([]Item, 0, b.Count(t))
	for i := range b.Items {
		if b.Items[i].Type == t {
			out = append(out, b.Items[i])
		}
	}
	return out
}

func (b *Backpack) RemoveAt(t ItemType, idx int) (Item, bool) {
	if t == ItemTreasure {
		return Item{}, false
	}
	if idx < 0 {
		return Item{}, false
	}
	seen := 0
	for i := 0; i < len(b.Items); i++ {
		if b.Items[i].Type != t {
			continue
		}
		if seen == idx {
			it := b.Items[i]
			copy(b.Items[i:], b.Items[i+1:])
			b.Items = b.Items[:len(b.Items)-1]
			return it, true
		}
		seen++
	}
	return Item{}, false
}

func (b *Backpack) FindByID(id int64) (Item, bool) {
	if id == 0 {
		return Item{}, false
	}
	for i := range b.Items {
		if b.Items[i].ID == id {
			return b.Items[i], true
		}
	}
	return Item{}, false
}

func (b *Backpack) RemoveByID(id int64) (Item, bool) {
	if id == 0 {
		return Item{}, false
	}
	for i := 0; i < len(b.Items); i++ {
		if b.Items[i].ID == id {
			it := b.Items[i]
			copy(b.Items[i:], b.Items[i+1:])
			b.Items = b.Items[:len(b.Items)-1]
			return it, true
		}
	}
	return Item{}, false
}
