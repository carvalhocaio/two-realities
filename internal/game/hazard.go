package game

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	SpikeHeight = 5
	spikeInset  = 1
)

func SpikeBox(p level.Point) geom.Rect {
	tile := level.TileRect(p)
	return geom.Rect{
		X: tile.X + spikeInset,
		Y: tile.Bottom() - SpikeHeight,
		W: tile.W - 2*spikeInset,
		H: SpikeHeight,
	}
}

func (s *Session) touchingSpikes() bool {
	for _, p := range s.level.Spikes[s.active] {
		if s.player.Box.Overlaps(SpikeBox(p)) {
			return true
		}
	}
	return false
}

func (s *Session) updateWatchers() {
	grey := s.level.Layer(world.Grey)
	for i := range s.watchers {
		w := &s.watchers[i]
		w.Patrol(grey)
		if s.active == world.Grey && w.Sees(s.player.Box, grey) {
			w.Alert++
		} else {
			w.Alert = 0
		}
		if w.Alert >= s.AlertLimit() {
			s.die(DeathCaught)
		}
	}
}
