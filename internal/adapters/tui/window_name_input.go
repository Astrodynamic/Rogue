package tui

import (
	"rogue/internal/service"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawNameInput() {
	wW, wH := w.screen.Size()

	boxWidth := 60
	boxHeight := 10
	boxX := (wW - boxWidth) / 2
	boxY := (wH - boxHeight) / 2

	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	textStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	labelStyle := tcell.StyleDefault.Foreground(tcell.ColorAqua)
	inputStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)

	topBorder := "+"
	for i := 0; i < boxWidth-2; i++ {
		topBorder += "-"
	}
	topBorder += "+"

	sideBorder := "|"

	for i := 0; i < boxHeight; i++ {
		for j := 0; j < boxWidth; j++ {
			var char rune = ' '
			style := borderStyle

			if i == 0 || i == boxHeight-1 {
				if j >= 0 && j < len(topBorder) {
					char = rune(topBorder[j])
				}
			} else if j == 0 || j == boxWidth-1 {
				char = rune(sideBorder[0])
			}

			if boxX+j >= 0 && boxX+j < wW && boxY+i >= 0 && boxY+i < wH {
				w.screen.SetContent(boxX+j, boxY+i, char, nil, style)
			}
		}
	}

	title := "Enter Player Name"
	titleX := boxX + (boxWidth-len(title))/2
	titleY := boxY + 1
	for i, r := range title {
		if titleX+i >= 0 && titleX+i < wW {
			w.screen.SetContent(titleX+i, titleY, r, nil, labelStyle.Bold(true))
		}
	}

	label := "Name:"
	labelX := boxX + 4
	labelY := boxY + 3
	for i, r := range label {
		if labelX+i >= 0 && labelX+i < wW {
			w.screen.SetContent(labelX+i, labelY, r, nil, labelStyle)
		}
	}

	inputAreaX := labelX + len(label) + 2
	inputAreaY := boxY + 3
	inputAreaWidth := boxWidth - (inputAreaX - boxX) - 4
	playerName := w.menuState.GetPlayerName()

	nameRunes := []rune(playerName)
	cursorPos := len(nameRunes)
	maxVisible := inputAreaWidth - 1
	startIdx := 0
	if cursorPos > maxVisible {
		startIdx = cursorPos - maxVisible
	}

	displayRunes := nameRunes[startIdx:]
	if len(displayRunes) > maxVisible {
		displayRunes = displayRunes[:maxVisible]
	}

	for i := 0; i < inputAreaWidth; i++ {
		var char rune = ' '
		style := textStyle
		if i < len(displayRunes) {
			char = displayRunes[i]
			style = inputStyle
		} else if i == len(displayRunes) && cursorPos < 20 {
			char = '_'
			style = tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
		}

		if inputAreaX+i >= 0 && inputAreaX+i < wW {
			w.screen.SetContent(inputAreaX+i, inputAreaY, char, nil, style)
		}
	}

	if cursorPos < 20 {
		cursorScreenX := inputAreaX + len(displayRunes)
		if cursorScreenX >= 0 && cursorScreenX < wW && cursorScreenX < inputAreaX+inputAreaWidth {
			cursorChar := '_'
			w.screen.SetContent(cursorScreenX, inputAreaY, cursorChar, nil, inputStyle.Bold(true).Reverse(true))
		}
	}

	instructions := []string{
		"Type your name and press ENTER to confirm",
		"Press ESC to cancel (will use 'Player')",
	}
	instY := boxY + 5
	for lineIdx, instruction := range instructions {
		instX := boxX + (boxWidth-len(instruction))/2
		for i, r := range instruction {
			if instX+i >= 0 && instX+i < wW && instY+lineIdx < wH {
				w.screen.SetContent(instX+i, instY+lineIdx, r, nil, textStyle)
			}
		}
	}
}

func (w *Window) HandleNameInput(cmd service.Command, char rune) {
	if cmd == service.CmdSelectConfirm {
		name := w.menuState.GetPlayerName()
		if name == "" {
			name = "Player"
		}
		w.menuState.SetPlayerName(name)
		return
	}

	if cmd == service.CmdSelectCancel || cmd == service.CmdQuit {
		w.menuState.SetPlayerName("Player")
		w.menuState.SetNameInput(false)
		return
	}

	if cmd == service.CmdBackToGame {
		w.menuState.DeleteLastChar()
		return
	}

	if char != 0 && (char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == ' ' || char == '_' || char == '-') {
		w.menuState.AppendToPlayerName(char)
	}
}
