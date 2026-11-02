package ebitenapp

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/i18n"
)

const (
	boxMargin      = 8
	boxPadding     = 4
	boxLineGap     = 3
	boxAlpha       = 0.92
	typewriterRate = 1
	boxClearance   = 6
)

type textbox struct {
	open bool
	top  bool
}

func (t *textbox) place(player geom.Rect, height float64) geom.Rect {
	if !t.open {
		t.open = true
		t.top = player.Y+player.H/2 > HUDHeight+ViewHeight/2
	}
	near := geom.Rect{X: player.X - boxClearance, Y: player.Y - boxClearance, W: player.W + 2*boxClearance, H: player.H + 2*boxClearance}
	box := t.rect(height)
	if box.Overlaps(near) {
		t.top = !t.top
		box = t.rect(height)
	}
	return box
}

func (t *textbox) rect(height float64) geom.Rect {
	y := float64(ScreenHeight - boxMargin - height)
	if t.top {
		y = HUDHeight + boxMargin
	}
	return geom.Rect{X: boxMargin, Y: y, W: ScreenWidth - 2*boxMargin, H: height}
}

func (t *textbox) draw(dst *ebiten.Image, s *game.Session, player geom.Rect, lang i18n.Lang) {
	view, ok := s.Transmission()
	if !ok {
		t.open = false
		return
	}
	pal := palettes[s.Active()]
	alpha := float32(1 - view.Fade)
	width := float64(ScreenWidth - 2*boxMargin)
	lines := wrapText(view.TextIn(lang), width-2*boxPadding)
	height := float64(len(lines))*(glyphHeight+boxLineGap) - boxLineGap + 2*boxPadding
	box := t.place(player, height)

	fillRect(dst, box, pal.panel, boxAlpha*alpha)
	drawOutline(dst, box, pal.tile)
	shown := view.Age / typewriterRate
	y := box.Y + boxPadding
	for _, line := range lines {
		runes := []rune(line)
		n := min(len(runes), max(shown, 0))
		drawText(dst, string(runes[:n]), box.X+boxPadding, y, pal.panelInk, alpha)
		shown -= len(runes) + 1
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
