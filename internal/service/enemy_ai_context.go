package service

import (
	"rogue/internal/domain"
)

type AICtxImpl struct {
	game *Game
}

func NewAICtx(game *Game) *AICtxImpl {
	return &AICtxImpl{game: game}
}

func (ctx *AICtxImpl) FindPathTo(from, to domain.Point, level *domain.Level) domain.Point {
	return ctx.game.findPathToPlayer(from, to, level)
}

func (ctx *AICtxImpl) MoveEnemy(oldPos, newPos domain.Point, level *domain.Level) bool {
	return level.MoveEnemy(oldPos, newPos)
}

func (ctx *AICtxImpl) GetRandomPositionInRoom(currentPos domain.Point, level *domain.Level) domain.Point {
	return ctx.game.findRandomTeleportPosition(level, currentPos)
}

func (ctx *AICtxImpl) FindDiagonalMove(from, to domain.Point, level *domain.Level, direction domain.Point, directions []domain.Point) domain.Point {
	newPos := from.Add(direction)
	if level.IsWalkableTile(newPos) && level.GetEnemy(newPos) == nil {
		return newPos
	}

	for _, dir := range directions {
		newPos := from.Add(dir)
		if level.IsWalkableTile(newPos) && level.GetEnemy(newPos) == nil {
			return newPos
		}
	}

	return domain.Point{X: -1, Y: -1}
}

func (ctx *AICtxImpl) GetRandomGenerator() domain.RandomGenerator {
	return ctx.game.rng
}

func (ctx *AICtxImpl) GetCombatResolver() domain.CombatResolver {
	return NewCombatRes(ctx.game)
}

func (ctx *AICtxImpl) GetPlayerActor() *domain.Actor {
	return &ctx.game.World.Player.Actor
}

func (ctx *AICtxImpl) OnEnemyAttack(enemy domain.Enemy, result domain.CombatResult) {
	ctx.game.RecordHitReceived()
	if result.Killed {
		ctx.game.handlePlayerDeath()
	}
}

func (ctx *AICtxImpl) OnPlayerDeath() {
	ctx.game.handlePlayerDeath()
}
