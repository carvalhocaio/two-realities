package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/level"
)

const blockedFillAlpha = 0.25

func drawBlockingTiles(dst *ebiten.Image, tiles []level.Point, strength float64) {
	for _, p := range tiles {
		tile := level.TileRect(p)
		fillRect(dst, tile, dangerColor, float32(strength*blockedFillAlpha))
		drawOutline(dst, tile, dangerColor, float32(strength))
	}
}
