package ebitenapp

import (
	"encoding/binary"
	"math"
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
	noiseSeed    = 0x2545F491
	noiseMul     = 1664525
	noiseInc     = 1013904223
	noiseDivisor = 1 << 31
)

func beatFrames() int {
	return int(int64(rhythm.BeatDuration(rhythm.BPM)) * sampleRate / int64(time.Second))
}

func synthBar() []byte {
	beat := beatFrames()
	mix := make([]float64, beat*beatsPerBar)
	for b := range beatsPerBar {
		addKick(mix[b*beat:])
		addHat(mix[b*beat+beat/2:])
	}
	return encodeStereoF32(mix)
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

func addHat(dst []float64) {
	seed := uint32(noiseSeed)
	prev := 0.0
	for i := range min(len(dst), int(hatLength*sampleRate)) {
		seed = seed*noiseMul + noiseInc
		noise := float64(seed)/noiseDivisor - 1
		t := float64(i) / sampleRate
		dst[i] += hatGain * math.Exp(-t*hatDecay) * (noise - prev)
		prev = noise
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
	data []byte
	pos  int
}

func (l *loopReader) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		c := copy(p[n:], l.data[l.pos:])
		n += c
		l.pos = (l.pos + c) % len(l.data)
	}
	return n, nil
}
