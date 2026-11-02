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
	eyeWidth          = 12
	eyeHeight         = 8
	pupilSize         = 3
	pupilReach        = 3
	railAlpha         = 0.35
	beamAlpha         = 0.45
	beamCoreAlpha     = 0.9
	ghostBeamAlpha    = 0.18
	telegraphMinAlpha = 0.25
	telegraphDash     = 3
	gateBarSpacing    = 3
	gateAlpha         = 0.85
	lampPulseTicks    = 20
	hitBlinkTicks     = 3
	pipSize           = 3
	pipGap            = 1
	pipSpaces         = 4
	unlitLampAlpha    = 0.4
)

var lampRows = []int{2, 4, 4, 2}

func drawWardenBack(dst *ebiten.Image, s *game.Session, pal palette, tick int) {
	w, ok := s.Warden()
	if !ok {
		return
	}
	drawGates(dst, s.ClosedGates(), pal)
	drawLamps(dst, s.Level().Boss.Lamps, w.Lit, pal, tick)
	drawBeam(dst, w.Beam, s.Active() == world.Grey, tick)
}

func drawWardenFront(dst *ebiten.Image, s *game.Session, tick int) {
	w, ok := s.Warden()
	if !ok {
		return
	}
	rail := s.Level().Boss.Rail
	drawRail(dst, rail)
	if w.Lit < level.LampCount {
		drawWardenEye(dst, w, rail, s.PlayerBox(), s.Active() == world.Grey, tick)
	}
}

func drawRail(dst *ebiten.Image, r level.Rail) {
	y := float64(r.Row*level.TileSize) + 1
	x0 := float64(r.From * level.TileSize)
	fillRect(dst, geom.Rect{X: x0, Y: y, W: float64((r.To - r.From + 1) * level.TileSize), H: 1}, dangerColor, railAlpha)
}

func drawWardenEye(dst *ebiten.Image, w game.WardenView, r level.Rail, target geom.Rect, active bool, tick int) {
	box := geom.Rect{X: math.Round(w.Beam.X - eyeWidth/2), Y: float64(r.Row*level.TileSize) + 1, W: eyeWidth, H: eyeHeight}
	if !active {
		drawDottedBox(dst, box, enemyColor)
		return
	}
	fillRect(dst, box, enemyColor, 1)
	if w.Hit > 0 && (tick/hitBlinkTicks)%2 == 0 {
		fillRect(dst, geom.Rect{X: box.X + 1, Y: box.Y + eyeHeight/2, W: eyeWidth - 2, H: 1}, deathFlash, 1)
		return
	}
	dx := max(-pupilReach, min(pupilReach, (target.X+target.W/2-w.Beam.X)/8))
	dy := max(0, min(eyeHeight-pupilSize-2, (target.Y-box.Y)/40))
	pupil := geom.Rect{
		X: math.Round(w.Beam.X - pupilSize/2 + dx),
		Y: math.Round(box.Y + 1 + dy),
		W: pupilSize, H: pupilSize,
	}
	fillRect(dst, pupil, dangerColor, 1)
}

func drawBeam(dst *ebiten.Image, b game.Beam, active bool, tick int) {
	switch b.Stage {
	case game.BeamFiring:
		if !active {
			fillRect(dst, b.Box, dangerColor, ghostBeamAlpha)
			return
		}
		fillRect(dst, b.Box, dangerColor, beamAlpha)
		core := geom.Rect{X: b.Box.X + 2, Y: b.Box.Y, W: b.Box.W - 4, H: b.Box.H}
		fillRect(dst, core, beamCoreColor, beamCoreAlpha)
	case game.BeamTelegraph:
		alpha := float32(telegraphMinAlpha + (1-telegraphMinAlpha)*b.Progress)
		if !active {
			alpha *= ghostBeamAlpha / beamAlpha
		}
		x := math.Round(b.X - 0.5)
		offset := float64(tick % (2 * telegraphDash))
		for y := b.Box.Y - offset; y < b.Box.Bottom(); y += 2 * telegraphDash {
			top := max(y, b.Box.Y)
			fillRect(dst, geom.Rect{X: x, Y: top, W: 1, H: min(y+telegraphDash, b.Box.Bottom()) - top}, dangerColor, alpha)
		}
	}
}

func drawLamps(dst *ebiten.Image, lamps []level.Point, lit int, pal palette, tick int) {
	for i, p := range lamps {
		tile := level.TileRect(p)
		bulb := geom.Rect{X: tile.X + 3, Y: tile.Y + 1, W: 4, H: float64(len(lampRows))}
		pole := geom.Rect{X: tile.X + 4, Y: bulb.Bottom(), W: 2, H: tile.Bottom() - bulb.Bottom()}
		fillRect(dst, pole, enemyColor, 1)
		switch {
		case i < lit:
			drawRows(dst, bulb, lampRows, fragmentColor, 1)
		case i == lit:
			pulse := float32(0.4 + 0.6*math.Abs(math.Sin(float64(tick)*math.Pi/lampPulseTicks)))
			drawRows(dst, bulb, lampRows, fragmentColor, pulse)
		default:
			drawRows(dst, bulb, lampRows, pal.text, unlitLampAlpha)
		}
	}
}

func drawGates(dst *ebiten.Image, gates []level.Point, pal palette) {
	for _, p := range gates {
		tile := level.TileRect(p)
		fillRect(dst, tile, pal.background, 1)
		for x := tile.X + 1; x < tile.Right(); x += gateBarSpacing {
			fillRect(dst, geom.Rect{X: x, Y: tile.Y, W: 1, H: tile.H}, pal.text, gateAlpha)
		}
	}
}

func drawWardenPips(dst *ebiten.Image, x, y float64, lit int) {
	for i := range level.LampCount {
		pip := geom.Rect{X: x + float64(i*(pipSize+pipGap)), Y: y + 1, W: pipSize, H: pipSize}
		if i < lit {
			fillRect(dst, pip, fragmentColor, 1)
			continue
		}
		fillRect(dst, pip, dangerColor, 1)
	}
}
