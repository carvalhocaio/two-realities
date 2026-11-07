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
	bobSpeed           = 0.08
	bobAmplitude       = 1.5
	secretTwinkleTicks = 10
	secretBurstRadius  = 18
	secretBurstSparks  = 8
	secretGlintPeriod  = 120
	secretGlintTicks   = 10
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

var secretRows = []int{1, 3, 6, 6, 3, 1}

func drawSecretGlint(dst *ebiten.Image, s *game.Session, tick int) {
	sec, ok := s.Secret()
	phase := tick % secretGlintPeriod
	if !ok || s.SecretTaken() || sec.Layer == s.Active() || phase >= secretGlintTicks {
		return
	}
	half := float64(secretGlintTicks) / 2
	alpha := float32(1 - math.Abs(float64(phase)-half)/half)
	c := centerOf(game.FragmentBox(sec.Point))
	x, y := math.Floor(c.X), math.Floor(c.Y)
	fillRect(dst, geom.Rect{X: x, Y: y, W: 1, H: 1}, secretColor, alpha)
	for _, d := range [][2]float64{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		fillRect(dst, geom.Rect{X: x + d[0], Y: y + d[1], W: 1, H: 1}, secretColor, alpha/2)
	}
}

func drawSecretBurst(dst *ebiten.Image, s *game.Session) {
	f := s.Feedback().SecretFound
	sec, ok := s.Secret()
	if f <= 0 || !ok {
		return
	}
	c := centerOf(game.FragmentBox(sec.Point))
	r := (1 - f) * secretBurstRadius
	for i := range secretBurstSparks {
		a := 2 * math.Pi * float64(i) / secretBurstSparks
		p := geom.Rect{X: math.Round(c.X + math.Cos(a)*r), Y: math.Round(c.Y + math.Sin(a)*r), W: 1, H: 1}
		fillRect(dst, p, fragmentColor, float32(f))
	}
}

func drawSecret(dst *ebiten.Image, s *game.Session, tick int) {
	box, ok := s.SecretBox()
	if !ok {
		return
	}
	box.Y += math.Round(math.Sin(float64(tick)*bobSpeed) * bobAmplitude)
	c := secretColor
	if (tick/secretTwinkleTicks)%2 == 0 {
		c = fragmentColor
	}
	drawRows(dst, box, secretRows, c, 1)
}
