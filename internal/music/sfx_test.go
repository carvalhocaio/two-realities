package music

import (
	"math"
	"testing"
)

func TestEverySoundRendersAShortBoundedClip(t *testing.T) {
	for s := range SoundCount {
		b := RenderSound(s)
		if len(b) == 0 || len(b) >= sfxMaxSeconds*SampleRate {
			t.Fatalf("sound %d: %d frames", s, len(b))
		}
		peak := 0.0
		for _, v := range b {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				t.Fatalf("sound %d: invalid sample", s)
			}
			peak = math.Max(peak, math.Abs(float64(v)))
		}
		if peak == 0 || peak > sfxPeak+1e-6 {
			t.Fatalf("sound %d: peak %v", s, peak)
		}
		if b[len(b)-1] != 0 {
			t.Fatalf("sound %d: does not end in silence", s)
		}
	}
}
