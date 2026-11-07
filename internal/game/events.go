package game

import "github.com/carvalhocaio/two-realities/internal/world"

type Event int

const (
	EventJump Event = iota
	EventLand
	EventShiftGrey
	EventShiftYellow
	EventShiftBlocked
	EventShiftQueued
	EventFragment
	EventSecret
	EventCheckpoint
	EventTransmission
	EventSpotted
	EventDeath
	EventLamp
	EventExit
	EventCount
)

func (s *Session) emit(e Event) { s.events = append(s.events, e) }

func (s *Session) emitShift(to world.World) {
	if to == world.Grey {
		s.emit(EventShiftGrey)
		return
	}
	s.emit(EventShiftYellow)
}

func (s *Session) Events() []Event {
	out := s.events
	s.events = nil
	return out
}
