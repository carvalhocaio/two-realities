package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/world"
)

func TestWatcherCorridorRoofsCannotBeWalked(t *testing.T) {
	tests := []struct {
		level   string
		hopX    float64
		roofTop float64
	}{
		{"04_watcher", 84, 130},
		{"05_planning", 84, 130},
		{"06_finale", 84, 120},
		{"07_crossfire", 84, 130},
		{"09_on_the_beat", 84, 130},
	}
	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			r := newRunner(t, tt.level)
			r.update(input.Intent{ShiftPressed: true})
			for r.s.Active() != world.Grey {
				r.update(idle)
			}
			r.until(r.reachX(56), right)
			r.update(jumpRight)
			for !r.s.PlayerGrounded() || r.x() < tt.hopX {
				r.update(right)
			}
			r.update(jumpRight)
			for range 120 {
				r.update(right)
				b := r.s.PlayerBox()
				if r.s.Death() != Alive {
					return
				}
				if r.s.PlayerGrounded() && b.Bottom() <= tt.roofTop && b.X > 110 {
					t.Fatalf("walking on the corridor roof at %+v", b)
				}
			}
		})
	}
}
