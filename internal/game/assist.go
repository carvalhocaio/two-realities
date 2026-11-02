package game

import (
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/physics"
)

const (
	assistAlertFactor = 2
	pushReach         = 3 * level.TileSize
)

func (s *Session) SetAssist(on bool) {
	s.assist = on
	if on {
		s.timer.assisted = true
		s.timer.levelAssisted = true
	}
}

func (s *Session) Assist() bool { return s.assist }

func (s *Session) AlertLimit() int {
	if s.assist {
		return s.baseAlert() * assistAlertFactor
	}
	return s.baseAlert()
}

func nearestFree(box geom.Rect, g level.Grid) (geom.Rect, bool) {
	for d := 1.0; d <= pushReach; d++ {
		for _, off := range []geom.Vec{{Y: -d}, {X: -d}, {X: d}, {Y: d}} {
			next := box
			next.X += off.X
			next.Y += off.Y
			if !physics.OverlapsSolid(next, g) {
				return next, true
			}
		}
	}
	return box, false
}
