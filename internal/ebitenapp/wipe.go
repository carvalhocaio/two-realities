package ebitenapp

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	wipeTicks       = 9
	wipeBandWidth   = 6
	wipeStripeSize  = 4
	wipeShadowAlpha = 0.35
)

type wipe struct {
	last, from *ebiten.Image
	active     world.World
	epoch      int
	start      int
	running    bool
}

func newWipe() *wipe {
	return &wipe{
		last: ebiten.NewImage(ScreenWidth, ViewHeight),
		from: ebiten.NewImage(ScreenWidth, ViewHeight),
	}
}

func (w *wipe) begin(active world.World, epoch, tick int) {
	shifted := active != w.active && epoch == w.epoch
	w.active, w.epoch = active, epoch
	if !shifted {
		return
	}
	w.from.Clear()
	w.from.DrawImage(w.last, nil)
	w.start, w.running = tick, true
}

func (w *wipe) remember(screen *ebiten.Image) {
	view := screen.SubImage(image.Rect(0, HUDHeight, ScreenWidth, ScreenHeight)).(*ebiten.Image)
	w.last.Clear()
	w.last.DrawImage(view, nil)
}

func (w *wipe) draw(screen *ebiten.Image, tick int) {
	if !w.running {
		return
	}
	p := float64(tick-w.start) / wipeTicks
	if p >= 1 || p < 0 {
		w.running = false
		return
	}
	x := int(p * ScreenWidth)
	rest := w.from.SubImage(image.Rect(x, 0, ScreenWidth, ViewHeight)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), HUDHeight)
	screen.DrawImage(rest, op)
	drawWipeBand(screen, float64(x))
}

func drawWipeBand(dst *ebiten.Image, x float64) {
	fillRect(dst, geom.Rect{X: x - wipeBandWidth, Y: HUDHeight, W: wipeBandWidth, H: ViewHeight}, enemyColor, wipeShadowAlpha)
	for y := 0; y < ViewHeight; y += wipeStripeSize {
		c := enemyColor
		if (y/wipeStripeSize)%2 == 0 {
			c = fragmentColor
		}
		fillRect(dst, geom.Rect{X: x, Y: float64(HUDHeight + y), W: wipeBandWidth, H: wipeStripeSize}, c, 1)
	}
}
