package ebitenapp

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
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

const storageNoticeMsg = "PROGRESS WON'T BE SAVED IN THIS BROWSER"

func drawTitleScreen(dst *ebiten.Image, storageNotice bool, tick int) {
	yellow, grey := palettes[world.Yellow], palettes[world.Grey]
	lines := []line{
		{text: "TWO REALITIES", y: 40, scale: titleScale, c: yellow.text},
		{text: "PRESS SHIFT TO START", y: 90, scale: 1, c: grey.background, blink: true},
		{text: "ARROWS OR A D: MOVE    SPACE: JUMP", y: 130, scale: 1, c: yellow.text},
		{text: "SHIFT: SWAP WORLDS    R: RESTART", y: 140, scale: 1, c: yellow.text},
		{text: "M: MUTE    T: TIMER    P OR ESC: PAUSE", y: 150, scale: 1, c: yellow.text},
	}
	if storageNotice {
		lines = append(lines, line{text: storageNoticeMsg, y: 168, scale: 1, c: warningColor})
	}
	drawLines(dst, yellow.background, lines, tick)
}

func drawEndScreen(dst *ebiten.Image, sum game.RunSummary, tick int) {
	grey := palettes[world.Grey]
	drawLines(dst, grey.background, []line{
		{text: "BOTH WORLDS CROSSED", y: 40, scale: subtitleScale, c: fragmentColor},
		{text: "THANKS FOR PLAYING", y: 64, scale: 1, c: grey.text},
		{text: runTimeText(sum), y: 84, scale: 1, c: grey.text},
		{text: bestRunText(sum), y: 96, scale: 1, c: grey.text},
		{text: fmt.Sprintf("SECRETS %d/%d", sum.Secrets, sum.SecretsTotal), y: 108, scale: 1, c: secretColor},
		{text: "PRESS SHIFT FOR TITLE", y: 136, scale: 1, c: grey.text, blink: true},
	}, tick)
}

func runTimeText(sum game.RunSummary) string {
	text := "RUN " + formatTicks(sum.Ticks)
	switch {
	case sum.Assisted:
		text += "  ASSIST"
	case !sum.Eligible:
		text += "  PRACTICE"
	case sum.NewBest:
		text += "  NEW BEST!"
	}
	return text
}

func bestRunText(sum game.RunSummary) string {
	if !sum.HasBest {
		return "BEST RUN -"
	}
	return "BEST RUN " + formatTicks(sum.Best)
}
