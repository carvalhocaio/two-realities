package physics

import (
	"math"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const edgeEpsilon = 1e-6

func solidAt(g level.Grid, x, y int) bool {
	if y >= g.Height() {
		return false
	}
	if x < 0 || x >= g.Width() || y < 0 {
		return true
	}
	return g.Solid(x, y)
}

func tileSpan(start, size float64) (first, last int) {
	first = int(math.Floor(start / level.TileSize))
	last = int(math.Floor((start + size - edgeEpsilon) / level.TileSize))
	return first, last
}

func OverlapsSolid(box geom.Rect, g level.Grid) bool {
	x0, x1 := tileSpan(box.X, box.W)
	y0, y1 := tileSpan(box.Y, box.H)
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if solidAt(g, x, y) {
				return true
			}
		}
	}
	return false
}
