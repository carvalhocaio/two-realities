package game

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

func (r *runner) w(i int) Watcher { return r.s.Watchers()[i] }

func solveCrossfire(r *runner) {
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.until(r.reachX(58), right)
	r.step(jumpRight)
	r.until(r.reachX(175), right)
	r.until(func() bool { return r.w(0).Dir < 0 }, idle)
	r.shiftTo(world.Yellow)
	r.until(func() bool { return r.w(0).Dir < 0 && r.w(0).Box.Right() < r.x()-level.TileSize }, idle)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.until(r.reachX(293), right)
	r.shiftTo(world.Yellow)
	r.until(r.reachX(311), right)
	r.until(func() bool { return r.w(1).Dir > 0 && r.w(1).Box.X > r.x()+60 }, idle)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.until(r.reachX(395), right)
	r.until(func() bool { return r.w(1).Dir < 0 }, idle)
	r.shiftTo(world.Yellow)
	r.until(func() bool { return r.w(1).Dir < 0 && r.w(1).Box.Right() < r.x()-level.TileSize }, idle)
	r.shiftTo(world.Grey)
	r.until(r.reachX(452), right)
	r.step(jumpRight)
	r.finish(right)
}

func solveDescent(r *runner) {
	r.until(func() bool { return r.s.PlayerBox().Y > 190 && r.s.player.Grounded() }, right)
	r.until(r.backTo(121), walkLeft)
	r.until(func() bool { return r.w(0).Dir > 0 && r.w(0).Box.X > 100 }, idle)
	r.step(walkLeft)
	r.step(walkLeft)
	r.until(func() bool { return r.s.PlayerBox().Y > 215 }, walkLeft)
	r.shiftTo(world.Grey)
	r.until(func() bool { return r.s.player.Grounded() }, idle)
	for r.w(0).Dir > 0 {
		if r.w(0).Box.X-r.s.PlayerBox().Right() > 16 {
			r.step(right)
		} else {
			r.step(idle)
		}
	}
	r.step(jumpRight)
	r.until(func() bool { return r.s.PlayerBox().Y > 280 }, right)
	r.shiftTo(world.Yellow)
	r.until(func() bool { return r.s.player.Grounded() }, idle)
	r.until(r.reachX(231), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
	r.finish(right)
}

func solveOffbeat(r *runner) {
	stops := []struct {
		x float64
		w world.World
	}{{82, world.Grey}, {162, world.Yellow}, {242, world.Grey}, {322, world.Yellow}}
	for _, st := range stops {
		r.until(r.reachX(st.x), right)
		r.queueAndWait(st.w)
	}
	r.until(r.reachX(402), right)
	r.until(func() bool { return r.w(0).Dir > 0 && r.w(0).Box.X > 480 }, idle)
	r.queueAndWait(world.Grey)
	r.until(r.reachX(412), right)
	r.step(jumpRight)
	r.until(r.hasFragments(1), right)
	r.until(r.reachX(535), right)
	r.until(func() bool { return r.w(0).Box.X > 555 }, idle)
	r.queueAndWait(world.Yellow)
	r.until(func() bool { return r.w(0).Dir < 0 && r.w(0).Box.Right() < r.x()-level.TileSize }, idle)
	r.queueAndWait(world.Grey)
	r.until(r.reachX(572), right)
	r.step(jumpRight)
	r.finish(right)
}

func (r *runner) freshWindow(w world.World) {
	r.t.Helper()
	r.until(func() bool { return r.s.Active() != w }, idle)
	r.until(func() bool { return r.s.Active() == w }, idle)
}

func (r *runner) hop(x float64) {
	r.t.Helper()
	r.step(jumpRight)
	r.until(r.reachX(x), right)
	r.until(func() bool { return r.s.player.Grounded() }, idle)
}

func predictWatcher(x float64, dir int, ticks int, lo, hi, speed float64) (float64, int) {
	for range ticks {
		next := x + float64(dir)*speed
		if next < lo || next > hi {
			dir = -dir
			next = x + float64(dir)*speed
		}
		x = next
	}
	return x, dir
}

func solveClockwork(r *runner) {
	r.until(r.reachX(82), right)
	r.freshWindow(world.Yellow)
	r.until(r.reachX(152), right)
	r.freshWindow(world.Grey)
	r.until(r.reachX(222), right)
	r.hop(241)
	const lo, hi = 250, 392
	for {
		r.freshWindow(world.Yellow)
		w := r.w(0)
		x72, d72 := predictWatcher(w.Box.X, w.Dir, 72, lo, hi, r.s.watcherSpeed())
		_, d144 := predictWatcher(w.Box.X, w.Dir, 144, lo, hi, r.s.watcherSpeed())
		if d72 > 0 && d144 > 0 && x72 > 315 {
			break
		}
	}
	r.until(r.reachX(258), right)
	r.until(func() bool { return r.s.Active() == world.Grey }, idle)
	r.until(r.reachX(301), right)
	r.until(func() bool { return r.s.Active() == world.Yellow }, idle)
	r.until(r.reachX(375), right)
	r.hop(401)
	r.freshWindow(world.Yellow)
	r.until(r.reachX(462), right)
	r.until(func() bool { return r.s.Active() == world.Grey }, idle)
	r.until(r.reachX(482), right)
	r.step(jumpRight)
	r.finish(right)
}

func TestRhythmMasteryWorksOnEveryBeatPhase(t *testing.T) {
	tests := []struct {
		name  string
		solve func(r *runner)
	}{
		{"11_offbeat", solveOffbeat},
		{"12_clockwork", solveClockwork},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for offset := range beatTicks {
				r := newRunnerAt(t, tt.name, tickDuration*time.Duration(offset))
				tt.solve(r)
			}
		})
	}
}

func TestMasteryNaiveAttemptsFail(t *testing.T) {
	tests := []struct {
		name  string
		level string
		setup func(r *runner)
		in    input.Intent
		want  DeathCause
	}{
		{"rushing the first watcher gets caught", "07_crossfire", func(r *runner) {
			r.shiftTo(world.Grey)
			r.until(r.reachX(58), right)
			r.step(jumpRight)
		}, right, DeathCaught},
		{"dropping in yellow lands on spikes", "08_descent", func(r *runner) {
			r.until(func() bool { return r.s.PlayerBox().Y > 190 && r.s.player.Grounded() }, right)
		}, walkLeft, DeathSpikes},
		{"walking against the pulse meets the spikes", "11_offbeat", func(*runner) {}, right, DeathSpikes},
		{"running against the clock falls", "12_clockwork", func(*runner) {}, right, DeathFall},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRunner(t, tt.level)
			tt.setup(r)
			for range maxTicks {
				if r.s.Death() != Alive {
					break
				}
				r.update(tt.in)
			}
			if got := r.s.Death(); got != tt.want {
				t.Fatalf("death = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrossfireHasExactFragments(t *testing.T) {
	lvl, err := level.Load("07_crossfire")
	if err != nil {
		t.Fatal(err)
	}
	const hidingPlaces = 3
	if got := lvl.StartFragments + len(lvl.Fragments); got != hidingPlaces {
		t.Fatalf("fragments available = %d, want exactly %d", got, hidingPlaces)
	}
}
