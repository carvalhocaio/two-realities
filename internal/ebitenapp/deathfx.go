package ebitenapp

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	shardSize    = 2
	shardSpeed   = 0.25
	shardLift    = 1.2
	shardGravity = 0.2
	shardJitter  = 0.08
	ashMotes     = 6
	ashRise      = 1.2
	sliceHeight  = 3
	sliceShift   = 3
	fallSparks   = 5
	fallRiseMax  = 24
	fallSpread   = 4
)

func deathTicks(p float64) float64 { return p * game.DeathFreezeTicks }

func shardCount(box geom.Rect) int {
	return int(box.W/shardSize) * int(math.Ceil(box.H/shardSize))
}

func shard(box geom.Rect, i int, p float64) geom.Rect {
	cols := int(box.W / shardSize)
	x0 := box.X + float64(i%cols)*shardSize
	y0 := box.Y + float64(i/cols)*shardSize
	cx, cy := box.X+box.W/2, box.Y+box.H/2
	vx := (x0+shardSize/2-cx)*shardSpeed + float64((i*7)%5-2)*shardJitter
	vy := (y0+shardSize/2-cy)*shardSpeed - shardLift - float64((i*3)%4)*shardJitter
	t := deathTicks(p)
	y := min(y0+vy*t+shardGravity*t*t/2, box.Bottom()-shardSize)
	return geom.Rect{X: math.Round(x0 + vx*t), Y: math.Round(y), W: shardSize, H: shardSize}
}

func ashMote(box geom.Rect, i int, p float64) geom.Rect {
	x := box.X + math.Mod(float64(i)*2.5, box.W)
	start := box.Bottom() - 1 - float64(i%3)*2
	return geom.Rect{X: math.Round(x), Y: math.Round(start - ashRise*deathTicks(p)*float64(1+i%2)), W: 1, H: 1}
}

func dissolvedBody(box geom.Rect, p float64) geom.Rect {
	h := math.Round(box.H * (1 - p))
	return geom.Rect{X: box.X, Y: box.Y, W: box.W, H: h}
}

func sliceOffset(i int, p float64) float64 {
	dir := 1.0
	if i%2 == 1 {
		dir = -1
	}
	return math.Round(dir * sliceShift * math.Sin(math.Pi*min(p*2, 1)) * float64(1+i%2))
}

func fallSpark(box, bounds geom.Rect, i int, p float64) geom.Rect {
	x := box.X + box.W/2 + float64(i-fallSparks/2)*fallSpread/2
	rise := fallRiseMax * (1 - (1-p)*(1-p)) * (2 + float64(i%2)) / 3
	y := bounds.Bottom() - 1 - rise
	return geom.Rect{X: math.Round(x), Y: math.Round(y), W: 1, H: 1}
}

func drawDeath(dst *ebiten.Image, s *game.Session, pal palette, reduced bool) {
	cause, p, box := s.Death(), s.DeathProgress(), s.PlayerBox()
	alpha := float32(1 - p)
	if reduced {
		if cause != game.DeathFall {
			fillRect(dst, box, pal.player, alpha)
		}
		return
	}
	switch cause {
	case game.DeathSpikes, game.DeathCrushed:
		for i := range shardCount(box) {
			fillRect(dst, shard(box, i, p), pal.player, alpha)
		}
	case game.DeathBeam:
		fillRect(dst, dissolvedBody(box, p), pal.player, 1)
		for i := range ashMotes {
			fillRect(dst, ashMote(box, i, p), ashColor(i, pal), alpha)
		}
	case game.DeathCaught:
		for i := 0; float64(i*sliceHeight) < box.H; i++ {
			y := box.Y + float64(i*sliceHeight)
			h := min(sliceHeight, box.H-float64(i*sliceHeight))
			slice := geom.Rect{X: box.X + sliceOffset(i, p), Y: y, W: box.W, H: h}
			fillRect(dst, slice, pal.player, alpha)
			fillRect(dst, geom.Rect{X: slice.X, Y: y, W: box.W, H: 1}, dangerColor, alpha*float32(min(p*2, 1)))
		}
	case game.DeathFall:
		for i := range fallSparks {
			fillRect(dst, fallSpark(box, s.Level().Bounds(), i, p), pal.text, alpha)
		}
	}
}

func ashColor(i int, pal palette) color.RGBA {
	if i%2 == 0 {
		return beamCoreColor
	}
	return pal.player
}
