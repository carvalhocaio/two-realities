package ebitenapp

import (
	"image/color"

	"github.com/carvalhocaio/two-realities/internal/world"
)

const ghostAlpha = 0.28

type palette struct {
	background color.RGBA
	tile       color.RGBA
	player     color.RGBA
	exit       color.RGBA
}

var palettes = [world.Count]palette{
	world.Yellow: {
		background: color.RGBA{0xFF, 0xF6, 0xDC, 0xFF},
		tile:       color.RGBA{0xE8, 0xB4, 0x3A, 0xFF},
		player:     color.RGBA{0x3A, 0x2A, 0x10, 0xFF},
		exit:       color.RGBA{0xFF, 0x8A, 0x00, 0xFF},
	},
	world.Grey: {
		background: color.RGBA{0x2E, 0x30, 0x36, 0xFF},
		tile:       color.RGBA{0x7A, 0x7D, 0x85, 0xFF},
		player:     color.RGBA{0xE6, 0xE6, 0xE6, 0xFF},
		exit:       color.RGBA{0xA2, 0x3B, 0x3B, 0xFF},
	},
}

var blockedFlash = color.RGBA{0xB0, 0x30, 0x30, 0xFF}
