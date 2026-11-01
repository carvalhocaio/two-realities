package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	ghostOutlineAlpha = 0.75
	ghostHatchAlpha   = 0.35
	hatchSpacing      = 4
)

func newGhostImage(g level.Grid, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(g.Width()*level.TileSize, g.Height()*level.TileSize)
	for y := range g.Height() {
		for x := range g.Width() {
			if g.Solid(x, y) {
				drawHatch(img, x, y, c)
				drawExposedEdges(img, g, x, y, c)
			}
		}
	}
	return img
}

func drawHatch(dst *ebiten.Image, tx, ty int, c color.Color) {
	for py := range level.TileSize {
		for px := range level.TileSize {
			gx, gy := tx*level.TileSize+px, ty*level.TileSize+py
			if (gx+gy)%hatchSpacing == 0 {
				fillRect(dst, geom.Rect{X: float64(gx), Y: float64(gy), W: 1, H: 1}, c, ghostHatchAlpha)
			}
		}
	}
}

func drawExposedEdges(dst *ebiten.Image, g level.Grid, x, y int, c color.Color) {
	r := level.TileRect(level.Point{X: x, Y: y})
	edges := []struct {
		exposed bool
		line    geom.Rect
	}{
		{!g.Solid(x, y-1), geom.Rect{X: r.X, Y: r.Y, W: r.W, H: 1}},
		{!g.Solid(x, y+1), geom.Rect{X: r.X, Y: r.Bottom() - 1, W: r.W, H: 1}},
		{!g.Solid(x-1, y), geom.Rect{X: r.X, Y: r.Y, W: 1, H: r.H}},
		{!g.Solid(x+1, y), geom.Rect{X: r.Right() - 1, Y: r.Y, W: 1, H: r.H}},
	}
	for _, e := range edges {
		if e.exposed {
			fillRect(dst, e.line, c, ghostOutlineAlpha)
		}
	}
}
