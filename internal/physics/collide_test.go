package physics

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

func gridFrom(rows ...string) level.Grid {
	solid := make([]bool, 0, len(rows)*len(rows[0]))
	for _, row := range rows {
		for _, ch := range row {
			solid = append(solid, ch == '#')
		}
	}
	return level.NewGrid(len(rows[0]), len(rows), solid)
}

var room = gridFrom(
	"......",
	"......",
	"...#..",
	"......",
	"######",
)

func TestMoveAndCollide(t *testing.T) {
	tests := []struct {
		name string
		box  geom.Rect
		vel  geom.Vec
		want Result
	}{
		{
			name: "free move",
			box:  geom.Rect{X: 5, Y: 5, W: 6, H: 9},
			vel:  geom.Vec{X: 1, Y: 1},
			want: Result{Box: geom.Rect{X: 6, Y: 6, W: 6, H: 9}},
		},
		{
			name: "lands on floor",
			box:  geom.Rect{X: 5, Y: 30, W: 6, H: 9},
			vel:  geom.Vec{Y: 3},
			want: Result{Box: geom.Rect{X: 5, Y: 31, W: 6, H: 9}, OnGround: true},
		},
		{
			name: "resting on floor stays grounded",
			box:  geom.Rect{X: 5, Y: 31, W: 6, H: 9},
			vel:  geom.Vec{Y: 0.25},
			want: Result{Box: geom.Rect{X: 5, Y: 31, W: 6, H: 9}, OnGround: true},
		},
		{
			name: "stopped by wall moving right",
			box:  geom.Rect{X: 23, Y: 20, W: 6, H: 9},
			vel:  geom.Vec{X: 2},
			want: Result{Box: geom.Rect{X: 24, Y: 20, W: 6, H: 9}, HitWall: true},
		},
		{
			name: "stopped by wall moving left",
			box:  geom.Rect{X: 41, Y: 20, W: 6, H: 9},
			vel:  geom.Vec{X: -2},
			want: Result{Box: geom.Rect{X: 40, Y: 20, W: 6, H: 9}, HitWall: true},
		},
		{
			name: "bumps ceiling",
			box:  geom.Rect{X: 32, Y: 31, W: 6, H: 9},
			vel:  geom.Vec{Y: -3},
			want: Result{Box: geom.Rect{X: 32, Y: 30, W: 6, H: 9}, HitCeiling: true},
		},
		{
			name: "diagonal move lands on block top",
			box:  geom.Rect{X: 22, Y: 10, W: 6, H: 9},
			vel:  geom.Vec{X: 3, Y: 3},
			want: Result{Box: geom.Rect{X: 25, Y: 11, W: 6, H: 9}, OnGround: true},
		},
		{
			name: "diagonal move into wall keeps falling",
			box:  geom.Rect{X: 23, Y: 15, W: 6, H: 9},
			vel:  geom.Vec{X: 3, Y: 2},
			want: Result{Box: geom.Rect{X: 24, Y: 17, W: 6, H: 9}, HitWall: true},
		},
		{
			name: "left bound is solid",
			box:  geom.Rect{X: 1, Y: 5, W: 6, H: 9},
			vel:  geom.Vec{X: -3},
			want: Result{Box: geom.Rect{X: 0, Y: 5, W: 6, H: 9}, HitWall: true},
		},
		{
			name: "top bound is solid",
			box:  geom.Rect{X: 5, Y: 1, W: 6, H: 9},
			vel:  geom.Vec{Y: -3},
			want: Result{Box: geom.Rect{X: 5, Y: 0, W: 6, H: 9}, HitCeiling: true},
		},
		{
			name: "bottom bound is open",
			box:  geom.Rect{X: 5, Y: 60, W: 6, H: 9},
			vel:  geom.Vec{Y: 5},
			want: Result{Box: geom.Rect{X: 5, Y: 65, W: 6, H: 9}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MoveAndCollide(tt.box, tt.vel, room); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestOverlapsSolid(t *testing.T) {
	tests := []struct {
		name string
		box  geom.Rect
		want bool
	}{
		{"empty air", geom.Rect{X: 5, Y: 5, W: 6, H: 9}, false},
		{"touching wall edge", geom.Rect{X: 24, Y: 20, W: 6, H: 9}, false},
		{"inside wall", geom.Rect{X: 27, Y: 20, W: 6, H: 9}, true},
		{"sunk into floor", geom.Rect{X: 5, Y: 32, W: 6, H: 9}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OverlapsSolid(tt.box, room); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
