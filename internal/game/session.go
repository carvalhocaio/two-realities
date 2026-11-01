package game

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const BlockedFeedbackTicks = 12

type Feedback struct {
	Shake        float64
	CounterFlash float64
}

type Session struct {
	level            level.Level
	active           world.World
	player           player.Player
	fragments        int
	collected        []bool
	solidFeedback    int
	fragmentFeedback int
	restartHint      bool
}

func NewSession(lvl level.Level) *Session {
	s := &Session{level: lvl}
	s.Restart()
	return s
}

func (s *Session) Restart() {
	s.active = world.Yellow
	s.player = player.New(s.level.Spawn)
	s.fragments = s.level.StartFragments
	s.collected = make([]bool, len(s.level.Fragments))
	s.solidFeedback = 0
	s.fragmentFeedback = 0
	s.restartHint = false
}

func (s *Session) Update(in input.Intent) {
	if in.RestartPressed {
		s.Restart()
		return
	}
	s.solidFeedback = max(s.solidFeedback-1, 0)
	s.fragmentFeedback = max(s.fragmentFeedback-1, 0)
	if in.ShiftPressed {
		s.TryShift()
	}

	s.player.Update(in, s.level.Layer(s.active))
	s.collectFragments()

	if s.reachedExit() || s.fellOut() {
		s.Restart()
	}
}

func (s *Session) Level() level.Level   { return s.level }
func (s *Session) Active() world.World  { return s.active }
func (s *Session) PlayerBox() geom.Rect { return s.player.Box }
func (s *Session) RestartHint() bool    { return s.restartHint }

func (s *Session) Feedback() Feedback {
	return Feedback{
		Shake:        float64(s.solidFeedback) / BlockedFeedbackTicks,
		CounterFlash: float64(s.fragmentFeedback) / BlockedFeedbackTicks,
	}
}

func (s *Session) reachedExit() bool {
	return s.player.Box.Overlaps(level.TileRect(s.level.Exit))
}

func (s *Session) fellOut() bool {
	return s.player.Box.Y > s.level.Layer(s.active).PixelHeight()
}
