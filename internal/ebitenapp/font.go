package ebitenapp

import (
	"image/color"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/font"
)

const (
	glyphWidth   = font.Width
	glyphHeight  = font.Height
	glyphSpacing = font.Spacing
)

func textWidth(msg string) float64 { return scaledTextWidth(msg, 1) }

func scaledTextWidth(msg string, scale float64) float64 {
	n := utf8.RuneCountInString(msg)
	if n == 0 {
		return 0
	}
	return float64(n*(glyphWidth+glyphSpacing)-glyphSpacing) * scale
}

func drawText(dst *ebiten.Image, msg string, x, y float64, c color.RGBA, alpha float32) {
	drawScaledText(dst, msg, x, y, 1, c, alpha)
}

func drawScaledText(dst *ebiten.Image, msg string, x, y, scale float64, c color.RGBA, alpha float32) {
	i := 0
	for _, ch := range msg {
		if img := glyphImage(ch); img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(x+float64(i*(glyphWidth+glyphSpacing))*scale, y-scale)
			scaleRGBA(&op.ColorScale, c)
			op.ColorScale.ScaleAlpha(alpha)
			dst.DrawImage(img, op)
		}
		i++
	}
}

const glyphCellHeight = glyphHeight + 2

var glyphImages = map[rune]*ebiten.Image{}

func glyphImage(ch rune) *ebiten.Image {
	if img, ok := glyphImages[ch]; ok {
		return img
	}
	g, ok := font.Lookup(ch)
	if !ok {
		glyphImages[ch] = nil
		return nil
	}
	img := ebiten.NewImage(glyphWidth, glyphCellHeight)
	img.WritePixels(glyphPixels(g))
	glyphImages[ch] = img
	return img
}

func glyphPixels(g font.Glyph) []byte {
	pix := make([]byte, 4*glyphWidth*glyphCellHeight)
	rows := append([]string{g.Above}, g.Rows[:]...)
	rows = append(rows, g.Below)
	for y, line := range rows {
		for x, px := range line {
			if px == '#' {
				copy(pix[4*(y*glyphWidth+x):], []byte{0xff, 0xff, 0xff, 0xff})
			}
		}
	}
	return pix
}

func runePrefix(s string, n int) string {
	if n <= 0 {
		return ""
	}
	for i := range s {
		if n == 0 {
			return s[:i]
		}
		n--
	}
	return s
}
