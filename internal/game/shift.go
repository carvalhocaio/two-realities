package game

import "github.com/carvalhocaio/two-realities/internal/physics"

func (s *Session) TryShift() bool {
	target := s.active.Other()
	if physics.OverlapsSolid(s.player.Box, s.level.Layer(target)) {
		s.feedback = BlockedFeedbackTicks
		return false
	}
	s.active = target
	return true
}
