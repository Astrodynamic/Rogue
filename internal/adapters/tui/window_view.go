package tui

import (
	"rogue/internal/domain"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawWorld(world *domain.World) {
	w.DrawLevel(world)
	w.DrawPlayer(world.Player)
}

func (w *Window) DrawLevel(world *domain.World) {
	level := world.Level
	for y := level.Y; y < level.Y+level.H; y++ {
		for x := level.X; x < level.X+level.W; x++ {
			tile := level.Tiles[y][x]

			if !tile.TestFlags(domain.TileExplored) {
				continue
			}

			var char rune
			var style tcell.Style

			if !tile.TestFlags(domain.TileVisible) {
				char, style = w.runeWall(level.Tiles, x, y)
				style = style.Foreground(tcell.ColorDarkSlateGray)
			} else {
				switch tile.Kind {
				case domain.TileWall:
					char, style = w.runeWall(level.Tiles, x, y)
				case domain.TileFloor:
					char, style = w.runeFloor()
				case domain.TileCorridor:
					char, style = w.runeCorridor()
				case domain.TileExit:
					char, style = w.runeExit()
				default:
					char, style = w.runeNone()
				}
			}

			w.screen.SetContent(x, y, char, nil, style)
		}
	}
}

func (w *Window) runeWall(tiles [][]domain.Tile, x, y int) (rune, tcell.Style) {
	wall := func(x1, y1 int) bool {
		if y1 < 0 || y1 >= len(tiles) || x1 < 0 || x1 >= len(tiles[y1]) {
			return true
		}
		return tiles[y1][x1].Kind == domain.TileWall
	}
	u := wall(x, y-1)
	d := wall(x, y+1)
	l := wall(x-1, y)
	r := wall(x+1, y)

	var char rune
	style := tcell.StyleDefault.Foreground(tcell.ColorGray)

	switch {
	case u && d && l && r:
		char = '┼'
	case u && d && l:
		char = '┤'
	case u && d && r:
		char = '├'
	case l && r && u:
		char = '┴'
	case l && r && d:
		char = '┬'
	case u && d:
		char = '│'
	case l && r:
		char = '─'
	case d && r:
		char = '┌'
	case d && l:
		char = '┐'
	case u && r:
		char = '└'
	case u && l:
		char = '┘'
	case u || d:
		char = '│'
	case l || r:
		char = '─'
	default:
		char = '█'
	}

	return char, style
}

func (w *Window) runeFloor() (rune, tcell.Style) {
	return '·', tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
}

func (w *Window) runeCorridor() (rune, tcell.Style) {
	return '∙', tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
}

func (w *Window) runeExit() (rune, tcell.Style) {
	return '◉', tcell.StyleDefault.Foreground(tcell.ColorYellow)
}

func (w *Window) runeNone() (rune, tcell.Style) {
	return ' ', tcell.StyleDefault
}

func (w *Window) DrawPlayer(player *domain.Player) {
	w.screen.SetContent(player.Point.X, player.Point.Y, '@', nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true))
}
