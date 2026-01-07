package tui

import (
	"fmt"
	"rogue/internal/domain"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) drawBox(rect domain.Rect, title string) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	for i := 0; i < rect.W; i++ {
		w.screen.SetContent(rect.X+i, rect.Y, '─', nil, style)
		w.screen.SetContent(rect.X+i, rect.Y+rect.H-1, '─', nil, style)
	}

	for i := 0; i < rect.H; i++ {
		w.screen.SetContent(rect.X, rect.Y+i, '│', nil, style)
		w.screen.SetContent(rect.X+rect.W-1, rect.Y+i, '│', nil, style)
	}

	w.screen.SetContent(rect.X, rect.Y, '┌', nil, style)
	w.screen.SetContent(rect.X+rect.W-1, rect.Y, '┐', nil, style)
	w.screen.SetContent(rect.X, rect.Y+rect.H-1, '└', nil, style)
	w.screen.SetContent(rect.X+rect.W-1, rect.Y+rect.H-1, '┘', nil, style)

	if title != "" {
		titleText := " " + title + " "
		titleX := rect.X + (rect.W-len(titleText))/2
		for i, r := range titleText {
			if titleX+i < rect.X+rect.W-1 {
				w.screen.SetContent(titleX+i, rect.Y, r, nil, style)
			}
		}
	}
}

func (w *Window) drawText(rect domain.Rect, line int, format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	w.drawTextRaw(rect, line, text, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}

func (w *Window) drawTextHighlighted(rect domain.Rect, line int, text string) {
	w.drawTextRaw(rect, line, text, tcell.StyleDefault.Reverse(true))
}

func (w *Window) drawTextRaw(rect domain.Rect, line int, text string, style tcell.Style) {
	width := rect.W - 2
	if len(text) > width {
		text = text[:width]
	}

	x := rect.X + 1
	y := rect.Y + line
	for i, r := range text {
		if i >= width {
			break
		}
		w.screen.SetContent(x+i, y, r, nil, style)
	}
}

func (w *Window) DrawMap(world *domain.World, rect domain.Rect) {
	level := world.Level
	player := world.Player

	offsetX := player.Point.X - rect.W/2
	offsetY := player.Point.Y - rect.H/2

	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}
	if offsetX+rect.W > level.W {
		offsetX = level.W - rect.W
	}
	if offsetY+rect.H > level.H {
		offsetY = level.H - rect.H
	}

	w.DrawWorld(world, rect, offsetX, offsetY)
}

func (w *Window) DrawWorld(world *domain.World, rect domain.Rect, offsetX, offsetY int) {
	w.DrawLevel(world, rect, offsetX, offsetY)
	w.DrawItems(world, rect, offsetX, offsetY)
	w.DrawEnemies(world, rect, offsetX, offsetY)
	w.DrawPlayer(world.Player, rect, offsetX, offsetY)
}

func (w *Window) DrawLevel(world *domain.World, rect domain.Rect, offsetX, offsetY int) {
	level := world.Level

	for y := 0; y < rect.H && y+rect.Y < w.layout.Screen.H; y++ {
		levelY := offsetY + y
		if levelY < 0 || levelY >= level.H {
			continue
		}

		for x := 0; x < rect.W && x+rect.X < w.layout.Screen.W; x++ {
			levelX := offsetX + x
			if levelX < 0 || levelX >= level.W {
				continue
			}

			tile := level.Tiles[levelY][levelX]

			if !tile.TestFlags(domain.TileExplored) {
				continue
			}

			var char rune
			var style tcell.Style

			if !tile.TestFlags(domain.TileVisible) {
				char, style = w.runeWall(level.Tiles, levelX, levelY)
				style = style.Foreground(tcell.ColorDarkSlateGray)
			} else {
				switch tile.Kind {
				case domain.TileWall:
					char, style = w.runeWall(level.Tiles, levelX, levelY)
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

			w.screen.SetContent(rect.X+x, rect.Y+y, char, nil, style)
		}
	}
}

func (w *Window) DrawItems(world *domain.World, rect domain.Rect, offsetX, offsetY int) {
	level := world.Level

	for pos, item := range level.Items {
		if !level.Contains(pos) {
			continue
		}

		tile := level.Tiles[pos.Y][pos.X]
		if !tile.TestFlags(domain.TileExplored) {
			continue
		}

		if !tile.TestFlags(domain.TileVisible) {
			continue
		}

		screenItemX := pos.X - offsetX + rect.X
		screenItemY := pos.Y - offsetY + rect.Y

		screenPos := domain.Point{X: screenItemX, Y: screenItemY}
		if !rect.Contains(screenPos) {
			continue
		}

		char, style := w.runeItem(item)
		w.screen.SetContent(screenItemX, screenItemY, char, nil, style)
	}
}

func (w *Window) DrawEnemies(world *domain.World, rect domain.Rect, offsetX, offsetY int) {
	level := world.Level

	for pos, enemy := range level.Enemies {
		if !level.Contains(pos) {
			continue
		}

		if !enemy.IsAlive() {
			continue
		}

		tile := level.Tiles[pos.Y][pos.X]
		if !tile.TestFlags(domain.TileExplored) {
			continue
		}

		if ghost, ok := enemy.(*domain.Ghost); ok {
			distance := domain.Manhattan(pos, world.Player.Point)
			inCombat := distance <= 1
			if !ghost.ShouldBeVisible(inCombat) {
				continue
			}
		}

		if !tile.TestFlags(domain.TileVisible) {
			continue
		}

		screenEnemyX := pos.X - offsetX + rect.X
		screenEnemyY := pos.Y - offsetY + rect.Y

		screenPos := domain.Point{X: screenEnemyX, Y: screenEnemyY}
		if !rect.Contains(screenPos) {
			continue
		}

		if mimic, ok := enemy.(*domain.Mimic); ok {
			if !mimic.IsRevealed() {

				char, style := w.runeMimicDisguise(mimic.DisguisedItem)
				w.screen.SetContent(screenEnemyX, screenEnemyY, char, nil, style)
				continue
			}
		}

		char, color := w.runeEnemy(enemy.GetEnemyType())
		style := tcell.StyleDefault.Foreground(color)
		w.screen.SetContent(screenEnemyX, screenEnemyY, char, nil, style)
	}
}

func (w *Window) DrawPlayer(player *domain.Player, rect domain.Rect, offsetX, offsetY int) {
	screenPlayerX := player.Point.X - offsetX + rect.X
	screenPlayerY := player.Point.Y - offsetY + rect.Y

	screenPos := domain.Point{X: screenPlayerX, Y: screenPlayerY}
	if !rect.Contains(screenPos) {
		return
	}

	char := '@'
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	w.screen.SetContent(screenPlayerX, screenPlayerY, char, nil, style)
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

func (w *Window) runeItem(item domain.Item) (rune, tcell.Style) {
	switch item.Type() {
	case domain.ItemFood:
		return 'f', tcell.StyleDefault.Foreground(tcell.ColorGreen)
	case domain.ItemElixir:
		return 'e', tcell.StyleDefault.Foreground(tcell.ColorBlue)
	case domain.ItemScroll:
		return 's', tcell.StyleDefault.Foreground(tcell.ColorPurple)
	case domain.ItemWeapon:
		return 'w', tcell.StyleDefault.Foreground(tcell.ColorRed)
	case domain.ItemArmor:
		return 'a', tcell.StyleDefault.Foreground(tcell.ColorLightCyan)
	case domain.ItemTreasure:
		return '$', tcell.StyleDefault.Foreground(tcell.ColorYellow)
	default:
		return '?', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	}
}

func (w *Window) runeEnemy(enemyType domain.EnemyType) (rune, tcell.Color) {
	switch enemyType {
	case domain.EnemyTypeZombie:
		return 'z', tcell.ColorGreen
	case domain.EnemyTypeVampire:
		return 'v', tcell.ColorRed
	case domain.EnemyTypeGhost:
		return 'g', tcell.ColorWhite
	case domain.EnemyTypeOgre:
		return 'O', tcell.ColorYellow
	case domain.EnemyTypeSnakeMage:
		return 's', tcell.ColorWhite
	case domain.EnemyTypeMimic:
		return 'm', tcell.ColorWhite
	default:
		return '?', tcell.ColorWhite
	}
}

func (w *Window) runeMimicDisguise(itemKind domain.ItemKind) (rune, tcell.Style) {
	switch itemKind {
	case domain.ItemFood:
		return 'f', tcell.StyleDefault.Foreground(tcell.ColorGreen)
	case domain.ItemElixir:
		return 'e', tcell.StyleDefault.Foreground(tcell.ColorBlue)
	case domain.ItemScroll:
		return 's', tcell.StyleDefault.Foreground(tcell.ColorPurple)
	case domain.ItemWeapon:
		return 'w', tcell.StyleDefault.Foreground(tcell.ColorRed)
	case domain.ItemTreasure:
		return '$', tcell.StyleDefault.Foreground(tcell.ColorYellow)
	default:
		return '?', tcell.StyleDefault.Foreground(tcell.ColorWhite)
	}
}
