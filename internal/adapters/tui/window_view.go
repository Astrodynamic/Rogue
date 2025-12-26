package tui

import (
	"rogue/internal/domain"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawWorld(world *domain.World) {
	w.DrawLevel(world.Level)
	w.DrawPlayer(world.Player)
}

func (w *Window) DrawLevel(level *domain.Level) {
	for y := 0; y < level.Height; y++ {
		for x := 0; x < level.Width; x++ {
			tile := level.Tiles[y][x]
			var char rune
			var style tcell.Style
			switch tile.Kind {
			case domain.TileWall:
				char = '#'
				style = tcell.StyleDefault.Foreground(tcell.ColorGray)
			case domain.TileFloor:
				char = '.'
				style = tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
			case domain.TileCorridor:
				char = '∙'
				style = tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
			}
			w.screen.SetContent(x, y, char, nil, style)
		}
	}
}

func (w *Window) DrawPlayer(player *domain.Player) {
	w.screen.SetContent(player.Point.X, player.Point.Y, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}
