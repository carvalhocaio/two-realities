package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
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

func drawScene(dst *ebiten.Image, s *game.Session, ghost *ebiten.Image, tick int) {
	active := s.Active()
	lvl := s.Level()
	pal := palettes[active]

	dst.Fill(pal.background)
	dst.DrawImage(ghost, nil)
	drawLayer(dst, lvl.Layer(active), pal.tile)
	fillRect(dst, level.TileRect(lvl.Exit), pal.exit, 1)
	drawFragments(dst, s.RemainingFragments(), tick, fragmentAlpha(active))
	fillRect(dst, s.PlayerBox(), pal.player, 1)
}

func fragmentAlpha(active world.World) float32 {
	if active == world.Grey {
		return 1
	}
	return inactiveFragmentAlpha
}

func drawFlash(dst *ebiten.Image, strength float64) {
	if strength <= 0 {
		return
	}
	b := dst.Bounds()
	area := geom.Rect{W: float64(b.Dx()), H: float64(b.Dy())}
	fillRect(dst, area, blockedFlash, float32(strength*0.35))
}
