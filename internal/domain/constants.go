package domain

var WorldConfig = struct {
	Depth  int
	Width  int
	Height int
	Radius int
}{
	Depth:  21,
	Width:  120,
	Height: 40,
	Radius: 7,
}

var PlayerConfig = struct {
	StartMaxHealth int
	StartHealth    int
	StartDexterity int
	StartStrength  int
}{
	StartMaxHealth: 150,
	StartHealth:    150,
	StartDexterity: 12,
	StartStrength:  8,
}

const (
	ZombieBaseHealth    = 25
	ZombieBaseDexterity = 3
	ZombieBaseStrength  = 4
	ZombieBaseHostility = 4
)

const (
	VampireBaseHealth         = 40
	VampireBaseDexterity      = 10
	VampireBaseStrength       = 6
	VampireBaseHostility      = 8
	VampireMaxHealthReduction = 3
)

const (
	GhostBaseHealth     = 15
	GhostBaseDexterity  = 12
	GhostBaseStrength   = 2
	GhostBaseHostility  = 3
	GhostTeleportChance = 25
	GhostInvisibleTurns = 2
)

const (
	OgreBaseHealth    = 60
	OgreBaseDexterity = 3
	OgreBaseStrength  = 12
	OgreBaseHostility = 5
	OgreMovesPerTurn  = 2
	OgreRestTurns     = 1
)

const (
	SnakeMageBaseHealth    = 20
	SnakeMageBaseDexterity = 14
	SnakeMageBaseStrength  = 3
	SnakeMageBaseHostility = 10
	SnakeMageSleepChance   = 20
)

const (
	MimicBaseHealth    = 30
	MimicBaseDexterity = 10
	MimicBaseStrength  = 3
	MimicBaseHostility = 2
)

var Combat = struct {
	HitChanceBase       float64
	HitChanceDexScale   float64
	DamageBase          int
	DamageStrengthScale float64
	AttackDistance      int
	PercentBase         int
}{
	HitChanceBase:       60.0,
	HitChanceDexScale:   3.0,
	DamageBase:          1,
	DamageStrengthScale: 1.0,
	AttackDistance:      1,
	PercentBase:         100,
}

var EnemyGeneration = struct {
	ScalingFactor   float64
	CountBase       int
	CountScaling    float64
	MaxCountPerRoom int
	SpawnChance     int
}{
	ScalingFactor:   0.05,
	CountBase:       1,
	CountScaling:    0.15,
	MaxCountPerRoom: 3,
	SpawnChance:     50,
}

const (
	TreasureBaseValue      = 5
	TreasureStatMultiplier = 2
	TreasureDepthBonus     = 3
)

const (
	ZombieMaxDepth    = 12
	VampireMinDepth   = 7
	GhostMinDepth     = 4
	OgreMinDepth      = 10
	SnakeMageMinDepth = 14
	MimicMinDepth     = 6
)

var ItemGeneration = struct {
	FoodBaseHealth     int
	FoodHealthPerDepth int
	ItemsPerRoomBase   int
	ItemSpawnChance    int
	FoodSpawnWeight    int
}{
	FoodBaseHealth:     20,
	FoodHealthPerDepth: 3,
	ItemsPerRoomBase:   4,
	ItemSpawnChance:    70,
	FoodSpawnWeight:    35,
}

var LevelGeneration = struct {
	RoomGridSize              int
	MinRoomWidth              int
	MinRoomHeight             int
	RoomPadding               int
	MaxStartPositionAttempts  int
	MaxRandomPositionAttempts int
}{
	RoomGridSize:              3,
	MinRoomWidth:              6,
	MinRoomHeight:             4,
	RoomPadding:               2,
	MaxStartPositionAttempts:  50,
	MaxRandomPositionAttempts: 20,
}

const (
	ItemElixirThresholdOffset   = 20
	ItemScrollThresholdOffset   = 15
	ItemTreasureThresholdOffset = 10
	ItemWeaponThresholdOffset   = 10
	ItemDifficultyBonus         = 10
	ItemMinPerRoom              = 2
	ItemDepthDivisor            = 7
	ElixirKindCount             = 4
	ScrollKindCount             = 5
	ElixirBaseAmount            = 5
	ElixirBaseDuration          = 8
	ElixirDurationDivisor       = 2
	ScrollBaseAmount            = 2
	ScrollAmountDivisor         = 2
	ScrollRegenDivisor          = 4
	TreasureBaseDropValue       = 10
	TreasureDepthMultiplier     = 5
	WeaponBaseStrength          = 2
	WeaponStrengthDivisor       = 3
	ArmorStatBase               = 1
	ArmorStatDivisor            = 4
	ArmorHealthBase             = 5
	ArmorHealthDivisor          = 2
	ArmorStatRollCount          = 3
)

const (
	EnemySpawnChanceDifficultyBonus   = 15
	EnemySpawnChanceDifficultyPenalty = 15
)

const (
	DifficultyHealthHighThreshold    = 0.7
	DifficultyHealthLowThreshold     = 0.3
	DifficultyHealthHighAdjustment   = 0.05
	DifficultyHealthLowAdjustment    = -0.05
	DifficultyHitRatioHighThreshold  = 0.6
	DifficultyHitRatioLowThreshold   = 0.4
	DifficultyHitRatioHighAdjustment = 0.03
	DifficultyHitRatioLowAdjustment  = -0.03
	DifficultyFoodPerLevelThreshold  = 3.0
	DifficultyFoodAdjustment         = -0.02
	DifficultyMinFactor              = 0.5
	DifficultyMaxFactor              = 1.5
)

const (
	UIMinWidthThreshold  = 80
	UIMinPanelWidth      = 20
	UIDefaultPanelWidth  = 28
	UIMinHeightThreshold = 30
	UIMinLogHeight       = 10
	UIDefaultLogHeight   = 15
	UIStatsHeight        = 6
	UIEquipmentHeight    = 8
)

const (
	MaxUniqueItems = 20
)
