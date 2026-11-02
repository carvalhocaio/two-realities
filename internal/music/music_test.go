package music

import (
	"math"
	"testing"
)

const testBPM = 100

func TestTracksLoopOnWholeBars(t *testing.T) {
	want := BarFrames(testBPM) * BarsPerLoop
	for tr := range TrackCount {
		got := Render(tr, testBPM)
		if len(got) != want {
			t.Errorf("track %d has %d frames, want %d", tr, len(got), want)
		}
	}
}

func TestTracksAreAudibleAndNeverClip(t *testing.T) {
	for tr := range TrackCount {
		peak, energy := 0.0, 0.0
		for _, v := range Render(tr, testBPM) {
			peak = math.Max(peak, math.Abs(float64(v)))
			energy += float64(v) * float64(v)
		}
		if peak > peakLimit+1e-6 || peak < 0.2 || energy == 0 {
			t.Errorf("track %d peak %.3f", tr, peak)
		}
	}
}

func TestKickMarksEveryBarDownbeat(t *testing.T) {
	beat := BeatFrames(testBPM)
	for tr := range TrackCount {
		b := Render(tr, testBPM)
		on, off := window(b, 0, beat/8), window(b, beat/4, beat/8)
		if on <= off {
			t.Errorf("track %d: downbeat energy %.3f not above off-beat %.3f", tr, on, off)
		}
	}
}

func TestRenderIsDeterministic(t *testing.T) {
	a, b := Render(ActThree, testBPM), Render(ActThree, testBPM)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("frame %d differs", i)
		}
	}
}

func window(b []float32, from, n int) float64 {
	sum := 0.0
	for _, v := range b[from : from+n] {
		sum += float64(v) * float64(v)
	}
	return sum
}
