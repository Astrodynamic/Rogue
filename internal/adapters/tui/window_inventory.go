package tui

import (
	"fmt"
	"rogue/internal/domain"
)

func (w *Window) formatItemProperties(item domain.Item) string {
	stats := item.GetStats()
	itemType := item.Type()

	switch itemType {
	case domain.ItemFood:
		if stats.Health > 0 {
			return fmt.Sprintf("+%dHP", stats.Health)
		}
	case domain.ItemElixir:
		if elixir, ok := item.(*domain.Elixir); ok {
			var statName string
			switch elixir.ElixirKind {
			case domain.ElixirHealth:
				statName = "HP"
			case domain.ElixirMaxHealth:
				statName = "MaxHP"
			case domain.ElixirDexterity:
				statName = "Dex"
			case domain.ElixirStrength:
				statName = "Str"
			default:
				statName = "?"
			}
			return fmt.Sprintf("+%d %s (%dt)", elixir.Amount, statName, elixir.Duration)
		}
	case domain.ItemScroll:
		if scroll, ok := item.(*domain.Scroll); ok {
			var statName string
			switch scroll.ScrollKind {
			case domain.ScrollHealth:
				statName = "HP"
			case domain.ScrollMaxHealth:
				statName = "MaxHP"
			case domain.ScrollDexterity:
				statName = "Dex"
			case domain.ScrollStrength:
				statName = "Str"
			default:
				statName = "?"
			}
			return fmt.Sprintf("+%d %s", scroll.Amount, statName)
		}
	case domain.ItemWeapon:
		if stats.Strength > 0 {
			return fmt.Sprintf("Str+%d", stats.Strength)
		}
	case domain.ItemArmor:
		var props []string
		if stats.Strength > 0 {
			props = append(props, fmt.Sprintf("Str+%d", stats.Strength))
		}
		if stats.Dexterity > 0 {
			props = append(props, fmt.Sprintf("Dex+%d", stats.Dexterity))
		}
		if stats.MaxHealth > 0 {
			props = append(props, fmt.Sprintf("MaxHP+%d", stats.MaxHealth))
		}
		if len(props) == 0 {
			return "No bonuses"
		}
		result := props[0]
		for i := 1; i < len(props); i++ {
			result += " " + props[i]
		}
		return result
	case domain.ItemTreasure:
		if treasure, ok := item.(*domain.Treasure); ok {
			return fmt.Sprintf("Value: %d", treasure.Value)
		}
	}
	return ""
}

func (w *Window) DrawInventory(player *domain.Player, rect domain.Rect) {
	w.drawBox(rect, "Inventory")

	line := 1
	itemKinds := []domain.ItemKind{
		domain.ItemFood,
		domain.ItemElixir,
		domain.ItemScroll,
		domain.ItemWeapon,
		domain.ItemArmor,
		domain.ItemTreasure,
	}

	for _, kind := range itemKinds {
		if kind == domain.ItemTreasure {
			totalTreasure := player.Backpack.GetTotalTreasure()
			if totalTreasure > 0 {
				if line >= rect.H-2 {
					return
				}
				w.drawText(rect, line, "Treasure: %d", totalTreasure)
				line++
			}
			continue
		}

		stacks := player.Backpack.GetStacks(kind)
		if len(stacks) == 0 {
			continue
		}

		for _, stack := range stacks {
			if line >= rect.H-2 {
				return
			}

			item := stack.Peek()
			if item == nil {
				continue
			}

			itemName := item.Name()
			properties := w.formatItemProperties(item)
			count := stack.Count

			if count > 1 {
				w.drawText(rect, line, "%s %s (x%d)", itemName, properties, count)
			} else {
				w.drawText(rect, line, "%s %s", itemName, properties)
			}
			line++
		}
	}
}
