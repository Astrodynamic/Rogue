package game

type EnemyType uint8

const (
	EnemyZombie EnemyType = iota
	EnemyVampire
	EnemyGhost
	EnemyOgre
	EnemySnakeMage
)

type Enemy struct {
	Type      EnemyType
	Health    int
	Dexterity int
	Strength  int
	Hostility int

	Pos Point

	// Type-specific state (kept in domain so it can be saved/restored).
	VampireFirstHitAlwaysMiss bool
	GhostInvisible            bool
	GhostTeleportCooldown     int
	GhostInvisToggle          int
	OgreRestTurns             int
	OgreGuaranteedHitNext     bool
	SnakeDx                   int
	SnakeDy                   int
}

func (t EnemyType) Name() string {
	switch t {
	case EnemyZombie:
		return "Zombie"
	case EnemyVampire:
		return "Vampire"
	case EnemyGhost:
		return "Ghost"
	case EnemyOgre:
		return "Ogre"
	case EnemySnakeMage:
		return "Snake-Mage"
	default:
		return "Enemy"
	}
}
