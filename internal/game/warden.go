package game

import (
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	WardenHitTicks    = 45
	assistWarning     = 1
	assistSweepFactor = 0.6
)

type WardenPhase struct {
	Shift     level.ShiftMode
	AutoEvery int
	Aim       Aim
	Sweep     float64
	Warning   int
}

var WardenPhases = [level.LampCount]WardenPhase{
	{Shift: level.ShiftFree, Aim: AimSweep, Sweep: 80},
	{Shift: level.ShiftQuantized, Aim: AimSweep, Sweep: 100},
	{Shift: level.ShiftAuto, AutoEvery: 2, Aim: AimTrack, Warning: 2},
}

type wardenState struct {
	lit       int
	start     int
	cycle     int
	fromRight bool
	beam      Beam
	hit       int
}

type WardenView struct {
	Lit  int
	Beam Beam
	Hit  float64
}

func (s *Session) Warden() (WardenView, bool) {
	if s.level.Boss == nil {
		return WardenView{}, false
	}
	return WardenView{
		Lit:  s.warden.lit,
		Beam: s.warden.beam,
		Hit:  float64(s.warden.hit) / WardenHitTicks,
	}, true
}

func (s *Session) wardenActive() bool {
	return s.level.Boss != nil && s.warden.lit < level.LampCount
}

func (s *Session) wardenPhase() WardenPhase { return WardenPhases[s.warden.lit] }

func (s *Session) shiftMode() level.ShiftMode {
	if s.wardenActive() {
		return s.wardenPhase().Shift
	}
	if s.level.Boss != nil {
		return level.ShiftFree
	}
	return s.level.Shift
}

func (s *Session) autoEvery() int {
	if s.wardenActive() {
		return s.wardenPhase().AutoEvery
	}
	return s.level.AutoEvery
}

func (s *Session) warning() int {
	w := s.wardenPhase().Warning
	if s.assist {
		w += assistWarning
	}
	return w
}

func (s *Session) sweepSpeed() float64 {
	v := s.wardenPhase().Sweep
	if s.assist {
		v *= assistSweepFactor
	}
	return v
}

func (s *Session) restBeats() int {
	if s.warden.lit == 0 {
		return 0
	}
	return RestBeats
}

func (s *Session) resetWarden(lit int) {
	beat := s.clock.Beat()
	rail := s.level.Boss.Rail
	railMid := float64((rail.From+rail.To+1)*level.TileSize) / 2
	s.warden = wardenState{
		lit:       lit,
		start:     beat,
		cycle:     -1,
		fromRight: s.player.Box.X+s.player.Box.W/2 < railMid,
		hit:       s.warden.hit,
	}
	s.beats = newBeatState(beat)
	s.closeGates()
}

func (s *Session) closeGates() {
	closed := s.ClosedGates()
	for w := range world.Count {
		s.layers[w] = s.level.Layer(world.World(w)).WithSolid(closed)
	}
}

func (s *Session) ClosedGates() []level.Point {
	if s.level.Boss == nil {
		return nil
	}
	var closed []level.Point
	for _, gate := range s.level.Boss.Gates[min(s.warden.lit, level.GateCount):] {
		closed = append(closed, gate...)
	}
	return closed
}

func (s *Session) updateWarden() {
	s.warden.hit = max(s.warden.hit-1, 0)
	if !s.wardenActive() || s.death != Alive {
		s.warden.beam.Stage = BeamIdle
		return
	}
	s.aimBeam()
	if s.lightLamp() {
		return
	}
	b := s.warden.beam
	if b.Stage == BeamFiring && s.active == world.Grey && s.player.Box.Overlaps(b.Box) {
		s.die(DeathBeam)
	}
}

func (s *Session) aimBeam() {
	rail := s.level.Boss.Rail
	rel := s.clock.Beat() - s.warden.start - s.restBeats()
	b := &s.warden.beam
	b.Progress = 0
	if s.wardenPhase().Aim == AimSweep {
		s.sweep(b, rel, rail)
	} else {
		s.track(b, rel, rail)
	}
	b.Box = beamBox(b.X, rail, s.layer(world.Grey))
}

func (s *Session) sweep(b *Beam, rel int, rail level.Rail) {
	b.Stage = BeamIdle
	if rel >= 0 {
		b.Stage = BeamFiring
		b.X = sweepX(float64(rel)+s.clock.Phase(), s.sweepSpeed(), rail, s.warden.fromRight)
	}
}

func (s *Session) track(b *Beam, rel int, rail level.Rail) {
	warning := s.warning()
	n, stage := beamCycle(rel, s.clock.Phase(), warning)
	if n != s.warden.cycle && n >= 0 {
		s.warden.cycle = n
		b.X = trackX(s.player.Box, rail)
	}
	b.Stage = stage
	if stage == BeamTelegraph {
		b.Progress = telegraphProgress(rel, s.clock.Phase(), warning)
	}
}

func (s *Session) lightLamp() bool {
	lamp := s.level.Boss.Lamps[s.warden.lit]
	if s.active != world.Grey || !s.player.Box.Overlaps(level.TileRect(lamp)) {
		return false
	}
	lit := s.warden.lit + 1
	s.checkpoint = &snapshot{
		index:     -1,
		point:     lamp,
		fragments: s.fragments,
		collected: append([]bool(nil), s.collected...),
		secret:    s.secretTaken,
		lit:       lit,
	}
	s.warden.hit = WardenHitTicks
	s.emit(EventLamp)
	s.resetWarden(lit)
	s.warden.beam.Stage = BeamIdle
	return true
}

func (s *Session) ExitOpen() bool { return s.level.Boss == nil || s.warden.lit == level.LampCount }
