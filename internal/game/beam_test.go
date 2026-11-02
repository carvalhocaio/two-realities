package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

func TestBeamCycle(t *testing.T) {
	tests := []struct {
		name      string
		rel       int
		phase     float64
		warning   int
		wantCycle int
		wantStage BeamStage
	}{
		{"resting before the first cycle", -1, 0.2, 2, -1, BeamIdle},
		{"first warning beat", 0, 0.1, 2, 0, BeamTelegraph},
		{"second warning beat", 1, 0.9, 2, 0, BeamTelegraph},
		{"fires at the start of the third beat", 2, 0.0, 2, 0, BeamFiring},
		{"still firing before half beat", 2, 0.49, 2, 0, BeamFiring},
		{"cools down after half beat", 2, 0.5, 2, 0, BeamIdle},
		{"next cycle aims again", 3, 0.0, 2, 1, BeamTelegraph},
		{"one-beat warning alternates", 1, 0.2, 1, 0, BeamFiring},
		{"one-beat warning next cycle", 2, 0.2, 1, 1, BeamTelegraph},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, stage := beamCycle(tt.rel, tt.phase, tt.warning)
			if n != tt.wantCycle || stage != tt.wantStage {
				t.Fatalf("got cycle %d stage %v, want %d %v", n, stage, tt.wantCycle, tt.wantStage)
			}
		})
	}
}

func TestSweepX(t *testing.T) {
	rail := level.Rail{Row: 1, From: 1, To: 4}
	tests := []struct {
		beats, speed float64
		fromRight    bool
		want         float64
	}{
		{0, 10, false, 15},
		{1, 10, false, 25},
		{2.5, 10, false, 40},
		{3, 10, false, 45},
		{4, 10, false, 35},
		{6, 10, false, 15},
		{7, 10, false, 25},
		{1, 25, false, 40},
		{0, 10, true, 45},
		{1, 10, true, 35},
		{4, 10, true, 25},
	}
	for _, tt := range tests {
		if got := sweepX(tt.beats, tt.speed, rail, tt.fromRight); got != tt.want {
			t.Errorf("beats %v speed %v right %v: x = %v, want %v", tt.beats, tt.speed, tt.fromRight, got, tt.want)
		}
	}
}

func TestTrackX(t *testing.T) {
	rail := level.Rail{Row: 1, From: 2, To: 8}
	tests := []struct {
		x, want float64
	}{
		{42, 45},
		{0, 25},
		{200, 85},
	}
	for _, tt := range tests {
		if got := trackX(geom.Rect{X: tt.x, W: 6, H: 9}, rail); got != tt.want {
			t.Errorf("x=%v beam x = %v, want %v", tt.x, got, tt.want)
		}
	}
}

func TestBeamStopsAtFirstGreySolid(t *testing.T) {
	g := level.NewGrid(3, 6, []bool{
		true, true, true,
		false, false, false,
		false, false, false,
		false, true, false,
		false, false, false,
		true, true, true,
	})
	rail := level.Rail{Row: 1, From: 0, To: 2}
	tests := []struct {
		x, want float64
	}{
		{5, 30},
		{15, 10},
		{10, 10},
	}
	for _, tt := range tests {
		box := beamBox(tt.x, rail, g)
		if box.Y != 20 || box.H != tt.want || box.X != tt.x-BeamWidth/2 {
			t.Errorf("beam at x=%v = %+v, want height %v from y=20", tt.x, box, tt.want)
		}
	}
}
