package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	testLevel        = "01_test"
	jumpForFragmentX = 135
	wallX            = level.TileSize * 22
	maxSolutionTicks = 30 * 60
)

type runner struct {
	t      *testing.T
	s      *Session
	shifts int
	ticks  int
}

func newRunner(t *testing.T) *runner {
	t.Helper()
	lvl, err := level.Load(testLevel)
	if err != nil {
		t.Fatal(err)
	}
	return &runner{t: t, s: NewSession(lvl)}
}

func (r *runner) step(in input.Intent) {
	before := r.s.Active()
	r.s.Update(in)
	r.ticks++
	if r.s.Active() != before {
		r.shifts++
	}
}

func (r *runner) walkTo(x float64, in input.Intent) {
	for range maxSolutionTicks {
		if r.s.PlayerBox().X >= x {
			return
		}
		r.step(in)
	}
	r.t.Fatalf("never reached x=%v, stuck at %+v", x, r.s.PlayerBox())
}

func (r *runner) crossToGreyWall(grabFragment bool) {
	r.step(shift)
	for range 10 {
		r.step(input.Intent{Right: true, JumpPressed: true})
	}
	if grabFragment {
		r.walkTo(jumpForFragmentX, right)
		r.step(input.Intent{Right: true, JumpPressed: true})
	}
	for range 200 {
		r.step(right)
	}
	if r.s.Active() != world.Grey || r.s.PlayerBox().X < wallX {
		r.t.Fatalf("expected to be against the grey wall, got %v at %+v", r.s.Active(), r.s.PlayerBox())
	}
}

func TestTestLevelNeedsShifting(t *testing.T) {
	r := newRunner(t)
	for range 600 {
		r.step(right)
		if r.s.PlayerBox().X > level.TileSize*12 {
			t.Fatal("crossed the pit without shifting")
		}
	}
}

func TestTestLevelNeedsFragment(t *testing.T) {
	r := newRunner(t)
	r.crossToGreyWall(false)
	if got := r.s.TryShift(); got != ShiftBlockedNoFragments {
		t.Fatalf("TryShift() = %v, want blocked by fragments", got)
	}
}

func TestTestLevelSolution(t *testing.T) {
	r := newRunner(t)
	exitX := level.TileRect(r.s.Level().Exit).X

	r.crossToGreyWall(true)
	if r.s.Fragments() != 1 {
		t.Fatalf("fragments = %d, want 1 before returning", r.s.Fragments())
	}
	r.step(shift)
	if r.shifts != 2 || r.s.Active() != world.Yellow || r.s.Fragments() != 0 {
		t.Fatalf("shifts=%d world=%v fragments=%d", r.shifts, r.s.Active(), r.s.Fragments())
	}

	for range maxSolutionTicks {
		lastX := r.s.PlayerBox().X
		r.step(right)
		if r.s.PlayerBox().X < lastX {
			if lastX+player.Width+player.RunSpeed <= exitX {
				t.Fatalf("restarted at x=%v before reaching exit", lastX)
			}
			if r.ticks > maxSolutionTicks {
				t.Fatalf("solution took %d ticks, want under %d", r.ticks, maxSolutionTicks)
			}
			return
		}
	}
	t.Fatalf("never reached the exit, stuck at %+v", r.s.PlayerBox())
}
