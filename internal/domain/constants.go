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
	ZombieBaseHP   = 25
	ZombieBaseDex  = 3
	ZombieBaseStr  = 4
	ZombieBaseHost = 4
)

const (
	VampireBaseHP   = 40
	VampireBaseDex  = 10
	VampireBaseStr  = 6
	VampireBaseHost = 8
	VampireMaxHPRed = 3
)

const (
	GhostBaseHP     = 15
	GhostBaseDex    = 12
	GhostBaseStr    = 2
	GhostBaseHost   = 3
	GhostTeleChance = 25
	GhostInvisTurns = 2
)

const (
	OgreBaseHP    = 60
	OgreBaseDex   = 3
	OgreBaseStr   = 12
	OgreBaseHost  = 5
	OgreMovesPerT = 2
	OgreRestTurns = 1
)

const (
	SnakeMageBaseHP   = 20
	SnakeMageBaseDex  = 14
	SnakeMageBaseStr  = 3
	SnakeMageBaseHost = 10
	SnakeMageSleepCh  = 20
)

const (
	MimicBaseHP   = 30
	MimicBaseDex  = 10
	MimicBaseStr  = 3
	MimicBaseHost = 2
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
	TreasureBaseVal  = 5
	TreasureStatMult = 2
	TreasureDepthBon = 3
)

const (
	ZombieMaxDepth    = 12
	VampireMinDepth   = 7
	GhostMinDepth     = 4
	OgreMinDepth      = 10
	SnakeMageMinDepth = 14
	MimicMinDepth     = 6
)

var ItemGen = struct {
	FoodBaseHP     int
	FoodHPPerDepth int
	ItemsPerRoom   int
	ItemSpawnCh    int
	FoodSpawnWgt   int
}{
	FoodBaseHP:     20,
	FoodHPPerDepth: 3,
	ItemsPerRoom:   4,
	ItemSpawnCh:    70,
	FoodSpawnWgt:   35,
}

var LevelGen = struct {
	RoomGridSize        int
	MinRoomW            int
	MinRoomH            int
	RoomPadding         int
	MaxStartPosAttempts int
	MaxRandPosAttempts  int
}{
	RoomGridSize:        3,
	MinRoomW:            6,
	MinRoomH:            4,
	RoomPadding:         2,
	MaxStartPosAttempts: 50,
	MaxRandPosAttempts:  20,
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
