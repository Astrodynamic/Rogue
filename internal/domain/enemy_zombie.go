package domain

type Zombie struct {
	BaseEnemy
}

func NewZombie(depth int) *Zombie {
	scaling := 1.0 + float64(depth)*EnemyScalingFactor
	health := int(float64(ZombieBaseHealth) * scaling)
	dexterity := int(float64(ZombieBaseDexterity) * scaling)
	strength := int(float64(ZombieBaseStrength) * scaling)
	hostility := int(float64(ZombieBaseHostility) * scaling)

	return &Zombie{
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
			EnemyType: EnemyTypeZombie,
			Hostility: hostility,
		},
	}
}

func (z *Zombie) GetEnemyType() EnemyType {
	return z.EnemyType
}

func (z *Zombie) Name() string {
	return "Zombie"
}

func (z *Zombie) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
	if !z.IsAlive() {
		return
	}

	enemyPos := z.Actor.Point
	distance := Manhattan(enemyPos, playerPos)

	if z.CanAttack(playerPos) {
		enemyActor := z.GetActor()
		playerActor := aiCtx.GetPlayerActor()
		if enemyActor.State != ActorStateSleep {
			resolver := aiCtx.GetCombatResolver()
			result := enemyActor.AttackWithResolver(playerActor, resolver)
			aiCtx.OnEnemyAttack(z, result)
		}
		return
	}

	if distance <= z.Hostility {
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
