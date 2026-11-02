package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const wardenArena = "fragments: 2\nboss: warden\n" +
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

const (
	lampOneCol   = 2
	exitCol      = 3
	coverCol     = 4
	lampTwoCol   = 7
	trackCol     = 8
	lampThreeCol = 10
	floorRow     = 6
)

func newWarden(t *testing.T) *beatDriver {
	t.Helper()
	return newBeatDriver(t, wardenArena)
}

func (d *beatDriver) place(col int) {
	d.s.player = player.New(level.Point{X: col, Y: floorRow})
}

func (d *beatDriver) survive(beat int) DeathCause {
	for d.clock.Beat() < beat && d.s.Death() == Alive {
		d.tick(idle)
	}
	return d.s.Death()
}

func (d *beatDriver) enterGrey(t *testing.T) {
	t.Helper()
	if d.s.Active() != world.Grey && d.s.TryShift() != ShiftDone {
		t.Fatal("could not shift to grey")
	}
}

func (d *beatDriver) light(t *testing.T, col int) {
	t.Helper()
	before := d.s.warden.lit
	d.enterGrey(t)
	d.place(col)
	d.tick(idle)
	if d.s.warden.lit != before+1 {
		t.Fatalf("lamp at column %d not lit (lit = %d)", col, d.s.warden.lit)
	}
}

func TestWardenPhasesDriveShiftModeAndGates(t *testing.T) {
	d := newWarden(t)
	gate := func(i int) bool {
		return d.s.layer(world.Yellow).Solid(6+3*i, 4) && d.s.layer(world.Grey).Solid(6+3*i, 4)
	}
	steps := []struct {
		lamp      int
		wantMode  level.ShiftMode
		wantGates [2]bool
	}{
		{0, level.ShiftFree, [2]bool{true, true}},
		{lampOneCol, level.ShiftQuantized, [2]bool{false, true}},
		{lampTwoCol, level.ShiftAuto, [2]bool{false, false}},
		{lampThreeCol, level.ShiftFree, [2]bool{false, false}},
	}
	for i, st := range steps {
		if st.lamp > 0 {
			d.light(t, st.lamp)
		}
		if got := d.s.shiftMode(); got != st.wantMode {
			t.Errorf("step %d mode = %v, want %v", i, got, st.wantMode)
		}
		if got := [2]bool{gate(0), gate(1)}; got != st.wantGates {
			t.Errorf("step %d gates closed = %v, want %v", i, got, st.wantGates)
		}
	}
	if _, ok := d.s.ActiveCheckpoint(); ok {
		t.Error("lamps must not light checkpoint flags")
	}
}

func TestLampsNeedGreyAndOrder(t *testing.T) {
	d := newWarden(t)
	d.place(lampOneCol)
	for range 10 {
		d.tick(idle)
	}
	if d.s.warden.lit != 0 {
		t.Fatal("lamp lit from the yellow world")
	}
	d.enterGrey(t)
	d.place(lampThreeCol)
	d.tick(idle)
	if d.s.warden.lit != 0 {
		t.Fatal("third lamp lit before the first")
	}
}

func TestExitLockedUntilWardenFalls(t *testing.T) {
	d := newWarden(t)
	d.place(exitCol)
	d.tick(idle)
	if d.s.Finished() {
		t.Fatal("exit open with the warden standing")
	}
	for _, col := range []int{lampOneCol, lampTwoCol, lampThreeCol} {
		d.light(t, col)
	}
	d.place(exitCol)
	d.tick(idle)
	if !d.s.Finished() {
		t.Fatal("exit still locked after the third lamp")
	}
}

func TestSweepHitsOnlyInOpenGrey(t *testing.T) {
	tests := []struct {
		name  string
		col   int
		grey  bool
		until int
		want  DeathCause
	}{
		{"light still far away", 1, true, 1, Alive},
		{"open grey under the light", 1, true, 3, DeathBeam},
		{"open grey nearer the start", 5, true, 2, DeathBeam},
		{"yellow under the light", 1, false, 3, Alive},
		{"under grey cover", coverCol, true, 3, Alive},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newWarden(t)
			if tt.grey {
				d.enterGrey(t)
			}
			d.place(tt.col)
			if got := d.survive(tt.until); got != tt.want {
				t.Fatalf("death = %v, want %v (beam %+v)", got, tt.want, d.s.warden.beam)
			}
		})
	}
}

func TestDeathKeepsLampsAndRestartDropsThem(t *testing.T) {
	d := newWarden(t)
	d.light(t, lampOneCol)
	d.s.die(DeathBeam)
	for range DeathFreezeTicks + 1 {
		d.tick(idle)
	}
	want := player.New(d.s.Level().Boss.Lamps[0]).Box
	if d.s.warden.lit != 1 || d.s.PlayerBox() != want || d.s.layer(world.Yellow).Solid(6, 4) {
		t.Fatalf("respawn lit=%d at %+v, want lit=1 at %+v with gate open", d.s.warden.lit, d.s.PlayerBox(), want)
	}
	d.tick(input.Intent{RestartPressed: true})
	if d.s.warden.lit != 0 || !d.s.layer(world.Yellow).Solid(6, 4) {
		t.Fatalf("restart kept lit=%d", d.s.warden.lit)
	}
}

func TestQueuedShiftResolvesBeforeTheLightArrives(t *testing.T) {
	for _, queue := range []bool{true, false} {
		d := newWarden(t)
		d.light(t, lampOneCol)
		start := d.s.warden.start
		d.place(trackCol)
		if queue {
			d.tick(input.Intent{ShiftPressed: true})
		}
		want := DeathBeam
		if queue {
			want = Alive
		}
		if got := d.survive(start + RestBeats + 2); got != want {
			t.Fatalf("queued=%v death = %v, want %v", queue, got, want)
		}
	}
}

func TestAssistSlowsTheSweep(t *testing.T) {
	for _, assist := range []bool{false, true} {
		d := newWarden(t)
		d.s.SetAssist(assist)
		d.enterGrey(t)
		d.place(5)
		want := DeathBeam
		if assist {
			want = Alive
		}
		if got := d.survive(1); got != want {
			t.Fatalf("assist=%v death = %v, want %v", assist, got, want)
		}
	}
}

func TestFinalPhaseTracksThePlayer(t *testing.T) {
	for _, assist := range []bool{false, true} {
		d := newWarden(t)
		d.s.SetAssist(assist)
		d.light(t, lampOneCol)
		d.light(t, lampTwoCol)
		start := d.s.warden.start
		d.place(trackCol)
		d.untilBeat(start + RestBeats)
		d.tick(idle)
		if b := d.s.warden.beam; b.Stage != BeamTelegraph || b.X != columnCenter(trackCol) {
			t.Fatalf("beam %+v, want telegraph on column %d", b, trackCol)
		}
		want := DeathBeam
		if assist {
			want = Alive
		}
		if got := d.survive(start + RestBeats + 2); got != want {
			t.Fatalf("assist=%v death = %v, want %v on the third beat", assist, got, want)
		}
	}
}
