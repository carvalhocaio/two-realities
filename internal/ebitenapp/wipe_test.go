package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/game"
)

func TestWipeStartsOnlyOnShiftEvents(t *testing.T) {
	tests := []struct {
		name   string
		events []game.Event
		want   bool
	}{
		{"no events", nil, false},
		{"jump and land", []game.Event{game.EventJump, game.EventLand}, false},
		{"blocked shift", []game.Event{game.EventShiftBlocked}, false},
		{"queued shift", []game.Event{game.EventShiftQueued}, false},
		{"death", []game.Event{game.EventDeath}, false},
		{"shift to grey", []game.Event{game.EventShiftGrey}, true},
		{"shift to yellow after a landing", []game.Event{game.EventLand, game.EventShiftYellow}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w wipe
			w.notice(tt.events)
			if w.shifted != tt.want {
				t.Fatalf("shifted %v, want %v", w.shifted, tt.want)
			}
		})
	}
}

func TestWipeWithoutShiftDoesNotRun(t *testing.T) {
	var w wipe
	w.begin(10)
	if w.running {
		t.Fatal("wipe ran without a shift")
	}
}
