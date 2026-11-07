package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

const threeSpots = "fragments: 0\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#S.S.S.#\n#......#\n########\n"

func TestSecretSpotIsPickedOnLevelEntry(t *testing.T) {
	lvl := mustParse(t, threeSpots)
	s := NewSession(&rhythm.Manual{}, lvl, lvl)
	picks := []int{2, 1}
	calls := 0
	s.UseSecretPicker(func(n int) int {
		if n != 3 {
			t.Fatalf("picker got %d spots, want 3", n)
		}
		calls++
		return picks[(calls-1)%len(picks)]
	})
	spot := func() int { sec, _ := s.Secret(); return sec.X }
	if spot() != 5 {
		t.Fatalf("spot x = %d, want 5", spot())
	}
	s.Update(input.Intent{RestartPressed: true})
	s.die(DeathSpikes)
	for s.Death() != Alive {
		s.Update(idle)
	}
	if spot() != 5 || calls != 1 {
		t.Fatalf("restart and death changed the spot: x=%d calls=%d", spot(), calls)
	}
	s.GoToLevel(1)
	if spot() != 3 || calls != 2 {
		t.Fatalf("entering a level should pick again: x=%d calls=%d", spot(), calls)
	}
}

func TestTakingTheSecretGivesFeedback(t *testing.T) {
	lvl := mustParse(t, "fragments: 0\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n#.S....#\n########\n")
	s := NewSession(&rhythm.Manual{}, lvl, lvl)
	s.Update(input.Intent{ShiftPressed: true})
	for range 20 {
		if s.SecretTaken() {
			break
		}
		s.Update(right)
	}
	if !s.SecretTaken() || s.Feedback().SecretFound != 1 {
		t.Fatalf("taken %v feedback %v, want the full pickup feedback", s.SecretTaken(), s.Feedback().SecretFound)
	}
	for range SecretFeedbackTicks {
		s.Update(idle)
	}
	if s.Feedback().SecretFound != 0 {
		t.Fatalf("feedback %v should fade out", s.Feedback().SecretFound)
	}
}

func TestFallingBeforeEverShiftingShowsTheShiftHint(t *testing.T) {
	r := newRunner(t, "01_first_shift")
	for range maxTicks {
		if r.s.Death() != Alive {
			break
		}
		r.update(jumpOverStep(r.s))
	}
	if r.s.Death() != DeathFall || !r.s.ShiftHint() {
		t.Fatalf("death %v hint %v, want a fall and the shift hint", r.s.Death(), r.s.ShiftHint())
	}
	for r.s.Death() != Alive {
		r.update(idle)
	}
	if !r.s.ShiftHint() {
		t.Fatal("the hint should stay after respawning")
	}
	r.update(input.Intent{ShiftPressed: true})
	if r.s.ShiftHint() {
		t.Fatal("shifting should clear the hint")
	}
}
