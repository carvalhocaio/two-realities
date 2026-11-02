package game

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	maxTicks      = 60 * 60
	maxLevelTicks = 90 * 60
	tickDuration  = time.Second / 60
)

var (
	jumpRight = input.Intent{Right: true, JumpPressed: true}
	walkLeft  = input.Intent{Left: true}
)

type runner struct {
	t     *testing.T
	s     *Session
	clock *rhythm.Manual
	ticks int
}

func newRunner(t *testing.T, name string) *runner {
	t.Helper()
	return newRunnerAt(t, name, 0)
}

func newRunnerAt(t *testing.T, name string, offset time.Duration) *runner {
	t.Helper()
	lvl, err := level.Load(name)
	if err != nil {
		t.Fatal(err)
	}
	clock := &rhythm.Manual{Elapsed: offset}
	return &runner{t: t, s: NewSession(clock, lvl, lvl), clock: clock}
}

func (r *runner) update(in input.Intent) {
	r.clock.Advance(tickDuration)
	r.s.Update(in)
}

func (r *runner) step(in input.Intent) {
	r.t.Helper()
	r.update(in)
	r.ticks++
	if r.s.Death() != Alive {
		r.t.Fatalf("died (%v) at %+v after %d ticks", r.s.Death(), r.s.PlayerBox(), r.ticks)
	}
}

func (r *runner) until(cond func() bool, in input.Intent) {
	r.t.Helper()
	for range maxTicks {
		if cond() {
			return
		}
		r.step(in)
	}
	r.t.Fatalf("condition never met; player at %+v", r.s.PlayerBox())
}

func (r *runner) shiftTo(w world.World) {
	r.t.Helper()
	r.step(input.Intent{ShiftPressed: true})
	if r.s.Active() != w {
		r.t.Fatalf("shift to %v failed at %+v", w, r.s.PlayerBox())
	}
}

func (r *runner) x() float64                   { return r.s.PlayerBox().X }
func (r *runner) reachX(x float64) func() bool { return func() bool { return r.x() >= x } }
func (r *runner) backTo(x float64) func() bool { return func() bool { return r.x() <= x } }
func (r *runner) hasFragments(n int) func() bool {
	return func() bool { return r.s.Fragments() == n }
}

func (r *runner) finish(in input.Intent) {
	r.t.Helper()
	r.until(func() bool { return r.s.LevelIndex() == 1 }, in)
	if r.ticks > maxLevelTicks {
		r.t.Fatalf("solution took %d ticks, want under %d", r.ticks, maxLevelTicks)
	}
}

func (r *runner) watcher() Watcher { return r.s.Watchers()[0] }

func (r *runner) dodgeCorridor() {
	const (
		pillarJumpX  = 58
		pocketX      = 155
		afterPocketX = 185
		gapJumpX     = 236
	)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.until(func() bool { return r.watcher().Dir > 0 && r.watcher().Box.X > afterPocketX }, idle)
	r.until(r.reachX(pillarJumpX), right)
	r.step(jumpRight)
	r.until(r.reachX(pocketX), right)
	r.shiftTo(world.Yellow)
	r.until(func() bool {
		w := r.watcher()
		return w.Dir < 0 && w.Box.Right() < r.x()-level.TileSize
	}, idle)
	r.shiftTo(world.Grey)
	r.until(r.reachX(gapJumpX), right)
	r.step(jumpRight)
}

func solveFirstShift(r *runner) {
	r.until(r.reachX(80), right)
	r.step(jumpRight)
	r.until(r.reachX(120), right)
	r.shiftTo(world.Grey)
	r.finish(right)
}

func solveFragments(r *runner) {
	r.until(r.reachX(80), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.until(r.reachX(254), right)
	r.shiftTo(world.Yellow)
	r.finish(right)
}

func solveSpikes(r *runner) {
	r.until(r.reachX(60), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
	r.until(r.reachX(232), right)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
	r.finish(right)
}

func solveWatcher(r *runner) {
	r.dodgeCorridor()
	r.finish(right)
}

func solvePlanning(r *runner) {
	r.dodgeCorridor()
	r.until(r.hasFragments(1), right)
	r.shiftTo(world.Yellow)
	r.until(r.reachX(430), right)
	r.shiftTo(world.Grey)
	r.finish(right)
}

func solveFinale(r *runner) {
	r.dodgeCorridor()
	r.until(func() bool { return r.s.PlayerBox().Y > 300 }, right)
	r.until(r.hasFragments(1), walkLeft)
	r.shiftTo(world.Yellow)
	r.until(r.backTo(430), walkLeft)
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), walkLeft)
	r.shiftTo(world.Yellow)
	r.finish(walkLeft)
}

func TestLevelSolutions(t *testing.T) {
	tests := []struct {
		name  string
		solve func(r *runner)
	}{
		{"01_first_shift", solveFirstShift},
		{"02_fragments", solveFragments},
		{"03_spikes", solveSpikes},
		{"04_watcher", solveWatcher},
		{"05_planning", solvePlanning},
		{"06_finale", solveFinale},
		{"07_on_the_beat", solveOnTheBeat},
		{"08_metronome", solveMetronome},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.solve(newRunner(t, tt.name))
		})
	}
}

func runUntilDeath(s *Session, inputs func(*Session) input.Intent) DeathCause {
	for range maxTicks {
		s.Update(inputs(s))
		if s.Death() != Alive {
			return s.Death()
		}
		if s.LevelIndex() != 0 {
			return Alive
		}
	}
	return Alive
}

func stayIn(w world.World, in input.Intent) func(*Session) input.Intent {
	return func(s *Session) input.Intent {
		if s.Active() != w {
			return input.Intent{ShiftPressed: true}
		}
		return in
	}
}

func jumpOverStep(s *Session) input.Intent {
	if x := s.PlayerBox().X; x >= 80 && x < 82 {
		return jumpRight
	}
	return right
}

func TestNaiveAttemptsFail(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		inputs func(*Session) input.Intent
		want   DeathCause
	}{
		{"never shifting falls into the pit", "01_first_shift", jumpOverStep, DeathFall},
		{"yellow spikes punish staying yellow", "03_spikes", stayIn(world.Yellow, right), DeathSpikes},
		{"grey spikes punish staying grey", "03_spikes", stayIn(world.Grey, right), DeathSpikes},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRunner(t, tt.level)
			if got := runUntilDeath(r.s, tt.inputs); got != tt.want {
				t.Fatalf("death = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFragmentsLevelBlocksStayingGrey(t *testing.T) {
	r := newRunner(t, "02_fragments")
	r.until(r.reachX(80), right)
	r.shiftTo(world.Grey)
	for range 600 {
		r.step(right)
	}
	if r.s.LevelIndex() != 0 || r.x() >= 260 {
		t.Fatalf("grey wall should stop the player, got x=%v", r.x())
	}
	if got := r.s.TryShift(); got != ShiftDone {
		t.Fatalf("the fragment picked on the bridge should pay the way back, got %v", got)
	}
}

type mechanics struct {
	fragments, spikes, watchers, checkpoints, quantized, auto bool
}

func mechanicsOf(l level.Level) mechanics {
	return mechanics{
		fragments:   l.StartFragments > 0 || len(l.Fragments) > 0,
		spikes:      len(l.Spikes[world.Yellow])+len(l.Spikes[world.Grey]) > 0,
		watchers:    len(l.Watchers) > 0,
		checkpoints: len(l.Checkpoints) > 0,
		quantized:   l.Shift == level.ShiftQuantized,
		auto:        l.Shift == level.ShiftAuto,
	}
}

func TestTeachingOrder(t *testing.T) {
	levels, err := level.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	want := []mechanics{
		{},
		{fragments: true},
		{fragments: true, spikes: true},
		{fragments: true, watchers: true},
		{fragments: true, spikes: true, watchers: true},
		{fragments: true, spikes: true, watchers: true, checkpoints: true},
		{fragments: true, watchers: true, quantized: true},
		{auto: true},
	}
	if len(levels) != len(want) {
		t.Fatalf("got %d levels, want %d", len(levels), len(want))
	}
	var seen mechanics
	for i, lvl := range levels {
		got := mechanicsOf(lvl)
		if got != want[i] {
			t.Errorf("level %d mechanics = %+v, want %+v", i+1, got, want[i])
		}
		introduced := 0
		for _, pair := range [][2]bool{
			{got.fragments, seen.fragments},
			{got.spikes, seen.spikes},
			{got.watchers, seen.watchers},
			{got.checkpoints, seen.checkpoints},
			{got.quantized, seen.quantized},
			{got.auto, seen.auto},
		} {
			if pair[0] && !pair[1] {
				introduced++
			}
		}
		if introduced > 1 {
			t.Errorf("level %d introduces %d mechanics at once", i+1, introduced)
		}
		seen = mechanics{
			fragments:   seen.fragments || got.fragments,
			spikes:      seen.spikes || got.spikes,
			watchers:    seen.watchers || got.watchers,
			checkpoints: seen.checkpoints || got.checkpoints,
			quantized:   seen.quantized || got.quantized,
			auto:        seen.auto || got.auto,
		}
	}
}

func TestRushingTheWatcherGetsCaught(t *testing.T) {
	r := newRunner(t, "04_watcher")
	r.shiftTo(world.Grey)
	r.until(r.hasFragments(1), right)
	r.until(func() bool { return r.watcher().Dir > 0 && r.watcher().Box.X > 185 }, idle)
	r.until(r.reachX(58), right)
	r.step(jumpRight)
	got := runUntilDeath(r.s, func(s *Session) input.Intent {
		if s.PlayerBox().X >= 236 {
			return jumpRight
		}
		return right
	})
	if got != DeathCaught {
		t.Fatalf("death = %v, want caught", got)
	}
}

func TestPlanningLevelHasExactFragments(t *testing.T) {
	lvl, err := level.Load("05_planning")
	if err != nil {
		t.Fatal(err)
	}
	const mandatoryReturns = 2
	if got := lvl.StartFragments + len(lvl.Fragments); got != mandatoryReturns {
		t.Fatalf("fragments available = %d, want exactly %d", got, mandatoryReturns)
	}
}

func TestFinaleRespawnsAtCheckpoint(t *testing.T) {
	r := newRunner(t, "06_finale")
	r.dodgeCorridor()
	r.until(func() bool { return r.s.PlayerBox().Y > 300 }, right)
	r.until(func() bool { _, ok := r.s.ActiveCheckpoint(); return ok && r.x() < 560 }, walkLeft)
	idx, _ := r.s.ActiveCheckpoint()
	cp := r.s.Level().Checkpoints[idx]

	walkUntilDead(r.s, walkLeft)
	if r.s.Death() != Alive {
		t.Fatalf("still dead: %v", r.s.Death())
	}
	want := player.New(cp).Box
	if r.s.PlayerBox() != want || r.s.Fragments() != 0 || len(r.s.RemainingFragments()) != 2 {
		t.Fatalf("respawn at %+v with %d fragments (%d remaining), want %+v",
			r.s.PlayerBox(), r.s.Fragments(), len(r.s.RemainingFragments()), want)
	}
}
