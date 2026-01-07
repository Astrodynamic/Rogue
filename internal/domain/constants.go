package domain

const (
	Depth  = 21
	Width  = 120
	Height = 40
	Radius = 7
)

const (
	PlayerStartMaxHealth = 150
	PlayerStartHealth    = 150
	PlayerStartDexterity = 12
	PlayerStartStrength  = 8
)

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

const (
	HitChanceBase       = 60.0
	HitChanceDexScale   = 3.0
	DamageBase          = 1
	DamageStrengthScale = 1.0
)

const (
	EnemyScalingFactor   = 0.05
	EnemyCountBase       = 1
	EnemyCountScaling    = 0.15
	EnemyMaxCountPerRoom = 3
)

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

const (
	FoodBaseHealth     = 20
	FoodHealthPerDepth = 3
	ItemsPerRoomBase   = 4
	ItemSpawnChance    = 70
	FoodSpawnWeight    = 35
	EnemySpawnChance   = 50
)

const (
	PercentBase    = 100
	AttackDistance = 1
)

const (
	MaxStartPositionAttempts  = 50
	MaxRandomPositionAttempts = 20
)

const (
	RoomGridSize  = 3
	MinRoomWidth  = 6
	MinRoomHeight = 4
	RoomPadding   = 2
)

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
