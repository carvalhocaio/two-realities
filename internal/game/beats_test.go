package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	quantizedOpen  = "fragments: 0\nshift: quantized\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n#......#\n########\n"
	quantizedWall  = "fragments: 0\nshift: quantized\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n##.....#\n########\n"
	autoOpen       = "fragments: 0\nshift: auto\nauto_every: 4\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n#......#\n########\n"
	autoEveryTwo   = "fragments: 0\nshift: auto\nauto_every: 2\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n#......#\n########\n"
	autoGreyWalled = "fragments: 0\nshift: auto\nauto_every: 2\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n##.....#\n########\n"
)

type beatDriver struct {
	s     *Session
	clock *rhythm.Manual
}

func newBeatDriver(t *testing.T, text string) *beatDriver {
	t.Helper()
	clock := &rhythm.Manual{}
	return &beatDriver{s: NewSession(clock, mustParse(t, text)), clock: clock}
}

func (d *beatDriver) tick(in input.Intent) {
	d.clock.Advance(tickDuration)
	d.s.Update(in)
}

func (d *beatDriver) untilBeat(b int) {
	for d.clock.Beat() < b {
		d.tick(idle)
	}
}

func TestQuantizedShift(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		play   func(d *beatDriver)
		world  world.World
		queued bool
		check  func(t *testing.T, s *Session)
	}{
		{"press queues without shifting", quantizedOpen,
			func(d *beatDriver) { d.tick(shift); d.tick(idle) },
			world.Yellow, true, nil},
		{"queued shift resolves on the next beat", quantizedOpen,
			func(d *beatDriver) { d.tick(shift); d.untilBeat(1) },
			world.Grey, false, nil},
		{"pressing again cancels", quantizedOpen,
			func(d *beatDriver) { d.tick(shift); d.tick(shift); d.untilBeat(1) },
			world.Yellow, false, nil},
		{"only one shift per press", quantizedOpen,
			func(d *beatDriver) { d.tick(shift); d.untilBeat(3) },
			world.Grey, false, nil},
		{"fragment cost applies at resolution", quantizedOpen,
			func(d *beatDriver) { d.tick(shift); d.untilBeat(1); d.tick(shift); d.untilBeat(2) },
			world.Grey, false,
			func(t *testing.T, s *Session) {
				if !s.RestartHint() || s.Feedback().CounterFlash == 0 {
					t.Fatal("want no-fragments feedback at resolution")
				}
			}},
		{"solid overlap blocks at resolution", quantizedWall,
			func(d *beatDriver) { d.tick(shift); d.untilBeat(1) },
			world.Yellow, false,
			func(t *testing.T, s *Session) {
				if s.Feedback().Shake == 0 {
					t.Fatal("want blocked-solid feedback at resolution")
				}
			}},
		{"restart clears the queue", quantizedOpen,
			func(d *beatDriver) { d.tick(shift); d.tick(restart); d.untilBeat(1) },
			world.Yellow, false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newBeatDriver(t, tt.level)
			tt.play(d)
			if d.s.Active() != tt.world || d.s.ShiftQueued() != tt.queued {
				t.Fatalf("world %v queued %v, want %v %v", d.s.Active(), d.s.ShiftQueued(), tt.world, tt.queued)
			}
			if tt.check != nil {
				tt.check(t, d.s)
			}
		})
	}
}

func TestFreeShiftIgnoresBeats(t *testing.T) {
	d := newBeatDriver(t, "fragments: 0\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n#......#\n########\n")
	d.tick(shift)
	if d.s.Active() != world.Grey || d.s.ShiftQueued() {
		t.Fatalf("free shift should be immediate, got %v queued %v", d.s.Active(), d.s.ShiftQueued())
	}
}

func TestAutoShift(t *testing.T) {
	tests := []struct {
		name  string
		level string
		play  func(d *beatDriver)
		world world.World
		death DeathCause
	}{
		{"no swap before N beats", autoEveryTwo, func(d *beatDriver) { d.untilBeat(1) }, world.Yellow, Alive},
		{"swaps every N beats", autoEveryTwo, func(d *beatDriver) { d.untilBeat(2) }, world.Grey, Alive},
		{"swaps back for free", autoEveryTwo, func(d *beatDriver) { d.untilBeat(4) }, world.Yellow, Alive},
		{"player shift is ignored", autoEveryTwo, func(d *beatDriver) { d.tick(shift); d.untilBeat(1) }, world.Yellow, Alive},
		{"swapping into a solid kills", autoGreyWalled, func(d *beatDriver) { d.untilBeat(2) }, world.Grey, DeathCrushed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newBeatDriver(t, tt.level)
			tt.play(d)
			if d.s.Active() != tt.world || d.s.Death() != tt.death {
				t.Fatalf("world %v death %v, want %v %v", d.s.Active(), d.s.Death(), tt.world, tt.death)
			}
			if d.s.Fragments() != 0 {
				t.Fatalf("auto swaps must not cost fragments, got %d", d.s.Fragments())
			}
		})
	}
}

func TestAutoCountdown(t *testing.T) {
	d := newBeatDriver(t, autoOpen)
	want := []int{0, 3, 2, 1, 0, 3, 2, 1, 0}
	for beat, w := range want {
		d.untilBeat(beat)
		if got := d.s.AutoCountdown(); got != w {
			t.Fatalf("beat %d: countdown %d, want %d", beat, got, w)
		}
	}
}

func TestAutoCountdownHiddenOutsideAuto(t *testing.T) {
	d := newBeatDriver(t, quantizedOpen)
	for beat := range 5 {
		d.untilBeat(beat)
		if got := d.s.AutoCountdown(); got != 0 {
			t.Fatalf("beat %d: countdown %d outside auto mode", beat, got)
		}
	}
}

func TestAutoRespawnRestartsTheCount(t *testing.T) {
	d := newBeatDriver(t, autoGreyWalled)
	d.untilBeat(2)
	if d.s.Death() != DeathCrushed {
		t.Fatalf("death = %v, want crushed", d.s.Death())
	}
	for range DeathFreezeTicks {
		d.tick(idle)
	}
	respawnBeat := d.clock.Beat()
	if d.s.Death() != Alive || d.s.Active() != world.Yellow {
		t.Fatalf("after freeze: death %v world %v", d.s.Death(), d.s.Active())
	}
	d.untilBeat(respawnBeat + 1)
	if d.s.Death() != Alive {
		t.Fatalf("swapped one beat after respawn, want %d beats", 2)
	}
	d.untilBeat(respawnBeat + 2)
	if d.s.Death() != DeathCrushed {
		t.Fatalf("death = %v, want crushed two beats after respawn", d.s.Death())
	}
}

func TestLevelShiftModeDefaultsToFree(t *testing.T) {
	levels, err := level.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	for i, lvl := range levels[:6] {
		if lvl.Shift != level.ShiftFree {
			t.Errorf("level %d shift = %v, want free", i+1, lvl.Shift)
		}
	}
}
