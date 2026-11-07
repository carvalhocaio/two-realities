package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	epilogueTitleY  = 34
	epilogueFirstY  = 70
	epilogueStep    = 16
	epiloguePromptY = 150
)

func drawEpilogue(dst *ebiten.Image, lines []string, done bool, lang i18n.Lang, pad bool, tick int) {
	pal := palettes[world.Grey]
	dst.Fill(pal.background)
	drawCenteredLine(dst, lang.T(i18n.EpilogueTitle), epilogueTitleY, subtitleScale, fragmentColor)
	for i, l := range lines {
		drawCenteredLine(dst, l, float64(epilogueFirstY+i*epilogueStep), 1, pal.text)
	}
	if done && (tick/promptBlink)%2 == 0 {
		drawCenteredLine(dst, hint(lang, i18n.HelpPrompt, pad), epiloguePromptY, 1, pal.text)
	}
}
