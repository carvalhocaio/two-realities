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

func spikeTeeth(p level.Point) []geom.Rect {
	box := game.SpikeBox(p)
	w := box.W / spikesPerTile
	teeth := make([]geom.Rect, spikesPerTile)
	for i := range teeth {
		teeth[i] = geom.Rect{X: box.X + float64(i)*w, Y: box.Y, W: w, H: box.H}
	}
	return teeth
}

func drawSpikes(dst *ebiten.Image, spikes []level.Point, c color.Color) {
	for _, p := range spikes {
		for _, tooth := range spikeTeeth(p) {
			drawRows(dst, tooth, spikeRows, c, 1)
		}
	}
}

func drawGhostSpikes(dst *ebiten.Image, spikes []level.Point, c color.Color) {
	for _, p := range spikes {
		for _, tooth := range spikeTeeth(p) {
			drawToothOutline(dst, tooth, c)
		}
	}
}

func drawToothOutline(dst *ebiten.Image, tooth geom.Rect, c color.Color) {
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
