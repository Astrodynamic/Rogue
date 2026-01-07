package domain

import "time"

type Statistics struct {
	TreasureCollected int `json:"treasure_collected"`
	DeepestLevel      int `json:"deepest_level"`
	EnemiesDefeated   int `json:"enemies_defeated"`
	FoodConsumed      int `json:"food_consumed"`
	ElixirsDrunk      int `json:"elixirs_drunk"`
	ScrollsRead       int `json:"scrolls_read"`
	HitsDealt         int `json:"hits_dealt"`
	HitsReceived      int `json:"hits_received"`
	TilesTraveled     int `json:"tiles_traveled"`
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
	Statistics `json:",inline"`
	Timestamp  time.Time `json:"timestamp"`
	PlayerName string    `json:"player_name"`
}

func NewPlayStats(session *Statistics, playerName string) *PlayStats {
	return &PlayStats{
		Statistics: *session,
		Timestamp:  time.Now(),
		PlayerName: playerName,
	}
}
