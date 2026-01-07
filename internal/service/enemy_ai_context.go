package service

import (
	"rogue/internal/domain"
)

type EnemyAIContextImpl struct {
	game *Game
}

func NewEnemyAIContext(game *Game) *EnemyAIContextImpl {
	return &EnemyAIContextImpl{game: game}
}

func (ctx *EnemyAIContextImpl) FindPathTo(from, to domain.Point, level *domain.Level) domain.Point {
	return ctx.game.findPathToPlayer(from, to, level)
}

func (ctx *EnemyAIContextImpl) MoveEnemy(oldPos, newPos domain.Point, level *domain.Level) bool {
	return level.MoveEnemy(oldPos, newPos)
}

func (ctx *EnemyAIContextImpl) GetRandomPositionInRoom(currentPos domain.Point, level *domain.Level) domain.Point {
	return ctx.game.findRandomTeleportPosition(level, currentPos)
}

func (ctx *EnemyAIContextImpl) FindDiagonalMove(from, to domain.Point, level *domain.Level, direction domain.Point, directions []domain.Point) domain.Point {
	newPos := from.Add(direction)
	if level.Contains(newPos) {
		tile := level.Tiles[newPos.Y][newPos.X]
		if (tile.Kind == domain.TileFloor || tile.Kind == domain.TileCorridor) && level.GetEnemy(newPos) == nil {
			return newPos
		}
	}

	for _, dir := range directions {
		newPos := from.Add(dir)
		if level.Contains(newPos) {
			tile := level.Tiles[newPos.Y][newPos.X]
			if (tile.Kind == domain.TileFloor || tile.Kind == domain.TileCorridor) && level.GetEnemy(newPos) == nil {
				return newPos
			}
		}
	}

	return domain.Point{X: -1, Y: -1}
}

func (ctx *EnemyAIContextImpl) GetRandomGenerator() domain.RandomGenerator {
	return ctx.game.rng
}

func (ctx *EnemyAIContextImpl) GetCombatResolver() domain.CombatResolver {
	return NewGameCombatResolver(ctx.game)
}

func (ctx *EnemyAIContextImpl) GetPlayerActor() *domain.Actor {
	return &ctx.game.World.Player.Actor
}

func (ctx *EnemyAIContextImpl) OnEnemyAttack(enemy domain.Enemy, result domain.CombatResult) {
	ctx.game.RecordHitReceived()
	if result.Killed {
		ctx.game.handlePlayerDeath()
	}
}

func (ctx *EnemyAIContextImpl) OnPlayerDeath() {
	ctx.game.handlePlayerDeath()
}
