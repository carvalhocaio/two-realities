package ebitenapp

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const spikesPerTile = 2

var spikeRows = []int{1, 1, 3, 3, 5}

func spikeTeeth(p level.Point) [spikesPerTile]geom.Rect {
	box := game.SpikeBox(p)
	w := box.W / spikesPerTile
	var teeth [spikesPerTile]geom.Rect
	for i := range teeth {
		teeth[i] = geom.Rect{X: box.X + float64(i)*w, Y: box.Y, W: w, H: box.H}
	}
	return teeth
}

func drawSpikes(dst *ebiten.Image, spikes []level.Point, c color.RGBA) {
	for _, p := range spikes {
		for _, tooth := range spikeTeeth(p) {
			drawRows(dst, tooth, spikeRows, c, 1)
		}
	}
}

func drawGhostSpikes(dst *ebiten.Image, spikes []level.Point, c color.RGBA) {
	for _, p := range spikes {
		for _, tooth := range spikeTeeth(p) {
			drawToothOutline(dst, tooth, c)
		}
	}
}

func drawToothOutline(dst *ebiten.Image, tooth geom.Rect, c color.RGBA) {
	for i, w := range spikeRows {
		left := tooth.X + float64(int(tooth.W)-w)/2
		y := tooth.Y + float64(i)
		if i == len(spikeRows)-1 {
			fillRect(dst, geom.Rect{X: left, Y: y, W: float64(w), H: 1}, c, ghostOutlineAlpha)
			continue
		}
		fillRect(dst, geom.Rect{X: left, Y: y, W: 1, H: 1}, c, ghostOutlineAlpha)
		fillRect(dst, geom.Rect{X: left + float64(w) - 1, Y: y, W: 1, H: 1}, c, ghostOutlineAlpha)
	}
}

const (
	glintPeriod = 100
	glintStep   = 3
	glintTicks  = 2
	glintAlpha  = 0.9
)

func drawSpikeGlints(dst *ebiten.Image, spikes []level.Point, tick int) {
	k := 0
	for _, p := range spikes {
		for _, tooth := range spikeTeeth(p) {
			if glinting(k, tick) {
				fillRect(dst, toothTip(tooth), beamCoreColor, glintAlpha)
			}
			k++
		}
	}
}

func glinting(tooth, tick int) bool {
	phase := ((tick-tooth*glintStep)%glintPeriod + glintPeriod) % glintPeriod
	return phase < glintTicks
}

func toothTip(tooth geom.Rect) geom.Rect {
	return geom.Rect{X: tooth.X + float64(int(tooth.W)-spikeRows[0])/2, Y: tooth.Y, W: float64(spikeRows[0]), H: 1}
}
