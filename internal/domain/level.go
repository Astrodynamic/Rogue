package domain

type EnemyWithPos struct {
	Enemy Enemy
	Pos   Point
}

type Level struct {
	Rect
	Tiles     [][]Tile
	Rooms     []Room
	Corridors []Corridor
	Exits     []Point
	Items     map[Point]Item
	Enemies   map[Point]Enemy
}

func NewLevel(width, height int) *Level {
	tiles := make([][]Tile, height)
	for i := range tiles {
		tiles[i] = make([]Tile, width)
	}

	return &Level{
		Rect:      NewRect(0, 0, width, height),
		Tiles:     tiles,
		Rooms:     make([]Room, 0),
		Corridors: make([]Corridor, 0),
		Exits:     make([]Point, 0),
		Items:     make(map[Point]Item),
		Enemies:   make(map[Point]Enemy),
	}
}

func (l *Level) Clear() {
	l.ClearTiles()
	l.Enemies = make(map[Point]Enemy)
	l.Items = make(map[Point]Item)
}

func (l *Level) ClearTiles() {
	for y := l.Y; y < l.Y+l.H; y++ {
		for x := l.X; x < l.X+l.W; x++ {
			l.Tiles[y][x].Clear()
		}
	}
}

func (l *Level) AddRoom(room Room) {
	l.Rooms = append(l.Rooms, room)
	l.FillRoom(room)
}

func (l *Level) FillRoom(room Room) {
	for y := room.Y; y < room.Y+room.H; y++ {
		for x := room.X; x < room.X+room.W; x++ {
			if x == room.X || y == room.Y || x == room.X+room.W-1 || y == room.Y+room.H-1 {
				l.Tiles[y][x].SetKind(TileWall)
			} else {
				l.Tiles[y][x].SetKind(TileFloor)
			}
		}
	}
}

func (l *Level) AddCorridor(corridor Corridor) {
	l.Corridors = append(l.Corridors, corridor)
	l.FillCorridor(corridor)
}

func (l *Level) FillCorridor(corridor Corridor) {
	for _, point := range corridor.Points {
		l.Tiles[point.Y][point.X].SetKind(TileCorridor)
	}

	for _, point := range corridor.Points {
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if x == 0 && y == 0 {
					continue
				}

				if point.Y+y < 0 || point.Y+y >= l.H || point.X+x < 0 || point.X+x >= l.W {
					continue
				}

				tile := &l.Tiles[point.Y+y][point.X+x]
				if tile.Kind == TileNone {
					(*tile).SetKind(TileWall)
				}
			}
		}
	}
}

func (l *Level) AddExit(exit Point) {
	l.Exits = append(l.Exits, exit)
	l.FillExit(exit)
}

func (l *Level) FillExit(exit Point) {
	l.Tiles[exit.Y][exit.X].SetKind(TileExit)
}

func (l *Level) BlocksSight(p Point) bool {
	if !l.Contains(p) {
		return true
	}
	return l.Tiles[p.Y][p.X].TestFlags(TileOpaque)
}

func (l *Level) SetVisible(p Point, visible bool) {
	if !l.Contains(p) {
		return
	}

	if visible {
		l.Tiles[p.Y][p.X].SetFlags(TileVisible, true)
	} else {
		l.Tiles[p.Y][p.X].SetFlags(TileVisible, false)
	}
}

func (l *Level) ClearVisible() {
	for y := l.Y; y < l.Y+l.H; y++ {
		for x := l.X; x < l.X+l.W; x++ {
			l.Tiles[y][x].SetFlags(TileVisible, false)
		}
	}
}

func (l *Level) UpdateExplored() {
	for y := l.Y; y < l.Y+l.H; y++ {
		for x := l.X; x < l.X+l.W; x++ {
			if l.Tiles[y][x].TestFlags(TileVisible) {
				l.Tiles[y][x].SetFlags(TileExplored, true)
			}
		}
	}
}

func (l *Level) IsWalkableTile(p Point) bool {
	if !l.Contains(p) {
		return false
	}
	tile := l.Tiles[p.Y][p.X]
	return tile.Kind == TileFloor || tile.Kind == TileCorridor
}

func (l *Level) AddItem(p Point, item Item) {
	if l.IsWalkableTile(p) {
		if l.Items[p] == nil {
			l.Items[p] = item
		}
	}
}

func (l *Level) RemoveItem(p Point) Item {
	item, exists := l.Items[p]
	if exists {
		delete(l.Items, p)
	}
	return item
}

func (l *Level) GetItem(p Point) Item {
	return l.Items[p]
}

func (l *Level) GetAdjacentDropPoint(from Point) Point {
	for _, dir := range Dirs8 {
		pos := from.Add(dir)
		if !l.Contains(pos) {
			continue
		}
		if !l.IsWalkableTile(pos) {
			continue
		}
		if l.GetItem(pos) != nil {
			continue
		}
		if l.GetEnemy(pos) != nil {
			continue
		}
		return pos
	}
	return Point{X: -1, Y: -1}
}

func (l *Level) AddEnemy(p Point, enemy Enemy) {
	if l.IsWalkableTile(p) {
		if enemy.GetActor() != nil {
			enemy.GetActor().Point = p
		}
		l.Enemies[p] = enemy
	}
}

func (l *Level) RemoveEnemy(p Point) Enemy {
	enemy, exists := l.Enemies[p]
	if exists {
		delete(l.Enemies, p)
	}
	return enemy
}

func (l *Level) GetEnemy(p Point) Enemy {
	return l.Enemies[p]
}

func (l *Level) MoveEnemy(oldPos, newPos Point) bool {
	enemy := l.Enemies[oldPos]
	if enemy == nil {
		return false
	}
	if !l.IsWalkableTile(newPos) {
		return false
	}
	if l.GetEnemy(newPos) != nil {
		return false
	}
	delete(l.Enemies, oldPos)
	l.Enemies[newPos] = enemy
	if enemy.GetActor() != nil {
		enemy.GetActor().Point = newPos
	}
	return true
}

func (l *Level) GetEnemiesInRoom(room Room) []EnemyWithPos {
	var enemies []EnemyWithPos
	for pos, enemy := range l.Enemies {
		if room.Contains(pos) {
			enemies = append(enemies, EnemyWithPos{
				Enemy: enemy,
				Pos:   pos,
			})
		}
	}
	return enemies
}

func (l *Level) GetAllEnemies() []EnemyWithPos {
	var enemies []EnemyWithPos
	for pos, enemy := range l.Enemies {
		enemies = append(enemies, EnemyWithPos{
			Enemy: enemy,
			Pos:   pos,
		})
	}
	return enemies
}
