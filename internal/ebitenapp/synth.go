package ebitenapp

import (
	"encoding/binary"
	"math"
	"sync/atomic"
	"time"

	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

const (
	sampleRate    = 44100
	beatsPerBar   = 4
	bytesPerFrame = 8

	kickLength   = 0.22
	kickStartHz  = 140.0
	kickEndHz    = 45.0
	kickSweep    = 30.0
	kickDecay    = 18.0
	kickGain     = 0.8
	hatLength    = 0.06
	hatDecay     = 70.0
	hatGain      = 0.18
	bassHz       = 55.0
	bassLength   = 0.18
	bassDecay    = 14.0
	bassGain     = 0.35
	snareLength  = 0.14
	snareDecay   = 28.0
	snareGain    = 0.3
	noiseSeed    = 0x2545F491
	noiseMul     = 1664525
	noiseInc     = 1013904223
	noiseDivisor = 1 << 31
)

func beatFrames() int {
	return int(int64(rhythm.BeatDuration(rhythm.BPM)) * sampleRate / int64(time.Second))
}

const intensityLevels = 4

func synthBars() [][]byte {
	bars := make([][]byte, intensityLevels)
	for i := range bars {
		bars[i] = synthBar(i)
	}
	return bars
}

func synthBar(intensity int) []byte {
	beat := beatFrames()
	mix := make([]float64, beat*beatsPerBar)
	for b := range beatsPerBar {
		start := b * beat
		addKick(mix[start:])
		addHat(mix[start+beat/2:], hatGain)
		if intensity >= 1 {
			addBass(mix[start+beat/2:])
		}
		if intensity >= 2 {
			addHat(mix[start+beat/4:], hatGain/2)
			addHat(mix[start+3*beat/4:], hatGain/2)
		}
		if intensity >= 3 && b%2 == 1 {
			addSnare(mix[start:])
		}
	}
	return encodeStereoF32(mix)
}

func addBass(dst []float64) {
	for i := range min(len(dst), int(bassLength*sampleRate)) {
		t := float64(i) / sampleRate
		dst[i] += bassGain * math.Exp(-t*bassDecay) * math.Sin(2*math.Pi*bassHz*t)
	}
}

func addSnare(dst []float64) {
	addNoise(dst, snareLength, snareDecay, snareGain)
}

func addKick(dst []float64) {
	phase := 0.0
	for i := range min(len(dst), int(kickLength*sampleRate)) {
		t := float64(i) / sampleRate
		freq := kickEndHz + (kickStartHz-kickEndHz)*math.Exp(-t*kickSweep)
		phase += 2 * math.Pi * freq / sampleRate
		dst[i] += kickGain * math.Exp(-t*kickDecay) * math.Sin(phase)
	}
}

func addHat(dst []float64, gain float64) {
	seed := uint32(noiseSeed)
	prev := 0.0
	for i := range min(len(dst), int(hatLength*sampleRate)) {
		seed = seed*noiseMul + noiseInc
		noise := float64(seed)/noiseDivisor - 1
		t := float64(i) / sampleRate
		dst[i] += gain * math.Exp(-t*hatDecay) * (noise - prev)
		prev = noise
	}
}

func addNoise(dst []float64, length, decay, gain float64) {
	seed := uint32(noiseSeed)
	for i := range min(len(dst), int(length*sampleRate)) {
		seed = seed*noiseMul + noiseInc
		noise := float64(seed)/noiseDivisor - 1
		t := float64(i) / sampleRate
		dst[i] += gain * math.Exp(-t*decay) * noise
	}
}

func encodeStereoF32(mono []float64) []byte {
	out := make([]byte, len(mono)*bytesPerFrame)
	for i, v := range mono {
		bits := math.Float32bits(float32(max(-1, min(1, v))))
		binary.LittleEndian.PutUint32(out[i*bytesPerFrame:], bits)
		binary.LittleEndian.PutUint32(out[i*bytesPerFrame+4:], bits)
	}
	return out
}

type loopReader struct {
	bars [][]byte
	cur  int
	next atomic.Int32
	pos  int
}

func (l *loopReader) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		data := l.bars[l.cur]
		c := copy(p[n:], data[l.pos:])
		n += c
		l.pos += c
		if l.pos == len(data) {
			l.pos = 0
			l.cur = int(l.next.Load())
		}
	}
	return n, nil
}
