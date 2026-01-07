package domain

type Ghost struct {
	BaseEnemy
	InvisibleTurns int
	InCombat       bool
}

func NewGhost(depth int) *Ghost {
	scaling := 1.0 + float64(depth)*EnemyScalingFactor
	health := int(float64(GhostBaseHealth) * scaling)
	dexterity := int(float64(GhostBaseDexterity) * scaling)
	strength := int(float64(GhostBaseStrength) * scaling)
	hostility := int(float64(GhostBaseHostility) * scaling)

	return &Ghost{
		BaseEnemy: BaseEnemy{
			Actor: Actor{
				Stats: Stats{
					MaxHealth: health,
					Health:    health,
					Dexterity: dexterity,
					Strength:  strength,
				},
				State:    ActorStateNormal,
				Backpack: NewBackpack(),
			},
			EnemyType: EnemyTypeGhost,
			Hostility: hostility,
		},
		InvisibleTurns: 0,
		InCombat:       false,
	}
}

func (g *Ghost) GetEnemyType() EnemyType {
	return g.EnemyType
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

	if g.CanAttack(playerPos) {
		enemyActor := g.GetActor()
		playerActor := aiCtx.GetPlayerActor()
		if enemyActor.State != ActorStateSleep {
			resolver := aiCtx.GetCombatResolver()
			result := enemyActor.AttackWithResolver(playerActor, resolver)
			aiCtx.OnEnemyAttack(g, result)
		}
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
