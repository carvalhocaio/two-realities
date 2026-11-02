package game

const MaxPressure = 5

type PressureRule int

const (
	FasterWatchers PressureRule = iota + 1
	NoCheckpoints
	BeatGhost
	DeathRestartsLevel
	DeathRestartsAct
)

const (
	pressureSpeedFactor = 1.5
	pressureAlertTicks  = 16
)

func (s *Session) SetPressure(p int) { s.pressure = min(max(p, 0), MaxPressure) }
func (s *Session) Pressure() int     { return s.pressure }

func (s *Session) Has(rule PressureRule) bool { return s.pressure >= int(rule) }

func (s *Session) watcherSpeed() float64 {
	if s.Has(FasterWatchers) {
		return WatcherSpeed * pressureSpeedFactor
	}
	return WatcherSpeed
}

func (s *Session) baseAlert() int {
	if s.Has(FasterWatchers) {
		return pressureAlertTicks
	}
	return AlertTicks
}

func (s *Session) CheckpointsActive() bool { return !s.Has(NoCheckpoints) }

func (s *Session) afterDeath() {
	switch {
	case s.Has(DeathRestartsAct):
		s.enterLevel(s.actStart())
	case s.Has(DeathRestartsLevel):
		s.Restart()
	default:
		s.respawn()
	}
}

func (s *Session) actStart() int {
	i := s.index
	for i > 0 && s.levels[i-1].Act == s.level.Act {
		i--
	}
	return i
}

func (s *Session) ranked() bool { return !s.timer.levelAssisted && s.pressure == 0 }
