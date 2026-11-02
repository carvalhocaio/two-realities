package game

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

func TestWatcherPatrol(t *testing.T) {
	walled := gridFrom(
		"#.....#",
		"#######",
	)
	ledged := gridFrom(
		"........",
		"..####..",
	)
	pocket := gridFrom(
		"#.#",
		"###",
	)
	tests := []struct {
		name    string
		grid    level.Grid
		start   level.Point
		dir     int
		ticks   int
		wantX   float64
		wantDir int
	}{
		{"walks forward", walled, level.Point{X: 4}, 1, 10, 46, 1},
		{"turns at wall", walled, level.Point{X: 4}, 1, 30, 48, -1},
		{"turns at right ledge", ledged, level.Point{X: 3}, 1, 50, 48, -1},
		{"turns at left ledge", ledged, level.Point{X: 3}, -1, 30, 24, 1},
		{"bounces inside narrow pocket", pocket, level.Point{X: 1}, 1, 10, 12, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWatcher(tt.start)
			w.Dir = tt.dir
			for range tt.ticks {
				w.Patrol(tt.grid, WatcherSpeed)
			}
			if w.Box.X != tt.wantX || w.Dir != tt.wantDir {
				t.Fatalf("x=%v dir=%d, want x=%v dir=%d", w.Box.X, w.Dir, tt.wantX, tt.wantDir)
			}
		})
	}
}

func TestWatcherSight(t *testing.T) {
	open := gridFrom(
		"#..................#",
		"#..................#",
		"####################",
	)
	blocked := gridFrom(
		"#..................#",
		"#....#.............#",
		"####################",
	)
	standing := func(x float64) geom.Rect { return geom.Rect{X: x, Y: 11, W: 6, H: 9} }

	tests := []struct {
		name   string
		grid   level.Grid
		dir    int
		target geom.Rect
		want   bool
	}{
		{"in front within range", open, 1, standing(60), true},
		{"just inside range", open, 1, standing(105), true},
		{"out of range", open, 1, standing(110), false},
		{"behind", open, 1, standing(5), false},
		{"facing left sees left", open, -1, standing(5), true},
		{"blocked by solid tile", blocked, 1, standing(70), false},
		{"in front of the block", blocked, 1, standing(40), true},
		{"above eye line", open, 1, geom.Rect{X: 60, Y: 0, W: 6, H: 9}, false},
		{"bottom exactly at eye line", open, 1, geom.Rect{X: 60, Y: 4, W: 6, H: 9}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWatcher(level.Point{X: 2, Y: 1})
			w.Dir = tt.dir
			if got := w.Sees(tt.target, tt.grid); got != tt.want {
				t.Fatalf("Sees() = %v, want %v (sight end %v)", got, tt.want, w.SightEnd(tt.grid))
			}
		})
	}
}
