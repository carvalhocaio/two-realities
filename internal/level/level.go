package level

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/world"
)

type Level struct {
	Layers [world.Count]Grid
	Spawn  Point
	Exit   Point
}

func (l Level) Layer(w world.World) Grid { return l.Layers[w] }

func TileRect(p Point) geom.Rect {
	return geom.Rect{X: float64(p.X * TileSize), Y: float64(p.Y * TileSize), W: TileSize, H: TileSize}
}
