package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	glyphWidth   = 3
	glyphHeight  = 5
	glyphSpacing = 1
)

var glyphs = map[rune][glyphHeight]string{
	'A': {".#.", "#.#", "###", "#.#", "#.#"},
	'B': {"##.", "#.#", "##.", "#.#", "##."},
	'C': {".##", "#..", "#..", "#..", ".##"},
	'D': {"##.", "#.#", "#.#", "#.#", "##."},
	'E': {"###", "#..", "##.", "#..", "###"},
	'F': {"###", "#..", "##.", "#..", "#.."},
	'G': {".##", "#..", "#.#", "#.#", ".##"},
	'H': {"#.#", "#.#", "###", "#.#", "#.#"},
	'I': {"###", ".#.", ".#.", ".#.", "###"},
	'J': {"..#", "..#", "..#", "#.#", ".#."},
	'K': {"#.#", "#.#", "##.", "#.#", "#.#"},
	'L': {"#..", "#..", "#..", "#..", "###"},
	'M': {"#.#", "###", "###", "#.#", "#.#"},
	'N': {"##.", "#.#", "#.#", "#.#", "#.#"},
	'O': {".#.", "#.#", "#.#", "#.#", ".#."},
	'P': {"##.", "#.#", "##.", "#..", "#.."},
	'Q': {".#.", "#.#", "#.#", "##.", ".##"},
	'R': {"##.", "#.#", "##.", "#.#", "#.#"},
	'S': {".##", "#..", ".#.", "..#", "##."},
	'T': {"###", ".#.", ".#.", ".#.", ".#."},
	'U': {"#.#", "#.#", "#.#", "#.#", "###"},
	'V': {"#.#", "#.#", "#.#", "#.#", ".#."},
	'W': {"#.#", "#.#", "#.#", "###", "#.#"},
	'X': {"#.#", "#.#", ".#.", "#.#", "#.#"},
	'Y': {"#.#", "#.#", ".#.", ".#.", ".#."},
	'Z': {"###", "..#", ".#.", "#..", "###"},
	'.': {"...", "...", "...", "...", ".#."},
	',': {"...", "...", "...", ".#.", "#.."},
	'-': {"...", "...", "###", "...", "..."},
	':': {"...", ".#.", "...", ".#.", "..."},
	'0': {"###", "#.#", "#.#", "#.#", "###"},
	'1': {".#.", "##.", ".#.", ".#.", "###"},
	'2': {"##.", "..#", ".#.", "#..", "###"},
	'3': {"##.", "..#", ".#.", "..#", "##."},
	'4': {"#.#", "#.#", "###", "..#", "..#"},
	'5': {"###", "#..", "##.", "..#", "##."},
	'6': {".##", "#..", "###", "#.#", "###"},
	'7': {"###", "..#", ".#.", ".#.", ".#."},
	'8': {"###", "#.#", "###", "#.#", "###"},
	'9': {"###", "#.#", "###", "..#", "##."},
	'x': {"...", "#.#", ".#.", "#.#", "..."},
	'/': {"..#", "..#", ".#.", "#..", "#.."},
	'!': {".#.", ".#.", ".#.", "...", ".#."},
}

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
		drawGlyph(dst, glyphs[ch], ox, y, scale, c, alpha)
	}
}

func drawGlyph(dst *ebiten.Image, g [glyphHeight]string, x, y, scale float64, c color.Color, alpha float32) {
	for row, line := range g {
		for col, px := range line {
			if px == '#' {
				fillRect(dst, geom.Rect{X: x + float64(col)*scale, Y: y + float64(row)*scale, W: scale, H: scale}, c, alpha)
			}
		}
	}
}
