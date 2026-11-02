package game

import (
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/physics"
)

const CountdownBeats = 3

type beatState struct {
	last   int
	start  int
	queued bool
}

func newBeatState(beat int) beatState {
	return beatState{last: beat, start: beat}
}

func (s *Session) requestShift() {
	switch s.level.Shift {
	case level.ShiftFree:
		s.TryShift()
	case level.ShiftQuantized:
		s.beats.queued = !s.beats.queued
	}
}

func (s *Session) onBeats() {
	from, to := s.beats.last, s.clock.Beat()
	if to <= from {
		return
	}
	s.beats.last = to
	switch s.level.Shift {
	case level.ShiftQuantized:
		if s.beats.queued {
			s.beats.queued = false
			s.TryShift()
		}
	case level.ShiftAuto:
		if s.nextAutoBeat(from) <= to {
			s.autoShift()
		}
	}
}

func (s *Session) nextAutoBeat(after int) int {
	every := s.level.AutoEvery
	return s.beats.start + ((after-s.beats.start)/every+1)*every
}

func (s *Session) autoShift() {
	s.active = s.active.Other()
	if physics.OverlapsSolid(s.player.Box, s.level.Layer(s.active)) {
		s.die(DeathCrushed)
	}
}

func (s *Session) ShiftQueued() bool { return s.beats.queued }

func (s *Session) AutoCountdown() int {
	if s.level.Shift != level.ShiftAuto || s.death != Alive {
		return 0
	}
	left := s.nextAutoBeat(s.beats.last) - s.beats.last
	if left > CountdownBeats {
		return 0
	}
	return left
}
