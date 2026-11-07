package game

import (
	"slices"
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

func threeLevels(t *testing.T) []level.Level {
	t.Helper()
	var out []level.Level
	for _, name := range []string{"a", "b", "c"} {
		l, err := level.Parse(exitLevel)
		if err != nil {
			t.Fatal(err)
		}
		l.Name = name
		out = append(out, l)
	}
	return out
}

func TestTitleStartsAtTheFurthestReachedLevel(t *testing.T) {
	tests := []struct {
		name    string
		data    save.Data
		start   int
		choices int
	}{
		{"new save", save.Data{}, 0, 1},
		{"reached the second", save.Data{Reached: map[string]bool{"b": true}}, 1, 2},
		{"old save with a level best", save.Data{BestLevel: map[string]int{"a": 100, "b": 200}}, 2, 3},
		{"campaign done starts from the top", save.Data{BestRun: 900}, 0, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &save.Memory{}
			_ = store.Save(tt.data)
			f := seenFlow(&rhythm.Manual{}, store, threeLevels(t)...)
			if f.StartLevel() != tt.start || f.StartChoices() != tt.choices {
				t.Fatalf("start %d of %d, want %d of %d", f.StartLevel(), f.StartChoices(), tt.start, tt.choices)
			}
		})
	}
}

func TestChoosingAStartLevel(t *testing.T) {
	store := &save.Memory{}
	_ = store.Save(save.Data{Reached: map[string]bool{"c": true}})
	f := seenFlow(&rhythm.Manual{}, store, threeLevels(t)...)
	f.Update(input.Intent{UpPressed: true})
	if f.StartLevel() != 2 {
		t.Fatalf("up past the furthest level: %d", f.StartLevel())
	}
	f.Update(input.Intent{DownPressed: true})
	f.Update(shift)
	if f.Session().LevelIndex() != 1 {
		t.Fatalf("started at %d, want 1", f.Session().LevelIndex())
	}
	if f.Session().timer.runEligible {
		t.Fatal("a run started past level 1 must not be eligible for the best run")
	}
}

func TestEnteringALevelSavesItAsReached(t *testing.T) {
	store := &save.Memory{}
	f := seenFlow(&rhythm.Manual{}, store, threeLevels(t)...)
	f.Update(shift)
	for range 20 {
		if f.Session().LevelIndex() == 1 {
			break
		}
		f.Update(right)
	}
	if f.Session().LevelIndex() != 1 {
		t.Fatal("expected to reach level 2")
	}
	again := seenFlow(&rhythm.Manual{}, store, threeLevels(t)...)
	if again.StartLevel() != 1 {
		t.Fatalf("continue from %d, want 1", again.StartLevel())
	}
}

func TestSkipLevelNeedsAssist(t *testing.T) {
	f := seenFlow(&rhythm.Manual{}, &save.Memory{}, threeLevels(t)...)
	f.Update(shift)
	f.Update(input.Intent{PausePressed: true})
	if slices.Contains(f.MenuItems(), MenuSkip) {
		t.Fatal("skip shown without assist")
	}
	f.menu = MenuAssist
	f.Update(input.Intent{ConfirmPressed: true})
	f.Update(input.Intent{DownPressed: true})
	if f.Menu() != MenuSkip {
		t.Fatalf("menu at %v after assist, want skip", f.Menu())
	}
	f.Update(input.Intent{ConfirmPressed: true})
	if f.Screen() != ScreenPlaying || f.Session().LevelIndex() != 1 {
		t.Fatalf("screen %v level %d after skip", f.Screen(), f.Session().LevelIndex())
	}
	f.Session().GoToLevel(2)
	f.Update(input.Intent{PausePressed: true})
	if slices.Contains(f.MenuItems(), MenuSkip) {
		t.Fatal("the last level cannot be skipped")
	}
}
