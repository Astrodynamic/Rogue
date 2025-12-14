package game

type LootEntry struct {
	Pos  Point
	Item Item
}

type RoomSnapshot struct {
	ID      int
	Bounds  Rect
	IsStart bool
	IsExit  bool
	Loot    []LootEntry
}

type LevelSnapshot struct {
	Depth     int
	Rooms     []RoomSnapshot
	Corridors []Corridor
	Tiles     [][]Tile
}

type Snapshot struct {
	Seed       int64
	LevelDepth int
	Level      LevelSnapshot

	Player    Character
	Backpack  Backpack
	PlayerPos Point
	Enemies   []Enemy

	Explored [][]bool
	Effects  []Effect
	Stats    Stats
	Messages []string
}

func (s *GameSession) Snapshot() Snapshot {
	rs := make([]RoomSnapshot, 0, len(s.Level.Rooms))
	for _, r := range s.Level.Rooms {
		loot := make([]LootEntry, 0, len(r.FloorLoot))
		for p, it := range r.FloorLoot {
			loot = append(loot, LootEntry{Pos: p, Item: it})
		}
		rs = append(rs, RoomSnapshot{
			ID:      r.ID,
			Bounds:  r.Bounds,
			IsStart: r.IsStart,
			IsExit:  r.IsExit,
			Loot:    loot,
		})
	}
	return Snapshot{
		Seed:       s.Seed,
		LevelDepth: s.LevelDepth,
		Level: LevelSnapshot{
			Depth:     s.Level.Depth,
			Rooms:     rs,
			Corridors: s.Level.Corridors,
			Tiles:     s.Level.Tiles,
		},
		Player:    s.Player,
		Backpack:  s.Backpack,
		PlayerPos: s.PlayerPos,
		Enemies:   append([]Enemy(nil), s.Enemies...),
		Explored:  s.Explored,
		Effects:   append([]Effect(nil), s.Effects...),
		Stats:     s.Stats,
		Messages:  append([]string(nil), s.Messages...),
	}
}

func FromSnapshot(sn Snapshot) *GameSession {
	rooms := make([]Room, 0, len(sn.Level.Rooms))
	for _, rs := range sn.Level.Rooms {
		loot := make(map[Point]Item, len(rs.Loot))
		for _, le := range rs.Loot {
			loot[le.Pos] = le.Item
		}
		rooms = append(rooms, Room{
			ID:        rs.ID,
			Bounds:    rs.Bounds,
			IsStart:   rs.IsStart,
			IsExit:    rs.IsExit,
			FloorLoot: loot,
		})
	}
	s := &GameSession{
		Seed:       sn.Seed,
		LevelDepth: sn.LevelDepth,
		Level: Level{
			Depth:     sn.Level.Depth,
			Rooms:     rooms,
			Corridors: sn.Level.Corridors,
			Tiles:     sn.Level.Tiles,
		},
		Player:     sn.Player,
		Backpack:   sn.Backpack,
		PlayerPos:  sn.PlayerPos,
		NextItemID: 1,
		Enemies:    sn.Enemies,
		Explored:   sn.Explored,
		Effects:    sn.Effects,
		Stats:      sn.Stats,
		Messages:   sn.Messages,
	}
	s.normalizeItemIDs()
	return s
}

func (s *GameSession) normalizeItemIDs() {
	maxID := int64(0)
	// Backpack items.
	for i := range s.Backpack.Items {
		if s.Backpack.Items[i].ID == 0 {
			maxID++
			s.Backpack.Items[i].ID = maxID
			continue
		}
		if s.Backpack.Items[i].ID > maxID {
			maxID = s.Backpack.Items[i].ID
		}
	}
	// Floor loot.
	for ri := range s.Level.Rooms {
		for p, it := range s.Level.Rooms[ri].FloorLoot {
			if it.ID == 0 {
				maxID++
				it.ID = maxID
				s.Level.Rooms[ri].FloorLoot[p] = it
				continue
			}
			if it.ID > maxID {
				maxID = it.ID
			}
		}
	}
	s.NextItemID = maxID + 1
	// If equipped weapon doesn't exist anymore, clear it.
	if wid := s.Player.Equipment.Get(SlotWeapon); wid != 0 {
		if it, ok := s.Backpack.FindByID(wid); !ok || it.Type != ItemWeapon {
			s.Player.Equipment.Clear(SlotWeapon)
		}
	}
}
