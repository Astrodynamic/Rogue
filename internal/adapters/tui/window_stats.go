package tui

import (
	"fmt"
	"rogue/internal/domain"

	"github.com/gdamore/tcell/v2"
)

func (w *Window) DrawStats(world *domain.World, rect domain.Rect) {
	w.drawBox(rect, "Stats")

	player := world.Player
	line := 1
	w.drawText(rect, line, "Health: %d/%d", player.Health, player.MaxHealth)
	line++
	w.drawText(rect, line, "Dexterity: %d", player.Dexterity)
	line++
	w.drawText(rect, line, "Strength: %d", player.Strength)
	line++
	w.drawText(rect, line, "Depth: %d", world.GetDepth())
}

func (w *Window) drawStatisticsView(playthroughs []*domain.PlayStats, rect domain.Rect) {
	w.drawBox(rect, "Statistics / Leaderboard")

	line := 1
	if len(playthroughs) == 0 {
		w.drawText(rect, line, "No playthroughs yet")
		w.drawText(rect, rect.H-2, "Press ESC to return")
		return
	}

	headerStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	separatorStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	colWidths := []int{3, 12, 10, 6, 8, 5, 8, 8, 5, 6}
	colHeaders := []string{"#", "Player", "Treasure", "Level", "Enemies", "Food", "Elixirs", "Scrolls", "Hits", "Tiles"}

	x := rect.X + 1
	y := rect.Y + line

	for i, header := range colHeaders {
		if i > 0 {
			w.screen.SetContent(x, y, '│', nil, separatorStyle)
			x++
		}
		headerText := fmt.Sprintf("%-*s", colWidths[i], header)
		for _, r := range headerText {
			if x >= rect.X+rect.W-1 {
				break
			}
			w.screen.SetContent(x, y, r, nil, headerStyle)
			x++
		}
	}
	line++

	x = rect.X + 1
	y = rect.Y + line
	for i := 0; i < len(colWidths); i++ {
		if i > 0 {
			w.screen.SetContent(x, y, '┼', nil, separatorStyle)
			x++
		}
		for j := 0; j < colWidths[i]; j++ {
			if x >= rect.X+rect.W-1 {
				break
			}
			w.screen.SetContent(x, y, '─', nil, separatorStyle)
			x++
		}
	}
	line++

	availableHeight := rect.H - 3
	if availableHeight < 1 {
		availableHeight = 1
	}

	maxDisplay := availableHeight
	if len(playthroughs) > maxDisplay {
		playthroughs = playthroughs[:maxDisplay]
	}

	for i, stats := range playthroughs {
		if line >= rect.H-1 {
			break
		}

		x = rect.X + 1
		y = rect.Y + line

		playerName := stats.PlayerName
		if playerName == "" {
			playerName = "Player"
		}
		if len(playerName) > colWidths[1]-1 {
			playerName = playerName[:colWidths[1]-1]
		}

		values := []interface{}{
			i + 1,
			playerName,
			stats.TreasureCollected,
			stats.DeepestLevel,
			stats.EnemiesDefeated,
			stats.FoodConsumed,
			stats.ElixirsDrunk,
			stats.ScrollsRead,
			stats.HitsDealt + stats.HitsReceived,
			stats.TilesTraveled,
		}

		style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
		if i == 0 {
			style = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
		}

		for j, val := range values {
			if j > 0 {
				w.screen.SetContent(x, y, '│', nil, separatorStyle)
				x++
			}
			var text string
			if j == 0 {
				text = fmt.Sprintf("%-*d", colWidths[j], val)
			} else if j == 1 {
				text = fmt.Sprintf("%-*s", colWidths[j], val)
			} else {
				text = fmt.Sprintf("%*d", colWidths[j], val)
			}
			for _, r := range text {
				if x >= rect.X+rect.W-1 {
					break
				}
				w.screen.SetContent(x, y, r, nil, style)
				x++
			}
		}
		line++
	}

	w.drawText(rect, rect.H-2, "Press ESC to return")
}
