package domain

type SnakeMage struct {
	BaseEnemy
	Direction      Point
	DirectionClock int
}

func NewSnakeMage(depth int) *SnakeMage {
	config := EnemyConfig{
		BaseHealth:    SnakeMageBaseHealth,
		BaseDexterity: SnakeMageBaseDexterity,
		BaseStrength:  SnakeMageBaseStrength,
		BaseHostility: SnakeMageBaseHostility,
	}
	stats := ScaleEnemyStats(config, depth)
	hostility := ScaleHostility(config.BaseHostility, depth)

	return &SnakeMage{
		BaseEnemy:      NewBaseEnemy(EnemyTypeSnakeMage, stats, hostility),
		Direction:      DirUR,
		DirectionClock: 0,
	}
}

func (s *SnakeMage) GetEnemyType() EnemyType {
	return s.EnemyType
}

func (s *SnakeMage) Name() string {
	return "Snake Mage"
}

func (s *SnakeMage) SwitchDirection() {
	diagonalDirs := []Point{DirUL, DirUR, DirDR, DirDL}
	s.DirectionClock = (s.DirectionClock + 1) % len(diagonalDirs)
	s.Direction = diagonalDirs[s.DirectionClock]
}

func (s *SnakeMage) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
	if !s.IsAlive() {
		return
	}

	enemyPos := s.Actor.Point
	distance := Manhattan(enemyPos, playerPos)

	if s.CanAttack(playerPos) {
		enemyActor := s.GetActor()
		playerActor := aiCtx.GetPlayerActor()
		if enemyActor.State != ActorStateSleep {
			resolver := aiCtx.GetCombatResolver()
			result := enemyActor.AttackWithResolver(playerActor, resolver)
			if result.Hit {
				rng := aiCtx.GetRandomGenerator()
				if rng.IntN(100) < SnakeMageSleepChance {
					playerActor.State = ActorStateSleep
				}
			}
			aiCtx.OnEnemyAttack(s, result)
		}
		return
	}

	diagonalDirs := []Point{DirUL, DirUR, DirDR, DirDL}
	if distance <= s.Hostility {
		newPos := aiCtx.FindDiagonalMove(enemyPos, playerPos, level, s.Direction, diagonalDirs)
		if newPos.X >= 0 && newPos.Y >= 0 {
			foundDir := newPos.Sub(enemyPos)
			for i, dir := range diagonalDirs {
				if dir == foundDir {
					s.DirectionClock = i
					s.Direction = dir
					break
				}
			}
			aiCtx.MoveEnemy(enemyPos, newPos, level)
		} else {
			newPos := aiCtx.FindPathTo(enemyPos, playerPos, level)
			if newPos.X >= 0 && newPos.Y >= 0 {
				aiCtx.MoveEnemy(enemyPos, newPos, level)
			}
		}
	} else {
		newPos := enemyPos.Add(s.Direction)
		if !level.Contains(newPos) || level.GetEnemy(newPos) != nil {
			s.SwitchDirection()
			newPos = aiCtx.FindDiagonalMove(enemyPos, Point{X: -1, Y: -1}, level, s.Direction, diagonalDirs)
		}
		if newPos.X >= 0 && newPos.Y >= 0 {
			aiCtx.MoveEnemy(enemyPos, newPos, level)
		}
	}
}
