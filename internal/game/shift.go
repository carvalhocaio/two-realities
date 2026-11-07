package game

import (
	"github.com/carvalhocaio/two-realities/internal/physics"
	"github.com/carvalhocaio/two-realities/internal/world"
)

type ShiftOutcome int

const (
	ShiftDone ShiftOutcome = iota
	ShiftBlockedSolid
	ShiftBlockedNoFragments
)

func shiftCost(from world.World) int {
	if from == world.Grey {
		return 1
	}
	return 0
}

func (s *Session) TryShift() ShiftOutcome {
	cost := shiftCost(s.active)
	target := s.active.Other()

	if s.fragments < cost {
		s.fragmentFeedback = BlockedFeedbackTicks
		s.restartHint = true
		s.emit(EventShiftBlocked)
		return ShiftBlockedNoFragments
	}
	if physics.OverlapsSolid(s.player.Box, s.layer(target)) {
		s.solidFeedback = BlockedFeedbackTicks
		s.emit(EventShiftBlocked)
		return ShiftBlockedSolid
	}

	s.fragments -= cost
	s.active = target
	s.shifted, s.shiftHint = true, false
	s.emitShift(target)
	return ShiftDone
}
