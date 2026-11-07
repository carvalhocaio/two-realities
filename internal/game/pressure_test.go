package game

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

func TestPressureRulesStack(t *testing.T) {
	s := NewSession(&rhythm.Manual{}, mustParse(t, quantizedOpen))
	rules := []PressureRule{FasterWatchers, NoCheckpoints, BeatGhost, DeathRestartsLevel, DeathRestartsAct}
	for p := range MaxPressure + 1 {
		s.SetPressure(p)
		for _, r := range rules {
			if got, want := s.Has(r), int(r) <= p; got != want {
				t.Errorf("pressure %d rule %d = %v, want %v", p, r, got, want)
			}
		}
	}
	s.SetPressure(9)
	if s.Pressure() != MaxPressure {
		t.Fatalf("pressure clamps to %d, got %d", MaxPressure, s.Pressure())
	}
}

func TestSolutionsSurviveFullPressure(t *testing.T) {
	for _, tt := range levelSolutions {
		t.Run(tt.name, func(t *testing.T) {
			r := newRunner(t, tt.name)
			r.s.SetPressure(MaxPressure)
			tt.solve(r)
		})
	}
}

func TestFasterWatchersAndShorterAlert(t *testing.T) {
	r := newRunner(t, "04_watcher")
	r.s.SetPressure(int(FasterWatchers))
	x := r.watcher().Box.X
	r.update(idle)
	if got := r.watcher().Box.X - x; got != WatcherSpeed*pressureSpeedFactor {
		t.Fatalf("watcher moved %v, want %v", got, WatcherSpeed*pressureSpeedFactor)
	}
	if r.s.AlertLimit() >= AlertTicks {
		t.Fatalf("alert limit %d not shorter than %d", r.s.AlertLimit(), AlertTicks)
	}
}

func TestNoCheckpointsUnderPressure(t *testing.T) {
	r := newRunner(t, "07_crossfire")
	r.s.SetPressure(int(NoCheckpoints))
	r.s.player = player.New(r.s.Level().Checkpoints[0])
	r.update(idle)
	if _, ok := r.s.ActiveCheckpoint(); ok || r.s.CheckpointsActive() {
		t.Fatal("checkpoint activated under pressure")
	}
}

func TestDeathUnderPressure(t *testing.T) {
	levels, err := level.LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	const descent = 7
	tests := []struct {
		name      string
		pressure  int
		wantLevel int
		wantCheck bool
	}{
		{"respawn at the checkpoint", 0, descent, true},
		{"restart the level", int(DeathRestartsLevel), descent, false},
		{"restart the act", int(DeathRestartsAct), 4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(&rhythm.Manual{}, levels...)
			s.SetPressure(tt.pressure)
			s.GoToLevel(descent)
			s.player = player.New(s.Level().Checkpoints[0])
			s.Update(idle)
			s.die(DeathFall)
			for range DeathFreezeTicks + 1 {
				s.Update(idle)
			}
			_, checked := s.ActiveCheckpoint()
			if s.LevelIndex() != tt.wantLevel || checked != tt.wantCheck {
				t.Fatalf("level %d checkpoint %v, want %d %v", s.LevelIndex(), checked, tt.wantLevel, tt.wantCheck)
			}
		})
	}
}

func TestPressureUnlocksAfterACleanRunAndRecordsClears(t *testing.T) {
	store := &save.Memory{}
	f := seenFlow(&rhythm.Manual{}, store, mustParse(t, quantizedOpen))
	f.Update(input.Intent{IncreasePressed: true})
	if f.Pressure() != 0 {
		t.Fatal("pressure chosen before it was unlocked")
	}
	finish := func() {
		f.Update(input.Intent{ShiftPressed: true})
		for range maxTicks {
			if f.Screen() == ScreenEnd {
				break
			}
			f.Update(right)
		}
		f.Update(input.Intent{ShiftPressed: true})
	}
	finish()
	if !f.Records().PressureUnlocked() {
		t.Fatal("finishing the campaign did not unlock pressure")
	}
	f.Update(input.Intent{IncreasePressed: true})
	f.Update(input.Intent{IncreasePressed: true})
	if f.Pressure() != 2 {
		t.Fatalf("pressure %d, want 2", f.Pressure())
	}
	best, _ := f.Records().BestRun()
	finish()
	if f.Records().PressureCleared() != 2 {
		t.Fatalf("cleared %d, want 2", f.Records().PressureCleared())
	}
	if after, _ := f.Records().BestRun(); after != best {
		t.Fatal("a pressure run changed the ranked best run")
	}
}

func TestRhythmSolutionsSurviveFullPressureOnEveryBeatPhase(t *testing.T) {
	tests := []struct {
		name  string
		solve func(r *runner)
	}{
		{"09_on_the_beat", solveOnTheBeat},
		{"11_offbeat", solveOffbeat},
		{"12_clockwork", solveClockwork},
		{wardenLevel, solveWarden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for offset := range beatTicks {
				r := newRunnerAt(t, tt.name, tickDuration*time.Duration(offset))
				r.s.SetPressure(MaxPressure)
				tt.solve(r)
			}
		})
	}
}

func TestOldSavesWithABestRunUnlockPressure(t *testing.T) {
	store := &save.Memory{}
	_ = store.Save(save.Data{BestRun: 9000})
	if !NewRecords(store).PressureUnlocked() {
		t.Fatal("a saved campaign best should unlock pressure")
	}
}
