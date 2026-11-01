package player

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
)

type step struct {
	grounded, pressed bool
}

func repeat(s step, n int) []step {
	out := make([]step, n)
	for i := range out {
		out[i] = s
	}
	return out
}

func concat(parts ...[]step) []step {
	var out []step
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

var (
	ground  = step{grounded: true}
	air     = step{}
	press   = step{pressed: true}
	landJmp = step{grounded: true, pressed: true}
)

func TestJumpWindow(t *testing.T) {
	tests := []struct {
		name  string
		steps []step
		want  bool
	}{
		{"press on ground", []step{landJmp}, true},
		{"never pressed", repeat(ground, 10), false},
		{"coyote: press right after leaving", concat([]step{ground}, []step{press}), true},
		{"coyote: press on last valid tick", concat([]step{ground}, repeat(air, CoyoteTicks-2), []step{press}), true},
		{"coyote: window expired", concat([]step{ground}, repeat(air, CoyoteTicks-1), []step{press}), false},
		{"coyote: never grounded", concat(repeat(air, 3), []step{press}), false},
		{"buffer: land right after press", []step{press, ground}, true},
		{"buffer: land on last valid tick", concat([]step{press}, repeat(air, BufferTicks-2), []step{ground}), true},
		{"buffer: window expired", concat([]step{press}, repeat(air, BufferTicks-1), []step{ground}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w jumpWindow
			got := false
			for _, s := range tt.steps {
				got = w.tick(s.grounded, s.pressed)
			}
			if got != tt.want {
				t.Fatalf("last tick jumped = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJumpWindowConsumesJump(t *testing.T) {
	var w jumpWindow
	if !w.tick(true, true) {
		t.Fatal("expected first jump")
	}
	if w.tick(false, true) {
		t.Fatal("double jump in air")
	}
}

func floorLevel(t *testing.T) level.Level {
	t.Helper()
	lvl, err := level.Parse("[yellow]\n......\n.P....\n######\n[grey]\n.....E\n......\n......\n")
	if err != nil {
		t.Fatal(err)
	}
	return lvl
}

func TestPlayerUpdate(t *testing.T) {
	lvl := floorLevel(t)
	grid := lvl.Layer(0)

	tests := []struct {
		name   string
		inputs []input.Intent
		check  func(p Player) bool
	}{
		{"settles on floor", make([]input.Intent, 3), func(p Player) bool {
			return p.Grounded() && p.Box.Bottom() == 20
		}},
		{"moves right", []input.Intent{{Right: true}}, func(p Player) bool {
			return p.Box.X > New(lvl.Spawn).Box.X
		}},
		{"moves left", []input.Intent{{Left: true}}, func(p Player) bool {
			return p.Box.X < New(lvl.Spawn).Box.X
		}},
		{"opposite keys cancel", []input.Intent{{Left: true, Right: true}}, func(p Player) bool {
			return p.Box.X == New(lvl.Spawn).Box.X
		}},
		{"jumps from ground", []input.Intent{{}, {JumpPressed: true}}, func(p Player) bool {
			return p.Vel.Y < 0 && !p.Grounded()
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lvl.Spawn)
			for _, in := range tt.inputs {
				p.Update(in, grid)
			}
			if !tt.check(p) {
				t.Fatalf("unexpected state: %+v", p)
			}
		})
	}
}
