package game

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	BlockedFeedbackTicks = 12
	SecretFeedbackTicks  = 90
	landSoundSpeed       = 2
)

type Feedback struct {
	Shake        float64
	CounterFlash float64
	SecretFound  float64
}

type Session struct {
	clock            rhythm.Clock
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
	secretFeedback   int
	restartHint      bool
	death            DeathCause
	deathTimer       int
	checkpoint       *snapshot
	finished         bool
	epoch            int
	beats            beatState
	signals          transmissionState
	secretTaken      bool
	secretSpot       int
	shifted          bool
	shiftHint        bool
	picker           func(n int) int
	timer            runTimer
	result           LevelResult
	resultTimer      int
	records          *Records
	assist           bool
	layers           [world.Count]level.Grid
	warden           wardenState
	pressure         int
	events           []Event
}

func NewSession(clock rhythm.Clock, levels ...level.Level) *Session {
	s := &Session{clock: clock, levels: levels, records: NewRecords(&save.Memory{}), picker: randomSpot}
	s.timer.runEligible = true
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
	s.secretTaken = start.secret
	s.signals = newTransmissionState(len(s.level.Transmissions))
	s.epoch++
	s.solidFeedback = 0
	s.fragmentFeedback = 0
	s.secretFeedback = 0
	s.restartHint = false
	s.death = Alive
	s.deathTimer = 0
	s.beats = newBeatState(s.clock.Beat())
	s.warden = wardenState{lit: start.lit}
	if s.level.Boss != nil {
		s.resetWarden(start.lit)
		return
	}
	s.closeGates()
}

func (s *Session) Update(in input.Intent) {
	s.events = s.events[:0]
	if s.finished {
		return
	}
	s.tickTimer(in)
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
	s.secretFeedback = max(s.secretFeedback-1, 0)
	if in.ShiftPressed {
		s.requestShift()
	}
	s.onBeats()
	if s.death != Alive {
		return
	}

	s.movePlayer(in)
	s.collectFragments()
	s.collectSecret()
	s.activateCheckpoints()
	s.updateTransmissions()
	s.updateWatchers()
	s.updateWarden()

	switch {
	case s.touchingSpikes():
		s.die(DeathSpikes)
	case s.fellOut():
		s.die(DeathFall)
	}
	if s.death == Alive && s.ExitOpen() && s.reachedExit() {
		s.emit(EventExit)
		s.completeLevel()
		s.nextLevel()
	}
}

func (s *Session) Level() level.Level    { return s.level }
func (s *Session) Active() world.World   { return s.active }
func (s *Session) PlayerBox() geom.Rect  { return s.player.Box }
func (s *Session) PlayerGrounded() bool  { return s.player.Grounded() }
func (s *Session) RestartHint() bool     { return s.restartHint }
func (s *Session) ShiftHint() bool       { return s.shiftHint }
func (s *Session) Watchers() []Watcher   { return s.watchers }
func (s *Session) Finished() bool        { return s.finished }
func (s *Session) Epoch() int            { return s.epoch }
func (s *Session) Levels() []level.Level { return s.levels }
func (s *Session) Records() *Records     { return s.records }

func (s *Session) UseRecords(r *Records) { s.records = r }

func (s *Session) Feedback() Feedback {
	return Feedback{
		Shake:        float64(s.solidFeedback) / BlockedFeedbackTicks,
		CounterFlash: float64(s.fragmentFeedback) / BlockedFeedbackTicks,
		SecretFound:  float64(s.secretFeedback) / SecretFeedbackTicks,
	}
}

func (s *Session) layer(w world.World) level.Grid { return s.layers[w] }

func (s *Session) reachedExit() bool {
	return s.player.Box.Overlaps(level.TileRect(s.level.Exit))
}

func (s *Session) fellOut() bool {
	return s.player.Box.Y > s.layer(s.active).PixelHeight()
}

func (s *Session) movePlayer(in input.Intent) {
	wasGrounded, fall := s.player.Grounded(), s.player.Vel.Y
	s.player.Update(in, s.layer(s.active))
	switch {
	case s.player.Jumped():
		s.emit(EventJump)
	case !wasGrounded && s.player.Grounded() && fall > landSoundSpeed:
		s.emit(EventLand)
	}
}
