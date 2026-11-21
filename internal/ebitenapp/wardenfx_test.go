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

const wardenFXArena = "fragments: 2\nboss: warden\n" +
	"[yellow]\n" +
	"##############\n" +
	"#............#\n" +
	"#............#\n" +
	"#.....1..2...#\n" +
	"#.....1..2...#\n" +
	"#.....1..2...#\n" +
	"#PLE..1L.2L..#\n" +
	"##############\n" +
	"[grey]\n" +
	"##############\n" +
	"#OOOOOOOOOOOO#\n" +
	"#............#\n" +
	"#............#\n" +
	"#...#........#\n" +
	"#............#\n" +
	"#............#\n" +
	"##############\n"

func wardenBoss(t *testing.T) *level.Boss {
	t.Helper()
	lvl, err := level.Parse(wardenFXArena)
	if err != nil {
		t.Fatal(err)
	}
	return lvl.Boss
}

func TestWardenFXLampLightsShockAndOpensGate(t *testing.T) {
	boss := wardenBoss(t)
	f := wardenFX{gate: noGate}
	f.step(boss, game.WardenView{Lit: 1})
	if f.shock != shockTicks || f.gate != 0 || f.opening != gateOpenTicks {
		t.Fatalf("after lamp 1: shock %d gate %d opening %d", f.shock, f.gate, f.opening)
	}
	tile := level.TileRect(boss.Lamps[0])
	if f.lamp != (geom.Vec{X: tile.X + 5, Y: tile.Y + 3}) {
		t.Fatalf("shock centered at %v, want lamp 1 bulb", f.lamp)
	}
	for range gateOpenTicks {
		f.step(boss, game.WardenView{Lit: 1})
	}
	if f.shock != 0 || f.opening != 0 || f.gateProgress() != 1 {
		t.Fatalf("effects outlived their ticks: shock %d opening %d", f.shock, f.opening)
	}
	f.step(boss, game.WardenView{Lit: 2})
	if f.gate != 1 {
		t.Fatalf("lamp 2 opened gate %d, want 1", f.gate)
	}
	f.step(boss, game.WardenView{Lit: 3})
	if f.gate != noGate || f.opening != 0 || f.shock != shockTicks {
		t.Fatalf("lamp 3 has no gate: gate %d opening %d shock %d", f.gate, f.opening, f.shock)
	}
}

func TestWardenFXRecoilsOnlyWhenATrackedShotFires(t *testing.T) {
	boss := wardenBoss(t)
	f := wardenFX{gate: noGate}
	fire := func(stage game.BeamStage) {
		f.step(boss, game.WardenView{Lit: 2, Beam: game.Beam{Stage: stage}})
	}
	fire(game.BeamFiring)
	fire(game.BeamFiring)
	if f.recoil != 0 {
		t.Fatal("a sweep that keeps firing recoiled")
	}
	fire(game.BeamTelegraph)
	fire(game.BeamFiring)
	if f.recoil != recoilTicks || f.eyeLift(false) != recoilPixels || f.eyeLift(true) != 0 {
		t.Fatalf("recoil %d, lift %v / reduced %v", f.recoil, f.eyeLift(false), f.eyeLift(true))
	}
	for range recoilTicks {
		fire(game.BeamFiring)
	}
	if f.eyeLift(false) != 0 {
		t.Fatal("recoil outlived its ticks")
	}
	if coreAlpha(true) != 1 || coreAlpha(false) != beamCoreAlpha {
		t.Fatal("beam core does not snap to full on the recoil")
	}
}

func TestWardenFXFollowsTheSession(t *testing.T) {
	lvl, err := level.Parse(wardenFXArena)
	if err != nil {
		t.Fatal(err)
	}
	clock := &rhythm.Manual{}
	s := game.NewSession(clock, lvl, lvl)
	var f wardenFX
	tick := func(in input.Intent) {
		clock.Advance(time.Second / 60)
		s.Update(in)
		f.update(s)
	}
	tick(input.Intent{ShiftPressed: true})
	for range 20 {
		if w, _ := s.Warden(); w.Lit == 1 {
			break
		}
		tick(input.Intent{Right: true})
	}
	if w, _ := s.Warden(); w.Lit != 1 || f.shock == 0 || f.gate != 0 {
		t.Fatalf("lit %d shock %d gate %d", w.Lit, f.shock, f.gate)
	}
	tick(input.Intent{RestartPressed: true})
	if f.shock != 0 || f.opening != 0 || f.lit != 0 {
		t.Fatalf("restart kept effects: shock %d opening %d lit %d", f.shock, f.opening, f.lit)
	}
	plain, err := level.Parse("fragments: 0\n[yellow]\n#....#\n#P..E#\n######\n[grey]\n#....#\n#....#\n######\n")
	if err != nil {
		t.Fatal(err)
	}
	f.update(game.NewSession(clock, plain))
	if f.session != nil || f.gate != noGate {
		t.Fatal("a level without the Warden kept Warden effects")
	}
}

func TestGateColumnsMergeStackedTiles(t *testing.T) {
	cols := gateColumns([]level.Point{{X: 6, Y: 3}, {X: 6, Y: 4}, {X: 6, Y: 5}, {X: 9, Y: 5}})
	want := []geom.Rect{{X: 60, Y: 30, W: 10, H: 30}, {X: 90, Y: 50, W: 10, H: 10}}
	if len(cols) != len(want) || cols[0] != want[0] || cols[1] != want[1] {
		t.Fatalf("columns %v, want %v", cols, want)
	}
}
