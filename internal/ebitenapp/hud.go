package ebitenapp

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	hudMargin   = 3
	hudTop      = 1
	hudPadding  = 1
	panelHeight = glyphHeight + 2*hudPadding
	swatchSize  = glyphHeight
	panelAlpha  = 0.9
	blinkTicks  = 3
)

type hudState struct {
	phase  float64
	muted  bool
	timer  bool
	notice string
	lang   i18n.Lang
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
	drawWorldIndicator(dst, s, pal, tick, h.lang)
	right := drawFragmentCounter(dst, s, pal, tick)
	if h.muted {
		muted := h.lang.T(i18n.HudMuted)
		drawPanelText(dst, muted, right-hudMargin-textWidth(muted)-2*hudPadding, pal, pal.text)
	}
	result, showResult := s.LevelResult()
	switch {
	case h.notice != "":
		drawCentered(dst, h.notice, pal, pal.text)
	case showResult:
		drawCentered(dst, resultText(result, h.lang), pal, pal.text)
	case s.RestartHint():
		drawCentered(dst, h.lang.T(i18n.HudRestart), pal, dangerColor)
	default:
		drawLevelLabel(dst, s, h, pal)
	}
}

func drawLevelLabel(dst *ebiten.Image, s *game.Session, h hudState, pal palette) {
	w, ok := s.Warden()
	if !ok {
		drawCentered(dst, levelLabel(s, h.timer, h.lang), pal, pal.text)
		return
	}
	label := h.lang.T(i18n.HudWarden) + strings.Repeat(" ", pipSpaces)
	if h.timer {
		label += "  " + formatTicks(s.LevelTicks())
	}
	x := (float64(dst.Bounds().Dx())-textWidth(label))/2 - hudPadding
	drawPanelText(dst, label, x, pal, pal.text)
	name := textWidth(h.lang.T(i18n.HudWarden) + " ")
	drawWardenPips(dst, x+hudPadding+name, hudTop+hudPadding, w.Lit)
}

func levelLabel(s *game.Session, timer bool, lang i18n.Lang) string {
	label := lang.T(i18n.HudLevel, s.LevelIndex()+1, s.LevelCount())
	if timer {
		label += "  " + formatTicks(s.LevelTicks())
	}
	return label
}

func resultText(r game.LevelResult, lang i18n.Lang) string {
	switch {
	case r.NewBest:
		return lang.T(i18n.ResultNewBest, formatTicks(r.Ticks))
	case r.HasBest:
		return lang.T(i18n.ResultTimeBest, formatTicks(r.Ticks), formatTicks(r.Best))
	default:
		return lang.T(i18n.ResultTime, formatTicks(r.Ticks))
	}
}

func worldName(w world.World, lang i18n.Lang) string {
	if w == world.Yellow {
		return lang.T(i18n.WorldYellow)
	}
	return lang.T(i18n.WorldGrey)
}

func drawWorldIndicator(dst *ebiten.Image, s *game.Session, pal palette, tick int, lang i18n.Lang) {
	shown := shownWorld(s, tick)
	label := worldName(shown, lang)
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
		return dangerColor
	}
	return pal.text
}
