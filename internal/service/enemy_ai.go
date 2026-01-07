package service

import (
	"rogue/internal/domain"
)

func (g *Game) ProcessEnemyTurns() {
	enemies := g.World.Level.GetAllEnemies()
	aiCtx := NewEnemyAIContext(g)
	playerPos := g.World.Player.Point
	level := g.World.Level

	for _, enemyWithPos := range enemies {
		if enemyWithPos.Enemy == nil || !enemyWithPos.Enemy.IsAlive() {
			continue
		}

		enemyWithPos.Enemy.ProcessTurn(aiCtx, level, playerPos)
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

		if domain.Manhattan(current, to) == 1 {
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
			tile := level.Tiles[next.Y][next.X]
			if tile.Kind != domain.TileFloor && tile.Kind != domain.TileCorridor {
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

	maxAttempts := 20
	for i := 0; i < maxAttempts; i++ {
		x := g.rng.IntN(currentRoom.W-2) + currentRoom.X + 1
		y := g.rng.IntN(currentRoom.H-2) + currentRoom.Y + 1
		pos := domain.Point{X: x, Y: y}

		if level.Tiles[y][x].Kind == domain.TileFloor {
			if level.GetEnemy(pos) == nil && pos != currentPos {
				return pos
			}
		}
	}

	return domain.Point{X: -1, Y: -1}
}
