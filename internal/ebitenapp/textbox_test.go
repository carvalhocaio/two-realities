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

func TestBoxUnfolds(t *testing.T) {
	box := geom.Rect{X: 8, Y: 20, W: 304, H: 30}
	if boxOpening(0, true) != 1 || typedRunes(0, true) != 0 || typedRunes(9, true) != 9 {
		t.Fatal("reduced effects should open at once and type from age 0")
	}
	if boxOpening(0, false) != 0 || boxOpening(boxOpenTicks/2, false) != 0.5 || boxOpening(boxOpenTicks+5, false) != 1 {
		t.Fatal("box should open linearly over boxOpenTicks")
	}
	if typedRunes(boxOpenTicks, false) != 0 || typedRunes(boxOpenTicks+4, false) != 4 {
		t.Fatal("typing should start once the box is open")
	}
	half := unfold(box, 0.5)
	if half.H != 15 || half.Y != 28 || half.X != box.X || half.W != box.W {
		t.Fatalf("half-open panel %+v from %+v", half, box)
	}
	if unfold(box, 0).H < 2 || unfold(box, 1) != box {
		t.Fatalf("closed %+v / open %+v", unfold(box, 0), unfold(box, 1))
	}
}
