package domain

type Mimic struct {
	BaseEnemy
	Revealed      bool
	DisguisedItem ItemKind
}

func NewMimic(depth int, disguisedItem ItemKind) *Mimic {
	scaling := 1.0 + float64(depth)*EnemyScalingFactor
	health := int(float64(MimicBaseHealth) * scaling)
	dexterity := int(float64(MimicBaseDexterity) * scaling)
	strength := int(float64(MimicBaseStrength) * scaling)
	hostility := int(float64(MimicBaseHostility) * scaling)

	return &Mimic{
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
			EnemyType: EnemyTypeMimic,
			Hostility: hostility,
		},
		Revealed:      false,
		DisguisedItem: disguisedItem,
	}
}

func (m *Mimic) GetEnemyType() EnemyType {
	return m.EnemyType
}

func (m *Mimic) Name() string {
	return "Mimic"
}

func (m *Mimic) IsRevealed() bool {
	return m.Revealed
}

func (m *Mimic) Reveal() {
	m.Revealed = true
}

func (m *Mimic) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
	if !m.IsAlive() {
		return
	}

	enemyPos := m.Actor.Point
	distance := Manhattan(enemyPos, playerPos)

	if distance <= 2 && !m.Revealed {
		m.Reveal()
	}

	if m.CanAttack(playerPos) {
		m.Reveal()
		enemyActor := m.GetActor()
		playerActor := aiCtx.GetPlayerActor()
		if enemyActor.State != ActorStateSleep {
			resolver := aiCtx.GetCombatResolver()
			result := enemyActor.AttackWithResolver(playerActor, resolver)
			aiCtx.OnEnemyAttack(m, result)
		}
		return
	}

	if m.Revealed && distance <= m.Hostility {
		newPos := aiCtx.FindPathTo(enemyPos, playerPos, level)
		if newPos.X >= 0 && newPos.Y >= 0 {
			aiCtx.MoveEnemy(enemyPos, newPos, level)
		}
	}
}
