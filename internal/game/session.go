package game

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	BlockedFeedbackTicks = 12
)

type Feedback struct {
	Shake        float64
	CounterFlash float64
}

type Session struct {
	levels           []level.Level
	index            int
	level            level.Level
	active           world.World
	player           player.Player
	watchers         []Watcher
	fragments        int
	collected        []bool
	solidFeedback    int
	fragmentFeedback int
	restartHint      bool
	death            DeathCause
	deathTimer       int
	checkpoint       *snapshot
	finished         bool
	epoch            int
}

func NewSession(levels ...level.Level) *Session {
	s := &Session{levels: levels}
	s.enterLevel(0)
	return s
}

func (s *Session) Restart() {
	s.checkpoint = nil
	s.respawn()
}

func (s *Session) respawn() {
	start := s.startState()
	s.active = world.Yellow
	s.player = player.New(start.point)
	s.watchers = make([]Watcher, len(s.level.Watchers))
	for i, p := range s.level.Watchers {
		s.watchers[i] = NewWatcher(p)
	}
	s.fragments = start.fragments
	s.collected = append([]bool(nil), start.collected...)
	s.epoch++
	s.solidFeedback = 0
	s.fragmentFeedback = 0
	s.restartHint = false
	s.death = Alive
	s.deathTimer = 0
}

func (s *Session) Update(in input.Intent) {
	if s.finished {
		return
	}
	if s.death != Alive {
		s.tickDeath()
		return
	}
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
	s.activateCheckpoints()
	s.updateWatchers()

	switch {
	case s.touchingSpikes():
		s.die(DeathSpikes)
	case s.fellOut():
		s.die(DeathFall)
	}
	if s.death == Alive && s.reachedExit() {
		s.nextLevel()
	}
}

func (s *Session) Level() level.Level   { return s.level }
func (s *Session) Active() world.World  { return s.active }
func (s *Session) PlayerBox() geom.Rect { return s.player.Box }
func (s *Session) RestartHint() bool    { return s.restartHint }
func (s *Session) Watchers() []Watcher  { return s.watchers }
func (s *Session) Finished() bool       { return s.finished }
func (s *Session) Epoch() int           { return s.epoch }

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
