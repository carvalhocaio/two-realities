package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	dangerColorAlpha = 0.35
	deathFlashAlpha  = 0.6
	deathBlinkTicks  = 4
	lockedExitAlpha  = 0.25
	exitRise         = 4
	exitFrame        = 2
	beatGhostFade    = 2.5
)

var pixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()

func fillRect(dst *ebiten.Image, r geom.Rect, c color.Color, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(r.W, r.H)
	op.GeoM.Translate(r.X, r.Y)
	op.ColorScale.ScaleWithColor(c)
	op.ColorScale.ScaleAlpha(alpha)
	dst.DrawImage(pixel, op)
}

func drawLayer(dst *ebiten.Image, g level.Grid, c color.Color) {
	for y := range g.Height() {
		for x := range g.Width() {
			if g.Solid(x, y) {
				fillRect(dst, level.TileRect(level.Point{X: x, Y: y}), c, 1)
			}
		}
	}
}

func drawScene(dst *ebiten.Image, s *game.Session, ghost *ebiten.Image, ghostAlpha float32, fx *playerFX, tick int) {
	active := s.Active()
	lvl := s.Level()
	pal := palettes[active]

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(ghostAlpha)
	dst.DrawImage(ghost, op)
	drawLayer(dst, lvl.Layer(active), pal.tile)
	drawSpikes(dst, lvl.Spikes[active], pal.spike)
	drawExit(dst, level.TileRect(lvl.Exit), pal, exitAlpha(s))
	drawCheckpoints(dst, s, pal)
	drawWardenBack(dst, s, pal, tick)
	drawFragments(dst, s.RemainingFragments(), tick, fragmentAlpha(active))
	drawSignals(dst, s, pal)
	drawSecret(dst, s, tick)
	drawWatchers(dst, s.Watchers(), lvl.Layer(world.Grey), active == world.Grey, s.AlertLimit())
	drawGlows(dst, s, tick)
	fx.draw(dst, s.PlayerBox(), pal, playerVisible(s, tick))
	drawWardenFront(dst, s, tick)
}

func exitAlpha(s *game.Session) float32 {
	if s.ExitOpen() {
		return 1
	}
	return lockedExitAlpha
}

func playerVisible(s *game.Session, tick int) bool {
	return s.Death() == game.Alive || (tick/deathBlinkTicks)%2 == 0
}

func fragmentAlpha(active world.World) float32 {
	if active == world.Grey {
		return 1
	}
	return inactiveFragmentAlpha
}

func drawFlash(dst *ebiten.Image, area geom.Rect, c color.Color, strength float64) {
	if strength <= 0 {
		return
	}
	fillRect(dst, area, c, float32(strength))
}

func drawExit(dst *ebiten.Image, tile geom.Rect, pal palette, alpha float32) {
	frame := geom.Rect{X: tile.X + 1, Y: tile.Y - exitRise, W: tile.W - 2, H: tile.H + exitRise}
	fillRect(dst, frame, pal.exit, alpha)
	door := geom.Rect{X: frame.X + exitFrame, Y: frame.Y + exitFrame, W: frame.W - 2*exitFrame, H: frame.H - exitFrame}
	fillRect(dst, door, secretColor, alpha)
}
