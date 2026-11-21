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

const dropLevel = "fragments: 0\n[yellow]\n#######\n#P...E#\n#.....#\n#.....#\n#.....#\n#######\n[grey]\n#######\n#.....#\n#.....#\n#.....#\n#.....#\n#######\n"

func TestLandingSquashesAndRaisesDust(t *testing.T) {
	lvl, err := level.Parse(dropLevel)
	if err != nil {
		t.Fatal(err)
	}
	clock := &rhythm.Manual{}
	s := game.NewSession(clock, lvl)
	var fx playerFX
	landed := false
	for range 120 {
		clock.Advance(time.Second / 60)
		s.Update(input.Intent{})
		fx.update(s, s.Events())
		if fx.squash == squashTicks {
			landed = true
			break
		}
	}
	if !landed || len(fx.dust) != dustPerLanding {
		t.Fatalf("landing: squash=%d dust=%d", fx.squash, len(fx.dust))
	}
	body := fx.body(s.PlayerBox(), false)
	if body.W <= s.PlayerBox().W || body.H >= s.PlayerBox().H || body.Bottom() != s.PlayerBox().Bottom() {
		t.Fatalf("squashed body %+v from %+v", body, s.PlayerBox())
	}
	for range dustTicks {
		s.Update(input.Intent{Left: true})
		fx.update(s, s.Events())
	}
	if fx.squash != 0 || len(fx.dust) != 0 || fx.facing != -1 {
		t.Fatalf("after settling: squash=%d dust=%d facing=%d", fx.squash, len(fx.dust), fx.facing)
	}
}

const (
	floorLevel   = "fragments: 0\n[yellow]\n##########\n#........#\n#........#\n#P......E#\n##########\n[grey]\n##########\n#........#\n#........#\n#........#\n##########\n"
	watchedLevel = "fragments: 0\n[yellow]\n#########\n#.......#\n#...P..E#\n#########\n[grey]\n#########\n#.......#\n#.W.....#\n#########\n"
)

type playerRun struct {
	clock *rhythm.Manual
	s     *game.Session
	fx    playerFX
}

func newPlayerRun(t *testing.T, text string) *playerRun {
	t.Helper()
	lvl, err := level.Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	r := &playerRun{clock: &rhythm.Manual{}}
	r.s = game.NewSession(r.clock, lvl)
	r.fx.update(r.s, nil)
	return r
}

func (r *playerRun) step(in input.Intent, n int) {
	for range n {
		r.clock.Advance(time.Second / 60)
		r.s.Update(in)
		r.fx.update(r.s, r.s.Events())
	}
}

func TestJumpStretchesThePlayer(t *testing.T) {
	r := newPlayerRun(t, floorLevel)
	r.step(input.Intent{}, 5)
	r.step(input.Intent{JumpPressed: true}, 1)
	if r.fx.stretch != stretchTicks {
		t.Fatalf("stretch %d after a jump, want %d", r.fx.stretch, stretchTicks)
	}
	box := r.s.PlayerBox()
	body := r.fx.body(box, false)
	if body.H <= box.H || body.W >= box.W || body.Bottom() != box.Bottom() {
		t.Fatalf("jump body %+v from %+v", body, box)
	}
	r.step(input.Intent{}, stretchTicks)
	if r.fx.stretch != 0 {
		t.Fatalf("stretch %d outlived %d ticks", r.fx.stretch, stretchTicks)
	}
}

func TestFastFallStretchesAndLandingEndsTheJump(t *testing.T) {
	var fx playerFX
	box := geom.Rect{X: 10, Y: 10, W: 6, H: 9}
	fx.dy = fallStretchSpeed + 0.5
	if body := fx.body(box, false); body.H != box.H+1 || body.Bottom() != box.Bottom() {
		t.Fatalf("fast fall body %+v from %+v", body, box)
	}
	fx.stretch = stretchTicks
	fx.land(box)
	if fx.stretch != 0 || fx.squash != squashTicks {
		t.Fatalf("landing left stretch %d squash %d", fx.stretch, fx.squash)
	}
}

func TestRunStepsAndSkidsOnTurns(t *testing.T) {
	r := newPlayerRun(t, floorLevel)
	r.step(input.Intent{}, 5)
	dips := map[bool]int{}
	for range 3 * stepTicks {
		r.step(input.Intent{Right: true}, 1)
		dips[r.fx.stepDip(false)]++
		if r.fx.stepDip(true) {
			t.Fatal("reduced effects still bobbed the run")
		}
	}
	if dips[true] == 0 || dips[false] == 0 {
		t.Fatalf("run bob never alternated: %v", dips)
	}
	if len(r.fx.dust) != 0 {
		t.Fatalf("running raised %d specks", len(r.fx.dust))
	}
	for range 10 {
		r.step(input.Intent{Left: true}, 1)
		if len(r.fx.dust) > 0 {
			break
		}
	}
	if len(r.fx.dust) != dustPerTurn {
		t.Fatalf("turn raised %d specks, want %d", len(r.fx.dust), dustPerTurn)
	}
	r.step(input.Intent{}, 2)
	if r.fx.run != 0 || r.fx.stepDip(false) {
		t.Fatal("standing still kept the run cycle")
	}
}

func TestShortRunTurnsWithoutSkid(t *testing.T) {
	r := newPlayerRun(t, floorLevel)
	r.step(input.Intent{}, 5)
	r.step(input.Intent{Right: true}, 2)
	r.step(input.Intent{Left: true}, 10)
	if len(r.fx.dust) != 0 {
		t.Fatalf("a short step back raised %d specks", len(r.fx.dust))
	}
}

func TestEyeBlinksLooksAndWidens(t *testing.T) {
	var fx playerFX
	for _, tt := range []struct {
		age  int
		open bool
	}{{0, true}, {eyeBlinkPeriod - eyeBlinkTicks - 1, true}, {eyeBlinkPeriod - 1, false}, {eyeBlinkPeriod, true}} {
		fx.age = tt.age
		if fx.eyeOpen() != tt.open {
			t.Fatalf("age %d: open %v, want %v", tt.age, fx.eyeOpen(), tt.open)
		}
	}
	for _, tt := range []struct{ dy, want float64 }{{-3, -1}, {-0.5, 0}, {0, 0}, {1, 0}, {2, 1}} {
		fx.dy = tt.dy
		if got := fx.eyeLook(); got != tt.want {
			t.Fatalf("dy %v: look %v, want %v", tt.dy, got, tt.want)
		}
	}
	if fx.eyeHeight() != 1 {
		t.Fatal("calm eye is not 1px")
	}
	fx.alarmed = true
	if fx.eyeHeight() != 2 {
		t.Fatal("alarmed eye did not widen")
	}
}

func TestEyeWidensOnlyWhileAWatcherSeesYou(t *testing.T) {
	r := newPlayerRun(t, watchedLevel)
	r.step(input.Intent{}, 3)
	if r.fx.alarmed {
		t.Fatal("alarmed in yellow")
	}
	r.step(input.Intent{ShiftPressed: true}, 1)
	r.step(input.Intent{}, 2)
	if !r.fx.alarmed {
		t.Fatal("not alarmed while the watcher fills its meter")
	}
	r.step(input.Intent{}, 60)
	if r.fx.alarmed {
		t.Fatal("still alarmed after the respawn")
	}
}

func TestShiftLeavesAnEchoAndRestartResets(t *testing.T) {
	r := newPlayerRun(t, floorLevel)
	r.step(input.Intent{}, 30)
	r.step(input.Intent{ShiftPressed: true}, 1)
	e := r.fx.echo
	if e.left != echoTicks || e.box != r.s.PlayerBox() {
		t.Fatalf("echo %+v after shifting to grey", e)
	}
	first, firstAlpha := r.fx.echoRing()
	r.step(input.Intent{}, echoTicks/2)
	later, laterAlpha := r.fx.echoRing()
	if !(later.W > first.W && first.W > e.box.W) || laterAlpha >= firstAlpha {
		t.Fatalf("pulse should grow and fade: %+v a%v then %+v a%v", first, firstAlpha, later, laterAlpha)
	}
	r.step(input.Intent{}, echoTicks-echoTicks/2-1)
	if !r.fx.showEcho(false) || r.fx.showEcho(true) {
		t.Fatal("echo should show, except with reduced effects")
	}
	r.step(input.Intent{}, 1)
	if r.fx.showEcho(false) {
		t.Fatal("echo outlived its ticks")
	}
	if r.fx.age == 0 {
		t.Fatal("blink clock did not run")
	}
	r.step(input.Intent{RestartPressed: true}, 1)
	if r.fx.age != 0 || r.fx.echo.left != 0 || r.fx.stretch != 0 {
		t.Fatalf("restart kept effects: age %d echo %d stretch %d", r.fx.age, r.fx.echo.left, r.fx.stretch)
	}
}
