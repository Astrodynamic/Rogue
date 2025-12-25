package tui

import (
	"rogue/internal/domain"
	"rogue/internal/service"

	"github.com/gdamore/tcell/v2"
)

type Window struct {
	screen tcell.Screen
}

func NewWindow() *Window {
	return &Window{}
}

func (w *Window) Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}
	w.screen = s
	return nil
}

func (w *Window) Close() {
	if w.screen != nil {
		w.screen.Fini()
	}
}

func (w *Window) Draw(world *domain.World) {
	w.screen.Clear()

	lvl := world.Level
	for y := 0; y < lvl.Height; y++ {
		for x := 0; x < lvl.Width; x++ {
			tile := lvl.Tiles[y][x]
			var char rune
			var style tcell.Style
			switch tile.Kind {
			case domain.TileWall:
				char = '#'
				style = tcell.StyleDefault.Foreground(tcell.ColorGray)
			case domain.TileFloor:
				char = '.'
				style = tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
			}
			w.screen.SetContent(x, y, char, nil, style)
		}
	}

	player := world.Player
	w.screen.SetContent(player.Point.X, player.Point.Y, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))

	w.screen.Show()
}

func (w *Window) Input() service.Command {
	ev := w.screen.PollEvent()

	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			return service.CmdQuit
		case tcell.KeyUp:
			return service.CmdMoveUp
		case tcell.KeyDown:
			return service.CmdMoveDown
		case tcell.KeyLeft:
			return service.CmdMoveLeft
		case tcell.KeyRight:
			return service.CmdMoveRight
		}

		switch ev.Rune() {
		case 'q', 'Q':
			return service.CmdQuit
		case 'w', 'W':
			return service.CmdMoveUp
		case 's', 'S':
			return service.CmdMoveDown
		case 'a', 'A':
			return service.CmdMoveLeft
		case 'd', 'D':
			return service.CmdMoveRight
		}
	}
	return service.CmdNone
}
