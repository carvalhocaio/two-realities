package ebitenapp

import (
	"math"

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
	flutterTicks    = 12
	hoistSparks     = 4
	hoistSparkRange = 5
)

func drawCheckpoints(dst *ebiten.Image, s *game.Session, pal palette, hoist float64, tick int) {
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
			drawActiveFlag(dst, pole, hoist, tick)
			continue
		}
		fillRect(dst, flag, pal.text, idleFlagAlpha)
	}
}

func drawActiveFlag(dst *ebiten.Image, pole geom.Rect, hoist float64, tick int) {
	y := pole.Y + math.Round((1-hoist)*(poleHeight-flagHeight))
	body := geom.Rect{X: pole.Right(), Y: y, W: flagWidth - 1, H: flagHeight}
	tip := geom.Rect{X: body.Right(), Y: y + flutterOffset(tick), W: 1, H: flagHeight - 1}
	fillRect(dst, body, fragmentColor, activeFlagAlpha)
	fillRect(dst, tip, fragmentColor, activeFlagAlpha)
	if hoist < 1 {
		top := geom.Vec{X: pole.X + 0.5, Y: pole.Y}
		drawSparkRing(dst, top, 1+hoist*hoistSparkRange, hoistSparks, float32(1-hoist))
	}
}

func flutterOffset(tick int) float64 {
	return float64((tick / flutterTicks) % 2)
}
