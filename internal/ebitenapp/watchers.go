package ebitenapp

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	sightAlpha  = 0.5
	visorWidth  = 4
	meterHeight = 1
	markGap     = 2
)

func drawWatchers(dst *ebiten.Image, watchers []game.Watcher, grey level.Grid, active bool, limit int) {
	for _, w := range watchers {
		if !active {
			drawDottedBox(dst, w.Box, enemyColor)
			continue
		}
		drawSight(dst, w, grey)
		fillRect(dst, w.Box, enemyColor, 1)
		drawEye(dst, w)
		drawAlert(dst, w, limit)
	}
}

func drawOutline(dst *ebiten.Image, r geom.Rect, c color.Color) {
	edges := []geom.Rect{
		{X: r.X, Y: r.Y, W: r.W, H: 1},
		{X: r.X, Y: r.Bottom() - 1, W: r.W, H: 1},
		{X: r.X, Y: r.Y, W: 1, H: r.H},
		{X: r.Right() - 1, Y: r.Y, W: 1, H: r.H},
	}
	for _, e := range edges {
		fillRect(dst, e, c, ghostOutlineAlpha)
	}
}

func drawSight(dst *ebiten.Image, w game.Watcher, grey level.Grid) {
	origin := w.SightOrigin()
	end := w.SightEnd(grey)
	x0, x1 := min(origin.X, end), max(origin.X, end)
	fillRect(dst, geom.Rect{X: x0, Y: origin.Y, W: x1 - x0, H: 1}, dangerColor, sightAlpha)
}

func drawEye(dst *ebiten.Image, w game.Watcher) {
	x := w.Box.Right() - visorWidth - 1
	if w.Dir < 0 {
		x = w.Box.X + 1
	}
	fillRect(dst, geom.Rect{X: math.Round(x), Y: w.EyeY(), W: visorWidth, H: 1}, dangerColor, 1)
}

func drawAlert(dst *ebiten.Image, w game.Watcher, limit int) {
	if w.Alert == 0 {
		return
	}
	fill := w.AlertLevel(limit)
	markX := w.Box.X + (w.Box.W-glyphWidth)/2
	markY := w.Box.Y - markGap - glyphHeight - meterHeight - 1
	drawText(dst, "!", markX, markY, dangerColor, float32(0.4+0.6*fill))

	meter := geom.Rect{X: w.Box.X, Y: w.Box.Y - markGap - meterHeight, W: w.Box.W, H: meterHeight}
	fillRect(dst, meter, enemyColor, 0.5)
	meter.W *= fill
	fillRect(dst, meter, dangerColor, 1)
}
