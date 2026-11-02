package ebitenapp

import (
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
)

func drawSignals(dst *ebiten.Image, s *game.Session, pal palette) {
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
		fillRect(dst, mast, pal.text, alpha)
		fillRect(dst, head, pal.exit, alpha)
	}
}
