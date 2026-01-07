package domain

type RandomGenerator interface {
	Float64() float64
	IntN(n int) int
	Shuffle(n int, swap func(i, j int))
}

type CombatResolver interface {
	GetWeaponStrength(actor *Actor) int
	GetRandomGenerator() RandomGenerator
	GetAttackModifiers(attacker *Actor, target *Actor) AttackModifiers
}

type AttackModifiers struct {
	GuaranteedHit bool
	FirstHitMiss  bool
}

type CombatResult struct {
	Hit    bool
	Damage int
	Killed bool
}

type attackOptions struct {
	WeaponStrength  int
	GuaranteedHit   bool
	FirstHitMiss    bool
	RandomGenerator RandomGenerator
}

func (a *Actor) Attack(target *Actor) CombatResult {
	return a.attackWithResolver(target, globalCombatResolver)
}

func (a *Actor) attackWithResolver(target *Actor, resolver CombatResolver) CombatResult {
	options := a.buildAttackOptions(target, resolver)
	return a.executeAttack(target, options)
}

func (a *Actor) buildAttackOptions(target *Actor, resolver CombatResolver) attackOptions {
	options := attackOptions{
		WeaponStrength:  0,
		GuaranteedHit:   false,
		FirstHitMiss:    false,
		RandomGenerator: nil,
	}

	if resolver != nil {
		options.WeaponStrength = resolver.GetWeaponStrength(a)
		options.RandomGenerator = resolver.GetRandomGenerator()
		modifiers := resolver.GetAttackModifiers(a, target)
		options.GuaranteedHit = modifiers.GuaranteedHit
		options.FirstHitMiss = modifiers.FirstHitMiss
	}

	return options
}

func (a *Actor) executeAttack(target *Actor, options attackOptions) CombatResult {
	if a.State == ActorStateSleep {
		return CombatResult{Hit: false, Damage: 0, Killed: false}
	}

	if options.FirstHitMiss {
		return CombatResult{Hit: false, Damage: 0, Killed: false}
	}

	var hit bool
	if options.GuaranteedHit {
		hit = true
	} else {
		hit = a.calculateHit(target, options.RandomGenerator)
	}

	if !hit {
		return CombatResult{Hit: false, Damage: 0, Killed: false}
	}

	damage := a.calculateDamage(options.WeaponStrength)
	target.AddHealth(-damage)
	killed := target.GetHealth() <= 0

	return CombatResult{
		Hit:    true,
		Damage: damage,
		Killed: killed,
	}
}

var globalCombatResolver CombatResolver

func SetGlobalCombatResolver(resolver CombatResolver) {
	globalCombatResolver = resolver
}

func (a *Actor) AttackWithResolver(target *Actor, resolver CombatResolver) CombatResult {
	return a.attackWithResolver(target, resolver)
}

func (a *Actor) calculateHit(target *Actor, rng RandomGenerator) bool {
	if rng == nil {
		return false
	}

	attackerDex := a.GetDexterity()
	targetDex := target.GetDexterity()

	dexDiff := float64(attackerDex - targetDex)
	hitChance := HitChanceBase + (dexDiff * HitChanceDexScale)

	if hitChance < 0 {
		hitChance = 0
	}
	if hitChance > 100 {
		hitChance = 100
	}

	roll := rng.Float64() * 100.0
	return roll < hitChance
}

func (a *Actor) calculateDamage(weaponStrength int) int {
	totalStrength := a.GetStrength() + weaponStrength
	damage := int(float64(DamageBase+totalStrength) * DamageStrengthScale)

	if damage < 1 {
		damage = 1
	}
	return damage
}
