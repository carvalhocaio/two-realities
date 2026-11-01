package player

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/physics"
)

type Player struct {
	Box      geom.Rect
	Vel      geom.Vec
	grounded bool
	jump     jumpWindow
}

func New(spawn level.Point) Player {
	tile := level.TileRect(spawn)
	return Player{Box: geom.Rect{
		X: tile.X + (tile.W-Width)/2,
		Y: tile.Bottom() - Height,
		W: Width,
		H: Height,
	}}
}

func (p *Player) Grounded() bool { return p.grounded }

func (p *Player) Update(in input.Intent, g level.Grid) {
	p.Vel.X = horizontal(in) * RunSpeed

	if p.jump.tick(p.grounded, in.JumpPressed) {
		p.Vel.Y = -JumpSpeed
	}
	p.Vel.Y = min(p.Vel.Y+Gravity, MaxFall)

	res := physics.MoveAndCollide(p.Box, p.Vel, g)
	p.Box = res.Box
	p.grounded = res.OnGround
	if res.OnGround || res.HitCeiling {
		p.Vel.Y = 0
	}
}

func horizontal(in input.Intent) float64 {
	switch {
	case in.Left && !in.Right:
		return -1
	case in.Right && !in.Left:
		return 1
	}
	return 0
}
