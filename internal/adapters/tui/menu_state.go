package tui

type MenuState struct {
	showStartMenu  bool
	showStatistics bool
	showNameInput  bool
	menuOption     int
	playerName     string
}

func NewMenuState() *MenuState {
	return &MenuState{
		showStartMenu:  true,
		showStatistics: false,
		showNameInput:  false,
		menuOption:     0,
		playerName:     "",
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

func (m *MenuState) IsNameInput() bool {
	return m.showNameInput
}

func (m *MenuState) SetNameInput(show bool) {
	m.showNameInput = show
	if !show {
		m.playerName = ""
	}
}

func (m *MenuState) GetPlayerName() string {
	return m.playerName
}

func (m *MenuState) SetPlayerName(name string) {
	m.playerName = name
}

func (m *MenuState) AppendToPlayerName(char rune) {
	runes := []rune(m.playerName)
	if len(runes) < 20 {
		m.playerName += string(char)
	}
}

func (m *MenuState) DeleteLastChar() {
	runes := []rune(m.playerName)
	if len(runes) > 0 {
		m.playerName = string(runes[:len(runes)-1])
	}
}

func (m *MenuState) Reset() {
	m.showStartMenu = true
	m.showStatistics = false
	m.showNameInput = false
	m.menuOption = 0
	m.playerName = ""
}
