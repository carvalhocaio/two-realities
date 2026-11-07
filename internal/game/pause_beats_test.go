package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
	"github.com/carvalhocaio/two-realities/internal/world"
)

var beatDuration = rhythm.BeatDuration(rhythm.BPM)

func TestPauseDoesNotBendTheBeat(t *testing.T) {
	tests := []struct {
		name  string
		level string
		setup func(f *Flow)
	}{
		{"a queued shift still waits for a beat after resuming", quantizedOpen, func(f *Flow) { f.Update(shift) }},
		{"auto swaps skipped while paused do not happen", autoEveryTwo, func(*Flow) {}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := &rhythm.Manual{}
			f := NewFlow(clock, &save.Memory{}, mustParse(t, tt.level))
			f.Update(shift)
			clock.Advance(beatDuration / 4)
			f.Update(idle)
			tt.setup(f)
			s := f.Session()
			active, countdown := s.Active(), s.AutoCountdown()
			f.Update(pauseKey)
			clock.Advance(7 * beatDuration)
			f.Update(pauseKey)
			f.Update(idle)
			if s.Active() != active || s.AutoCountdown() != countdown {
				t.Fatalf("after resume: %v countdown %d, want %v countdown %d", s.Active(), s.AutoCountdown(), active, countdown)
			}
			changed := false
			for range beatTicks * 3 {
				clock.Advance(tickDuration)
				f.Update(idle)
				changed = changed || s.Active() != world.Yellow
			}
			if !changed {
				t.Fatal("the shift should still happen on a later beat")
			}
		})
	}
}
