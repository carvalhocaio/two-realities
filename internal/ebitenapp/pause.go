package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	pauseDimAlpha = 0.7
	pauseTitleY   = 50
	pauseFirstY   = 80
	pauseStep     = 14
)

func onOff(on bool) string {
	if on {
		return "ON"
	}
	return "OFF"
}

func menuLabels(f *game.Flow) [game.MenuCount]string {
	return [game.MenuCount]string{
		game.MenuResume: "RESUME",
		game.MenuAssist: "ASSIST MODE: " + onOff(f.Assist()),
		game.MenuMute:   "MUTE: " + onOff(f.Muted()),
		game.MenuTitle:  "BACK TO TITLE",
	}
}

func drawPauseMenu(dst *ebiten.Image, f *game.Flow) {
	pal := palettes[f.Session().Active()]
	fillRect(dst, geom.Rect{W: ScreenWidth, H: ScreenHeight}, pal.background, pauseDimAlpha)
	drawCenteredLine(dst, "PAUSED", pauseTitleY, subtitleScale, pal.text)
	for i, label := range menuLabels(f) {
		c := color.Color(pal.text)
		if game.MenuItem(i) == f.Menu() {
			c = pal.exit
			label = "- " + label + " -"
		}
		drawCenteredLine(dst, label, float64(pauseFirstY+i*pauseStep), 1, c)
	}
}

func drawCenteredLine(dst *ebiten.Image, msg string, y, scale float64, c color.Color) {
	w := scaledTextWidth(msg, scale)
	drawScaledText(dst, msg, (ScreenWidth-w)/2, y, scale, c, 1)
}
