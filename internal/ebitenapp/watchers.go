package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	sightAlpha  = 0.5
	eyeSize     = 2
	meterHeight = 1
	markGap     = 2
)

func drawWatchers(dst *ebiten.Image, watchers []game.Watcher, grey level.Grid, active bool) {
	for _, w := range watchers {
		if !active {
			drawGhostBox(dst, w.Box, watcherColor)
			continue
		}
		drawSight(dst, w, grey)
		fillRect(dst, w.Box, watcherColor, 1)
		drawEye(dst, w)
		drawAlert(dst, w)
	}
}

func drawGhostBox(dst *ebiten.Image, r geom.Rect, c color.Color) {
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
	fillRect(dst, geom.Rect{X: x0, Y: origin.Y, W: x1 - x0, H: 1}, watcherColor, sightAlpha)
}

func drawEye(dst *ebiten.Image, w game.Watcher) {
	x := w.Box.Right() - eyeSize - 1
	if w.Dir < 0 {
		x = w.Box.X + 1
	}
	fillRect(dst, geom.Rect{X: x, Y: w.EyeY() - 1, W: eyeSize, H: eyeSize}, watcherEyeColor, 1)
}

func drawAlert(dst *ebiten.Image, w game.Watcher) {
	if w.Alert == 0 {
		return
	}
	fill := w.AlertLevel()
	markX := w.Box.X + (w.Box.W-glyphWidth)/2
	markY := w.Box.Y - markGap - glyphHeight - meterHeight - 1
	drawText(dst, "!", markX, markY, warningColor, float32(0.4+0.6*fill))

	meter := geom.Rect{X: w.Box.X, Y: w.Box.Y - markGap - meterHeight, W: w.Box.W, H: meterHeight}
	fillRect(dst, meter, watcherColor, 0.5)
	meter.W *= fill
	fillRect(dst, meter, warningColor, 1)
}
