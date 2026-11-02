package music

import "math"

const (
	SampleRate  = 44100
	BeatsPerBar = 4
	BarsPerLoop = 4
	peakLimit   = 0.85
)

func BeatFrames(bpm int) int { return SampleRate * 60 / bpm }

func BarFrames(bpm int) int { return BeatFrames(bpm) * BeatsPerBar }

type buffer []float32

func (b buffer) add(i int, v float64) { b[((i%len(b))+len(b))%len(b)] += float32(v) }

func midiHz(note int) float64 { return 440 * math.Pow(2, float64(note-69)/12) }

type wave func(phase float64) float64

func sine(p float64) float64     { return math.Sin(2 * math.Pi * p) }
func triangle(p float64) float64 { return 4*math.Abs(p-math.Floor(p+0.5)) - 1 }
func saw(p float64) float64      { return 2*(p-math.Floor(p)) - 1 }
func square(p float64) float64 {
	if p-math.Floor(p) < 0.5 {
		return 1
	}
	return -1
}

type voice struct {
	wave    wave
	gain    float64
	attack  float64
	decay   float64
	sustain float64
	release float64
	cutoff  float64
	detune  float64
}

func (v voice) envelope(t, length float64) float64 {
	var e float64
	switch {
	case t < v.attack:
		e = t / v.attack
	case t < v.attack+v.decay:
		e = 1 - (1-v.sustain)*(t-v.attack)/v.decay
	default:
		e = v.sustain
	}
	if t > length {
		e *= math.Max(0, 1-(t-length)/v.release)
	}
	return e
}

func (v voice) play(b buffer, start, frames int, note int) {
	freq := midiHz(note)
	total := frames + int(v.release*SampleRate)
	length := float64(frames) / SampleRate
	alpha := 1.0
	if v.cutoff > 0 {
		alpha = 1 - math.Exp(-2*math.Pi*v.cutoff/SampleRate)
	}
	low := 0.0
	for i := range total {
		t := float64(i) / SampleRate
		x := v.wave(freq * t)
		if v.detune != 0 {
			x = (x + v.wave(freq*(1+v.detune)*t)) / 2
		}
		low += alpha * (x - low)
		b.add(start+i, v.gain*v.envelope(t, length)*low)
	}
}

func kick(b buffer, start int, gain float64) {
	phase := 0.0
	for i := range int(0.22 * SampleRate) {
		t := float64(i) / SampleRate
		freq := 45 + 95*math.Exp(-t*30)
		phase += 2 * math.Pi * freq / SampleRate
		b.add(start+i, gain*math.Exp(-t*18)*math.Sin(phase))
	}
}

type noise uint32

func (n *noise) next() float64 {
	*n = *n*1664525 + 1013904223
	return float64(uint32(*n))/(1<<31) - 1
}

func hat(b buffer, start int, gain float64) {
	n, prev := noise(0x2545F491), 0.0
	for i := range int(0.05 * SampleRate) {
		x := n.next()
		b.add(start+i, gain*math.Exp(-float64(i)/SampleRate*80)*(x-prev))
		prev = x
	}
}

func snare(b buffer, start int, gain float64) {
	n := noise(0x9E3779B9)
	for i := range int(0.16 * SampleRate) {
		t := float64(i) / SampleRate
		body := 0.5 * math.Sin(2*math.Pi*185*t) * math.Exp(-t*30)
		b.add(start+i, gain*(math.Exp(-t*22)*n.next()*0.7+body))
	}
}

func normalize(b buffer) buffer {
	peak := 0.0
	for _, v := range b {
		peak = math.Max(peak, math.Abs(float64(v)))
	}
	if peak > peakLimit {
		k := float32(peakLimit / peak)
		for i := range b {
			b[i] *= k
		}
	}
	return b
}
