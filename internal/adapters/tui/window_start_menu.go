package tui

import (
	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawStartMenu(hasSave bool, currentOption int) {
	w.screen.Clear()
	wW, wH := w.screen.Size()

	titleLines := []string{
		"          _____                   _______                   _____                    _____                    _____          ",
		"         /\\    \\                 /::\\    \\                 /\\    \\                  /\\    \\                  /\\    \\         ",
		"        /::\\    \\               /::::\\    \\               /::\\    \\                /::\\____\\                /::\\    \\        ",
		"       /::::\\    \\             /::::::\\    \\             /::::\\    \\              /:::/    /               /::::\\    \\       ",
		"      /::::::\\    \\           /::::::::\\    \\           /::::::\\    \\            /:::/    /               /::::::\\    \\      ",
		"     /:::/\\:::\\    \\         /:::/~~\\:::\\    \\         /:::/\\:::\\    \\          /:::/    /               /:::/\\:::\\    \\     ",
		"    /:::/__\\:::\\    \\       /:::/    \\:::\\    \\       /:::/  \\:::\\    \\        /:::/    /               /:::/__\\:::\\    \\    ",
		"   /::::\\   \\:::\\    \\     /:::/    / \\:::\\    \\     /:::/    \\:::\\    \\      /:::/    /               /::::\\   \\:::\\    \\   ",
		"  /::::::\\   \\:::\\    \\   /:::/____/   \\:::\\____\\   /:::/    / \\:::\\    \\    /:::/    /      _____    /::::::\\   \\:::\\    \\  ",
		" /:::/\\:::\\   \\:::\\____\\ |:::|    |     |:::|    | /:::/    /   \\:::\\ ___\\  /:::/____/      /\\    \\  /:::/\\:::\\   \\:::\\    \\ ",
		"/:::/  \\:::\\   \\:::|    ||:::|____|     |:::|    |/:::/____/  ___\\:::|    ||:::|    /      /::\\____\\/:::/__\\:::\\   \\:::\\____\\",
		"\\::/   |::::\\  /:::|____| \\:::\\    \\   /:::/    / \\:::\\    \\ /\\  /:::|____||:::|____\\     /:::/    /\\:::\\   \\:::\\   \\::/    /",
		" \\/____|:::::\\/:::/    /   \\:::\\    \\ /:::/    /   \\:::\\    /::\\ \\::/    /  \\:::\\    \\   /:::/    /  \\:::\\   \\:::\\   \\/____/ ",
		"       |:::::::::/    /     \\:::\\    /:::/    /     \\:::\\   \\:::\\ \\/____/    \\:::\\    \\ /:::/    /    \\:::\\   \\:::\\    \\     ",
		"       |::|\\::::/    /       \\:::\\__/:::/    /       \\:::\\   \\:::\\____\\       \\:::\\    /:::/    /      \\:::\\   \\:::\\____\\    ",
		"       |::| \\::/____/         \\::::::::/    /         \\:::\\  /:::/    /        \\:::\\__/:::/    /        \\:::\\   \\::/    /    ",
		"       |::|  ~|                \\::::::/    /           \\:::\\/:::/    /          \\::::::::/    /          \\:::\\   \\/____/     ",
		"       |::|   |                 \\::::/    /             \\::::::/    /            \\::::::/    /            \\:::\\    \\         ",
		"       \\::|   |                  \\::/____/               \\::::/    /              \\::::/    /              \\:::\\____\\        ",
		"        \\:|   |                   ~~                      \\::/____/                \\::/____/                \\::/    /        ",
		"         \\|___|                                                                     ~~                       \\/____/         ",
	}

	titleWidth := len(titleLines[0])
	titleHeight := len(titleLines)
	titleX := (wW - titleWidth) / 2
	titleY := (wH-titleHeight)/2 - 8

	titleStyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	for i, line := range titleLines {
		if titleY+i >= 0 && titleY+i < wH {
			for j, r := range line {
				if titleX+j >= 0 && titleX+j < wW {
					w.screen.SetContent(titleX+j, titleY+i, r, nil, titleStyle)
				}
			}
		}
	}

	menuOptions := []string{
		"NEW   GAME",
		"LOAD  GAME",
		"SCOREBOARD",
		"EXIT  GAME",
	}

	if !hasSave {
		menuOptions = []string{
			"NEW   GAME",
			"SCOREBOARD",
			"EXIT  GAME",
		}
	}

	menuWidth := 32
	menuHeight := len(menuOptions) + 4
	menuX := (wW - menuWidth) / 2
	menuY := titleY + titleHeight + 3

	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	menuStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)

	titleLine := "           GAME  MENU           "
	menuTitleX := menuX + (menuWidth-len(titleLine))/2
	for i, r := range titleLine {
		if menuTitleX+i >= 0 && menuTitleX+i < wW {
			w.screen.SetContent(menuTitleX+i, menuY, r, nil, borderStyle)
		}
	}

	topBorder := "+------------------------------+"
	for i, r := range topBorder {
		if menuX+i >= 0 && menuX+i < wW {
			w.screen.SetContent(menuX+i, menuY+1, r, nil, borderStyle)
		}
	}

	for i := 0; i < len(menuOptions); i++ {
		emptyLine := "|                              |"
		for j, r := range emptyLine {
			if menuX+j >= 0 && menuX+j < wW {
				w.screen.SetContent(menuX+j, menuY+2+i, r, nil, borderStyle)
			}
		}

		optionText := "|          " + menuOptions[i] + "          |"
		optionStyle := menuStyle
		if i == currentOption {
			optionStyle = selectedStyle
		}
		for j, r := range optionText {
			if menuX+j >= 0 && menuX+j < wW {
				w.screen.SetContent(menuX+j, menuY+2+i, r, nil, optionStyle)
			}
		}

		if i == currentOption {
			arrowLeftX := menuX + 5
			arrowRightX := menuX + menuWidth - 5
			arrowY := menuY + 2 + i
			if arrowY >= 0 && arrowY < wH {
				w.screen.SetContent(arrowLeftX, arrowY, '<', nil, selectedStyle)
				w.screen.SetContent(arrowLeftX+1, arrowY, '<', nil, selectedStyle)
				w.screen.SetContent(arrowLeftX+2, arrowY, '<', nil, selectedStyle)
				w.screen.SetContent(arrowRightX-2, arrowY, '>', nil, selectedStyle)
				w.screen.SetContent(arrowRightX-1, arrowY, '>', nil, selectedStyle)
				w.screen.SetContent(arrowRightX, arrowY, '>', nil, selectedStyle)
			}
		}
	}

	bottomBorder := "+------------------------------+"
	for i, r := range bottomBorder {
		if menuX+i >= 0 && menuX+i < wW {
			w.screen.SetContent(menuX+i, menuY+len(menuOptions)+2, r, nil, borderStyle)
		}
	}

	instructionsStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
	instructionsY := menuY + menuHeight + 2

	instructions := []string{
		"Use W/S to navigate, ENTER to select",
		"Press Q to quit",
	}

	for i, instruction := range instructions {
		if instructionsY+i >= 0 && instructionsY+i < wH {
			instX := (wW - len(instruction)) / 2
			for j, r := range instruction {
				if instX+j >= 0 && instX+j < wW {
					w.screen.SetContent(instX+j, instructionsY+i, r, nil, instructionsStyle)
				}
			}
		}
	}

	w.screen.Show()
}
