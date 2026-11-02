package ebitenapp

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	towerSeed       = 0x9E3779B9
	towerMinWidth   = 24
	towerMaxWidth   = 56
	towerMinHeight  = 50
	towerMaxHeight  = 150
	towerMinGap     = 6
	towerMaxGap     = 40
	vignetteInner   = 0.55
	vignetteOuter   = 1.2
	vignetteDarkest = 0.55
)

type towerLayer struct {
	parallax float64
	c        color.RGBA
	towers   []geom.Rect
}

var towerColors = [2]color.RGBA{
	{0x2F, 0x31, 0x38, 0xFF},
	{0x34, 0x37, 0x3F, 0xFF},
}

type rng uint32

func (r *rng) next(lo, hi int) float64 {
	*r = *r*1664525 + 1013904223
	return float64(lo + int(uint32(*r)>>8)%(hi-lo+1))
}

func newBackdrop(seed int, levelW, levelH float64) []towerLayer {
	layers := []towerLayer{{parallax: 0.25, c: towerColors[0]}, {parallax: 0.5, c: towerColors[1]}}
	r := rng(towerSeed + uint32(seed)*7919)
	for i := range layers {
		l := &layers[i]
		span := levelW*l.parallax + ScreenWidth
		for x := -r.next(0, towerMaxWidth); x < span; {
			w := r.next(towerMinWidth, towerMaxWidth)
			h := r.next(towerMinHeight, towerMaxHeight) * (1 + float64(i)*0.2)
			l.towers = append(l.towers, geom.Rect{X: x, Y: levelH - h, W: w, H: h})
			x += w + r.next(towerMinGap, towerMaxGap)
		}
	}
	return layers
}

func drawBackdrop(dst *ebiten.Image, layers []towerLayer, camX float64) {
	for _, l := range layers {
		shift := math.Round(camX * (1 - l.parallax))
		for _, t := range l.towers {
			t.X += shift
			fillRect(dst, t, l.c, 1)
		}
	}
}

var vignette = func() *ebiten.Image {
	img := ebiten.NewImage(ScreenWidth, ViewHeight)
	pix := make([]byte, ScreenWidth*ViewHeight*4)
	for y := range ViewHeight {
		for x := range ScreenWidth {
			dx := (float64(x)+0.5)/ScreenWidth*2 - 1
			dy := (float64(y)+0.5)/ViewHeight*2 - 1
			d := math.Hypot(dx, dy)
			t := min(max((d-vignetteInner)/(vignetteOuter-vignetteInner), 0), 1)
			a := byte(t * t * (3 - 2*t) * vignetteDarkest * 255)
			pix[(y*ScreenWidth+x)*4+3] = a
		}
	}
	img.WritePixels(pix)
	return img
}()

func drawVignette(dst *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(vignette, op)
}
