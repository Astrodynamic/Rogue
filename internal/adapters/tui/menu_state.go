package tui

type MenuState struct {
	showStartMenu  bool
	showStatistics bool
	menuOption     int
}

func NewMenuState() *MenuState {
	return &MenuState{
		showStartMenu:  true,
		showStatistics: false,
		menuOption:     0,
	}
}

func (m *MenuState) IsStartMenu() bool {
	return m.showStartMenu
}

func (m *MenuState) IsStatistics() bool {
	return m.showStatistics
}

func (m *MenuState) GetMenuOption() int {
	return m.menuOption
}

func (m *MenuState) SetStartMenu(show bool) {
	m.showStartMenu = show
}

func (m *MenuState) SetStatistics(show bool) {
	m.showStatistics = show
}

func (m *MenuState) SetMenuOption(option int) {
	m.menuOption = option
}

func (m *MenuState) MoveUp(maxOptions int) {
	if m.menuOption > 0 {
		m.menuOption--
	}
}

func (m *MenuState) MoveDown(maxOptions int) {
	if m.menuOption < maxOptions-1 {
		m.menuOption++
	}
}

func (m *MenuState) Reset() {
	m.showStartMenu = true
	m.showStatistics = false
	m.menuOption = 0
}
