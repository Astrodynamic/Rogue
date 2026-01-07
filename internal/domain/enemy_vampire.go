package domain

type Vampire struct {
	BaseEnemy
	FirstHitMissed bool
}

func NewVampire(depth int) *Vampire {
	scaling := 1.0 + float64(depth)*EnemyScalingFactor
	health := int(float64(VampireBaseHealth) * scaling)
	dexterity := int(float64(VampireBaseDexterity) * scaling)
	strength := int(float64(VampireBaseStrength) * scaling)
	hostility := int(float64(VampireBaseHostility) * scaling)

	return &Vampire{
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
			EnemyType: EnemyTypeVampire,
			Hostility: hostility,
		},
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
		dirs := Dirs4
		rng := aiCtx.GetRandomGenerator()
		for i := len(dirs) - 1; i > 0; i-- {
			j := rng.IntN(i + 1)
			dirs[i], dirs[j] = dirs[j], dirs[i]
		}
		for _, dir := range dirs {
			newPos := enemyPos.Add(dir)
			if aiCtx.MoveEnemy(enemyPos, newPos, level) {
				break
			}
		}
	}
}
