package domain

type Ogre struct {
	BaseEnemy
	Resting   bool
	RestTurns int
	MovesLeft int
}

func NewOgre(depth int) *Ogre {
	scaling := 1.0 + float64(depth)*EnemyScalingFactor
	health := int(float64(OgreBaseHealth) * scaling)
	dexterity := int(float64(OgreBaseDexterity) * scaling)
	strength := int(float64(OgreBaseStrength) * scaling)
	hostility := int(float64(OgreBaseHostility) * scaling)

	return &Ogre{
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
			EnemyType: EnemyTypeOgre,
			Hostility: hostility,
		},
		Resting:   false,
		RestTurns: 0,
		MovesLeft: OgreMovesPerTurn,
	}
}

func (o *Ogre) GetEnemyType() EnemyType {
	return o.EnemyType
}

func (o *Ogre) Name() string {
	return "Ogre"
}

func (o *Ogre) StartRest() {
	o.Resting = true
	o.RestTurns = OgreRestTurns
}

func (o *Ogre) TickRest() {
	if o.Resting {
		o.RestTurns--
		if o.RestTurns <= 0 {
			o.Resting = false
		}
	}
}

func (o *Ogre) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
	if !o.IsAlive() {
		return
	}

	o.TickRest()

	enemyPos := o.Actor.Point
	distance := Manhattan(enemyPos, playerPos)

	if o.CanAttack(playerPos) {
		enemyActor := o.GetActor()
		playerActor := aiCtx.GetPlayerActor()
		if enemyActor.State != ActorStateSleep {
			resolver := aiCtx.GetCombatResolver()
			result := enemyActor.AttackWithResolver(playerActor, resolver)
			aiCtx.OnEnemyAttack(o, result)
			o.StartRest()
		}
		return
	}

	if distance <= o.Hostility && !o.Resting && o.MovesLeft > 0 {
		newPos := aiCtx.FindPathTo(enemyPos, playerPos, level)
		if newPos.X >= 0 && newPos.Y >= 0 {
			aiCtx.MoveEnemy(enemyPos, newPos, level)
			o.MovesLeft--
			if o.MovesLeft <= 0 {
				o.MovesLeft = OgreMovesPerTurn
			}

			if o.MovesLeft > 0 {
				enemyPos = newPos
				newPos2 := aiCtx.FindPathTo(enemyPos, playerPos, level)
				if newPos2.X >= 0 && newPos2.Y >= 0 {
					aiCtx.MoveEnemy(enemyPos, newPos2, level)
					o.MovesLeft--
				}
			}
		}
	}
}
