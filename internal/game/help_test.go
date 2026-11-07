package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

var helpKey = input.Intent{HelpPressed: true}

func helpFlow(t *testing.T, store save.Storage) *Flow {
	t.Helper()
	return NewFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
}

func TestFirstStartShowsTheHelpCard(t *testing.T) {
	tests := []struct {
		name string
		data save.Data
		want Screen
	}{
		{"new save", save.Data{}, ScreenHelp},
		{"seen before", save.Data{HelpSeen: true}, ScreenPlaying},
		{"old save with a level best", save.Data{BestLevel: map[string]int{"x": 1}}, ScreenPlaying},
		{"old save with a reached level", save.Data{Reached: map[string]bool{"x": true}}, ScreenPlaying},
		{"finished campaign", save.Data{BestRun: 10}, ScreenPlaying},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := helpFlow(t, &save.Memory{Stored: tt.data})
			f.Update(shift)
			if f.Screen() != tt.want {
				t.Fatalf("screen %v, want %v", f.Screen(), tt.want)
			}
		})
	}
}

func TestClosingTheFirstHelpStartsAndIsSaved(t *testing.T) {
	store := &save.Memory{}
	f := helpFlow(t, store)
	f.Update(shift)
	f.Update(shift)
	if f.Screen() != ScreenPlaying || f.Session() == nil {
		t.Fatalf("screen %v after closing the first help", f.Screen())
	}
	again := helpFlow(t, store)
	again.Update(shift)
	if again.Screen() != ScreenPlaying {
		t.Fatal("the help card should show only once per save")
	}
}

func TestHelpFromTheTitle(t *testing.T) {
	tests := []struct {
		name  string
		close input.Intent
	}{
		{"shift", shift},
		{"confirm", input.Intent{ConfirmPressed: true}},
		{"pause key", input.Intent{PausePressed: true}},
		{"h again", helpKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := helpFlow(t, &save.Memory{})
			f.Update(helpKey)
			if f.Screen() != ScreenHelp {
				t.Fatalf("h on the title: screen %v", f.Screen())
			}
			f.Update(tt.close)
			if f.Screen() != ScreenTitle || f.Session() != nil {
				t.Fatalf("closing returns to the title, got %v", f.Screen())
			}
		})
	}
}

func TestHelpIgnoresOtherKeys(t *testing.T) {
	f := helpFlow(t, &save.Memory{})
	f.Update(helpKey)
	for _, in := range []input.Intent{idle, right, left, {JumpPressed: true}, {UpPressed: true}, {DecreasePressed: true}, {LanguagePressed: true}, {FocusLost: true}, restart} {
		f.Update(in)
		if f.Screen() != ScreenHelp {
			t.Fatalf("%+v closed the help card", in)
		}
	}
}

func TestHelpFromThePauseMenu(t *testing.T) {
	f := helpFlow(t, &save.Memory{Stored: save.Data{HelpSeen: true}})
	f.Update(shift)
	f.Update(input.Intent{PausePressed: true})
	f.Update(input.Intent{DownPressed: true})
	if f.Menu() != MenuHelp {
		t.Fatalf("menu at %v, want help", f.Menu())
	}
	f.Update(input.Intent{ConfirmPressed: true})
	if f.Screen() != ScreenHelp {
		t.Fatalf("screen %v, want help", f.Screen())
	}
	before := f.Session().PlayerBox()
	f.Update(right)
	f.Update(shift)
	if f.Screen() != ScreenPaused || f.Session().PlayerBox() != before {
		t.Fatalf("closing returns to the paused game untouched, got %v", f.Screen())
	}
}

func TestHelpWorksWithBrokenStorage(t *testing.T) {
	f := helpFlow(t, save.NewFallback(unavailableStorage{}))
	f.Update(shift)
	f.Update(shift)
	if f.Screen() != ScreenPlaying {
		t.Fatalf("screen %v, want playing", f.Screen())
	}
}
