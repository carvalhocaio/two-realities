package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

func TestSessionAccessors(t *testing.T) {
	s := newSession(t, timedLevel)
	if len(s.Levels()) != 2 || s.Records() == nil {
		t.Fatal("levels or records missing")
	}
	if _, ok := s.Warden(); ok {
		t.Fatal("a normal level has no warden")
	}
	if s.DeathProgress() != 0 {
		t.Fatal("alive means no death progress")
	}
	play(s, ticks(right, 5))
	if s.RunTicks() != 5 {
		t.Fatalf("run ticks %d, want 5", s.RunTicks())
	}
	s.die(DeathSpikes)
	s.die(DeathFall)
	if s.Death() != DeathSpikes {
		t.Fatalf("a second death replaced the first: %v", s.Death())
	}
	s.Update(idle)
	if p := s.DeathProgress(); p <= 0 || p >= 1 {
		t.Fatalf("death progress %v, want between 0 and 1", p)
	}
}

func TestWatcherAccessors(t *testing.T) {
	w := NewWatcher(level.Point{X: 2, Y: 1})
	w.Alert = 6
	if got := w.AlertLevel(12); got != 0.5 {
		t.Fatalf("alert level %v, want 0.5", got)
	}
	if o := w.SightOrigin(); o.X != w.Box.Right() || o.Y != w.EyeY() {
		t.Fatalf("sight origin %v", o)
	}
	w.Dir = -1
	if o := w.SightOrigin(); o.X != w.Box.X {
		t.Fatalf("sight origin facing left %v", o)
	}
}

func TestTrappedWatcherStandsStill(t *testing.T) {
	lvl, err := level.Parse("fragments: 0\n[yellow]\n#....#\n#P..E#\n######\n[grey]\n#....#\n##W###\n######\n")
	if err != nil {
		t.Fatal(err)
	}
	w := NewWatcher(lvl.Watchers[0])
	before := w.Box
	w.Patrol(lvl.Layer(1), 2)
	if w.Box != before {
		t.Fatalf("trapped watcher moved to %v", w.Box)
	}
}

func TestNearestFreeGivesUpWhenBuried(t *testing.T) {
	lvl, err := level.Parse("fragments: 0\n[yellow]\n#########\n#########\n#########\n####P####\n####E####\n#########\n#########\n#########\n[grey]\n#########\n#########\n#########\n#########\n#########\n#########\n#########\n#########\n")
	if err != nil {
		t.Skip("level parser rejects a buried spawn:", err)
	}
	box := geom.Rect{X: 40, Y: 30, W: 6, H: 9}
	if _, ok := nearestFree(box, lvl.Layer(1)); ok {
		t.Fatal("found room inside solid rock")
	}
}

func TestPressureGoesDownAndStopsAtZero(t *testing.T) {
	f := seenFlow(&rhythm.Manual{}, &save.Memory{Stored: save.Data{PressureUnlocked: true}}, mustParse(t, timedLevel))
	for _, in := range then(ticks(moreKey, 2), ticks(lessKey, 3)) {
		f.Update(in)
	}
	if f.Pressure() != 0 {
		t.Fatalf("pressure %d, want 0", f.Pressure())
	}
	f.Update(moreKey)
	if f.Pressure() != 1 {
		t.Fatalf("pressure %d, want 1", f.Pressure())
	}
}

func TestWardenView(t *testing.T) {
	r := newRunner(t, "13_warden")
	w, ok := r.s.Warden()
	if !ok || w.Lit != 0 || w.Hit != 0 {
		t.Fatalf("warden view %+v %v at the start of the fight", w, ok)
	}
}
