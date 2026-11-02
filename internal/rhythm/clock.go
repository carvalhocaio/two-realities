package rhythm

import (
	"math"
	"time"
)

const BPM = 100

type Clock interface {
	Beat() int
	Phase() float64
}

func BeatDuration(bpm int) time.Duration {
	return time.Minute / time.Duration(bpm)
}

type Position struct {
	beat  int
	phase float64
}

func (p Position) Beat() int      { return p.beat }
func (p Position) Phase() float64 { return p.phase }

func At(elapsed time.Duration, bpm int) Position {
	beats := float64(elapsed) / float64(BeatDuration(bpm))
	whole := math.Floor(beats)
	return Position{beat: int(whole), phase: beats - whole}
}

type Manual struct {
	Elapsed time.Duration
}

func (m *Manual) Advance(d time.Duration) { m.Elapsed += d }
func (m *Manual) Beat() int               { return At(m.Elapsed, BPM).Beat() }
func (m *Manual) Phase() float64          { return At(m.Elapsed, BPM).Phase() }
