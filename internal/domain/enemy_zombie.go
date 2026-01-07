package domain

type Zombie struct {
	BaseEnemy
}

func NewZombie(depth int) *Zombie {
	config := EnemyConfig{
		BaseHealth:    ZombieBaseHealth,
		BaseDexterity: ZombieBaseDexterity,
		BaseStrength:  ZombieBaseStrength,
		BaseHostility: ZombieBaseHostility,
	}
	stats := ScaleEnemyStats(config, depth)
	hostility := ScaleHostility(config.BaseHostility, depth)

	return &Zombie{
		BaseEnemy: NewBaseEnemy(EnemyTypeZombie, stats, hostility),
	}
}

func (z *Zombie) GetEnemyType() EnemyType {
	return z.EnemyType
}

func (z *Zombie) Name() string {
	return "Zombie"
}

func (z *Zombie) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
	ProcessStandardTurn(z, aiCtx, level, playerPos)
}
