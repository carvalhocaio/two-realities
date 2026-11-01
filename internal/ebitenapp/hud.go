package ebitenapp

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/basicfont"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	hudMargin      = 4
	hudPadding     = 3
	swatchSize     = 7
	lineHeight     = 13
	blinkTicks     = 3
	restartHintMsg = "R to restart"
)

var hudFace = text.NewGoXFace(basicfont.Face7x13)

func drawText(dst *ebiten.Image, msg string, x, y float64, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, msg, hudFace, op)
}

func textWidth(msg string) float64 {
	w, _ := text.Measure(msg, hudFace, lineHeight)
	return w
}

func drawPanel(dst *ebiten.Image, r geom.Rect, pal palette) {
	fillRect(dst, r, pal.background, 0.85)
}

func drawHUD(dst *ebiten.Image, s *game.Session, tick int) {
	pal := palettes[s.Active()]
	drawWorldIndicator(dst, s, pal)
	drawFragmentCounter(dst, s, pal, tick)
	if s.RestartHint() {
		drawRestartHint(dst, pal)
	}
}

func drawWorldIndicator(dst *ebiten.Image, s *game.Session, pal palette) {
	label := strings.ToUpper(s.Active().String())
	panel := geom.Rect{X: hudMargin, Y: hudMargin, W: swatchSize + 3*hudPadding + textWidth(label), H: lineHeight + 2*hudPadding}
	drawPanel(dst, panel, pal)

	swatch := geom.Rect{X: panel.X + hudPadding, Y: panel.Y + (panel.H-swatchSize)/2, W: swatchSize, H: swatchSize}
	fillRect(dst, swatch, pal.tile, 1)
	drawText(dst, label, swatch.Right()+hudPadding, panel.Y+hudPadding, pal.text)
}

func drawFragmentCounter(dst *ebiten.Image, s *game.Session, pal palette, tick int) {
	label := fmt.Sprintf("x%d", s.Fragments())
	w := float64(game.FragmentSize) + 3*hudPadding + textWidth(label)
	b := dst.Bounds()
	panel := geom.Rect{X: float64(b.Dx()) - hudMargin - w, Y: hudMargin, W: w, H: lineHeight + 2*hudPadding}
	drawPanel(dst, panel, pal)

	icon := geom.Rect{X: panel.X + hudPadding, Y: panel.Y + (panel.H-game.FragmentSize)/2, W: game.FragmentSize, H: game.FragmentSize}
	drawDiamond(dst, icon, 1)
	drawText(dst, label, icon.Right()+hudPadding, panel.Y+hudPadding, counterColor(s.Feedback().CounterFlash, pal, tick))
}

func counterColor(flash float64, pal palette, tick int) color.Color {
	if flash > 0 && (tick/blinkTicks)%2 == 0 {
		return warningColor
	}
	return pal.text
}

func drawRestartHint(dst *ebiten.Image, pal palette) {
	b := dst.Bounds()
	w := textWidth(restartHintMsg) + 2*hudPadding
	panel := geom.Rect{X: (float64(b.Dx()) - w) / 2, Y: 3 * lineHeight, W: w, H: lineHeight + 2*hudPadding}
	drawPanel(dst, panel, pal)
	drawText(dst, restartHintMsg, panel.X+hudPadding, panel.Y+hudPadding, warningColor)
}
