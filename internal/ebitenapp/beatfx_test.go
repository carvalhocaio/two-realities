package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/world"
)

func TestCountdownPop(t *testing.T) {
	tests := []struct {
		name    string
		phase   float64
		reduced bool
		want    float64
	}{
		{"beat onset", 0, false, countdownLift},
		{"mid beat", 0.5, false, 0},
		{"end of beat", 0.99, false, 0},
		{"reduced effects", 0, true, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countdownPop(tt.phase, tt.reduced); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountdownColor(t *testing.T) {
	for w := range world.Count {
		pal := palettes[w]
		if countdownColor(3, pal) != pal.text || countdownColor(2, pal) != pal.text {
			t.Fatalf("world %d: 3 and 2 should use the text color", w)
		}
		if countdownColor(1, pal) != pal.exit {
			t.Fatalf("world %d: 1 should use the exit color", w)
		}
	}
}
