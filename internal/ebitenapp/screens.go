package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	titleScale    = 3
	subtitleScale = 2
	promptBlink   = 30
)

type line struct {
	text  string
	y     float64
	scale float64
	c     color.Color
	blink bool
}

func drawLines(dst *ebiten.Image, bg color.Color, lines []line, tick int) {
	dst.Fill(bg)
	for _, l := range lines {
		if l.blink && (tick/promptBlink)%2 == 1 {
			continue
		}
		x := (ScreenWidth - scaledTextWidth(l.text, l.scale)) / 2
		drawScaledText(dst, l.text, x, l.y, l.scale, l.c, 1)
	}
}

func drawTitleScreen(dst *ebiten.Image, tick int) {
	yellow, grey := palettes[world.Yellow], palettes[world.Grey]
	drawLines(dst, yellow.background, []line{
		{text: "TWO REALITIES", y: 40, scale: titleScale, c: yellow.text},
		{text: "PRESS SHIFT TO START", y: 90, scale: 1, c: grey.background, blink: true},
		{text: "ARROWS OR A D: MOVE    SPACE: JUMP", y: 130, scale: 1, c: yellow.text},
		{text: "SHIFT: SWAP WORLDS    R: RESTART", y: 140, scale: 1, c: yellow.text},
		{text: "M: MUTE", y: 150, scale: 1, c: yellow.text},
	}, tick)
}

func drawEndScreen(dst *ebiten.Image, tick int) {
	grey := palettes[world.Grey]
	drawLines(dst, grey.background, []line{
		{text: "BOTH WORLDS CROSSED", y: 60, scale: subtitleScale, c: fragmentColor},
		{text: "THANKS FOR PLAYING", y: 90, scale: 1, c: grey.text},
		{text: "PRESS SHIFT FOR TITLE", y: 120, scale: 1, c: grey.text, blink: true},
	}, tick)
}
