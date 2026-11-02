package ebitenapp

import (
	"image/color"

	"github.com/carvalhocaio/two-realities/internal/world"
)

const inactiveFragmentAlpha = 0.35

type palette struct {
	background color.RGBA
	tile       color.RGBA
	player     color.RGBA
	exit       color.RGBA
	text       color.RGBA
	spike      color.RGBA
	panel      color.RGBA
	panelInk   color.RGBA
}

var palettes = [world.Count]palette{
	world.Yellow: {
		background: color.RGBA{0xFF, 0xF3, 0xD6, 0xFF},
		tile:       color.RGBA{0xE8, 0xA3, 0x3D, 0xFF},
		player:     color.RGBA{0x3B, 0x2A, 0x16, 0xFF},
		exit:       color.RGBA{0xFF, 0x8A, 0x1E, 0xFF},
		text:       color.RGBA{0x3B, 0x2A, 0x16, 0xFF},
		spike:      dangerColor,
		panel:      color.RGBA{0x3B, 0x2A, 0x16, 0xFF},
		panelInk:   color.RGBA{0xFF, 0xF3, 0xD6, 0xFF},
	},
	world.Grey: {
		background: color.RGBA{0x2B, 0x2D, 0x33, 0xFF},
		tile:       color.RGBA{0x56, 0x5A, 0x64, 0xFF},
		player:     color.RGBA{0xE4, 0xE2, 0xDC, 0xFF},
		exit:       fragmentColor,
		text:       color.RGBA{0xE4, 0xE2, 0xDC, 0xFF},
		spike:      dangerColor,
		panel:      enemyColor,
		panelInk:   color.RGBA{0xE4, 0xE2, 0xDC, 0xFF},
	},
}

var (
	dangerColor   = color.RGBA{0xD9, 0x43, 0x4F, 0xFF}
	enemyColor    = color.RGBA{0x11, 0x12, 0x16, 0xFF}
	fragmentColor = color.RGBA{0xFF, 0xD2, 0x3F, 0xFF}
	secretColor   = color.RGBA{0xFF, 0xF8, 0xE0, 0xFF}
	beamCoreColor = color.RGBA{0xFF, 0xD6, 0xD6, 0xFF}
	deathFlash    = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
)
