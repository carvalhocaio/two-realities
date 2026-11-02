package game

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/world"
)

var beatTicks = int(rhythm.BeatDuration(rhythm.BPM) / tickDuration)

const (
	beatPocketX  = 203
	beatWallJump = 58
	beatGapJump  = 336
)

func (r *runner) queueShift(w world.World) {
	r.t.Helper()
	r.step(shift)
	if !r.s.ShiftQueued() && r.s.Active() != w {
		r.t.Fatalf("shift not queued at %+v", r.s.PlayerBox())
	}
	r.until(func() bool { return r.s.Active() == w }, idle)
}

func (r *runner) waitFor(w world.World) {
	r.t.Helper()
	r.until(func() bool { return r.s.Active() == w }, idle)
}

func (r *runner) reachPocketBehindWatcher() {
	r.t.Helper()
	r.queueShift(world.Grey)
	r.until(r.hasFragments(1), right)
	r.until(func() bool { w := r.watcher(); return w.Dir > 0 && w.Box.X > 200 && w.Box.X < 250 }, idle)
	r.until(r.reachX(beatWallJump), right)
	r.step(jumpRight)
	r.until(r.reachX(beatPocketX), right)
}

func solveOnTheBeat(r *runner) {
	r.reachPocketBehindWatcher()
	r.queueShift(world.Yellow)
	r.until(func() bool {
		w := r.watcher()
		return w.Dir < 0 && w.Box.Right() < r.x()-level.TileSize
	}, idle)
	r.queueShift(world.Grey)
	r.until(r.reachX(beatGapJump), right)
	r.step(jumpRight)
	r.finish(right)
}

func solveMetronome(r *runner) {
	r.until(r.reachX(158), right)
	r.waitFor(world.Grey)
	r.until(r.reachX(285), right)
	r.waitFor(world.Yellow)
	r.until(r.reachX(415), right)
	r.waitFor(world.Grey)
	r.finish(right)
}

func TestQueuingEarlyWorksOnAnyBeatPhase(t *testing.T) {
	for offset := range beatTicks {
		r := newRunnerAt(t, "09_on_the_beat", tickDuration*time.Duration(offset))
		solveOnTheBeat(r)
	}
}

func TestReactingOnSightIsTooLateOnSomePhases(t *testing.T) {
	caught := 0
	for offset := range beatTicks {
		r := newRunnerAt(t, "09_on_the_beat", tickDuration*time.Duration(offset))
		r.reachPocketBehindWatcher()
		for range maxTicks {
			if r.watcher().Alert > 0 {
				break
			}
			r.step(idle)
		}
		r.update(shift)
		for r.s.Death() == Alive && r.s.Active() != world.Yellow {
			r.update(idle)
		}
		if r.s.Death() == DeathCaught {
			caught++
		}
	}
	t.Logf("caught on %d of %d phases", caught, beatTicks)
	if caught == 0 {
		t.Fatal("reacting only after being seen should get caught when the next beat is far")
	}
}

func TestMetronomeMistakes(t *testing.T) {
	tests := []struct {
		name string
		play func(r *runner)
		want DeathCause
	}{
		{"running without waiting falls", func(r *runner) {}, DeathFall},
		{"standing under the grey block when it flips is crushed", func(r *runner) {
			r.until(r.reachX(158), right)
			r.waitFor(world.Grey)
			r.until(r.reachX(285), right)
			r.waitFor(world.Yellow)
			r.until(r.reachX(350), right)
		}, DeathCrushed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRunner(t, "10_metronome")
			tt.play(r)
			in := right
			if r.ticks > 0 {
				in = idle
			}
			for range maxTicks {
				if r.s.Death() != Alive {
					break
				}
				r.update(in)
			}
			if r.s.Death() != tt.want {
				t.Fatalf("death = %v, want %v", r.s.Death(), tt.want)
			}
		})
	}
}
