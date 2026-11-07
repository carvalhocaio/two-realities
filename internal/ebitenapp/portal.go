package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	exitRise    = 4
	exitPost    = 2
	exitMotes   = 3
	exitMoteGap = 13
	exitCycle   = 39
)

func drawExit(dst *ebiten.Image, tile geom.Rect, pal palette, alpha float32, open bool, tick int) {
	top := tile.Y - exitRise
	height := tile.H + exitRise
	fillRect(dst, geom.Rect{X: tile.X - 1, Y: tile.Bottom() - 1, W: tile.W + 2, H: 1}, pal.exit, alpha)
	fillRect(dst, geom.Rect{X: tile.X + 1, Y: top + 1, W: exitPost, H: height - 1}, pal.exit, alpha)
	fillRect(dst, geom.Rect{X: tile.Right() - 1 - exitPost, Y: top + 1, W: exitPost, H: height - 1}, pal.exit, alpha)
	fillRect(dst, geom.Rect{X: tile.X + 2, Y: top, W: tile.W - 4, H: exitPost}, pal.exit, alpha)
	door := geom.Rect{X: tile.X + 1 + exitPost, Y: top + exitPost, W: tile.W - 2 - 2*exitPost, H: height - exitPost - 1}
	fillRect(dst, door, pal.door, alpha)
	if !open {
		return
	}
	fillRect(dst, geom.Rect{X: door.X, Y: door.Bottom() - 1, W: door.W, H: 1}, secretColor, 0.8)
	for i := range exitMotes {
		phase := float64((tick+i*exitMoteGap)%exitCycle) / exitCycle
		x := door.X + float64((i*3+1)%int(door.W))
		y := door.Bottom() - 2 - phase*(door.H-2)
		fillRect(dst, geom.Rect{X: x, Y: y, W: 1, H: 1}, secretColor, float32(1-phase))
	}
}
