package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/font"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	glyphWidth   = font.Width
	glyphHeight  = font.Height
	glyphSpacing = font.Spacing
)

func textWidth(msg string) float64 { return scaledTextWidth(msg, 1) }

func scaledTextWidth(msg string, scale float64) float64 {
	n := len([]rune(msg))
	if n == 0 {
		return 0
	}
	return float64(n*(glyphWidth+glyphSpacing)-glyphSpacing) * scale
}

func drawText(dst *ebiten.Image, msg string, x, y float64, c color.Color, alpha float32) {
	drawScaledText(dst, msg, x, y, 1, c, alpha)
}

func drawScaledText(dst *ebiten.Image, msg string, x, y, scale float64, c color.Color, alpha float32) {
	for i, ch := range []rune(msg) {
		ox := x + float64(i*(glyphWidth+glyphSpacing))*scale
		if g, ok := font.Lookup(ch); ok {
			drawGlyph(dst, g, ox, y, scale, c, alpha)
		}
	}
}

func drawGlyph(dst *ebiten.Image, g font.Glyph, x, y, scale float64, c color.Color, alpha float32) {
	drawGlyphRow(dst, g.Above, x, y-scale, scale, c, alpha)
	for row, line := range g.Rows {
		drawGlyphRow(dst, line, x, y+float64(row)*scale, scale, c, alpha)
	}
	drawGlyphRow(dst, g.Below, x, y+glyphHeight*scale, scale, c, alpha)
}

func drawGlyphRow(dst *ebiten.Image, line string, x, y, scale float64, c color.Color, alpha float32) {
	for col, px := range line {
		if px == '#' {
			fillRect(dst, geom.Rect{X: x + float64(col)*scale, Y: y, W: scale, H: scale}, c, alpha)
		}
	}
}
