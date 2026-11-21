package ebitenapp

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

func TestVisorSlidesAcrossWhenTurning(t *testing.T) {
	w := game.Watcher{Box: geom.Rect{X: 20, Y: 10, W: 8, H: 10}, Dir: 1}
	if v := visorRect(w, 0); v.X != w.Box.Right()-visorWidth-1 || v.W != visorWidth {
		t.Fatalf("facing right visor %+v", v)
	}
	w.Dir = -1
	if v := visorRect(w, 0); v.X != w.Box.X+1 {
		t.Fatalf("facing left visor %+v", v)
	}
	mid := visorRect(w, visorTurnTicks/2)
	if mid.W >= visorWidth || mid.X <= w.Box.X+1 || mid.Right() >= w.Box.Right()-1 {
		t.Fatalf("mid-turn visor %+v should be narrower and between the sides", mid)
	}
}

func TestWatcherFXTracksTurns(t *testing.T) {
	lvl, err := level.Parse(watchedLevel)
	if err != nil {
		t.Fatal(err)
	}
	clock := &rhythm.Manual{}
	s := game.NewSession(clock, lvl)
	var f watcherFX
	f.update(s)
	turned := false
	for range 300 {
		clock.Advance(time.Second / 60)
		before := s.Watchers()[0].Dir
		s.Update(input.Intent{})
		f.update(s)
		if s.Watchers()[0].Dir != before {
			turned = true
			break
		}
	}
	if !turned || f.turning(0) != visorTurnTicks {
		t.Fatalf("turned %v, turning %d", turned, f.turning(0))
	}
	s.Update(input.Intent{RestartPressed: true})
	f.update(s)
	if f.turning(0) != 0 || f.turning(5) != 0 {
		t.Fatal("restart kept a turn, or a missing watcher reported one")
	}
}

func TestTreadSightAndJitter(t *testing.T) {
	a := treadNotch(geom.Rect{X: 0, Y: 0, W: 8, H: 10})
	b := treadNotch(geom.Rect{X: treadStride, Y: 0, W: 8, H: 10})
	if a.X == b.X-treadStride || a.Y != 9 || b.Y != 9 {
		t.Fatalf("tread did not alternate: %+v %+v", a, b)
	}
	for tick := range 40 {
		x, ok := sightPulse(10, 50, 1, tick)
		if !ok || x < 10 || x+sightPulseLen > 50 {
			t.Fatalf("tick %d pulse at %v outside the ray", tick, x)
		}
		back, _ := sightPulse(10, 50, -1, tick)
		if back != 50-sightPulseLen-(x-10) {
			t.Fatalf("tick %d: left-facing pulse %v does not mirror %v", tick, back, x)
		}
	}
	if _, ok := sightPulse(10, 11, 1, 0); ok {
		t.Fatal("pulse on a ray shorter than the pulse")
	}
	if alertJitter(0.5, 0, false) != 0 || alertJitter(0.9, 0, true) != 0 {
		t.Fatal("jitter below the level or with reduced effects")
	}
	if alertJitter(0.9, 0, false) == alertJitter(0.9, jitterTicks, false) {
		t.Fatal("high alert did not jitter")
	}
}

func TestSpikeGlintSweepsTeethInOrder(t *testing.T) {
	for tooth := range 5 {
		at := tooth * glintStep
		if !glinting(tooth, at) || glinting(tooth, at+glintTicks) || !glinting(tooth, at+glintPeriod) {
			t.Fatalf("tooth %d glint timing wrong", tooth)
		}
	}
	tip := toothTip(geom.Rect{X: 0, Y: 5, W: 5, H: 5})
	if tip.W != 1 || tip.X != 2 || tip.Y != 5 {
		t.Fatalf("tip %+v, want the 1px top center", tip)
	}
}
