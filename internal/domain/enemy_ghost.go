package domain

type Ghost struct {
	BaseEnemy
	InvisibleTurns int
	InCombat       bool
}

func NewGhost(depth int) *Ghost {
	config := EnemyConfig{
		BaseHealth:    GhostBaseHealth,
		BaseDexterity: GhostBaseDexterity,
		BaseStrength:  GhostBaseStrength,
		BaseHostility: GhostBaseHostility,
	}
	stats := ScaleEnemyStats(config, depth)
	hostility := ScaleHostility(config.BaseHostility, depth)

	return &Ghost{
		BaseEnemy:      NewBaseEnemy(EnemyTypeGhost, stats, hostility),
		InvisibleTurns: 0,
		InCombat:       false,
	}
}

func (g *Ghost) GetEnemyType() EnemyType {
	return g.EnemyType
}

func (g *Ghost) Name() string {
	return "Ghost"
}

func (g *Ghost) ShouldBeVisible(inCombat bool) bool {
	if inCombat {
		return true
	}
	return g.InvisibleTurns == 0
}

func (g *Ghost) TickInvisibility() {
	if g.InvisibleTurns > 0 {
		g.InvisibleTurns--
	}
}

func (g *Ghost) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
	if !g.IsAlive() {
		return
	}

	g.TickInvisibility()

	enemyPos := g.Actor.Point
	distance := Manhattan(enemyPos, playerPos)

	if ProcessAttack(g, aiCtx) {
		return
	}

	rng := aiCtx.GetRandomGenerator()
	if rng.IntN(100) < GhostTeleportChance {
		teleportPos := aiCtx.GetRandomPositionInRoom(enemyPos, level)
		if teleportPos.X >= 0 {
			aiCtx.MoveEnemy(enemyPos, teleportPos, level)
			return
		}
	}

	if distance <= g.Hostility {
		newPos := aiCtx.FindPathTo(enemyPos, playerPos, level)
		if newPos.X >= 0 && newPos.Y >= 0 {
			aiCtx.MoveEnemy(enemyPos, newPos, level)
		}
	}
}
