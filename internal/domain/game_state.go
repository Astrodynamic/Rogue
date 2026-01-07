package domain

type GameState struct {
	Statistics       Statistics
	Depth            int
	DifficultyFactor float64
}

func NewGameState() *GameState {
	return &GameState{
		Statistics:       *NewStatistics(),
		Depth:            1,
		DifficultyFactor: 1.0,
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

func (gs *GameState) ToPlaythroughStatistics(playerName string) *PlaythroughStatistics {
	return NewPlaythroughStatistics(&gs.Statistics, playerName)
}

func (gs *GameState) AdjustDifficulty(playerHealth, playerMaxHealth int) {

	hitRatio := 0.0
	if gs.Statistics.HitsDealt+gs.Statistics.HitsReceived > 0 {
		hitRatio = float64(gs.Statistics.HitsDealt) / float64(gs.Statistics.HitsDealt+gs.Statistics.HitsReceived)
	}

	healthRatio := float64(playerHealth) / float64(playerMaxHealth)

	adjustment := 0.0

	if healthRatio > DifficultyHealthHighThreshold {
		adjustment += DifficultyHealthHighAdjustment
	} else if healthRatio < DifficultyHealthLowThreshold {
		adjustment += DifficultyHealthLowAdjustment
	}

	if hitRatio > DifficultyHitRatioHighThreshold {
		adjustment += DifficultyHitRatioHighAdjustment
	} else if hitRatio < DifficultyHitRatioLowThreshold {
		adjustment += DifficultyHitRatioLowAdjustment
	}

	foodPerLevel := float64(gs.Statistics.FoodConsumed) / float64(gs.Depth)
	if foodPerLevel > DifficultyFoodPerLevelThreshold {
		adjustment += DifficultyFoodAdjustment
	}

	gs.DifficultyFactor += adjustment

	if gs.DifficultyFactor < DifficultyMinFactor {
		gs.DifficultyFactor = DifficultyMinFactor
	}
	if gs.DifficultyFactor > DifficultyMaxFactor {
		gs.DifficultyFactor = DifficultyMaxFactor
	}
}
