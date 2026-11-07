package music

import "math"

type Sound int

const (
	SoundJump Sound = iota
	SoundLand
	SoundShiftGrey
	SoundShiftYellow
	SoundBlocked
	SoundQueued
	SoundFragment
	SoundSecret
	SoundCheckpoint
	SoundTransmission
	SoundSpotted
	SoundDeath
	SoundLamp
	SoundExit
	SoundCount
)

const (
	sfxPeak       = 0.5
	sfxMaxSeconds = 2
)

var (
	blip  = voice{wave: square, gain: 0.25, attack: 0.002, decay: 0.04, sustain: 0.3, release: 0.03, cutoff: 3000}
	chime = voice{wave: sine, gain: 0.5, attack: 0.002, decay: 0.15, sustain: 0.2, release: 0.25}
	bell  = voice{wave: triangle, gain: 0.45, attack: 0.003, decay: 0.3, sustain: 0.3, release: 0.5, detune: 0.004}
	buzz  = voice{wave: saw, gain: 0.35, attack: 0.002, decay: 0.05, sustain: 0.8, release: 0.04, cutoff: 900}
)

type tone struct {
	at, length float64
	pitch      int
	v          voice
}

func RenderSound(s Sound) []float32 {
	b := make(buffer, sfxMaxSeconds*SampleRate)
	switch s {
	case SoundJump:
		sweep(b, 0.07, 520, 880, 0.12, square, 2500)
	case SoundLand:
		thud(b, 0.06, 0.25)
	case SoundShiftGrey:
		sweep(b, 0.16, 660, 220, 0.18, triangle, 0)
		hiss(b, 0.12, 0.05)
	case SoundShiftYellow:
		sweep(b, 0.16, 330, 990, 0.18, triangle, 0)
		hiss(b, 0.12, 0.05)
	case SoundBlocked:
		notes(b, tone{0, 0.06, 40, buzz}, tone{0.08, 0.08, 39, buzz})
	case SoundQueued:
		notes(b, tone{0, 0.02, 88, blip})
	case SoundFragment:
		notes(b, tone{0, 0.06, 88, chime}, tone{0.06, 0.12, 95, chime})
	case SoundSecret:
		notes(b, tone{0, 0.07, 81, chime}, tone{0.07, 0.07, 85, chime}, tone{0.14, 0.07, 88, chime}, tone{0.21, 0.3, 93, chime})
	case SoundCheckpoint:
		notes(b, tone{0, 0.1, 69, bell}, tone{0.1, 0.25, 76, bell})
	case SoundTransmission:
		notes(b, tone{0, 0.03, 84, blip}, tone{0.06, 0.03, 84, blip}, tone{0.12, 0.05, 91, blip})
		hiss(b, 0.2, 0.03)
	case SoundSpotted:
		notes(b, tone{0, 0.05, 86, blip}, tone{0.07, 0.08, 91, blip})
	case SoundDeath:
		sweep(b, 0.35, 440, 55, 0.22, saw, 1200)
		hiss(b, 0.25, 0.18)
	case SoundLamp:
		notes(b, tone{0, 0.5, 60, bell}, tone{0.02, 0.5, 67, bell}, tone{0.04, 0.5, 72, bell}, tone{0.06, 0.6, 76, bell})
	case SoundExit:
		notes(b, tone{0, 0.08, 69, bell}, tone{0.08, 0.08, 72, bell}, tone{0.16, 0.08, 76, bell}, tone{0.24, 0.3, 81, bell})
	}
	b = trim(b)
	fadeOut(b)
	return limit(b)
}

func trim(b buffer) buffer {
	end := len(b)
	for end > 0 && math.Abs(float64(b[end-1])) < 1e-4 {
		end--
	}
	return b[:end+1]
}

func notes(b buffer, ns ...tone) {
	for _, n := range ns {
		n.v.play(b, int(n.at*SampleRate), int(n.length*SampleRate), n.pitch)
	}
}

func sweep(b buffer, seconds, from, to, gain float64, w wave, cutoff float64) {
	frames := int(seconds * SampleRate)
	alpha := 1.0
	if cutoff > 0 {
		alpha = 1 - math.Exp(-2*math.Pi*cutoff/SampleRate)
	}
	phase, low := 0.0, 0.0
	for i := range frames {
		t := float64(i) / float64(frames)
		phase += from * math.Pow(to/from, t) / SampleRate
		low += alpha * (w(phase) - low)
		env := math.Min(1, t*40) * (1 - t)
		b.add(i, gain*env*low)
	}
}

func hiss(b buffer, seconds, gain float64) {
	n := noise(0x51ED270B)
	frames := int(seconds * SampleRate)
	for i := range frames {
		t := float64(i) / float64(frames)
		b.add(i, gain*(1-t)*(1-t)*n.next())
	}
}

func thud(b buffer, seconds, gain float64) {
	n := noise(0x1B873593)
	frames := int(seconds * SampleRate)
	low := 0.0
	for i := range frames {
		t := float64(i) / SampleRate
		low += 0.08 * (n.next() - low)
		env := math.Exp(-t * 60)
		b.add(i, gain*env*(math.Sin(2*math.Pi*(70+40*math.Exp(-t*50))*t)+2*low))
	}
}

func fadeOut(b buffer) {
	tail := min(len(b), SampleRate/100)
	for i := range tail {
		b[len(b)-1-i] *= float32(i) / float32(tail)
	}
}

func limit(b buffer) []float32 {
	peak := 0.0
	for _, v := range b {
		peak = math.Max(peak, math.Abs(float64(v)))
	}
	if peak > sfxPeak {
		k := float32(sfxPeak / peak)
		for i := range b {
			b[i] *= k
		}
	}
	return b
}
