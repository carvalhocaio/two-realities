package ebitenapp

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	pickupBurstTicks  = 10
	pickupBurstSparks = 6
	pickupBurstRadius = 9
	pickupBurstStart  = 2
	hoistTicks        = 8
	noCheckpoint      = -1
)

type burst struct {
	pos geom.Vec
	age int
}

type pickupFX struct {
	session    *game.Session
	level      int
	epoch      int
	collected  []bool
	bursts     []burst
	checkpoint int
	hoist      int
}

func (f *pickupFX) update(s *game.Session) {
	if s != f.session || s.LevelIndex() != f.level || s.Epoch() != f.epoch {
		f.sync(s)
		return
	}
	f.age()
	for i, p := range s.Level().Fragments {
		if s.FragmentCollected(i) && !f.collected[i] {
			f.collected[i] = true
			f.bursts = append(f.bursts, burst{pos: centerOf(game.FragmentBox(p))})
		}
	}
	f.hoist = max(f.hoist-1, 0)
	if cp := activeCheckpoint(s); cp != f.checkpoint {
		f.checkpoint = cp
		if cp != noCheckpoint {
			f.hoist = hoistTicks
		}
	}
}

func (f *pickupFX) sync(s *game.Session) {
	f.session, f.level, f.epoch = s, s.LevelIndex(), s.Epoch()
	f.collected = f.collected[:0]
	for i := range s.Level().Fragments {
		f.collected = append(f.collected, s.FragmentCollected(i))
	}
	f.bursts = f.bursts[:0]
	f.checkpoint, f.hoist = activeCheckpoint(s), 0
}

func (f *pickupFX) age() {
	live := f.bursts[:0]
	for _, b := range f.bursts {
		if b.age++; b.age < pickupBurstTicks {
			live = append(live, b)
		}
	}
	f.bursts = live
}

func (f *pickupFX) hoistProgress() float64 {
	return 1 - float64(f.hoist)/hoistTicks
}

func activeCheckpoint(s *game.Session) int {
	if i, ok := s.ActiveCheckpoint(); ok {
		return i
	}
	return noCheckpoint
}

func (f *pickupFX) draw(dst *ebiten.Image) {
	for _, b := range f.bursts {
		t := float64(b.age) / pickupBurstTicks
		drawSparkRing(dst, b.pos, pickupBurstStart+t*(pickupBurstRadius-pickupBurstStart), pickupBurstSparks, float32(1-t))
	}
}

func drawSparkRing(dst *ebiten.Image, c geom.Vec, r float64, sparks int, alpha float32) {
	for i := range sparks {
		a := 2 * math.Pi * float64(i) / float64(sparks)
		p := geom.Rect{X: math.Round(c.X + math.Cos(a)*r), Y: math.Round(c.Y + math.Sin(a)*r), W: 1, H: 1}
		fillRect(dst, p, fragmentColor, alpha)
	}
}
