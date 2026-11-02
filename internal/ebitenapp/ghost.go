package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	ghostOutlineAlpha = 0.75
	ghostDotAlpha     = 0.9
	dotSpacing        = 2
)

func newGhostImage(g level.Grid, spikes []level.Point, c, spike color.Color) *ebiten.Image {
	img := ebiten.NewImage(g.Width()*level.TileSize, g.Height()*level.TileSize)
	for y := range g.Height() {
		for x := range g.Width() {
			if g.Solid(x, y) {
				drawDottedEdges(img, g, x, y, c)
			}
		}
	}
	drawGhostSpikes(img, spikes, spike)
	return img
}

func drawDot(dst *ebiten.Image, x, y float64, c color.Color) {
	if (int(x)+int(y))%dotSpacing == 0 {
		fillRect(dst, geom.Rect{X: x, Y: y, W: 1, H: 1}, c, ghostDotAlpha)
	}
}

func drawDottedEdges(dst *ebiten.Image, g level.Grid, x, y int, c color.Color) {
	r := level.TileRect(level.Point{X: x, Y: y})
	for i := range level.TileSize {
		o := float64(i)
		if !g.Solid(x, y-1) {
			drawDot(dst, r.X+o, r.Y, c)
		}
		if !g.Solid(x, y+1) {
			drawDot(dst, r.X+o, r.Bottom()-1, c)
		}
		if !g.Solid(x-1, y) {
			drawDot(dst, r.X, r.Y+o, c)
		}
		if !g.Solid(x+1, y) {
			drawDot(dst, r.Right()-1, r.Y+o, c)
		}
	}
}

func drawDottedBox(dst *ebiten.Image, r geom.Rect, c color.Color) {
	r = geom.Rect{X: float64(int(r.X)), Y: float64(int(r.Y)), W: r.W, H: r.H}
	for x := r.X; x < r.Right(); x++ {
		drawDot(dst, x, r.Y, c)
		drawDot(dst, x, r.Bottom()-1, c)
	}
	for y := r.Y; y < r.Bottom(); y++ {
		drawDot(dst, r.X, y, c)
		drawDot(dst, r.Right()-1, y, c)
	}
}
