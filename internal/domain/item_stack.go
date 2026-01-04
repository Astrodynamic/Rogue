package domain

type ItemStack struct {
	Item  Item
	Count int
}

func NewItemStack(item Item) *ItemStack {
	return &ItemStack{
		Item:  item,
		Count: 1,
	}
}

func (s *ItemStack) CanAdd(item Item) bool {
	if !s.Item.Stackable() {
		return false
	}
	if !s.Item.Equals(item) {
		return false
	}
	maxStack := s.Item.MaxStack()
	if maxStack == 0 {
		return true
	}
	return s.Count < maxStack
}

func (s *ItemStack) Add(item Item) bool {
	if !s.CanAdd(item) {
		return false
	}
	s.Count++
	return true
}

func (s *ItemStack) Remove() Item {
	if s.Count <= 0 {
		return nil
	}
	s.Count--
	if s.Count == 0 {
		return s.Item
	}
	return nil
}

func (s *ItemStack) IsEmpty() bool {
	return s.Count <= 0
}

func (s *ItemStack) Peek() Item {
	if s.IsEmpty() {
		return nil
	}
	return s.Item
}
