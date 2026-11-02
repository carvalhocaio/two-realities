package game

import (
	"math"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

type BeamStage int

const (
	BeamIdle BeamStage = iota
	BeamTelegraph
	BeamFiring
)

type Aim int

const (
	AimSweep Aim = iota
	AimTrack
)

const (
	RestBeats    = 2
	FireFraction = 0.5
	BeamWidth    = 6
)

type Beam struct {
	Stage    BeamStage
	X        float64
	Progress float64
	Box      geom.Rect
}

func beamCycle(rel int, phase float64, warning int) (int, BeamStage) {
	if rel < 0 {
		return -1, BeamIdle
	}
	length := warning + 1
	n, k := rel/length, rel%length
	switch {
	case k < warning:
		return n, BeamTelegraph
	case phase < FireFraction:
		return n, BeamFiring
	default:
		return n, BeamIdle
	}
}

func telegraphProgress(rel int, phase float64, warning int) float64 {
	k := rel % (warning + 1)
	return min((float64(k)+phase)/float64(warning), 1)
}

func sweepX(beats, speed float64, r level.Rail, fromRight bool) float64 {
	length := float64((r.To - r.From) * level.TileSize)
	d := 0.0
	if length > 0 && beats > 0 {
		d = math.Mod(beats*speed, 2*length)
		if d > length {
			d = 2*length - d
		}
	}
	if fromRight {
		return columnCenter(r.To) - d
	}
	return columnCenter(r.From) + d
}

func trackX(target geom.Rect, r level.Rail) float64 {
	col := tileOf(target.X + target.W/2)
	return columnCenter(min(max(col, r.From), r.To))
}

func columnCenter(col int) float64 { return float64(col*level.TileSize) + level.TileSize/2 }

func beamBox(x float64, r level.Rail, g level.Grid) geom.Rect {
	left, right := tileOf(x-BeamWidth/2), tileOf(x+BeamWidth/2-1e-6)
	top := r.Row + 1
	bottom := top
	for bottom < g.Height() && !g.Solid(left, bottom) && !g.Solid(right, bottom) {
		bottom++
	}
	return geom.Rect{
		X: x - BeamWidth/2,
		Y: float64(top * level.TileSize),
		W: BeamWidth,
		H: float64((bottom - top) * level.TileSize),
	}
}
