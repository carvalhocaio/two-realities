package ebitenapp

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	pauseDimAlpha = 0.7
	pauseTitleY   = 50
	pauseFirstY   = 80
	pauseStep     = 14
	pauseHintGap  = 6
	beatDotSize   = 3
	beatDotGap    = 6
	offsetHint    = "LEFT/RIGHT OR -/= TO ADJUST, MATCH THE DOT TO THE KICK"
)

func onOff(on bool) string {
	if on {
		return "ON"
	}
	return "OFF"
}

func menuLabels(f *game.Flow) [game.MenuCount]string {
	return [game.MenuCount]string{
		game.MenuResume:      "RESUME",
		game.MenuAssist:      "ASSIST MODE: " + onOff(f.Assist()),
		game.MenuMute:        "MUTE: " + onOff(f.Muted()),
		game.MenuAudioOffset: "AUDIO OFFSET: " + offsetLabel(f.AudioOffsetMs()),
		game.MenuTitle:       "BACK TO TITLE",
	}
}

func offsetLabel(ms int) string {
	if ms > 0 {
		return fmt.Sprintf("+%dMS", ms)
	}
	return fmt.Sprintf("%dMS", ms)
}

func drawPauseMenu(dst *ebiten.Image, f *game.Flow, phase float64) {
	pal := palettes[f.Session().Active()]
	fillRect(dst, geom.Rect{W: ScreenWidth, H: ScreenHeight}, pal.background, pauseDimAlpha)
	drawCenteredLine(dst, "PAUSED", pauseTitleY, subtitleScale, pal.text)
	for i, label := range menuLabels(f) {
		c := color.Color(pal.text)
		if game.MenuItem(i) == f.Menu() {
			c = pal.exit
			label = "- " + label + " -"
		}
		y := float64(pauseFirstY + i*pauseStep)
		drawCenteredLine(dst, label, y, 1, c)
		if game.MenuItem(i) == game.MenuAudioOffset {
			drawBeatDot(dst, (ScreenWidth+textWidth(label))/2+beatDotGap, y+1, phase, pal)
		}
	}
	if f.Menu() == game.MenuAudioOffset {
		drawCenteredLine(dst, offsetHint, float64(pauseFirstY+int(game.MenuCount)*pauseStep+pauseHintGap), 1, pal.text)
	}
}

func drawBeatDot(dst *ebiten.Image, x, y, phase float64, pal palette) {
	alpha := float32(math.Pow(1-phase, pulseSharpness))
	fillRect(dst, geom.Rect{X: x, Y: y, W: beatDotSize, H: beatDotSize}, pal.exit, alpha)
}

func drawCenteredLine(dst *ebiten.Image, msg string, y, scale float64, c color.Color) {
	w := scaledTextWidth(msg, scale)
	drawScaledText(dst, msg, (ScreenWidth-w)/2, y, scale, c, 1)
}
