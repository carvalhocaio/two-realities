package ebitenapp

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	signalHeardAlpha = 0.35
	signalMastWidth  = 2
	signalMastHeight = 6
	signalHeadWidth  = 6
	signalWavePeriod = 45
	signalWaveTicks  = 18
	signalWaveRise   = 3
	signalWaveGrow   = 4
)

func drawSignals(dst *ebiten.Image, s *game.Session, pal palette, tick int) {
	for i, tr := range s.Level().Transmissions {
		if tr.Layer != s.Active() {
			continue
		}
		alpha := float32(1)
		if s.Heard(i) {
			alpha = signalHeardAlpha
		}
		tile := level.TileRect(tr.Point)
		mast := geom.Rect{X: tile.X + (tile.W-signalMastWidth)/2, Y: tile.Bottom() - signalMastHeight, W: signalMastWidth, H: signalMastHeight}
		head := geom.Rect{X: tile.X + (tile.W-signalHeadWidth)/2, Y: mast.Y - 2, W: signalHeadWidth, H: 2}
		fillRect(dst, mast, pal.tile, alpha)
		fillRect(dst, head, pal.text, alpha)
		if wave, a, ok := signalWave(head, tick+i*signalWavePeriod/3); ok && !s.Heard(i) {
			fillRect(dst, wave, pal.text, a)
		}
	}
}

func signalWave(head geom.Rect, tick int) (geom.Rect, float32, bool) {
	age := tick % signalWavePeriod
	if age >= signalWaveTicks {
		return geom.Rect{}, 0, false
	}
	t := float64(age) / signalWaveTicks
	w := 2 + 2*math.Round(t*signalWaveGrow/2)
	y := head.Y - 2 - math.Round(t*signalWaveRise)
	return geom.Rect{X: head.X + (head.W-w)/2, Y: y, W: w, H: 1}, float32(1 - t), true
}
