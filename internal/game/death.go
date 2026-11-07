package game

type DeathCause int

const (
	Alive DeathCause = iota
	DeathSpikes
	DeathCaught
	DeathFall
	DeathCrushed
	DeathBeam
)

const DeathFreezeTicks = 30

func (s *Session) die(cause DeathCause) {
	if s.death != Alive {
		return
	}
	s.death = cause
	s.emit(EventDeath)
	if cause == DeathFall && !s.shifted {
		s.shiftHint = true
	}
	s.deathTimer = DeathFreezeTicks
}

func (s *Session) Death() DeathCause { return s.death }

func (s *Session) DeathProgress() float64 {
	if s.death == Alive {
		return 0
	}
	return 1 - float64(s.deathTimer)/DeathFreezeTicks
}

func (s *Session) tickDeath() {
	s.deathTimer--
	if s.deathTimer <= 0 {
		s.afterDeath()
	}
}
