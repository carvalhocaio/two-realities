package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
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

func drawLayer(dst *ebiten.Image, g level.Grid, c color.Color, alpha float32) {
	for y := range g.Height() {
		for x := range g.Width() {
			if g.Solid(x, y) {
				fillRect(dst, level.TileRect(level.Point{X: x, Y: y}), c, alpha)
			}
		}
	}
}

func drawScene(dst *ebiten.Image, s *game.Session) {
	active := s.Active()
	inactive := active.Other()
	lvl := s.Level()
	pal := palettes[active]

	dst.Fill(pal.background)
	drawLayer(dst, lvl.Layer(inactive), palettes[inactive].tile, ghostAlpha)
	drawLayer(dst, lvl.Layer(active), pal.tile, 1)
	fillRect(dst, level.TileRect(lvl.Exit), pal.exit, 1)
	fillRect(dst, s.PlayerBox(), pal.player, 1)
}

func drawFlash(dst *ebiten.Image, strength float64) {
	if strength <= 0 {
		return
	}
	b := dst.Bounds()
	area := geom.Rect{W: float64(b.Dx()), H: float64(b.Dy())}
	fillRect(dst, area, blockedFlash, float32(strength*0.35))
}
