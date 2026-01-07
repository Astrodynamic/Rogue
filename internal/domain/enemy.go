package domain

type EnemyType int

const (
	EnemyTypeZombie EnemyType = iota
	EnemyTypeVampire
	EnemyTypeGhost
	EnemyTypeOgre
	EnemyTypeSnakeMage
)

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
}

type BaseEnemy struct {
	Actor
	EnemyType EnemyType
	Hostility int
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
	return distance <= 1
}
