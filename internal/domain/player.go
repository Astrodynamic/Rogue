package domain

type Player struct {
	Actor
	Equipment *Equipment
}

func NewPlayer(name string) *Player {
	if name == "" {
		name = "Player"
	}
	return &Player{
		Actor: Actor{
			Point: Point{X: 0, Y: 0},
			Stats: Stats{
				MaxHealth: PlayerStartMaxHealth,
				Health:    PlayerStartHealth,
				Dexterity: PlayerStartDexterity,
				Strength:  PlayerStartStrength,
			},
			Name:     name,
			State:    ActorStateNormal,
			Backpack: NewBackpack(),
		},
		Equipment: NewEquipment(),
	}
}

func (p *Player) EquipItem(part ActorPart, item Item) Item {
	if item == nil {
		return nil
	}
	itemPart, canEquip := GetEquipPart(item)
	if !canEquip || itemPart != part {
		return nil
	}
	return p.Equipment.Equip(part, item)
}

func (p *Player) UnequipItem(part ActorPart) Item {
	return p.Equipment.Unequip(part)
}

func (p *Player) UseItem(item Item) ItemUseResult {
	part, canEquip := GetEquipPart(item)
	if canEquip {
		oldItem := p.EquipItem(part, item)
		if oldItem != nil {
			if !p.Backpack.Add(oldItem) {
				return ItemUseResult{
					Success: false,
					Message: "Cannot add old item to backpack",
				}
			}
		}
		return ItemUseResult{
			Success:  true,
			Consumed: true,
			Message:  item.Name() + " equipped",
		}
	}
	return p.Actor.UseItem(item)
}

func (p *Player) GetEffectiveStats() Stats {
	equipment := p.Equipment.GetTotalStats()
	effective := p.Actor.GetEffectiveStats()
	effective.AddMaxHealth(equipment.MaxHealth)
	effective.AddHealth(equipment.Health)
	effective.AddDexterity(equipment.Dexterity)
	effective.AddStrength(equipment.Strength)
	return effective
}

func (p *Player) DropItem(itemKind ItemKind, stackIndex int) Item {
	return p.Backpack.Remove(itemKind, stackIndex)
}

func (p *Player) DropEquipment(part ActorPart) Item {
	return p.Equipment.Unequip(part)
}

func (p *Player) HandleOldItemOnEquip(oldItem Item, level *Level, playerPos Point) bool {
	if oldItem == nil {
		return true
	}
	if !p.Backpack.Add(oldItem) {
		dropPos := level.GetAdjacentDropPoint(playerPos)
		if dropPos.X >= 0 {
			level.AddItem(dropPos, oldItem)
		} else {
			return false
		}
	}
	return true
}
