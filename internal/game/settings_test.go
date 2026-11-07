package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

func pausedAt(t *testing.T, store save.Storage, item MenuItem) *Flow {
	t.Helper()
	f := seenFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
	f.Update(shift)
	f.Update(input.Intent{PausePressed: true})
	for range menuDowns(item) {
		f.Update(input.Intent{DownPressed: true})
	}
	if f.Menu() != item {
		t.Fatalf("menu at %v, want %v", f.Menu(), item)
	}
	return f
}

func TestVolumes(t *testing.T) {
	tests := []struct {
		name         string
		item         MenuItem
		inputs       []input.Intent
		music, sound int
	}{
		{"starts full", MenuMusic, nil, VolumeSteps, VolumeSteps},
		{"music down", MenuMusic, ticks(lessKey, 3), VolumeSteps - 3, VolumeSteps},
		{"music stops at zero", MenuMusic, ticks(lessKey, 40), 0, VolumeSteps},
		{"music stops at full", MenuMusic, ticks(moreKey, 5), VolumeSteps, VolumeSteps},
		{"sounds down and up", MenuSounds, then(ticks(lessKey, 4), ticks(moreKey, 1)), VolumeSteps, VolumeSteps - 3},
		{"other items ignore left and right", MenuEffects, ticks(lessKey, 3), VolumeSteps, VolumeSteps},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := pausedAt(t, &save.Memory{}, tt.item)
			for _, in := range tt.inputs {
				f.Update(in)
			}
			if f.MusicVolume() != tt.music || f.SoundVolume() != tt.sound {
				t.Fatalf("music %d sound %d, want %d %d", f.MusicVolume(), f.SoundVolume(), tt.music, tt.sound)
			}
			if f.Screen() != ScreenPaused {
				t.Fatal("adjusting a volume should keep the menu open")
			}
		})
	}
}

func TestVolumesPersistAndMuteStaysApart(t *testing.T) {
	store := &save.Memory{}
	f := pausedAt(t, store, MenuMusic)
	for _, in := range then(ticks(lessKey, 2), ticks(downKey, 1), ticks(lessKey, 5), one(input.Intent{MutePressed: true})) {
		f.Update(in)
	}
	again := seenFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
	if again.MusicVolume() != VolumeSteps-2 || again.SoundVolume() != VolumeSteps-5 {
		t.Fatalf("reloaded music %d sound %d", again.MusicVolume(), again.SoundVolume())
	}
	if !f.Muted() || again.Muted() {
		t.Fatal("mute is a session toggle that does not change or persist the volumes")
	}
}

func TestBadSavedVolumesAreClamped(t *testing.T) {
	f := seenFlow(&rhythm.Manual{}, &save.Memory{Stored: save.Data{MusicDown: -7, SoundDown: 99}}, mustParse(t, timedLevel))
	if f.MusicVolume() != VolumeSteps || f.SoundVolume() != 0 {
		t.Fatalf("music %d sound %d", f.MusicVolume(), f.SoundVolume())
	}
}

func TestReduceEffects(t *testing.T) {
	tests := []struct {
		name      string
		saved     string
		preferred bool
		toggle    bool
		want      bool
	}{
		{"off by default", "", false, false, false},
		{"follows the system preference", "", true, false, true},
		{"a saved full choice beats the system", effectsFull, true, false, false},
		{"a saved reduced choice beats the system", effectsReduced, false, false, true},
		{"garbage falls back to the system", "sparkly", true, false, true},
		{"the menu toggles it", "", false, true, true},
		{"the menu toggles it back off", "", true, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &save.Memory{Stored: save.Data{Effects: tt.saved}}
			f := seenFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
			f.UseReducedMotion(tt.preferred)
			if tt.toggle {
				f.Update(shift)
				f.Update(input.Intent{PausePressed: true})
				f.menu = MenuEffects
				f.Update(input.Intent{ConfirmPressed: true})
			}
			if f.ReduceEffects() != tt.want {
				t.Fatalf("reduce effects %v, want %v", f.ReduceEffects(), tt.want)
			}
			if tt.toggle {
				again := seenFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
				again.UseReducedMotion(!tt.want)
				if again.ReduceEffects() != tt.want {
					t.Fatal("the menu choice should be saved and win over the system")
				}
			}
		})
	}
}

func TestSettingsWithBrokenStorage(t *testing.T) {
	f := pausedAt(t, save.NewFallback(unavailableStorage{}), MenuMusic)
	f.Update(lessKey)
	f.menu = MenuEffects
	f.Update(input.Intent{ConfirmPressed: true})
	if f.MusicVolume() != VolumeSteps-1 || !f.ReduceEffects() {
		t.Fatalf("music %d reduced %v", f.MusicVolume(), f.ReduceEffects())
	}
}

func TestRecordsSkipUnchangedSettings(t *testing.T) {
	store := &countingStore{Memory: save.Memory{Stored: save.Data{MusicDown: 2, Effects: effectsReduced}}}
	r := NewRecords(store)
	r.setVolumes(VolumeSteps-2, VolumeSteps)
	r.setReducedEffects(true)
	if store.saves != 0 {
		t.Fatalf("%d writes for unchanged settings", store.saves)
	}
	r.setReducedEffects(false)
	if on, set := r.ReducedEffects(); on || !set || store.saves != 1 {
		t.Fatalf("reduced %v set %v writes %d", on, set, store.saves)
	}
}
