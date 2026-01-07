package domain

const (
	Depth  = 21
	Width  = 120
	Height = 40
	Radius = 7
)

const (
	PlayerStartMaxHealth = 100
	PlayerStartHealth    = 100
	PlayerStartDexterity = 10
	PlayerStartStrength  = 10
)

const (
	ZombieBaseHealth    = 50
	ZombieBaseDexterity = 5
	ZombieBaseStrength  = 8
	ZombieBaseHostility = 5
)

const (
	VampireBaseHealth         = 80
	VampireBaseDexterity      = 15
	VampireBaseStrength       = 10
	VampireBaseHostility      = 10
	VampireMaxHealthReduction = 5
)

const (
	GhostBaseHealth     = 30
	GhostBaseDexterity  = 15
	GhostBaseStrength   = 4
	GhostBaseHostility  = 3
	GhostTeleportChance = 30
	GhostInvisibleTurns = 3
)

const (
	OgreBaseHealth    = 100
	OgreBaseDexterity = 4
	OgreBaseStrength  = 20
	OgreBaseHostility = 6
	OgreMovesPerTurn  = 2
	OgreRestTurns     = 1
)

const (
	SnakeMageBaseHealth    = 40
	SnakeMageBaseDexterity = 20
	SnakeMageBaseStrength  = 6
	SnakeMageBaseHostility = 12
	SnakeMageSleepChance   = 30
)

const (
	HitChanceBase       = 50.0
	HitChanceDexScale   = 2.0
	DamageBase          = 1
	DamageStrengthScale = 1.5
)

const (
	EnemyScalingFactor   = 0.1
	EnemyCountBase       = 1
	EnemyCountScaling    = 0.3
	EnemyMaxCountPerRoom = 5
)

const (
	TreasureBaseValue      = 5
	TreasureStatMultiplier = 2
	TreasureDepthBonus     = 3
)

const (
	ZombieMaxDepth    = 10
	VampireMinDepth   = 5
	GhostMinDepth     = 3
	OgreMinDepth      = 8
	SnakeMageMinDepth = 10
)
