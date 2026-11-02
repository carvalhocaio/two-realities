package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	squashTicks    = 7
	squashWiden    = 2
	squashLower    = 2
	landingSpeed   = 2.0
	dustTicks      = 14
	dustSpeed      = 0.6
	dustLift       = 0.15
	dustSize       = 1
	dustAlpha      = 0.8
	playerEyeY     = 2
	playerEyeWidth = 2
	playerEyeInset = 1
	dustPerLanding = 3
)

type dust struct {
	pos, vel geom.Vec
	age      int
}

type playerFX struct {
	prev     geom.Rect
	fallDY   float64
	grounded bool
	facing   int
	squash   int
	dust     []dust
	epoch    int
	session  *game.Session
}

func (f *playerFX) update(s *game.Session) {
	box := s.PlayerBox()
	if s != f.session || s.Epoch() != f.epoch {
		*f = playerFX{prev: box, grounded: true, facing: 1, epoch: s.Epoch(), session: s}
		return
	}
	switch dx := box.X - f.prev.X; {
	case dx > 0:
		f.facing = 1
	case dx < 0:
		f.facing = -1
	}
	f.squash = max(f.squash-1, 0)
	grounded := s.PlayerGrounded()
	if grounded && !f.grounded && f.fallDY >= landingSpeed {
		f.land(box)
	}
	f.fallDY = box.Y - f.prev.Y
	f.grounded, f.prev = grounded, box
	f.ageDust()
}

func (f *playerFX) land(box geom.Rect) {
	f.squash = squashTicks
	feet := geom.Vec{X: box.X + box.W/2, Y: box.Bottom() - 1}
	for i := range dustPerLanding {
		dir := float64(i - 1)
		if dir == 0 {
			dir = float64(f.facing) * 0.5
		}
		f.dust = append(f.dust, dust{pos: feet, vel: geom.Vec{X: dir * dustSpeed, Y: -dustLift}})
	}
}

func (f *playerFX) ageDust() {
	kept := f.dust[:0]
	for _, d := range f.dust {
		d.age++
		d.pos.X += d.vel.X
		d.pos.Y += d.vel.Y
		if d.age < dustTicks {
			kept = append(kept, d)
		}
	}
	f.dust = kept
}

func (f *playerFX) body(box geom.Rect) geom.Rect {
	if f.squash == 0 {
		return box
	}
	k := float64(f.squash) / squashTicks
	w, h := box.W+squashWiden*k*2, box.H-squashLower*k*2
	return geom.Rect{X: box.X + (box.W-w)/2, Y: box.Bottom() - h, W: w, H: h}
}

func (f *playerFX) draw(dst *ebiten.Image, box geom.Rect, pal palette, visible bool) {
	for _, d := range f.dust {
		alpha := float32(dustAlpha * (1 - float64(d.age)/dustTicks))
		fillRect(dst, geom.Rect{X: d.pos.X, Y: d.pos.Y, W: dustSize, H: dustSize}, pal.tile, alpha)
	}
	if !visible {
		return
	}
	body := f.body(box)
	body = geom.Rect{X: roundHalf(body.X), Y: roundHalf(body.Y), W: roundHalf(body.W), H: roundHalf(body.H)}
	fillRect(dst, body, pal.player, 1)
	eyeX := body.Right() - playerEyeInset - playerEyeWidth
	if f.facing < 0 {
		eyeX = body.X + playerEyeInset
	}
	fillRect(dst, geom.Rect{X: eyeX, Y: body.Y + playerEyeY, W: playerEyeWidth, H: 1}, pal.background, 1)
}

func roundHalf(v float64) float64 { return float64(int(v + 0.5)) }
