package game

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const FragmentSize = 6

func FragmentBox(p level.Point) geom.Rect {
	tile := level.TileRect(p)
	inset := (level.TileSize - FragmentSize) / 2.0
	return geom.Rect{X: tile.X + inset, Y: tile.Y + inset, W: FragmentSize, H: FragmentSize}
}

func (s *Session) Fragments() int { return s.fragments }

func (s *Session) RemainingFragments() []level.Point {
	var out []level.Point
	for i, p := range s.level.Fragments {
		if !s.collected[i] {
			out = append(out, p)
		}
	}
	return out
}

func (s *Session) collectFragments() {
	if s.active != world.Grey {
		return
	}
	for i, p := range s.level.Fragments {
		if !s.collected[i] && s.player.Box.Overlaps(FragmentBox(p)) {
			s.collected[i] = true
			s.fragments++
			s.restartHint = false
			s.emit(EventFragment)
		}
	}
}
