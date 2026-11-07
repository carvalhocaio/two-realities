package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

var (
	focusLost = input.Intent{FocusLost: true}
	pauseKey  = input.Intent{PausePressed: true}
	enterKey  = input.Intent{PausePressed: true, ConfirmPressed: true}
	downKey   = input.Intent{DownPressed: true}
)

func TestAutoPause(t *testing.T) {
	tests := []struct {
		name   string
		inputs []input.Intent
		screen Screen
		assist bool
	}{
		{"focus loss pauses play", then(one(shift), ticks(right, 3), one(focusLost)), ScreenPaused, false},
		{"focus loss is ignored on the title", one(focusLost), ScreenTitle, false},
		{"focus loss keeps the menu open", then(one(shift), one(focusLost), one(focusLost)), ScreenPaused, false},
		{"never resumes by itself", then(one(shift), one(focusLost), ticks(idle, 120), ticks(right, 30)), ScreenPaused, false},
		{"pause key resumes", then(one(shift), one(focusLost), one(pauseKey)), ScreenPlaying, false},
		{"enter pauses during play", then(one(shift), one(enterKey)), ScreenPaused, false},
		{"enter confirms the selected item", then(one(shift), one(focusLost), ticks(downKey, menuDowns(MenuAssist)), one(enterKey)), ScreenPaused, true},
		{"enter on resume resumes", then(one(shift), one(focusLost), one(enterKey)), ScreenPlaying, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := seenFlow(&rhythm.Manual{}, &save.Memory{}, mustParse(t, timedLevel))
			for _, in := range tt.inputs {
				f.Update(in)
			}
			if f.Screen() != tt.screen || f.Assist() != tt.assist {
				t.Fatalf("screen %v assist %v, want %v %v", f.Screen(), f.Assist(), tt.screen, tt.assist)
			}
		})
	}
}

func TestFocusLossOnEndScreenIsIgnored(t *testing.T) {
	f := seenFlow(&rhythm.Manual{}, &save.Memory{}, mustParse(t, timedLevel))
	f.Update(shift)
	for f.Screen() == ScreenPlaying {
		f.Update(right)
	}
	passEpilogue(t, f)
	f.Update(focusLost)
	if f.Screen() != ScreenEnd {
		t.Fatalf("screen %v, want end", f.Screen())
	}
}

type unavailableStorage struct{}

func (unavailableStorage) Load() (save.Data, error) { return save.Data{}, save.ErrUnavailable }
func (unavailableStorage) Save(save.Data) error     { return save.ErrUnavailable }

func TestStorageNotice(t *testing.T) {
	tests := []struct {
		name   string
		store  save.Storage
		notice bool
	}{
		{"working storage stays quiet", save.NewFallback(&save.Memory{}), false},
		{"blocked storage warns", save.NewFallback(unavailableStorage{}), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := seenFlow(&rhythm.Manual{}, tt.store, mustParse(t, timedLevel))
			if f.StorageNotice() != tt.notice {
				t.Fatalf("notice %v, want %v", f.StorageNotice(), tt.notice)
			}
			f.Update(shift)
			for f.Screen() == ScreenPlaying {
				f.Update(right)
			}
			passEpilogue(t, f)
			f.Update(shift)
			if f.Screen() != ScreenTitle || f.StorageNotice() {
				t.Fatal("the notice shows only once, before the first run")
			}
		})
	}
}

func TestBrokenStorageNeverStopsTheGame(t *testing.T) {
	lvl := mustParse(t, timedLevel)
	f := seenFlow(&rhythm.Manual{}, save.NewFallback(unavailableStorage{}), lvl, lvl)
	f.Update(shift)
	for range 200 {
		f.Update(right)
	}
	passEpilogue(t, f)
	if f.Screen() != ScreenEnd {
		t.Fatalf("screen %v, want the run to finish", f.Screen())
	}
	if best, ok := f.Records().BestRun(); !ok || best == 0 {
		t.Fatal("session records should still work in memory")
	}
}
