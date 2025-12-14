package game

import "math/rand"

func (s *GameSession) enemyTurn(rng *rand.Rand) {
	for i := 0; i < len(s.Enemies); i++ {
		e := &s.Enemies[i]

		// Ghost periodic invisibility unless actively chasing in hostility range.
		if e.Type == EnemyGhost {
			e.GhostInvisToggle++
			if Manhattan(e.Pos, s.PlayerPos) <= e.Hostility {
				e.GhostInvisible = false
			} else if e.GhostInvisToggle%7 == 0 {
				e.GhostInvisible = !e.GhostInvisible
			}
		}

		// Ogre rest after attack.
		if e.Type == EnemyOgre && e.OgreRestTurns > 0 {
			e.OgreRestTurns--
			continue
		}

		// Chase if in range and path exists.
		inRange := Manhattan(e.Pos, s.PlayerPos) <= e.Hostility
		occ := s.occupiedMap(i)
		if inRange {
			if next, ok := ShortestPathNext(s.Level.Tiles, occ, e.Pos, s.PlayerPos); ok {
				if next == s.PlayerPos {
					s.enemyAttack(rng, e)
				} else {
					e.Pos = next
				}
				// Ogre moves two tiles per turn within the room (also while chasing).
				if e.Type == EnemyOgre {
					if roomID, inRoom := s.roomAt(e.Pos); inRoom {
						b := s.Level.Rooms[roomID].Bounds
						occ2 := s.occupiedMap(i)
						if next2, ok2 := ShortestPathNext(s.Level.Tiles, occ2, e.Pos, s.PlayerPos); ok2 {
							if next2 == s.PlayerPos {
								s.enemyAttack(rng, e)
							} else if b.Contains(next2) {
								e.Pos = next2
							}
						}
					}
				}
				continue
			}
		}

		// No path or not in range -> type-specific movement.
		s.enemyWander(rng, e, occ)
	}
}

func (s *GameSession) enemyWander(rng *rand.Rand, e *Enemy, occ map[Point]bool) {
	roomID, inRoom := s.roomAt(e.Pos)
	var bounds Rect
	if inRoom {
		bounds = s.Level.Rooms[roomID].Bounds
	}
	tryMove := func(np Point) bool {
		w := len(s.Level.Tiles[0])
		h := len(s.Level.Tiles)
		if !inBounds(np, w, h) {
			return false
		}
		if s.Level.Tiles[np.Y][np.X].BlocksMovement() {
			return false
		}
		if occ[np] {
			return false
		}
		if inRoom && !bounds.Contains(np) {
			return false
		}
		e.Pos = np
		return true
	}

	switch e.Type {
	case EnemyGhost:
		if inRoom {
			// Teleport within room.
			e.Pos = randomFloorInRoom(rng, s.Level.Tiles, bounds)
			return
		}
	case EnemySnakeMage:
		if e.SnakeDx == 0 && e.SnakeDy == 0 {
			e.SnakeDx = 1
			e.SnakeDy = 1
		}
		np := Point{X: e.Pos.X + e.SnakeDx, Y: e.Pos.Y + e.SnakeDy}
		if !tryMove(np) {
			// Bounce by flipping one component.
			if rng.Intn(2) == 0 {
				e.SnakeDx = -e.SnakeDx
			} else {
				e.SnakeDy = -e.SnakeDy
			}
			np2 := Point{X: e.Pos.X + e.SnakeDx, Y: e.Pos.Y + e.SnakeDy}
			_ = tryMove(np2)
		}
		return
	case EnemyOgre:
		// Move two tiles per turn within room when wandering.
		for step := 0; step < 2; step++ {
			d := Dirs4[rng.Intn(len(Dirs4))]
			np := e.Pos.Add(d)
			if !tryMove(np) {
				return
			}
		}
		return
	}

	// Default random 4-dir step.
	for attempts := 0; attempts < 4; attempts++ {
		d := Dirs4[rng.Intn(len(Dirs4))]
		np := e.Pos.Add(d)
		if tryMove(np) {
			return
		}
	}
}
