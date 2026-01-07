package service

import (
	"rogue/internal/domain"
)

func (g *Game) ProcessEnemyTurns() {
	enemies := g.World.Level.GetAllEnemies()
	aiCtx := NewAICtx(g)
	playerPos := g.World.Player.Point
	level := g.World.Level

	for _, enemy := range enemies {
		if enemy == nil || !enemy.IsAlive() {
			continue
		}

		enemy.ProcessTurn(aiCtx, level, playerPos)
	}
}

func (g *Game) findPathToPlayer(from, to domain.Point, level *domain.Level) domain.Point {
	queue := []domain.Point{from}
	visited := make(map[domain.Point]bool)
	visited[from] = true
	parent := make(map[domain.Point]domain.Point)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if domain.Manhattan(current, to) == domain.Combat.AttackDistance {
			path := []domain.Point{}
			pos := current
			for {
				if _, exists := parent[pos]; !exists {
					break
				}
				path = append([]domain.Point{pos}, path...)
				pos = parent[pos]
				if pos == from {
					break
				}
			}
			if len(path) > 0 {
				return path[0]
			}
			return current
		}

		for _, dir := range domain.Dirs4 {
			next := current.Add(dir)
			if !level.Contains(next) {
				continue
			}
			if visited[next] {
				continue
			}
			if !level.IsWalkableTile(next) {
				continue
			}
			if level.GetEnemy(next) != nil {
				continue
			}

			visited[next] = true
			parent[next] = current
			queue = append(queue, next)
		}
	}

	return domain.Point{X: -1, Y: -1}
}

func (g *Game) findRandomTeleportPosition(level *domain.Level, currentPos domain.Point) domain.Point {
	var currentRoom *domain.Room
	for i := range level.Rooms {
		if level.Rooms[i].Contains(currentPos) {
			currentRoom = &level.Rooms[i]
			break
		}
	}

	if currentRoom == nil {
		return domain.Point{X: -1, Y: -1}
	}

	return g.generator.findRandomPositionInRoom(level, *currentRoom, func(pos domain.Point) bool {
		return level.Tiles[pos.Y][pos.X].Kind == domain.TileFloor &&
			level.GetEnemy(pos) == nil && pos != currentPos
	})
}
