package domain

type Mimic struct {
	BaseEnemy
	Revealed      bool
	DisguisedItem ItemKind
}

func NewMimic(depth int, disguisedItem ItemKind) *Mimic {
	config := EnemyConfig{
		BaseHealth:    MimicBaseHealth,
		BaseDexterity: MimicBaseDexterity,
		BaseStrength:  MimicBaseStrength,
		BaseHostility: MimicBaseHostility,
	}
	stats := ScaleEnemyStats(config, depth)
	hostility := ScaleHostility(config.BaseHostility, depth)

	return &Mimic{
		BaseEnemy:     NewBaseEnemy(EnemyTypeMimic, stats, hostility),
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

	if m.CanAttack(playerPos) {
		if !m.Revealed {
			m.Reveal()
		}
		ProcessAttack(m, aiCtx)
		return
	}

	if m.Revealed && distance <= m.Hostility {
		newPos := aiCtx.FindPathTo(enemyPos, playerPos, level)
		if newPos.X >= 0 && newPos.Y >= 0 {
			aiCtx.MoveEnemy(enemyPos, newPos, level)
		}
	}
}
