package game

import (
	"errors"
	"math/rand"
)

func NewSession(seed int64, cfg Config) *GameSession {
	cfg = cfg.Normalize()
	rng := rand.New(rand.NewSource(seed))
	gen := NewGenerator(rng)
	lvl, startIdx, _ := gen.Generate(1, cfg)

	startPos := randomFloorInRoom(rng, lvl.Tiles, lvl.Rooms[startIdx].Bounds)

	s := &GameSession{
		Seed:       seed,
		LevelDepth: 1,
		Level:      lvl,
		Player:     NewCharacter(25, 10, 8),
		PlayerPos:  startPos,
		NextItemID: 1,
		Enemies:    nil,
		Explored:   make2D(cfg.Width, cfg.Height, false),
		Effects:    nil,
		Stats:      Stats{DeepestLevelReached: 1},
		Messages:   nil,
	}

	s.populate(rng)
	s.refreshFog()
	return s
}

func (s *GameSession) Step(rng *rand.Rand, a Action) (advanced bool, gameOver bool, won bool, err error) {
	if rng == nil {
		return false, false, false, errors.New("nil rng")
	}

	// Sleep blocks player action for one turn.
	if s.hasSleep() && a.Kind != ActionNone {
		s.addMsg("You are asleep.")
		s.tickEffects()
		s.enemyTurn(rng)
		s.tickEffects()
		s.refreshFog()
		return true, s.Player.IsDead(), false, nil
	}

	switch a.Kind {
	case ActionNone:
		// no-op
		return false, s.Player.IsDead(), false, nil
	case ActionMove:
		s.doMove(rng, a.Dx, a.Dy)
		advanced = true
	case ActionUseFood:
		s.useFood(a.Idx)
		advanced = true
	case ActionUseElixir:
		s.useElixir(rng, a.Idx)
		advanced = true
	case ActionUseScroll:
		s.useScroll(a.Idx)
		advanced = true
	case ActionEquipWeapon:
		s.equipWeapon(rng, a.Idx)
		advanced = true
	case ActionUnequipWeapon:
		s.Player.Equipment.Clear(SlotWeapon)
		s.addMsg("Weapon unequipped.")
		advanced = true
	default:
		return false, false, false, errors.New("unknown action")
	}

	if advanced {
		s.enemyTurn(rng)
		s.tickEffects()
		s.refreshFog()
	}

	if s.Player.IsDead() {
		return advanced, true, false, nil
	}
	if s.LevelDepth > MaxLevels {
		return advanced, false, true, nil
	}
	return advanced, false, false, nil
}

func (s *GameSession) refreshFog() {
	vis := s.ComputeVisible()
	w := len(s.Level.Tiles[0])
	h := len(s.Level.Tiles)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if vis[y][x] {
				s.Explored[y][x] = true
			}
		}
	}
}

func (s *GameSession) addMsg(m string) {
	const cap = 6
	if m == "" {
		return
	}
	s.Messages = append(s.Messages, m)
	if len(s.Messages) > cap {
		s.Messages = s.Messages[len(s.Messages)-cap:]
	}
}
