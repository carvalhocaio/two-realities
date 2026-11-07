package ebitenapp

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/i18n"
)

const (
	pauseDimAlpha = 0.7
	pauseTitleY   = 40
	pauseFirstY   = 58
	pauseStep     = 11
	pauseHintGap  = 6
	beatDotSize   = 3
	beatDotGap    = 6
)

func onOff(on bool, lang i18n.Lang) string {
	if on {
		return lang.T(i18n.On)
	}
	return lang.T(i18n.Off)
}

func menuLabels(f *game.Flow) [game.MenuCount]string {
	lang := f.Language()
	return [game.MenuCount]string{
		game.MenuSkip:        lang.T(i18n.MenuSkip),
		game.MenuHelp:        lang.T(i18n.MenuHelp),
		game.MenuResume:      lang.T(i18n.MenuResume),
		game.MenuAssist:      lang.T(i18n.MenuAssist, onOff(f.Assist(), lang)),
		game.MenuMusic:       lang.T(i18n.MenuMusic, volumeLabel(f.MusicVolume())),
		game.MenuSounds:      lang.T(i18n.MenuSounds, volumeLabel(f.SoundVolume())),
		game.MenuEffects:     lang.T(i18n.MenuEffects, onOff(f.ReduceEffects(), lang)),
		game.MenuAudioOffset: lang.T(i18n.MenuAudioOffset, offsetLabel(f.AudioOffsetMs())),
		game.MenuLanguage:    lang.T(i18n.MenuLanguage),
		game.MenuTitle:       lang.T(i18n.MenuTitle),
	}
}

func volumeLabel(level int) string { return fmt.Sprintf("%d/%d", level, game.VolumeSteps) }

func offsetLabel(ms int) string {
	if ms > 0 {
		return fmt.Sprintf("+%dMS", ms)
	}
	return fmt.Sprintf("%dMS", ms)
}

func drawPauseMenu(dst *ebiten.Image, f *game.Flow, phase float64) {
	pal := palettes[f.Session().Active()]
	fillRect(dst, geom.Rect{W: ScreenWidth, H: ScreenHeight}, pal.background, pauseDimAlpha)
	drawCenteredLine(dst, f.Language().T(i18n.Paused), pauseTitleY, subtitleScale, pal.text)
	labels, items := menuLabels(f), f.MenuItems()
	for i, item := range items {
		label, c := labels[item], color.Color(pal.text)
		if item == f.Menu() {
			c = pal.exit
			label = "- " + label + " -"
		}
		y := float64(pauseFirstY + i*pauseStep)
		drawCenteredLine(dst, label, y, 1, c)
		if item == game.MenuAudioOffset {
			drawBeatDot(dst, (ScreenWidth+textWidth(label))/2+beatDotGap, y+1, phase, pal)
		}
	}
	if f.Menu() == game.MenuAudioOffset {
		drawCenteredLine(dst, f.Language().T(i18n.OffsetHint), float64(pauseFirstY+len(items)*pauseStep+pauseHintGap), 1, pal.text)
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
