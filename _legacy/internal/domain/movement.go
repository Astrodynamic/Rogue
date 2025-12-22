package domain

import "math/rand"

func (s *GameSession) occupiedMap(exceptEnemyIdx int) map[Point]bool {
	occ := make(map[Point]bool, len(s.Enemies)+1)
	occ[s.PlayerPos] = true
	for i := range s.Enemies {
		if i == exceptEnemyIdx {
			continue
		}
		occ[s.Enemies[i].Pos] = true
	}
	return occ
}

func (s *GameSession) enemyIndexAt(p Point) (int, bool) {
	for i := range s.Enemies {
		if s.Enemies[i].Pos == p {
			return i, true
		}
	}
	return -1, false
}

func (s *GameSession) doMove(rng *rand.Rand, dx, dy int) {
	if dx == 0 && dy == 0 {
		return
	}
	to := Point{X: s.PlayerPos.X + dx, Y: s.PlayerPos.Y + dy}
	w := len(s.Level.Tiles[0])
	h := len(s.Level.Tiles)
	if !inBounds(to, w, h) {
		s.addMsg("You bump into the edge.")
		return
	}
	if s.Level.Tiles[to.Y][to.X].BlocksMovement() {
		s.addMsg("You bump into a wall.")
		return
	}
	if ei, ok := s.enemyIndexAt(to); ok {
		s.playerAttack(rng, ei)
		return
	}

	s.PlayerPos = to
	s.Stats.TilesTraveled++

	// Pick up item if any.
	for ri := range s.Level.Rooms {
		if !s.Level.Rooms[ri].Bounds.Contains(to) {
			continue
		}
		if s.Level.Rooms[ri].IsStart {
			// Start room may still have dropped weapon; allowed to pick up.
		}
		if it, ok := s.Level.Rooms[ri].FloorLoot[to]; ok {
			// Backward compatibility: assign an ID if missing.
			if it.ID == 0 {
				it.ID = s.nextItemID()
			}
			if s.Backpack.Add(it) {
				delete(s.Level.Rooms[ri].FloorLoot, to)
				s.addMsg("Picked up: " + it.DisplayName())
			} else {
				s.addMsg("Backpack is full for that item type.")
			}
		}
		break
	}

	if s.Level.Tiles[to.Y][to.X] == TileExit {
		s.advanceLevel(rng)
	}
}

func (s *GameSession) doAttack(rng *rand.Rand) {
	// Attack enemy in adjacent cell (4 directions). If multiple, attack first found.
	for _, d := range Dirs4 {
		adj := s.PlayerPos.Add(d)
		if ei, ok := s.enemyIndexAt(adj); ok {
			s.playerAttack(rng, ei)
			return
		}
	}
	s.addMsg("No enemy nearby to attack.")
}
