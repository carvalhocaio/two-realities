package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const airShiftTicks = 14

type secretLedges struct {
	jumpFromX  float64
	ledgeTop   float64
	ledgeEdgeX float64
	greyTop    float64
}

func ledgesOf(sec level.Marker) secretLedges {
	ts := float64(level.TileSize)
	return secretLedges{
		jumpFromX:  float64(sec.X-3)*ts - 2*ts,
		ledgeTop:   float64(sec.Y+3) * ts,
		ledgeEdgeX: float64(sec.X-1)*ts - player.Width - 1,
		greyTop:    float64(sec.Y+1) * ts,
	}
}

func (r *runner) ticksToBeat() float64 {
	return (1 - r.clock.Phase()) * float64(beatTicks)
}

func (r *runner) grabSecret() {
	r.t.Helper()
	l := ledgesOf(*r.s.Level().Secret)
	r.until(r.reachX(l.jumpFromX), right)
	r.step(jumpRight)
	r.until(func() bool { return r.s.PlayerBox().Bottom() == l.ledgeTop && r.x() >= l.ledgeEdgeX }, right)

	switch r.s.Level().Shift {
	case level.ShiftFree:
		r.step(jumpRight)
		r.until(func() bool { return r.s.PlayerBox().Bottom() <= l.greyTop }, right)
		r.step(input.Intent{Right: true, ShiftPressed: true})
	case level.ShiftQuantized:
		r.until(func() bool { return r.ticksToBeat() <= airShiftTicks }, idle)
		r.step(input.Intent{Right: true, JumpPressed: true, ShiftPressed: true})
	case level.ShiftAuto:
		r.until(func() bool { return r.s.AutoCountdown() == 1 && r.ticksToBeat() <= airShiftTicks }, idle)
		r.step(jumpRight)
	}
	r.until(func() bool { return r.s.Active() == world.Grey }, right)
	r.until(r.s.SecretTaken, right)
}

func toLeftWall(r *runner) { r.until(r.backTo(10.5), walkLeft) }

func spikesUntilLastFragment(r *runner) {
	r.until(r.reachX(60), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
	r.until(r.reachX(232), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
}

func TestSecretsCanBeFound(t *testing.T) {
	tests := []struct {
		name    string
		prelude func(r *runner)
	}{
		{"01_first_shift", toLeftWall},
		{"02_fragments", toLeftWall},
		{"03_spikes", spikesUntilLastFragment},
		{"04_watcher", toLeftWall},
		{"05_planning", toLeftWall},
		{"06_finale", toLeftWall},
		{"07_crossfire", toLeftWall},
		{"08_descent", toLeftWall},
		{"09_on_the_beat", toLeftWall},
		{"10_metronome", func(*runner) {}},
		{"11_offbeat", toLeftWall},
		{"12_clockwork", toLeftWall},
	}
	levels, err := level.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	if want := len(levels) - bossLevels(levels); len(tests) != want {
		t.Fatalf("%d secret scripts for %d levels with secrets", len(tests), want)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRunner(t, tt.name)
			if r.s.Level().Secret == nil {
				t.Fatal("level has no secret")
			}
			tt.prelude(r)
			before := r.s.Fragments()
			r.grabSecret()
			if r.s.Fragments() != before+1 {
				t.Fatalf("fragments %d, want %d", r.s.Fragments(), before+1)
			}
		})
	}
}

func TestSecretsNeedTheShift(t *testing.T) {
	levels, err := level.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, lvl := range levels {
		if lvl.Boss != nil {
			continue
		}
		sec := *lvl.Secret
		g := lvl.Layer(sec.Layer)
		other := lvl.Layer(sec.Layer.Other())
		if !g.Solid(sec.X, sec.Y+1) || other.Solid(sec.X, sec.Y+1) {
			t.Errorf("%s: secret should rest on a ledge that exists only in its own layer", lvl.Name)
		}
		for x := sec.X - 3; x <= sec.X+2; x++ {
			for y := sec.Y + 2; y <= sec.Y+3; y++ {
				if g.Solid(x, y) {
					t.Errorf("%s: %v has a stepping stone at %d,%d", lvl.Name, sec.Layer, x, y)
				}
			}
		}
	}
}

func bossLevels(levels []level.Level) int {
	n := 0
	for _, l := range levels {
		if l.Boss != nil {
			n++
		}
	}
	return n
}
