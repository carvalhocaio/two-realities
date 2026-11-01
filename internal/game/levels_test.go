package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const maxTicks = 60 * 60

var jumpRight = input.Intent{Right: true, JumpPressed: true}

type runner struct {
	t     *testing.T
	s     *Session
	start int
}

func newRunner(t *testing.T, name string) *runner {
	t.Helper()
	lvl, err := level.Load(name)
	if err != nil {
		t.Fatal(err)
	}
	return &runner{t: t, s: NewSession(lvl, lvl)}
}

func (r *runner) until(cond func() bool, in input.Intent) {
	r.t.Helper()
	for range maxTicks {
		if cond() {
			return
		}
		r.step(in)
	}
	r.t.Fatalf("condition never met; player at %+v", r.s.PlayerBox())
}

func (r *runner) step(in input.Intent) {
	r.t.Helper()
	r.s.Update(in)
	if r.s.Death() != Alive {
		r.t.Fatalf("died (%v) at %+v", r.s.Death(), r.s.PlayerBox())
	}
}

func (r *runner) shiftTo(w world.World) {
	r.t.Helper()
	r.step(shift)
	if r.s.Active() != w {
		r.t.Fatalf("shift to %v failed", w)
	}
}

func (r *runner) playerX() float64 { return r.s.PlayerBox().X }

func (r *runner) finished() bool { return r.s.LevelIndex() == 1 }

func runUntilDeath(s *Session, inputs func(*Session) input.Intent) DeathCause {
	for range maxTicks {
		s.Update(inputs(s))
		if s.Death() != Alive {
			return s.Death()
		}
		if s.LevelIndex() != 0 {
			return Alive
		}
	}
	return Alive
}

func loadPair(t *testing.T, name string) *Session {
	t.Helper()
	lvl, err := level.Load(name)
	if err != nil {
		t.Fatal(err)
	}
	return NewSession(lvl, lvl)
}

func TestSpikesLevelSolution(t *testing.T) {
	r := newRunner(t, "01_spikes")
	r.shiftTo(world.Grey)
	r.until(func() bool { return r.s.Fragments() == 1 }, right)
	r.shiftTo(world.Yellow)
	r.until(r.finished, right)
}

func TestSpikesLevelNeedsBothWorlds(t *testing.T) {
	tests := []struct {
		name   string
		inputs func(*Session) input.Intent
	}{
		{"stay in yellow", func(*Session) input.Intent { return right }},
		{"stay in grey", func(s *Session) input.Intent {
			if s.Active() == world.Yellow {
				return shift
			}
			return right
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runUntilDeath(loadPair(t, "01_spikes"), tt.inputs); got != DeathSpikes {
				t.Fatalf("death = %v, want spikes", got)
			}
		})
	}
}

const (
	pillarJumpX  = 58
	pocketX      = 155
	gapJumpX     = 236
	afterPocketX = 185
)

func (r *runner) watcher() Watcher { return r.s.Watchers()[0] }

func (r *runner) prepareWatcherRun() {
	r.shiftTo(world.Grey)
	r.until(func() bool { return r.s.Fragments() == 1 }, right)
	r.until(func() bool { return r.watcher().Dir > 0 && r.watcher().Box.X > afterPocketX }, idle)
	r.until(func() bool { return r.playerX() >= pillarJumpX }, right)
	r.step(jumpRight)
}

func TestWatcherLevelSolution(t *testing.T) {
	r := newRunner(t, "02_watcher")
	r.prepareWatcherRun()
	r.until(func() bool { return r.playerX() >= pocketX }, right)
	r.shiftTo(world.Yellow)

	passed := func() bool {
		w := r.watcher()
		return w.Dir < 0 && w.Box.Right() < r.playerX()-level.TileSize
	}
	r.until(passed, idle)
	r.shiftTo(world.Grey)
	r.until(func() bool { return r.playerX() >= gapJumpX }, right)
	r.step(jumpRight)
	r.until(r.finished, right)
}

func TestWatcherLevelCatchesRusher(t *testing.T) {
	r := newRunner(t, "02_watcher")
	r.prepareWatcherRun()
	got := runUntilDeath(r.s, func(s *Session) input.Intent {
		if s.PlayerBox().X >= gapJumpX {
			return jumpRight
		}
		return right
	})
	if got != DeathCaught {
		t.Fatalf("death = %v, want caught", got)
	}
}
