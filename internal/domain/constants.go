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
