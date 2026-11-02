package ebitenapp

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	boxMargin      = 8
	boxPadding     = 4
	boxLineGap     = 3
	boxAlpha       = 0.92
	typewriterRate = 1
)

func drawTransmission(dst *ebiten.Image, s *game.Session) {
	view, ok := s.Transmission()
	if !ok {
		return
	}
	pal := palettes[s.Active()]
	alpha := float32(1 - view.Fade)
	width := float64(ScreenWidth - 2*boxMargin)
	lines := wrapText(view.Text, width-2*boxPadding)
	height := float64(len(lines))*(glyphHeight+boxLineGap) - boxLineGap + 2*boxPadding
	box := geom.Rect{X: boxMargin, Y: ScreenHeight - boxMargin - height, W: width, H: height}

	fillRect(dst, box, pal.background, boxAlpha*alpha)
	drawGhostBox(dst, box, pal.text)
	shown := view.Age / typewriterRate
	y := box.Y + boxPadding
	for _, line := range lines {
		n := min(len(line), max(shown, 0))
		drawText(dst, line[:n], box.X+boxPadding, y, pal.text, alpha)
		shown -= len(line) + 1
		y += glyphHeight + boxLineGap
	}
}

func wrapText(text string, width float64) []string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(text) {
		next := word
		if line != "" {
			next = line + " " + word
		}
		if line != "" && textWidth(next) > width {
			lines = append(lines, line)
			next = word
		}
		line = next
	}
	return append(lines, line)
}
