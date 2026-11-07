package geom

import "testing"

func TestRect(t *testing.T) {
	a := Rect{X: 0, Y: 0, W: 10, H: 10}
	if a.Right() != 10 || a.Bottom() != 10 {
		t.Fatalf("right %v bottom %v", a.Right(), a.Bottom())
	}
	tests := []struct {
		name string
		b    Rect
		want bool
	}{
		{"inside", Rect{X: 2, Y: 2, W: 2, H: 2}, true},
		{"partial", Rect{X: 8, Y: 8, W: 5, H: 5}, true},
		{"touching edge", Rect{X: 10, Y: 0, W: 5, H: 5}, false},
		{"touching bottom", Rect{X: 0, Y: 10, W: 5, H: 5}, false},
		{"left of it", Rect{X: -6, Y: 0, W: 5, H: 5}, false},
		{"above it", Rect{X: 0, Y: -6, W: 5, H: 5}, false},
	}
	for _, tt := range tests {
		if a.Overlaps(tt.b) != tt.want || tt.b.Overlaps(a) != tt.want {
			t.Errorf("%s: overlaps want %v", tt.name, tt.want)
		}
	}
}
