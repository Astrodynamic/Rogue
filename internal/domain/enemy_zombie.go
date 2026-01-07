package domain

type Zombie struct {
	BaseEnemy
}

func NewZombie(depth int) *Zombie {
	config := EnemyConfig{
		BaseHealth:    ZombieBaseHP,
		BaseDexterity: ZombieBaseDex,
		BaseStrength:  ZombieBaseStr,
		BaseHostility: ZombieBaseHost,
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
	ProcessStdTurn(z, aiCtx, level, playerPos)
}
