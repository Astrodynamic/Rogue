package service

type SelectionModel struct {
	items         []SelectionItem
	selectedIndex int
	handler       SelectionHandler
}

type SelectionItem interface {
	DisplayText() string
	Data() interface{}
}

type SelectionHandler interface {
	OnConfirm(item SelectionItem) bool
	OnCancel()
}

func NewSelectionModel(items []SelectionItem, handler SelectionHandler) *SelectionModel {
	return &SelectionModel{
		items:         items,
		selectedIndex: 0,
		handler:       handler,
	}
}

func (m *SelectionModel) IsActive() bool {
	return m != nil && len(m.items) > 0
}

func (m *SelectionModel) Count() int {
	if m == nil {
		return 0
	}
	return len(m.items)
}

func (m *SelectionModel) SelectedIndex() int {
	if m == nil {
		return -1
	}
	return m.selectedIndex
}

func (m *SelectionModel) SelectedItem() SelectionItem {
	if m == nil || m.selectedIndex < 0 || m.selectedIndex >= len(m.items) {
		return nil
	}
	return m.items[m.selectedIndex]
}

func (m *SelectionModel) Items() []SelectionItem {
	if m == nil {
		return nil
	}
	return m.items
}

func (m *SelectionModel) MoveUp() {
	if m == nil || m.selectedIndex <= 0 {
		return
	}
	m.selectedIndex--
}

func (m *SelectionModel) MoveDown() {
	if m == nil || m.selectedIndex >= len(m.items)-1 {
		return
	}
	m.selectedIndex++
}

func (m *SelectionModel) Confirm() bool {
	if m == nil || m.handler == nil {
		return false
	}
	item := m.SelectedItem()
	if item == nil {
		return false
	}
	return m.handler.OnConfirm(item)
}

func (m *SelectionModel) Cancel() {
	if m == nil || m.handler == nil {
		return
	}
	m.handler.OnCancel()
}

func (m *SelectionModel) Reset() {
	if m == nil {
		return
	}
	m.items = nil
	m.selectedIndex = 0
	m.handler = nil
}
