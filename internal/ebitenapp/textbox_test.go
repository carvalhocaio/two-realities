package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/geom"
)

func TestTextboxAvoidsThePlayer(t *testing.T) {
	const height = 13
	low := geom.Rect{X: 40, Y: 150, W: 6, H: 9}
	high := geom.Rect{X: 40, Y: 30, W: 6, H: 9}
	tests := []struct {
		name    string
		player  []geom.Rect
		wantTop bool
	}{
		{"player on the floor gets the box on top", []geom.Rect{low}, true},
		{"player high up gets the box at the bottom", []geom.Rect{high}, false},
		{"box stays put while the player moves around", []geom.Rect{low, {X: 40, Y: 80, W: 6, H: 9}}, true},
		{"box moves away when the player walks under it", []geom.Rect{low, high}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tb textbox
			var box geom.Rect
			for _, p := range tt.player {
				box = tb.place(p, height)
			}
			if tb.top != tt.wantTop || box.Overlaps(tt.player[len(tt.player)-1]) {
				t.Fatalf("top = %v box %+v, want top = %v clear of the player", tb.top, box, tt.wantTop)
			}
		})
	}
}
