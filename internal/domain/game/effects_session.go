package game

func (s *GameSession) hasSleep() bool {
	for _, ef := range s.Effects {
		if ef.Type == EffectSleep && ef.TurnsLeft > 0 {
			return true
		}
	}
	return false
}

func (s *GameSession) addEffect(e Effect) {
	if e.TurnsLeft <= 0 {
		return
	}
	s.Effects = append(s.Effects, e)
}

func (s *GameSession) tickEffects() {
	kept := s.Effects[:0]
	for _, ef := range s.Effects {
		if ef.TurnsLeft > 0 {
			ef.TurnsLeft--
		}
		if ef.TurnsLeft <= 0 {
			// Expire: revert deltas for temporary buffs.
			switch ef.Type {
			case EffectDexBuff:
				s.Player.Dexterity -= ef.Delta
			case EffectStrBuff:
				s.Player.Strength -= ef.Delta
			case EffectMaxHPBuff:
				s.Player.MaxHealth -= ef.Delta
				if s.Player.Health > s.Player.MaxHealth {
					s.Player.Health = s.Player.MaxHealth
				}
				if s.Player.Health <= 0 {
					s.Player.Health = 1
				}
			}
			continue
		}
		kept = append(kept, ef)
	}
	s.Effects = kept
}
