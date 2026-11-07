package game

import (
	"fmt"
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const airShiftTicks = 14

type secretLedges struct {
	jumpFromX  float64
	ledgeTop   float64
	ledgeEdgeX float64
	greyTop    float64
}

func ledgesOf(lvl level.Level, sec level.Marker) secretLedges {
	ts := float64(level.TileSize)
	grey, yellow := lvl.Layer(sec.Layer), lvl.Layer(world.Yellow)
	start := sec.X
	for grey.Solid(start-1, sec.Y+1) {
		start--
	}
	row := sec.Y + 1
	for !yellow.Solid(start-2, row) {
		row++
	}
	return secretLedges{
		jumpFromX:  float64(start-3)*ts - 2*ts,
		ledgeTop:   float64(row) * ts,
		ledgeEdgeX: float64(start-1)*ts - player.Width - 1,
		greyTop:    float64(sec.Y+1) * ts,
	}
}

func (r *runner) ticksToBeat() float64 {
	return (1 - r.clock.Phase()) * float64(beatTicks)
}

func (r *runner) grabSecret() {
	r.t.Helper()
	sec, _ := r.s.Secret()
	l := ledgesOf(r.s.Level(), sec)
	r.until(r.reachX(l.jumpFromX), right)
	r.step(jumpRight)
	r.until(func() bool { return r.s.PlayerBox().Bottom() == l.ledgeTop && r.x() >= l.ledgeEdgeX }, right)

	switch r.s.Level().Shift {
	case level.ShiftFree:
		r.step(jumpRight)
		r.until(func() bool { return r.s.PlayerBox().Bottom() <= l.greyTop }, right)
		r.step(input.Intent{Right: true, ShiftPressed: true})
	case level.ShiftQuantized:
		r.until(func() bool { return r.ticksToBeat() <= airShiftTicks }, idle)
		r.step(input.Intent{Right: true, JumpPressed: true, ShiftPressed: true})
	case level.ShiftAuto:
		r.until(func() bool { return r.s.AutoCountdown() == 1 && r.ticksToBeat() <= airShiftTicks }, idle)
		r.step(jumpRight)
	}
	r.until(func() bool { return r.s.Active() == world.Grey }, right)
	r.until(r.s.SecretTaken, right)
}

func toLeftWall(r *runner) { r.until(r.backTo(10.5), walkLeft) }

func spikesUntilLastFragment(r *runner) {
	r.until(r.reachX(60), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
	r.until(r.reachX(232), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
}

func (r *runner) onFloor(bottom float64) func() bool {
	return func() bool { return r.s.player.Grounded() && r.s.PlayerBox().Bottom() == bottom }
}

func (r *runner) spendFuel() {
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
}

func finishRight(r *runner)  { r.finish(right) }
func finishLeft(r *runner)   { r.finish(walkLeft) }
func nothingFirst(r *runner) {}
func (r *runner) dropToFloorThenLeft(bottom float64) {
	r.until(r.onFloor(bottom), right)
	r.finish(walkLeft)
}

var secretRoutes = []struct {
	level   string
	spot    int
	prelude func(r *runner)
	after   func(r *runner)
}{
	{"01_first_shift", 0, toLeftWall, func(r *runner) {
		r.until(func() bool { return r.x() >= 80 && r.onFloor(150)() }, right)
		r.step(jumpRight)
		r.finish(right)
	}},
	{"02_fragments", 0, toLeftWall, func(r *runner) {
		r.until(r.reachX(254), right)
		r.shiftTo(world.Yellow)
		r.finish(right)
	}},
	{"02_fragments", 1, func(r *runner) {
		r.until(r.reachX(80), right)
		r.shiftTo(world.Grey)
		r.until(r.hasFragments(1), right)
		r.until(r.reachX(254), right)
		r.shiftTo(world.Yellow)
	}, finishRight},
	{"03_spikes", 0, spikesUntilLastFragment, func(r *runner) { r.dropToFloorThenLeft(150) }},
	{"04_watcher", 0, func(r *runner) { r.dodgeCorridor(); r.spendFuel() }, finishRight},
	{"04_watcher", 1, func(r *runner) { r.dodgeCorridor(); r.spendFuel() }, finishRight},
	{"05_planning", 0, func(r *runner) { r.dodgeCorridor(); r.spendFuel() }, finishRight},
	{"05_planning", 1, func(r *runner) {
		r.dodgeCorridor()
		r.spendFuel()
		r.until(r.reachX(430), right)
		r.shiftTo(world.Grey)
		r.spendFuel()
	}, func(r *runner) { r.dropToFloorThenLeft(150) }},
	{"06_finale", 0, func(r *runner) { finaleUntilLastYellow(r); r.until(r.backTo(45), walkLeft) }, finishLeft},
	{"06_finale", 1, func(r *runner) { finaleUntilLastYellow(r); r.until(r.backTo(140), walkLeft) }, finishLeft},
	{"07_crossfire", 0, func(r *runner) {
		crossfireUntilLastJump(r)
		r.until(r.hasFragments(1), right)
		r.until(r.onFloor(150), idle)
		r.shiftTo(world.Yellow)
		r.until(r.backTo(470), walkLeft)
	}, finishRight},
	{"08_descent", 0, nothingFirst, func(r *runner) {
		r.until(r.onFloor(150), right)
		r.shiftTo(world.Yellow)
		solveDescent(r)
	}},
	{"09_on_the_beat", 0, func(r *runner) {
		onTheBeatUntilGapJump(r)
		r.until(r.hasFragments(1), right)
		r.queueShift(world.Yellow)
		r.until(r.backTo(350), walkLeft)
	}, finishRight},
	{"10_metronome", 0, nothingFirst, func(r *runner) {
		r.waitFor(world.Yellow)
		r.until(func() bool { return r.s.player.Grounded() }, idle)
		solveMetronome(r)
	}},
	{"10_metronome", 1, func(r *runner) {
		r.until(r.reachX(158), right)
		r.waitFor(world.Grey)
		r.until(r.reachX(285), right)
		r.freshWindow(world.Yellow)
	}, finishRight},
	{"11_offbeat", 0, func(r *runner) {
		for _, st := range []struct {
			x float64
			w world.World
		}{{82, world.Grey}, {162, world.Yellow}, {242, world.Grey}, {322, world.Yellow}} {
			r.until(r.reachX(st.x), right)
			r.queueAndWait(st.w)
		}
	}, func(r *runner) {
		r.until(func() bool { return r.w(0).Dir > 0 && r.w(0).Box.X > 480 }, idle)
		r.until(r.reachX(451), right)
		r.until(r.onFloor(150), idle)
		r.until(r.reachX(535), right)
		r.until(func() bool { return r.w(0).Box.X > 555 }, idle)
		r.queueAndWait(world.Yellow)
		r.until(func() bool { return r.w(0).Dir < 0 && r.w(0).Box.Right() < r.x()-level.TileSize }, idle)
		r.queueAndWait(world.Grey)
		r.until(r.reachX(572), right)
		r.step(jumpRight)
		r.finish(right)
	}},
	{"12_clockwork", 0, toLeftWall, solveClockwork},
}

func TestEverySecretSpotCanBeTakenAndTheLevelFinished(t *testing.T) {
	levels, err := level.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	spots := 0
	for _, l := range levels {
		spots += len(l.Secrets)
	}
	if len(secretRoutes) != spots {
		t.Fatalf("%d secret routes for %d spots", len(secretRoutes), spots)
	}
	for _, tt := range secretRoutes {
		t.Run(fmt.Sprintf("%s/%d", tt.level, tt.spot), func(t *testing.T) {
			r := newRunner(t, tt.level)
			r.s.UseSecretPicker(func(int) int { return tt.spot })
			tt.prelude(r)
			before := r.s.Fragments()
			r.grabSecret()
			if r.s.Fragments() != before+1 {
				t.Fatalf("fragments %d, want %d", r.s.Fragments(), before+1)
			}
			tt.after(r)
		})
	}
}

func TestSecretsNeedTheShift(t *testing.T) {
	levels, err := level.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, lvl := range levels {
		if lvl.Boss != nil {
			continue
		}
		for _, sec := range lvl.Secrets {
			g := lvl.Layer(sec.Layer)
			other := lvl.Layer(sec.Layer.Other())
			if !g.Solid(sec.X, sec.Y+1) || other.Solid(sec.X, sec.Y+1) {
				t.Errorf("%s: secret should rest on a ledge that exists only in its own layer", lvl.Name)
			}
			for x := sec.X - 3; x <= sec.X+2; x++ {
				for y := sec.Y + 2; y <= sec.Y+3; y++ {
					if g.Solid(x, y) && x > 0 && x < g.Width()-1 {
						t.Errorf("%s: %v has a stepping stone at %d,%d", lvl.Name, sec.Layer, x, y)
					}
				}
			}
		}
	}
}

func bossLevels(levels []level.Level) int {
	n := 0
	for _, l := range levels {
		if l.Boss != nil {
			n++
		}
	}
	return n
}
