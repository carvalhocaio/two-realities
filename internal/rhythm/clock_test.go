package rhythm

import (
	"math"
	"testing"
	"time"
)

func TestAt(t *testing.T) {
	beat := BeatDuration(BPM)
	tests := []struct {
		name    string
		elapsed time.Duration
		bpm     int
		beat    int
		phase   float64
	}{
		{"start", 0, BPM, 0, 0},
		{"half beat", beat / 2, BPM, 0, 0.5},
		{"exactly one beat", beat, BPM, 1, 0},
		{"many beats", 10*beat + beat/4, BPM, 10, 0.25},
		{"before start", -beat / 4, BPM, -1, 0.75},
		{"faster tempo", time.Second, 120, 2, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := At(tt.elapsed, tt.bpm)
			if got.Beat() != tt.beat || math.Abs(got.Phase()-tt.phase) > 1e-9 {
				t.Fatalf("At(%v) = beat %d phase %v, want %d %v", tt.elapsed, got.Beat(), got.Phase(), tt.beat, tt.phase)
			}
		})
	}
}

func TestBeatDurationAtDefaultTempo(t *testing.T) {
	if got := BeatDuration(BPM); got != 600*time.Millisecond {
		t.Fatalf("beat = %v, want 600ms", got)
	}
}

func TestManualAdvances(t *testing.T) {
	var m Manual
	m.Advance(BeatDuration(BPM) * 3 / 2)
	if m.Beat() != 1 || math.Abs(m.Phase()-0.5) > 1e-9 {
		t.Fatalf("beat %d phase %v, want 1 0.5", m.Beat(), m.Phase())
	}
}
