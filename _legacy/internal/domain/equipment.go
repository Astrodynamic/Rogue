package domain

type EquipSlot string

const (
	SlotWeapon EquipSlot = "weapon"
)

// Equipment is intentionally generic: a map of slot -> Item.ID.
// This allows future expansion (armor/rings/etc.) without changing Character fields.
type Equipment struct {
	Slots map[EquipSlot]int64
}

func (e *Equipment) Get(slot EquipSlot) int64 {
	if e == nil || e.Slots == nil {
		return 0
	}
	return e.Slots[slot]
}

func (e *Equipment) Set(slot EquipSlot, id int64) {
	if e.Slots == nil {
		e.Slots = map[EquipSlot]int64{}
	}
	e.Slots[slot] = id
}

func (e *Equipment) Clear(slot EquipSlot) {
	if e == nil || e.Slots == nil {
		return
	}
	delete(e.Slots, slot)
}
