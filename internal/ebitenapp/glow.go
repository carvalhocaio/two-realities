package ebitenapp

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	glowTextureSize = 64
	glowPulseSpeed  = 0.06
	glowPulseDepth  = 0.15
	fragmentGlowR   = 11
	fragmentGlow    = 0.45
	exitGlowR       = 18
	exitGlow        = 0.55
	lampGlowR       = 14
	lampGlow        = 0.6
	flagGlowR       = 7
	flagGlow        = 0.35
)

var glowTexture = func() *ebiten.Image {
	img := ebiten.NewImage(glowTextureSize, glowTextureSize)
	pix := make([]byte, glowTextureSize*glowTextureSize*4)
	half := glowTextureSize / 2.0
	for y := range glowTextureSize {
		for x := range glowTextureSize {
			d := math.Hypot(float64(x)+0.5-half, float64(y)+0.5-half) / half
			a := math.Pow(max(0, 1-d), 2)
			v := byte(a * 255)
			i := (y*glowTextureSize + x) * 4
			pix[i], pix[i+1], pix[i+2], pix[i+3] = v, v, v, v
		}
	}
	img.WritePixels(pix)
	return img
}()

func drawGlow(dst *ebiten.Image, center geom.Vec, radius float64, c color.Color, strength float64) {
	op := &ebiten.DrawImageOptions{Blend: ebiten.BlendLighter}
	scale := 2 * radius / glowTextureSize
	op.GeoM.Translate(-glowTextureSize/2, -glowTextureSize/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(center.X, center.Y)
	op.ColorScale.ScaleWithColor(c)
	op.ColorScale.ScaleAlpha(float32(strength))
	dst.DrawImage(glowTexture, op)
}

func centerOf(r geom.Rect) geom.Vec { return geom.Vec{X: r.X + r.W/2, Y: r.Y + r.H/2} }

func pulse(tick int, phase float64) float64 {
	return 1 - glowPulseDepth + glowPulseDepth*math.Sin(float64(tick)*glowPulseSpeed+phase)
}

func drawGlows(dst *ebiten.Image, s *game.Session, tick int) {
	if s.Active() != world.Grey {
		return
	}
	for i, p := range s.RemainingFragments() {
		drawGlow(dst, centerOf(game.FragmentBox(p)), fragmentGlowR, fragmentColor, fragmentGlow*pulse(tick, float64(i)))
	}
	if box, ok := s.SecretBox(); ok {
		drawGlow(dst, centerOf(box), fragmentGlowR, secretColor, fragmentGlow*pulse(tick, 0))
	}
	if s.ExitOpen() {
		drawGlow(dst, centerOf(level.TileRect(s.Level().Exit)), exitGlowR, fragmentColor, exitGlow*pulse(tick, 1))
	}
	if i, ok := s.ActiveCheckpoint(); ok {
		tile := level.TileRect(s.Level().Checkpoints[i])
		drawGlow(dst, geom.Vec{X: tile.X + 5, Y: tile.Y + 3}, flagGlowR, fragmentColor, flagGlow)
	}
	if w, ok := s.Warden(); ok {
		for i, p := range s.Level().Boss.Lamps[:min(w.Lit+1, len(s.Level().Boss.Lamps))] {
			strength := lampGlow
			if i == w.Lit {
				strength *= pulse(tick*3, 0)
			}
			tile := level.TileRect(p)
			drawGlow(dst, geom.Vec{X: tile.X + 5, Y: tile.Y + 3}, lampGlowR, fragmentColor, strength)
		}
	}
}
