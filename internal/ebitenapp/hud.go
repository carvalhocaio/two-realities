package ebitenapp

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	hudMargin      = 3
	hudTop         = 1
	hudPadding     = 1
	panelHeight    = glyphHeight + 2*hudPadding
	swatchSize     = glyphHeight
	panelAlpha     = 0.9
	blinkTicks     = 3
	restartHintMsg = "R TO RESTART"
	mutedMsg       = "MUTE"
)

type hudState struct {
	phase  float64
	muted  bool
	notice string
}

var hudDiamondRows = []int{1, 3, 5, 3, 1}

func panelAt(x, w float64) geom.Rect {
	return geom.Rect{X: x, Y: hudTop, W: w + 2*hudPadding, H: panelHeight}
}

func drawPanel(dst *ebiten.Image, r geom.Rect, pal palette) {
	fillRect(dst, r, pal.background, panelAlpha)
}

func drawHUD(dst *ebiten.Image, s *game.Session, tick int, h hudState) {
	pal := palettes[s.Active()]
	fillRect(dst, geom.Rect{W: ScreenWidth, H: HUDHeight}, pal.tile, 1)
	drawBeatPulse(dst, h.phase, pal)
	drawWorldIndicator(dst, s, pal, tick)
	right := drawFragmentCounter(dst, s, pal, tick)
	if h.muted {
		drawPanelText(dst, mutedMsg, right-hudMargin-textWidth(mutedMsg)-2*hudPadding, pal, pal.text)
	}
	switch {
	case h.notice != "":
		drawCentered(dst, h.notice, pal, pal.text)
	case s.RestartHint():
		drawCentered(dst, restartHintMsg, pal, warningColor)
	default:
		drawCentered(dst, fmt.Sprintf("LV %d/%d", s.LevelIndex()+1, s.LevelCount()), pal, pal.text)
	}
}

func drawWorldIndicator(dst *ebiten.Image, s *game.Session, pal palette, tick int) {
	shown := shownWorld(s, tick)
	label := strings.ToUpper(shown.String())
	panel := panelAt(hudMargin, swatchSize+hudPadding+textWidth(label))
	drawPanel(dst, panel, pal)

	swatch := geom.Rect{X: panel.X + hudPadding, Y: panel.Y + hudPadding, W: swatchSize, H: swatchSize}
	fillRect(dst, swatch, palettes[shown].tile, 1)
	drawText(dst, label, swatch.Right()+hudPadding, panel.Y+hudPadding, pal.text, 1)
}

func drawPanelText(dst *ebiten.Image, msg string, x float64, pal palette, c color.Color) {
	panel := panelAt(x, textWidth(msg))
	drawPanel(dst, panel, pal)
	drawText(dst, msg, panel.X+hudPadding, panel.Y+hudPadding, c, 1)
}

func drawFragmentCounter(dst *ebiten.Image, s *game.Session, pal palette, tick int) float64 {
	label := fmt.Sprintf("x%d", s.Fragments())
	iconWidth := float64(hudDiamondRows[len(hudDiamondRows)/2])
	w := iconWidth + hudPadding + textWidth(label)
	panel := panelAt(float64(dst.Bounds().Dx())-hudMargin-w-2*hudPadding, w)
	drawPanel(dst, panel, pal)

	icon := geom.Rect{X: panel.X + hudPadding, Y: panel.Y + hudPadding, W: iconWidth, H: glyphHeight}
	drawRows(dst, icon, hudDiamondRows, fragmentColor, 1)
	drawText(dst, label, icon.Right()+hudPadding, panel.Y+hudPadding, counterColor(s.Feedback().CounterFlash, pal, tick), 1)
	return panel.X
}

func drawCentered(dst *ebiten.Image, msg string, pal palette, c color.Color) {
	drawPanelText(dst, msg, (float64(dst.Bounds().Dx())-textWidth(msg))/2-hudPadding, pal, c)
}

func counterColor(flash float64, pal palette, tick int) color.Color {
	if flash > 0 && (tick/blinkTicks)%2 == 0 {
		return warningColor
	}
	return pal.text
}
