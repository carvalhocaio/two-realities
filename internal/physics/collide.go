package physics

import (
	"math"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

type Result struct {
	Box        geom.Rect
	HitWall    bool
	OnGround   bool
	HitCeiling bool
}

func MoveAndCollide(box geom.Rect, vel geom.Vec, g level.Grid) Result {
	var res Result

	box.X += vel.X
	if vel.X != 0 && OverlapsSolid(box, g) {
		box.X = snap(box.X, box.W, vel.X)
		res.HitWall = true
	}

	box.Y += vel.Y
	if vel.Y != 0 && OverlapsSolid(box, g) {
		box.Y = snap(box.Y, box.H, vel.Y)
		res.OnGround = vel.Y > 0
		res.HitCeiling = vel.Y < 0
	}

	res.Box = box
	return res
}

func snap(pos, size, vel float64) float64 {
	if vel > 0 {
		return math.Floor((pos+size)/level.TileSize)*level.TileSize - size
	}
	return (math.Floor(pos/level.TileSize) + 1) * level.TileSize
}
