package ebitenapp

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	squashTicks      = 7
	squashWiden      = 2
	squashLower      = 2
	stretchTicks     = 6
	stretchTaller    = 2
	stretchNarrower  = 1
	fallStretchSpeed = 3.0
	stepTicks        = 6
	turnRunTicks     = 8
	dustPerTurn      = 2
	eyeBlinkPeriod   = 200
	eyeBlinkTicks    = 4
	lookUpSpeed      = -1.0
	lookDownSpeed    = 1.5
	echoTicks        = 12
	echoAlpha        = 0.7
	echoGrow         = 3
	dustTicks        = 14
	dustSpeed        = 0.6
	dustLift         = 0.15
	dustSize         = 1
	dustAlpha        = 0.8
	playerEyeY       = 2
	playerEyeWidth   = 2
	playerEyeInset   = 1
	dustPerLanding   = 3
)

type dust struct {
	pos, vel geom.Vec
	age      int
}

type echo struct {
	box  geom.Rect
	left int
}

type playerFX struct {
	prev    geom.Rect
	dy      float64
	facing  int
	squash  int
	stretch int
	run     int
	runDir  int
	age     int
	alarmed bool
	echo    echo
	dust    []dust
	epoch   int
	session *game.Session
}

func (f *playerFX) update(s *game.Session, events []game.Event) {
	box := s.PlayerBox()
	if s != f.session || s.Epoch() != f.epoch {
		*f = playerFX{prev: box, facing: 1, epoch: s.Epoch(), session: s, dust: f.dust[:0]}
		return
	}
	dx := box.X - f.prev.X
	f.dy = box.Y - f.prev.Y
	f.face(dx)
	f.age++
	f.squash = max(f.squash-1, 0)
	f.stretch = max(f.stretch-1, 0)
	f.echo.left = max(f.echo.left-1, 0)
	f.track(box, dx, s.PlayerGrounded())
	for _, e := range events {
		f.react(e, box)
	}
	f.alarmed = watched(s)
	f.prev = box
	f.ageDust()
}

func (f *playerFX) face(dx float64) {
	switch {
	case dx > 0:
		f.facing = 1
	case dx < 0:
		f.facing = -1
	}
}

func (f *playerFX) track(box geom.Rect, dx float64, grounded bool) {
	dir := sign(dx)
	if !grounded || dir == 0 {
		f.run, f.runDir = 0, 0
		return
	}
	if dir != f.runDir {
		if f.runDir != 0 && f.run >= turnRunTicks {
			f.skid(box, dir)
		}
		f.run, f.runDir = 0, dir
	}
	f.run++
}

func (f *playerFX) react(e game.Event, box geom.Rect) {
	switch e {
	case game.EventJump:
		f.stretch, f.squash = stretchTicks, 0
	case game.EventLand:
		f.land(box)
	case game.EventShiftGrey, game.EventShiftYellow:
		f.echo = echo{box: box, left: echoTicks}
	}
}

func watched(s *game.Session) bool {
	if s.Active() != world.Grey {
		return false
	}
	for _, w := range s.Watchers() {
		if w.Alert > 0 {
			return true
		}
	}
	return false
}

func sign(v float64) int {
	switch {
	case v > 0:
		return 1
	case v < 0:
		return -1
	}
	return 0
}

func (f *playerFX) land(box geom.Rect) {
	f.squash, f.stretch = squashTicks, 0
	feet := geom.Vec{X: box.X + box.W/2, Y: box.Bottom() - 1}
	for i := range dustPerLanding {
		dir := float64(i - 1)
		if dir == 0 {
			dir = float64(f.facing) * 0.5
		}
		f.dust = append(f.dust, dust{pos: feet, vel: geom.Vec{X: dir * dustSpeed, Y: -dustLift}})
	}
}

func (f *playerFX) skid(box geom.Rect, dir int) {
	feet := geom.Vec{X: box.X + box.W/2 - float64(dir)*box.W/2, Y: box.Bottom() - 1}
	for i := range dustPerTurn {
		speed := dustSpeed * float64(i+1) / dustPerTurn
		f.dust = append(f.dust, dust{pos: feet, vel: geom.Vec{X: -float64(dir) * speed, Y: -dustLift}})
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

func (f *playerFX) body(box geom.Rect, reduced bool) geom.Rect {
	w, h := box.W, box.H
	switch {
	case f.squash > 0:
		k := float64(f.squash) / squashTicks
		w, h = box.W+squashWiden*k*2, box.H-squashLower*k*2
	case f.stretch > 0:
		k := float64(f.stretch) / stretchTicks
		w, h = box.W-stretchNarrower*k*2, box.H+stretchTaller*k
	case f.dy > fallStretchSpeed:
		h++
	}
	if f.stepDip(reduced) {
		h--
	}
	return geom.Rect{X: box.X + (box.W-w)/2, Y: box.Bottom() - h, W: w, H: h}
}

func (f *playerFX) stepDip(reduced bool) bool {
	return !reduced && f.run > 0 && (f.run/stepTicks)%2 == 1
}

func (f *playerFX) eyeOpen() bool {
	return f.age%eyeBlinkPeriod < eyeBlinkPeriod-eyeBlinkTicks
}

func (f *playerFX) eyeLook() float64 {
	switch {
	case f.dy < lookUpSpeed:
		return -1
	case f.dy > lookDownSpeed:
		return 1
	}
	return 0
}

func (f *playerFX) eyeHeight() float64 {
	if f.alarmed {
		return 2
	}
	return 1
}

func (f *playerFX) showEcho(reduced bool) bool {
	return !reduced && f.echo.left > 0
}

func (f *playerFX) echoRing() (geom.Rect, float32) {
	t := 1 - float64(f.echo.left)/echoTicks
	grow := math.Round(1 + t*(echoGrow-1))
	b := f.echo.box
	ring := geom.Rect{X: b.X - grow, Y: b.Y - grow, W: b.W + 2*grow, H: b.H + 2*grow}
	return ring, float32(echoAlpha * (1 - t))
}

func (f *playerFX) draw(dst *ebiten.Image, box geom.Rect, pal palette, visible, reduced bool) {
	for _, d := range f.dust {
		alpha := float32(dustAlpha * (1 - float64(d.age)/dustTicks))
		fillRect(dst, geom.Rect{X: d.pos.X, Y: d.pos.Y, W: dustSize, H: dustSize}, pal.tile, alpha)
	}
	if f.showEcho(reduced) {
		ring, alpha := f.echoRing()
		drawOutline(dst, ring, pal.text, alpha)
	}
	if !visible {
		return
	}
	body := f.body(box, reduced)
	body = geom.Rect{X: roundHalf(body.X), Y: roundHalf(body.Y), W: roundHalf(body.W), H: roundHalf(body.H)}
	fillRect(dst, body, pal.player, 1)
	if !f.eyeOpen() && !f.alarmed {
		return
	}
	eyeX := body.Right() - playerEyeInset - playerEyeWidth
	if f.facing < 0 {
		eyeX = body.X + playerEyeInset
	}
	eye := geom.Rect{X: eyeX, Y: body.Y + playerEyeY + f.eyeLook(), W: playerEyeWidth, H: f.eyeHeight()}
	fillRect(dst, eye, pal.background, 1)
}

func roundHalf(v float64) float64 { return float64(int(v + 0.5)) }
