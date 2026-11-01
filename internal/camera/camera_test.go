package camera

import (
	"math"
	"testing"

	"github.com/carvalhocaio/two-realities/internal/geom"
)

var (
	wide  = geom.Rect{W: 1000, H: 500}
	small = geom.Rect{W: 200, H: 100}
)

func box(cx, cy float64) geom.Rect { return geom.Rect{X: cx - 3, Y: cy - 4.5, W: 6, H: 9} }

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestSnap(t *testing.T) {
	tests := []struct {
		name         string
		bounds       geom.Rect
		target       geom.Rect
		wantX, wantY float64
	}{
		{"centers on target", wide, box(500, 250), 340, 165},
		{"clamped at top-left", wide, box(10, 10), 0, 0},
		{"clamped at bottom-right", wide, box(990, 490), 680, 330},
		{"level smaller than view is centered", small, box(10, 10), -60, -35},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(320, 170)
			c.Snap(tt.target, tt.bounds)
			if !near(c.X, tt.wantX) || !near(c.Y, tt.wantY) {
				t.Fatalf("camera = (%v,%v), want (%v,%v)", c.X, c.Y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestFollow(t *testing.T) {
	tests := []struct {
		name         string
		start        geom.Vec
		target       geom.Rect
		wantX, wantY float64
	}{
		{"inside dead zone does not move", geom.Vec{X: 340, Y: 165}, box(520, 260), 340, 165},
		{"dead zone right edge does not move", geom.Vec{X: 340, Y: 165}, box(524, 250), 340, 165},
		{"past right edge eases toward target", geom.Vec{X: 340, Y: 165}, box(544, 250), 340 + 20*Smoothing, 165},
		{"past left edge eases toward target", geom.Vec{X: 340, Y: 165}, box(456, 250), 340 - 20*Smoothing, 165},
		{"past bottom edge eases down", geom.Vec{X: 340, Y: 165}, box(500, 290), 340, 165 + 24*Smoothing},
		{"easing still clamps to bounds", geom.Vec{X: 0, Y: 0}, box(10, 10), 0, 0},
		{"clamps at right bound", geom.Vec{X: 680, Y: 330}, box(999, 499), 680, 330},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Camera{X: tt.start.X, Y: tt.start.Y, ViewW: 320, ViewH: 170}
			c.Follow(tt.target, wide)
			if !near(c.X, tt.wantX) || !near(c.Y, tt.wantY) {
				t.Fatalf("camera = (%v,%v), want (%v,%v)", c.X, c.Y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestFollowConverges(t *testing.T) {
	c := Camera{ViewW: 320, ViewH: 170}
	target := box(700, 250)
	for range 200 {
		c.Follow(target, wide)
	}
	zoneRight := c.X + (c.ViewW+DeadZoneW)/2
	if math.Abs(zoneRight-700) > 0.01 {
		t.Fatalf("target should rest on dead zone edge, edge=%v", zoneRight)
	}
}
