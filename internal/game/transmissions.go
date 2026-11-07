package game

import (
	"math"

	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	TransmissionTicks     = 240
	TransmissionFadeTicks = 30
	TransmissionRange     = 12 * level.TileSize
)

type TransmissionView struct {
	Source level.Transmission
	Age    int
	Fade   float64
}

func (v TransmissionView) TextIn(l i18n.Lang) string { return v.Source.TextIn(l) }

type transmissionState struct {
	heard   []bool
	current int
	age     int
	fading  int
}

func newTransmissionState(count int) transmissionState {
	return transmissionState{heard: make([]bool, count), current: -1}
}

func (s *Session) updateTransmissions() {
	t := &s.signals
	for i, tr := range s.level.Transmissions {
		if !t.heard[i] && tr.Layer == s.active && s.player.Box.Overlaps(level.TileRect(tr.Point)) {
			t.heard[i] = true
			t.current, t.age, t.fading = i, 0, 0
			s.emit(EventTransmission)
		}
	}
	if t.current < 0 {
		return
	}
	t.age++
	if t.fading == 0 && (t.age >= TransmissionTicks || s.farFrom(s.level.Transmissions[t.current].Point)) {
		t.fading = 1
	}
	if t.fading > 0 {
		t.fading++
		if t.fading > TransmissionFadeTicks {
			t.current = -1
		}
	}
}

func (s *Session) farFrom(p level.Point) bool {
	tile := level.TileRect(p)
	box := s.player.Box
	dx := (box.X + box.W/2) - (tile.X + tile.W/2)
	dy := (box.Y + box.H/2) - (tile.Y + tile.H/2)
	return math.Hypot(dx, dy) > TransmissionRange
}

func (s *Session) Transmission() (TransmissionView, bool) {
	t := s.signals
	if t.current < 0 {
		return TransmissionView{}, false
	}
	return TransmissionView{
		Source: s.level.Transmissions[t.current],
		Age:    t.age,
		Fade:   float64(max(t.fading-1, 0)) / TransmissionFadeTicks,
	}, true
}

func (s *Session) Heard(i int) bool { return s.signals.heard[i] }
