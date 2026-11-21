package ebitenapp

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

const (
	pickupLevel     = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#.Y..#\n######\n"
	checkpointLevel = "fragments: 0\n[yellow]\n########\n#P.C..E#\n########\n[grey]\n########\n#......#\n########\n"
)

type fxRun struct {
	t     *testing.T
	clock *rhythm.Manual
	s     *game.Session
	fx    pickupFX
}

func newFXRun(t *testing.T, text string) *fxRun {
	t.Helper()
	lvl, err := level.Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	r := &fxRun{t: t, clock: &rhythm.Manual{}}
	r.s = game.NewSession(r.clock, lvl, lvl)
	r.fx.update(r.s)
	return r
}

func (r *fxRun) step(in input.Intent, n int) {
	for range n {
		r.clock.Advance(time.Second / 60)
		r.s.Update(in)
		r.fx.update(r.s)
	}
}

func TestFragmentPickupBursts(t *testing.T) {
	r := newFXRun(t, pickupLevel)
	r.step(input.Intent{ShiftPressed: true}, 1)
	for range 30 {
		r.step(input.Intent{Right: true}, 1)
		if len(r.fx.bursts) > 0 {
			break
		}
	}
	if !r.s.FragmentCollected(0) || len(r.fx.bursts) != 1 {
		t.Fatalf("collected %v, bursts %d", r.s.FragmentCollected(0), len(r.fx.bursts))
	}
	r.step(input.Intent{}, pickupBurstTicks)
	if len(r.fx.bursts) != 0 {
		t.Fatalf("burst outlived %d ticks", pickupBurstTicks)
	}
}

func TestRestartClearsBurstsWithoutNewOnes(t *testing.T) {
	r := newFXRun(t, pickupLevel)
	r.step(input.Intent{ShiftPressed: true}, 1)
	r.step(input.Intent{Right: true}, 10)
	r.step(input.Intent{RestartPressed: true}, 1)
	if len(r.fx.bursts) != 0 || r.fx.collected[0] {
		t.Fatalf("after restart: bursts %d, collected %v", len(r.fx.bursts), r.fx.collected[0])
	}
	r.step(input.Intent{}, 2)
	if len(r.fx.bursts) != 0 {
		t.Fatal("restart spawned a burst")
	}
}

func TestCheckpointHoist(t *testing.T) {
	r := newFXRun(t, checkpointLevel)
	if r.fx.hoistProgress() != 1 {
		t.Fatalf("idle hoist progress %v, want 1", r.fx.hoistProgress())
	}
	for range 40 {
		r.step(input.Intent{Right: true}, 1)
		if r.fx.hoist > 0 {
			break
		}
	}
	if r.fx.hoist != hoistTicks || r.fx.hoistProgress() != 0 {
		t.Fatalf("touching the checkpoint: hoist %d, progress %v", r.fx.hoist, r.fx.hoistProgress())
	}
	r.step(input.Intent{}, hoistTicks)
	if r.fx.hoistProgress() != 1 {
		t.Fatalf("hoist progress %v after %d ticks, want 1", r.fx.hoistProgress(), hoistTicks)
	}
}

func TestLevelFadeAlpha(t *testing.T) {
	tests := []struct {
		age  int
		want float32
	}{
		{-1, 0},
		{0, 1},
		{levelFadeTicks / 2, 0.5},
		{levelFadeTicks, 0},
	}
	for _, tt := range tests {
		if got := levelFadeAlpha(tt.age); got != tt.want {
			t.Fatalf("age %d: alpha %v, want %v", tt.age, got, tt.want)
		}
	}
}

func TestFlagFlutterAlternates(t *testing.T) {
	if flutterOffset(0) != 0 || flutterOffset(flutterTicks) != 1 || flutterOffset(2*flutterTicks) != 0 {
		t.Fatal("flag tip should alternate every flutterTicks")
	}
}
