package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

func newSolidImage(lvl level.Level, w world.World) *ebiten.Image {
	b := lvl.Bounds()
	img := ebiten.NewImage(int(b.W), int(b.H))
	drawLayer(img, lvl.Layer(w), palettes[w].tile)
	drawSpikes(img, lvl.Spikes[w], palettes[w].spike)
	return img
}

func deallocate(images ...*ebiten.Image) {
	for _, img := range images {
		if img != nil {
			img.Deallocate()
		}
	}
}
