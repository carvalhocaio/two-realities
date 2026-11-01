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
	'E': {"###", "#..", "##.", "#..", "###"},
	'G': {".##", "#..", "#.#", "#.#", ".##"},
	'L': {"#..", "#..", "#..", "#..", "###"},
	'O': {".#.", "#.#", "#.#", "#.#", ".#."},
	'R': {"##.", "#.#", "##.", "#.#", "#.#"},
	'S': {".##", "#..", ".#.", "..#", "##."},
	'T': {"###", ".#.", ".#.", ".#.", ".#."},
	'V': {"#.#", "#.#", "#.#", "#.#", ".#."},
	'W': {"#.#", "#.#", "#.#", "###", "#.#"},
	'Y': {"#.#", "#.#", ".#.", ".#.", ".#."},
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

func textWidth(msg string) float64 {
	n := len([]rune(msg))
	if n == 0 {
		return 0
	}
	return float64(n*(glyphWidth+glyphSpacing) - glyphSpacing)
}

func drawText(dst *ebiten.Image, msg string, x, y float64, c color.Color, alpha float32) {
	for i, ch := range []rune(msg) {
		ox := x + float64(i*(glyphWidth+glyphSpacing))
		drawGlyph(dst, glyphs[ch], ox, y, c, alpha)
	}
}

func drawGlyph(dst *ebiten.Image, g [glyphHeight]string, x, y float64, c color.Color, alpha float32) {
	for row, line := range g {
		for col, px := range line {
			if px == '#' {
				fillRect(dst, geom.Rect{X: x + float64(col), Y: y + float64(row), W: 1, H: 1}, c, alpha)
			}
		}
	}
}
