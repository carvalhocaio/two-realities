package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	helpTitleY   = 14
	helpFirstY   = 42
	helpStep     = 18
	helpIconX    = 14
	helpTextX    = 44
	helpPromptY  = 160
	helpBeat     = 36
	helpBeatDots = 4
)

type helpRow struct {
	key  i18n.Key
	icon func(dst *ebiten.Image, x, y float64, tick int)
}

var helpRows = []helpRow{
	{i18n.HelpWorlds, drawWorldsIcon},
	{i18n.HelpFragments, drawFragmentIcon},
	{i18n.HelpSpikes, drawSpikeIcon},
	{i18n.HelpWatchers, drawWatcherIcon},
	{i18n.HelpBeat, drawBeatIcon},
	{i18n.HelpRestart, drawRestartIcon},
}

func drawHelp(dst *ebiten.Image, lang i18n.Lang, tick int) {
	pal := palettes[world.Yellow]
	dst.Fill(pal.background)
	drawCenteredLine(dst, lang.T(i18n.HelpTitle), helpTitleY, subtitleScale, pal.text)
	for i, row := range helpRows {
		y := float64(helpFirstY + i*helpStep)
		row.icon(dst, helpIconX, y, tick)
		drawText(dst, lang.T(row.key), helpTextX, y+2, pal.text, 1)
	}
	if (tick/promptBlink)%2 == 0 {
		drawCenteredLine(dst, lang.T(i18n.HelpPrompt), helpPromptY, 1, pal.text)
	}
}

func iconTile(x, y float64) geom.Rect { return geom.Rect{X: x, Y: y, W: 10, H: 10} }

func drawWorldsIcon(dst *ebiten.Image, x, y float64, _ int) {
	fillRect(dst, iconTile(x, y), palettes[world.Yellow].tile, 1)
	drawDottedBox(dst, iconTile(x+13, y), palettes[world.Grey].tile)
}

func drawFragmentIcon(dst *ebiten.Image, x, y float64, _ int) {
	fillRect(dst, geom.Rect{X: x + 5, Y: y - 1, W: 12, H: 12}, palettes[world.Grey].background, 1)
	drawRows(dst, geom.Rect{X: x + 8, Y: y + 2, W: 6, H: 6}, diamondRows, fragmentColor, 1)
}

func drawSpikeIcon(dst *ebiten.Image, x, y float64, _ int) {
	for i := range 2 {
		drawRows(dst, geom.Rect{X: x + 6 + float64(i)*5, Y: y + 5, W: 5, H: 5}, spikeRows, dangerColor, 1)
	}
}

func drawWatcherIcon(dst *ebiten.Image, x, y float64, _ int) {
	body := geom.Rect{X: x + 7, Y: y, W: 8, H: 10}
	fillRect(dst, body, enemyColor, 1)
	fillRect(dst, geom.Rect{X: body.Right() - 4, Y: y + 2, W: 4, H: 2}, dangerColor, 1)
}

func drawBeatIcon(dst *ebiten.Image, x, y float64, tick int) {
	pal := palettes[world.Yellow]
	current := (tick / helpBeat) % helpBeatDots
	for i := range helpBeatDots {
		dot := geom.Rect{X: x + 4 + float64(i)*4, Y: y + 4, W: 3, H: 3}
		if i == current {
			fillRect(dst, dot, pal.exit, 1)
			continue
		}
		fillRect(dst, dot, pal.text, 0.35)
	}
}

func drawRestartIcon(dst *ebiten.Image, x, y float64, _ int) {
	pal := palettes[world.Yellow]
	box := geom.Rect{X: x + 6, Y: y - 1, W: 11, H: 11}
	drawPanel(dst, box, pal)
	drawText(dst, "R", box.X+4, box.Y+3, pal.panelInk, 1)
}
