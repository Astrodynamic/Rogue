package service

func (g *Game) RecordTreasureCollected(amount int) {
	g.World.GameState.RecordTreasureCollected(amount)
}

func (g *Game) RecordEnemyDefeated() {
	g.World.GameState.RecordEnemyDefeated()
}

func (g *Game) RecordFoodConsumed() {
	g.World.GameState.RecordFoodConsumed()
}

func (g *Game) RecordElixirDrunk() {
	g.World.GameState.RecordElixirDrunk()
}

func (g *Game) RecordScrollRead() {
	g.World.GameState.RecordScrollRead()
}

func (g *Game) RecordHitDealt() {
	g.World.GameState.RecordHitDealt()
}

func (g *Game) RecordHitReceived() {
	g.World.GameState.RecordHitReceived()
}

func (g *Game) RecordTileTraveled() {
	g.World.GameState.RecordTileTraveled()
}
