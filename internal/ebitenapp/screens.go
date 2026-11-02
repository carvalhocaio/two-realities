package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/i18n"
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

func drawTitleScreen(dst *ebiten.Image, lang i18n.Lang, storageNotice bool, tick int) {
	yellow, grey := palettes[world.Yellow], palettes[world.Grey]
	lines := []line{
		{text: lang.T(i18n.Title), y: 40, scale: titleScale, c: yellow.text},
		{text: lang.T(i18n.TitlePrompt), y: 90, scale: 1, c: grey.background, blink: true},
		{text: lang.T(i18n.TitleLanguage), y: 108, scale: 1, c: yellow.exit},
		{text: lang.T(i18n.TitleControlsMove), y: 130, scale: 1, c: yellow.text},
		{text: lang.T(i18n.TitleControlsWorld), y: 140, scale: 1, c: yellow.text},
		{text: lang.T(i18n.TitleControlsExtra), y: 150, scale: 1, c: yellow.text},
	}
	if storageNotice {
		lines = append(lines, line{text: lang.T(i18n.StorageNotice), y: 168, scale: 1, c: dangerColor})
	}
	drawLines(dst, yellow.background, lines, tick)
}

func drawEndScreen(dst *ebiten.Image, sum game.RunSummary, lang i18n.Lang, tick int) {
	grey := palettes[world.Grey]
	drawLines(dst, grey.background, []line{
		{text: lang.T(i18n.EndTitle), y: 40, scale: subtitleScale, c: fragmentColor},
		{text: lang.T(i18n.EndThanks), y: 64, scale: 1, c: grey.text},
		{text: runTimeText(sum, lang), y: 84, scale: 1, c: grey.text},
		{text: bestRunText(sum, lang), y: 96, scale: 1, c: grey.text},
		{text: lang.T(i18n.EndSecrets, sum.Secrets, sum.SecretsTotal), y: 108, scale: 1, c: secretColor},
		{text: lang.T(i18n.EndPrompt), y: 136, scale: 1, c: grey.text, blink: true},
	}, tick)
}

func runTimeText(sum game.RunSummary, lang i18n.Lang) string {
	text := lang.T(i18n.EndRun, formatTicks(sum.Ticks))
	switch {
	case sum.Assisted:
		text += "  " + lang.T(i18n.EndAssist)
	case !sum.Eligible:
		text += "  " + lang.T(i18n.EndPractice)
	case sum.NewBest:
		text += "  " + lang.T(i18n.EndNewBest)
	}
	return text
}

func bestRunText(sum game.RunSummary, lang i18n.Lang) string {
	if !sum.HasBest {
		return lang.T(i18n.EndNoBest)
	}
	return lang.T(i18n.EndBestRun, formatTicks(sum.Best))
}
