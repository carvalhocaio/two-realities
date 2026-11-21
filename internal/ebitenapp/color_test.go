package ebitenapp

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func gameColors() []color.RGBA {
	out := []color.RGBA{dangerColor, enemyColor, fragmentColor, secretColor, beamCoreColor, deathFlash, {}, {0x80, 0x40, 0x20, 0x80}}
	for _, p := range palettes {
		out = append(out, p.background, p.tile, p.player, p.exit, p.text, p.spike, p.panel, p.panelInk, p.door)
	}
	return out
}

func TestScaleRGBAMatchesScaleWithColor(t *testing.T) {
	for _, c := range gameColors() {
		var want, got ebiten.ColorScale
		want.ScaleWithColor(c)
		want.ScaleAlpha(0.37)
		scaleRGBA(&got, c)
		got.ScaleAlpha(0.37)
		if got != want {
			t.Fatalf("color %v: got %v, want %v", c, got, want)
		}
	}
}

func TestScaleRGBADoesNotAllocate(t *testing.T) {
	var cs ebiten.ColorScale
	pal := palettes[1]
	if n := testing.AllocsPerRun(1000, func() { scaleRGBA(&cs, pal.tile) }); n != 0 {
		t.Fatalf("scaleRGBA allocated %v times per call", n)
	}
}
