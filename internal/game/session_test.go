package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	openLevel    = "[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#....#\n######\n"
	blockedLevel = "[yellow]\n#...E#\n#P...#\n######\n[grey]\n#....#\n##...#\n######\n"
	exitLevel    = "[yellow]\n#....#\n#PE..#\n######\n[grey]\n#....#\n#....#\n######\n"
	pitLevel     = "[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#....#\n#....#\n"
)

var (
	idle    = input.Intent{}
	shift   = input.Intent{ShiftPressed: true}
	right   = input.Intent{Right: true}
	restart = input.Intent{RestartPressed: true}
)

func ticks(in input.Intent, n int) []input.Intent {
	out := make([]input.Intent, n)
	for i := range out {
		out[i] = in
	}
	return out
}

func then(parts ...[]input.Intent) []input.Intent {
	var out []input.Intent
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func TestSession(t *testing.T) {
	tests := []struct {
		name         string
		level        string
		inputs       []input.Intent
		wantWorld    world.World
		wantFeedback float64
		wantAtSpawn  bool
	}{
		{"starts in yellow", openLevel, nil, world.Yellow, 0, true},
		{"shift allowed", openLevel, []input.Intent{shift}, world.Grey, 0, true},
		{"shift twice returns", openLevel, []input.Intent{shift, shift}, world.Yellow, 0, true},
		{"shift blocked by solid tile", blockedLevel, []input.Intent{shift}, world.Yellow, 1, true},
		{"blocked feedback fades", blockedLevel, then([]input.Intent{shift}, ticks(idle, BlockedFeedbackTicks/2)), world.Yellow, 0.5, true},
		{"blocked feedback ends", blockedLevel, then([]input.Intent{shift}, ticks(idle, BlockedFeedbackTicks)), world.Yellow, 0, true},
		{"moving away leaves spawn", openLevel, ticks(right, 3), world.Yellow, 0, false},
		{"restart key resets", openLevel, then([]input.Intent{shift}, ticks(right, 3), []input.Intent{restart}), world.Yellow, 0, true},
		{"reaching exit restarts", exitLevel, ticks(right, 2), world.Yellow, 0, true},
		{"falling out restarts", pitLevel, then([]input.Intent{shift}, ticks(idle, 40)), world.Yellow, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl, err := level.Parse(tt.level)
			if err != nil {
				t.Fatal(err)
			}
			s := NewSession(lvl)
			for _, in := range tt.inputs {
				s.Update(in)
			}
			if s.Active() != tt.wantWorld {
				t.Errorf("active = %v, want %v", s.Active(), tt.wantWorld)
			}
			if s.Feedback() != tt.wantFeedback {
				t.Errorf("feedback = %v, want %v", s.Feedback(), tt.wantFeedback)
			}
			atSpawn := s.PlayerBox() == player.New(lvl.Spawn).Box
			if atSpawn != tt.wantAtSpawn {
				t.Errorf("at spawn = %v, want %v (box %+v)", atSpawn, tt.wantAtSpawn, s.PlayerBox())
			}
		})
	}
}

func TestTestLevelNeedsShifting(t *testing.T) {
	s := NewSession(mustLoad(t, "01_test"))
	for range 600 {
		s.Update(right)
		if s.PlayerBox().X > level.TileSize*12 {
			t.Fatal("crossed the pit without shifting")
		}
	}
}

func TestTestLevelIsSolvableWithTwoShifts(t *testing.T) {
	s := NewSession(mustLoad(t, "01_test"))
	exitX := level.TileRect(s.Level().Exit).X
	shifts := 0
	run := func(n int, in input.Intent) {
		for range n {
			before := s.Active()
			s.Update(in)
			if s.Active() != before {
				shifts++
			}
		}
	}

	run(1, shift)
	run(10, input.Intent{Right: true, JumpPressed: true})
	run(200, right)
	if s.Active() != world.Grey || s.PlayerBox().X < level.TileSize*22 {
		t.Fatalf("expected to be against the grey wall, got %v at %+v", s.Active(), s.PlayerBox())
	}
	run(1, shift)
	if shifts != 2 {
		t.Fatalf("shifts = %d, want 2", shifts)
	}
	for range 300 {
		lastX := s.PlayerBox().X
		run(1, right)
		if s.PlayerBox().X < lastX {
			if lastX+player.Width+player.RunSpeed <= exitX {
				t.Fatalf("restarted at x=%v before reaching exit", lastX)
			}
			return
		}
	}
	t.Fatalf("never reached the exit, stuck at %+v", s.PlayerBox())
}

func mustLoad(t *testing.T, name string) level.Level {
	t.Helper()
	lvl, err := level.Load(name)
	if err != nil {
		t.Fatal(err)
	}
	return lvl
}
