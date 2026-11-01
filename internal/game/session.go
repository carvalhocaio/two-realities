package game

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const BlockedFeedbackTicks = 12

type Session struct {
	level    level.Level
	active   world.World
	player   player.Player
	feedback int
}

func NewSession(lvl level.Level) *Session {
	s := &Session{level: lvl}
	s.Restart()
	return s
}

func (s *Session) Restart() {
	s.active = world.Yellow
	s.player = player.New(s.level.Spawn)
	s.feedback = 0
}

func (s *Session) Update(in input.Intent) {
	if in.RestartPressed {
		s.Restart()
		return
	}
	s.feedback = max(s.feedback-1, 0)
	if in.ShiftPressed {
		s.TryShift()
	}

	s.player.Update(in, s.level.Layer(s.active))

	if s.reachedExit() || s.fellOut() {
		s.Restart()
	}
}

func (s *Session) Level() level.Level   { return s.level }
func (s *Session) Active() world.World  { return s.active }
func (s *Session) PlayerBox() geom.Rect { return s.player.Box }

func (s *Session) Feedback() float64 {
	return float64(s.feedback) / BlockedFeedbackTicks
}

func (s *Session) reachedExit() bool {
	return s.player.Box.Overlaps(level.TileRect(s.level.Exit))
}

func (s *Session) fellOut() bool {
	return s.player.Box.Y > s.level.Layer(s.active).PixelHeight()
}
