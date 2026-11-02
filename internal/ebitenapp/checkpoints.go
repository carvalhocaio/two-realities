package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	poleHeight      = 8
	flagWidth       = 4
	flagHeight      = 3
	idleFlagAlpha   = 0.5
	activeFlagAlpha = 1
)

func drawCheckpoints(dst *ebiten.Image, s *game.Session, pal palette) {
	if !s.CheckpointsActive() {
		return
	}
	active, ok := s.ActiveCheckpoint()
	for i, p := range s.Level().Checkpoints {
		tile := level.TileRect(p)
		pole := geom.Rect{X: tile.X + 3, Y: tile.Bottom() - poleHeight, W: 1, H: poleHeight}
		flag := geom.Rect{X: pole.Right(), Y: pole.Y, W: flagWidth, H: flagHeight}
		fillRect(dst, pole, pal.text, idleFlagAlpha)
		if ok && i == active {
			fillRect(dst, flag, fragmentColor, activeFlagAlpha)
			continue
		}
		fillRect(dst, flag, pal.text, idleFlagAlpha)
	}
}
