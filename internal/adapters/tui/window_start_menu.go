package tui

import (
	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawStartMenu(hasSave bool) {
	w.screen.Clear()
	wW, wH := w.screen.Size()

	// Title with decorative box
	title := "ROGUE"
	subtitle := "A Dungeon Crawler Adventure"
	
	titleBoxWidth := len(title) + 4
	titleBoxHeight := 5
	titleBoxX := (wW - titleBoxWidth) / 2
	titleBoxY := wH/2 - 8

	// Draw title box
	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkRed)
	
	// Top border
	for i := 0; i < titleBoxWidth; i++ {
		if i == 0 {
			w.screen.SetContent(titleBoxX+i, titleBoxY, '╔', nil, borderStyle)
		} else if i == titleBoxWidth-1 {
			w.screen.SetContent(titleBoxX+i, titleBoxY, '╗', nil, borderStyle)
		} else {
			w.screen.SetContent(titleBoxX+i, titleBoxY, '═', nil, borderStyle)
		}
	}
	
	// Side borders and title
	for i := 1; i < titleBoxHeight-1; i++ {
		w.screen.SetContent(titleBoxX, titleBoxY+i, '║', nil, borderStyle)
		w.screen.SetContent(titleBoxX+titleBoxWidth-1, titleBoxY+i, '║', nil, borderStyle)
		
		if i == 2 {
			// Title text
			titleX := titleBoxX + (titleBoxWidth-len(title))/2
			for j, r := range title {
				w.screen.SetContent(titleX+j, titleBoxY+i, r, nil, titleStyle)
			}
		}
	}
	
	// Bottom border
	for i := 0; i < titleBoxWidth; i++ {
		if i == 0 {
			w.screen.SetContent(titleBoxX+i, titleBoxY+titleBoxHeight-1, '╚', nil, borderStyle)
		} else if i == titleBoxWidth-1 {
			w.screen.SetContent(titleBoxX+i, titleBoxY+titleBoxHeight-1, '╝', nil, borderStyle)
		} else {
			w.screen.SetContent(titleBoxX+i, titleBoxY+titleBoxHeight-1, '═', nil, borderStyle)
		}
	}

	// Subtitle
	subtitleX := (wW - len(subtitle)) / 2
	subtitleY := titleBoxY + titleBoxHeight + 1
	subtitleStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)
	for i, r := range subtitle {
		w.screen.SetContent(subtitleX+i, subtitleY, r, nil, subtitleStyle)
	}

	// Menu options
	menuY := wH/2 + 2
	keyStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	optionStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)

	// Continue option
	if hasSave {
		option1 := "[ENTER]"
		option1Text := " Continue saved game"
		option1X := (wW - len(option1+option1Text)) / 2
		
		for i, r := range option1 {
			w.screen.SetContent(option1X+i, menuY, r, nil, keyStyle)
		}
		for i, r := range option1Text {
			w.screen.SetContent(option1X+len(option1)+i, menuY, r, nil, selectedStyle)
		}
		menuY += 2
	}

	// New game option
	option2 := "[ESC]"
	option2Text := " Start new game"
	option2X := (wW - len(option2+option2Text)) / 2
	
	styleToUse := optionStyle
	if !hasSave {
		styleToUse = selectedStyle
	}
	
	for i, r := range option2 {
		w.screen.SetContent(option2X+i, menuY, r, nil, keyStyle)
	}
	for i, r := range option2Text {
		w.screen.SetContent(option2X+len(option2)+i, menuY, r, nil, styleToUse)
	}
	menuY += 3

	// Instructions
	instructions := []string{
		"Movement: W/A/S/D",
		"Items: H/R/J/K/E",
		"Equipment: I/U",
		"Drop: X",
		"Quit: Q",
	}
	
	instructionsStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
	instructionsTitle := "Controls:"
	instructionsTitleX := (wW - len(instructionsTitle)) / 2
	for i, r := range instructionsTitle {
		w.screen.SetContent(instructionsTitleX+i, menuY, r, nil, instructionsStyle)
	}
	menuY++

	for _, instruction := range instructions {
		instX := (wW - len(instruction)) / 2
		for i, r := range instruction {
			w.screen.SetContent(instX+i, menuY, r, nil, instructionsStyle)
		}
		menuY++
	}

	// Footer
	footer := "Press Q to quit at any time"
	footerX := (wW - len(footer)) / 2
	footerY := wH - 2
	footerStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
	for i, r := range footer {
		w.screen.SetContent(footerX+i, footerY, r, nil, footerStyle)
	}

	w.screen.Show()
}
