package ebitenapp

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	bobSpeed     = 0.08
	bobAmplitude = 1.5
)

var diamondRows = []int{2, 4, 6, 6, 4, 2}

func drawRows(dst *ebiten.Image, box geom.Rect, rows []int, c color.Color, alpha float32) {
	for i, w := range rows {
		row := geom.Rect{
			X: box.X + math.Floor((box.W-float64(w))/2),
			Y: box.Y + float64(i),
			W: float64(w),
			H: 1,
		}
		fillRect(dst, row, c, alpha)
	}
}

func drawFragments(dst *ebiten.Image, fragments []level.Point, tick int, alpha float32) {
	for i, p := range fragments {
		box := game.FragmentBox(p)
		box.Y += math.Round(math.Sin(float64(tick)*bobSpeed+float64(i)) * bobAmplitude)
		drawRows(dst, box, diamondRows, fragmentColor, alpha)
	}
}
