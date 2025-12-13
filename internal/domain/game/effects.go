package game

type EffectType uint8

const (
	EffectDexBuff EffectType = iota
	EffectStrBuff
	EffectMaxHPBuff
	EffectSleep
)

type Effect struct {
	Type      EffectType
	TurnsLeft int
	Delta     int
}
