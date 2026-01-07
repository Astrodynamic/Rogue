package domain

type EnemyType int

const (
	EnemyTypeZombie EnemyType = iota
	EnemyTypeVampire
	EnemyTypeGhost
	EnemyTypeOgre
	EnemyTypeSnakeMage
	EnemyTypeMimic
)

type EnemyConfig struct {
	BaseHealth    int
	BaseDexterity int
	BaseStrength  int
	BaseHostility int
}

func ScaleEnemyStats(config EnemyConfig, depth int) Stats {
	scaling := 1.0 + float64(depth)*EnemyGeneration.ScalingFactor
	return Stats{
		MaxHealth: int(float64(config.BaseHealth) * scaling),
		Health:    int(float64(config.BaseHealth) * scaling),
		Dexterity: int(float64(config.BaseDexterity) * scaling),
		Strength:  int(float64(config.BaseStrength) * scaling),
	}
}

func ScaleHostility(baseHostility int, depth int) int {
	scaling := 1.0 + float64(depth)*EnemyGeneration.ScalingFactor
	return int(float64(baseHostility) * scaling)
}

func NewBaseEnemy(enemyType EnemyType, stats Stats, hostility int) BaseEnemy {
	return BaseEnemy{
		Actor: Actor{
			Stats:    stats,
			State:    ActorStateNormal,
			Backpack: NewBackpack(),
		},
		EnemyType: enemyType,
		Hostility: hostility,
	}
}

type EnemyAIContext interface {
	FindPathTo(from, to Point, level *Level) Point
	MoveEnemy(oldPos, newPos Point, level *Level) bool
	GetRandomPositionInRoom(currentPos Point, level *Level) Point
	FindDiagonalMove(from, to Point, level *Level, direction Point, directions []Point) Point
	GetRandomGenerator() RandomGenerator
	GetCombatResolver() CombatResolver
	GetPlayerActor() *Actor
	OnEnemyAttack(enemy Enemy, result CombatResult)
	OnPlayerDeath()
}

type Enemy interface {
	GetEnemyType() EnemyType
	GetActor() *Actor
	GetHostility() int
	CanAttack(playerPos Point) bool
	IsAlive() bool
	ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point)
	Name() string
}

type BaseEnemy struct {
	Actor
	EnemyType EnemyType
	Hostility int
}

func (be *BaseEnemy) GetEnemyType() EnemyType {
	return be.EnemyType
}

func (be *BaseEnemy) GetActor() *Actor {
	return &be.Actor
}

func (be *BaseEnemy) GetHostility() int {
	return be.Hostility
}

func (be *BaseEnemy) IsAlive() bool {
	return be.Actor.Health > 0
}

func (be *BaseEnemy) CanAttack(playerPos Point) bool {
	if !be.IsAlive() {
		return false
	}
	distance := Manhattan(be.Actor.Point, playerPos)
	return distance <= Combat.AttackDistance
}

func (be *BaseEnemy) ProcessTurn(aiCtx EnemyAIContext, level *Level, playerPos Point) {
}

func ProcessStdTurn(enemy Enemy, aiCtx EnemyAIContext, level *Level, playerPos Point) {
	if !enemy.IsAlive() {
		return
	}

	enemyActor := enemy.GetActor()
	enemyPos := enemyActor.Point
	distance := Manhattan(enemyPos, playerPos)

	if enemy.CanAttack(playerPos) {
		playerActor := aiCtx.GetPlayerActor()
		if enemyActor.State != ActorStateSleep {
			resolver := aiCtx.GetCombatResolver()
			result := enemyActor.AttackWithResolver(playerActor, resolver)
			aiCtx.OnEnemyAttack(enemy, result)
		}
		return
	}

	if distance <= enemy.GetHostility() {
		newPos := aiCtx.FindPathTo(enemyPos, playerPos, level)
		if newPos.X >= 0 && newPos.Y >= 0 {
			aiCtx.MoveEnemy(enemyPos, newPos, level)
		}
	} else {
		ProcessRandMove(aiCtx, level, enemyPos)
	}
}

func ProcessRandMove(aiCtx EnemyAIContext, level *Level, enemyPos Point) {
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

func ProcessAttack(enemy Enemy, aiCtx EnemyAIContext) bool {
	if !enemy.IsAlive() || !enemy.CanAttack(aiCtx.GetPlayerActor().Point) {
		return false
	}

	enemyActor := enemy.GetActor()
	playerActor := aiCtx.GetPlayerActor()
	if enemyActor.State != ActorStateSleep {
		resolver := aiCtx.GetCombatResolver()
		result := enemyActor.AttackWithResolver(playerActor, resolver)
		aiCtx.OnEnemyAttack(enemy, result)
		return true
	}
	return false
}
