package domain

type GameState struct {
	Statistics Statistics
	Depth      int
}

func NewGameState() *GameState {
	return &GameState{
		Statistics: *NewStatistics(),
		Depth:      0,
	}
}

func (gs *GameState) RecordTreasureCollected(amount int) {
	gs.Statistics.RecordTreasureCollected(amount)
}

func (gs *GameState) RecordEnemyDefeated() {
	gs.Statistics.RecordEnemyDefeated()
}

func (gs *GameState) RecordFoodConsumed() {
	gs.Statistics.RecordFoodConsumed()
}

func (gs *GameState) RecordElixirDrunk() {
	gs.Statistics.RecordElixirDrunk()
}

func (gs *GameState) RecordScrollRead() {
	gs.Statistics.RecordScrollRead()
}

func (gs *GameState) RecordHitDealt() {
	gs.Statistics.RecordHitDealt()
}

func (gs *GameState) RecordHitReceived() {
	gs.Statistics.RecordHitReceived()
}

func (gs *GameState) RecordTileTraveled() {
	gs.Statistics.RecordTileTraveled()
}

func (gs *GameState) AdvanceLevel() {
	gs.Depth++
	gs.Statistics.UpdateDeepestLevel(gs.Depth)
}

func (gs *GameState) GetDeepestLevel() int {
	return gs.Statistics.DeepestLevel
}

func (gs *GameState) ToPlaythroughStatistics() *PlaythroughStatistics {
	return NewPlaythroughStatistics(&gs.Statistics)
}
