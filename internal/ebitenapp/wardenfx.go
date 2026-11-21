package ebitenapp

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	gateOpenTicks    = 14
	gateSparks       = 3
	shockTicks       = 12
	shockSparks      = 12
	shockStartRadius = 2
	shockEndRadius   = 22
	shockGlow        = 0.8
	recoilTicks      = 3
	recoilPixels     = 1
	impactSparks     = 4
	impactSpread     = 2
	noGate           = -1
)

type wardenFX struct {
	session *game.Session
	level   int
	epoch   int
	lit     int
	stage   game.BeamStage
	gate    int
	opening int
	lamp    geom.Vec
	shock   int
	recoil  int
}

func (f *wardenFX) update(s *game.Session) {
	w, ok := s.Warden()
	if !ok {
		*f = wardenFX{gate: noGate}
		return
	}
	if s != f.session || s.LevelIndex() != f.level || s.Epoch() != f.epoch {
		*f = wardenFX{session: s, level: s.LevelIndex(), epoch: s.Epoch(), lit: w.Lit, stage: w.Beam.Stage, gate: noGate}
		return
	}
	f.step(s.Level().Boss, w)
}

func (f *wardenFX) step(boss *level.Boss, w game.WardenView) {
	f.opening = max(f.opening-1, 0)
	f.shock = max(f.shock-1, 0)
	f.recoil = max(f.recoil-1, 0)
	if w.Lit > f.lit {
		f.light(boss, w.Lit-1)
	}
	f.lit = w.Lit
	if f.stage == game.BeamTelegraph && w.Beam.Stage == game.BeamFiring {
		f.recoil = recoilTicks
	}
	f.stage = w.Beam.Stage
}

func (f *wardenFX) light(boss *level.Boss, lamp int) {
	tile := level.TileRect(boss.Lamps[lamp])
	f.lamp = geom.Vec{X: tile.X + 5, Y: tile.Y + 3}
	f.shock = shockTicks
	f.gate, f.opening = noGate, 0
	if lamp < level.GateCount {
		f.gate, f.opening = lamp, gateOpenTicks
	}
}

func (f *wardenFX) gateProgress() float64 {
	return 1 - float64(f.opening)/gateOpenTicks
}

func (f *wardenFX) eyeLift(reduced bool) float64 {
	if reduced || f.recoil == 0 {
		return 0
	}
	return recoilPixels
}

func (f *wardenFX) drawOpening(dst *ebiten.Image, boss *level.Boss, pal palette, tick int) {
	if f.opening == 0 || f.gate == noGate {
		return
	}
	p := f.gateProgress()
	for _, col := range gateColumns(boss.Gates[f.gate]) {
		h := math.Round(col.H * (1 - p))
		for x := col.X + 1; x < col.Right(); x += gateBarSpacing {
			fillRect(dst, geom.Rect{X: x, Y: col.Y, W: 1, H: h}, pal.text, gateAlpha)
		}
		for i := range gateSparks {
			y := col.Y + h + float64((tick+i*3)%4)
			x := col.X + 1 + float64(i*gateBarSpacing)
			fillRect(dst, geom.Rect{X: x, Y: y, W: 1, H: 1}, fragmentColor, float32(1-p))
		}
	}
}

func (f *wardenFX) drawShock(dst *ebiten.Image, active world.World, reduced bool) {
	if f.shock == 0 || reduced {
		return
	}
	t := 1 - float64(f.shock)/shockTicks
	r := shockStartRadius + t*(shockEndRadius-shockStartRadius)
	drawSparkRing(dst, f.lamp, r, shockSparks, float32(1-t))
	if active == world.Grey {
		drawGlow(dst, f.lamp, r, fragmentColor, shockGlow*(1-t))
	}
}

func gateColumns(points []level.Point) []geom.Rect {
	var cols []geom.Rect
	for _, p := range points {
		tile := level.TileRect(p)
		merged := false
		for i, c := range cols {
			if c.X == tile.X {
				top := min(c.Y, tile.Y)
				cols[i] = geom.Rect{X: c.X, Y: top, W: c.W, H: max(c.Bottom(), tile.Bottom()) - top}
				merged = true
				break
			}
		}
		if !merged {
			cols = append(cols, tile)
		}
	}
	return cols
}

func drawBeamImpact(dst *ebiten.Image, box geom.Rect, tick int, reduced bool) {
	y := box.Bottom() - 1
	fillRect(dst, geom.Rect{X: box.X - impactSpread, Y: y, W: box.W + 2*impactSpread, H: 1}, beamCoreColor, 1)
	if reduced {
		return
	}
	span := int(box.W) + 2*impactSpread
	for i := range impactSparks {
		x := box.X - impactSpread + float64((tick*(i+2)+i*5)%span)
		lift := float64((tick + i) % 3)
		c := dangerColor
		if i%2 == 0 {
			c = beamCoreColor
		}
		fillRect(dst, geom.Rect{X: x, Y: y - 1 - lift, W: 1, H: 1}, c, 1)
	}
}
