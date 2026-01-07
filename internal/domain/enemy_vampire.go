package domain

type Vampire struct {
	BaseEnemy
	FirstHitMissed bool
}

func NewVampire(depth int) *Vampire {
	config := EnemyConfig{
		BaseHealth:    VampireBaseHealth,
		BaseDexterity: VampireBaseDexterity,
		BaseStrength:  VampireBaseStrength,
		BaseHostility: VampireBaseHostility,
	}
	stats := ScaleEnemyStats(config, depth)
	hostility := ScaleHostility(config.BaseHostility, depth)

	return &Vampire{
		BaseEnemy:       NewBaseEnemy(EnemyTypeVampire, stats, hostility),
		FirstHitMissed: false,
	}
}

func (v *Vampire) GetEnemyType() EnemyType {
	return v.EnemyType
}

func (v *Vampire) Name() string {
	return "Vampire"
}

func (v *Vampire) MarkFirstHitMissed() {
	v.FirstHitMissed = true
}

func (v *Vampire) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
	if !v.IsAlive() {
		return
	}

	enemyPos := v.Actor.Point
	distance := Manhattan(enemyPos, playerPos)

	if v.CanAttack(playerPos) {
		enemyActor := v.GetActor()
		playerActor := aiCtx.GetPlayerActor()
		if enemyActor.State != ActorStateSleep {
			resolver := aiCtx.GetCombatResolver()
			result := enemyActor.AttackWithResolver(playerActor, resolver)
			if result.Hit {
				reductionEffect := NewMaxHealthEffect(-VampireMaxHealthReduction, -1)
				playerActor.AddEffect(reductionEffect)
			}
			aiCtx.OnEnemyAttack(v, result)
		}
		return
	}

	if distance <= v.Hostility {
		newPos := aiCtx.FindPathTo(enemyPos, playerPos, level)
		if newPos.X >= 0 && newPos.Y >= 0 {
			aiCtx.MoveEnemy(enemyPos, newPos, level)
		}
	} else {
		ProcessRandomMove(aiCtx, level, enemyPos)
	}
}
