package tui

import (
	"rogue/internal/domain"
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
	w.drawText(rect, line, "Depth: %d", world.Depth)
}
