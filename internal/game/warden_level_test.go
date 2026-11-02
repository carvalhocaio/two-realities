package game

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const wardenLevel = "13_warden"

func (r *runner) lit() int { return r.s.warden.lit }

func (r *runner) litAt(n int) func() bool { return func() bool { return r.lit() >= n } }

const (
	planHorizon = 120
	planStride  = 4
)

func (r *runner) lightAhead() []geom.Rect {
	s := r.s
	rail := s.level.Boss.Rail
	grey := s.layer(world.Grey)
	boxes := make([]geom.Rect, planHorizon+1)
	for k := range boxes {
		pos := rhythm.At(r.clock.Elapsed+time.Duration(k+1)*tickDuration, rhythm.BPM)
		rel := pos.Beat() - s.warden.start - s.restBeats()
		if rel < 0 {
			continue
		}
		x := sweepX(float64(rel)+pos.Phase(), s.sweepSpeed(), rail, s.warden.fromRight)
		boxes[k] = beamBox(x, rail, grey)
	}
	return boxes
}

func planSafe(box geom.Rect, dir float64, stop int, light []geom.Rect) bool {
	for k, beam := range light {
		b := box
		b.X += dir * player.RunSpeed * float64(min(k+1, stop))
		if beam.W > 0 && b.Overlaps(beam) {
			return false
		}
	}
	return true
}

func (r *runner) carefully(in input.Intent) input.Intent {
	if r.s.Active() != world.Grey || r.s.warden.beam.Stage == BeamIdle && r.s.restBeats() == 0 && r.lit() >= level.LampCount {
		return in
	}
	dir := 1.0
	back := walkLeft
	if in.Left {
		dir, back = -1, right
	}
	light := r.lightAhead()
	box := r.s.PlayerBox()
	for stop := planHorizon; stop > 0; stop -= planStride {
		if planSafe(box, dir, stop, light) {
			return in
		}
	}
	if planSafe(box, 0, 0, light) {
		return idle
	}
	for stop := planStride; stop <= planHorizon; stop += planStride {
		if planSafe(box, -dir, stop, light) {
			return back
		}
	}
	return in
}

func (r *runner) sneak(cond func() bool, in input.Intent) {
	r.t.Helper()
	for range maxTicks {
		if cond() {
			return
		}
		r.step(r.carefully(in))
	}
	r.t.Fatalf("condition never met; player at %+v", r.s.PlayerBox())
}

func (r *runner) queueAndWait(w world.World) {
	r.t.Helper()
	if r.s.Active() != w {
		r.step(input.Intent{ShiftPressed: true})
	}
	r.until(func() bool { return r.s.Active() == w }, idle)
}

func wardenPhaseOne(r *runner) {
	r.until(r.reachX(62), right)
	r.shiftTo(world.Grey)
	r.sneak(r.litAt(1), right)
}

func (r *runner) lightX(k int) float64 {
	light := r.lightAhead()
	return light[k].X + BeamWidth/2
}

func (r *runner) lightLeaving(minX float64) func() bool {
	return func() bool {
		now, later := r.lightX(0), r.lightX(6)
		return now < later && now > minX
	}
}

func (r *runner) lightReturning(maxX float64) func() bool {
	return func() bool {
		now, later := r.lightX(0), r.lightX(6)
		return now > later && now < maxX
	}
}

func (r *runner) land() {
	r.t.Helper()
	r.until(func() bool { return r.s.player.Grounded() }, idle)
}

func wardenPhaseTwo(r *runner) {
	r.until(r.lightReturning(225), idle)
	r.until(r.reachX(262), right)
	r.step(jumpRight)
	r.until(func() bool { return r.s.player.Grounded() && r.s.PlayerBox().Bottom() == 140 }, right)
	r.until(r.backTo(281), walkLeft)
	r.step(input.Intent{Left: true, JumpPressed: true})
	r.until(func() bool { return r.s.player.Grounded() && r.s.PlayerBox().Bottom() == 120 }, walkLeft)
	r.until(r.backTo(241), walkLeft)
	r.step(input.Intent{Left: true, JumpPressed: true})
	r.until(func() bool { return r.s.player.Grounded() && r.s.PlayerBox().Bottom() == 100 }, walkLeft)
	r.until(r.hasFragments(2), walkLeft)
	r.queueAndWait(world.Yellow)
	r.until(r.backTo(82), walkLeft)
	r.until(r.lightLeaving(110), idle)
	r.queueAndWait(world.Grey)
	r.sneak(r.litAt(2), walkLeft)
}

func (r *runner) dodge(lo, hi float64) input.Intent {
	b := r.s.warden.beam
	box := r.s.PlayerBox()
	if b.Stage == BeamIdle || box.Right() <= b.Box.X || box.X >= b.Box.Right() {
		return idle
	}
	if box.X+box.W/2 < b.X && box.X-player.RunSpeed >= lo || box.X+player.RunSpeed > hi {
		return walkLeft
	}
	return right
}

func (r *runner) holdUntil(cond func() bool, lo, hi float64) {
	r.t.Helper()
	for range maxTicks {
		if cond() {
			return
		}
		r.step(r.dodge(lo, hi))
	}
	r.t.Fatalf("condition never met; player at %+v", r.s.PlayerBox())
}

func (r *runner) nextFlip(w world.World, lo, hi float64) {
	r.t.Helper()
	r.holdUntil(func() bool { return r.s.Active() != w }, lo, hi)
	r.holdUntil(func() bool { return r.s.Active() == w }, lo, hi)
}

func (r *runner) pathClear(dir float64) bool {
	b := r.s.warden.beam
	box := r.s.PlayerBox()
	behind := dir > 0 && b.Box.Right() <= box.X || dir < 0 && b.Box.X >= box.Right()
	return behind || b.Stage == BeamIdle || b.Stage == BeamTelegraph && b.Progress < 0.5
}

func (r *runner) jumpTo(in input.Intent, top float64) {
	r.t.Helper()
	dir := 1.0
	if in.Left {
		dir = -1
	}
	r.until(func() bool { return r.pathClear(dir) }, idle)
	jump := in
	jump.JumpPressed = true
	r.step(jump)
	r.until(func() bool { return r.s.player.Grounded() && r.s.PlayerBox().Bottom() == top }, in)
}

func wardenPhaseThree(r *runner) {
	r.until(r.backTo(36), walkLeft)
	r.jumpTo(walkLeft, 80)
	r.until(r.reachX(22), right)
	r.nextFlip(world.Grey, 10, 24)
	r.jumpTo(right, 60)
	r.until(r.reachX(53), right)
	r.jumpTo(right, 40)
	r.until(r.reachX(78), right)
	r.nextFlip(world.Yellow, 70, 84)
	r.until(r.reachX(146), right)
	r.nextFlip(world.Yellow, 140, 154)
	r.until(r.reachX(221), right)
	r.holdUntil(r.litAt(3), 214, 226)
	r.finish(right)
}

func solveWarden(r *runner) {
	wardenPhaseOne(r)
	wardenPhaseTwo(r)
	wardenPhaseThree(r)
}

func TestWardenSolutionWorksOnEveryBeatPhase(t *testing.T) {
	for offset := range beatTicks {
		r := newRunnerAt(t, wardenLevel, tickDuration*time.Duration(offset))
		solveWarden(r)
	}
}

func TestWardenNaiveAttemptsFail(t *testing.T) {
	tests := []struct {
		name  string
		setup func(r *runner)
		in    input.Intent
	}{
		{"rushing the bridge meets the light", func(r *runner) {
			r.until(r.reachX(62), right)
			r.shiftTo(world.Grey)
		}, right},
		{"standing still in the last phase gets tracked", func(r *runner) {
			wardenPhaseOne(r)
			wardenPhaseTwo(r)
			r.until(r.backTo(36), walkLeft)
			r.jumpTo(walkLeft, 80)
		}, idle},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for offset := range beatTicks {
				r := newRunnerAt(t, wardenLevel, tickDuration*time.Duration(offset))
				tt.setup(r)
				for range maxTicks {
					if r.s.Death() != Alive {
						break
					}
					r.update(tt.in)
				}
				if got := r.s.Death(); got != DeathBeam {
					t.Fatalf("offset %d: death = %v, want beam", offset, got)
				}
			}
		})
	}
}
