package domain

import "time"

type Statistics struct {
	TreasureCollected int
	DeepestLevel      int
	EnemiesDefeated   int
	FoodConsumed      int
	ElixirsDrunk      int
	ScrollsRead       int
	HitsDealt         int
	HitsReceived      int
	TilesTraveled     int
}

func NewStatistics() *Statistics {
	return &Statistics{
		DeepestLevel: 1,
	}
}

func (s *Statistics) RecordTreasureCollected(amount int) {
	s.TreasureCollected += amount
}

func (s *Statistics) RecordEnemyDefeated() {
	s.EnemiesDefeated++
}

func (s *Statistics) RecordFoodConsumed() {
	s.FoodConsumed++
}

func (s *Statistics) RecordElixirDrunk() {
	s.ElixirsDrunk++
}

func (s *Statistics) RecordScrollRead() {
	s.ScrollsRead++
}

func (s *Statistics) RecordHitDealt() {
	s.HitsDealt++
}

func (s *Statistics) RecordHitReceived() {
	s.HitsReceived++
}

func (s *Statistics) RecordTileTraveled() {
	s.TilesTraveled++
}

func (s *Statistics) UpdateDeepestLevel(level int) {
	if level > s.DeepestLevel {
		s.DeepestLevel = level
	}
}

type PlayStats struct {
	Statistics
	Timestamp  time.Time
	PlayerName string
}

func NewPlayStats(session *Statistics, playerName string) *PlayStats {
	return &PlayStats{
		Statistics: *session,
		Timestamp:  time.Now(),
		PlayerName: playerName,
	}
}
