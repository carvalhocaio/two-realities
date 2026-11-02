package game

import (
	"math"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/physics"
)

const (
	WatcherWidth  = 8
	WatcherHeight = 10
	WatcherSpeed  = 0.5
	WatcherEyeY   = 3
	SightRange    = 8
	AlertTicks    = 24
)

type Watcher struct {
	Box   geom.Rect
	Dir   int
	Alert int
}

func NewWatcher(p level.Point) Watcher {
	tile := level.TileRect(p)
	return Watcher{
		Box: geom.Rect{X: tile.X + (tile.W-WatcherWidth)/2, Y: tile.Bottom() - WatcherHeight, W: WatcherWidth, H: WatcherHeight},
		Dir: 1,
	}
}

func (w Watcher) AlertLevel(limit int) float64 { return float64(w.Alert) / float64(limit) }

func (w *Watcher) Patrol(g level.Grid) {
	if w.blockedAhead(g) {
		w.Dir = -w.Dir
		if w.blockedAhead(g) {
			return
		}
	}
	w.Box.X += float64(w.Dir) * WatcherSpeed
}

func (w Watcher) blockedAhead(g level.Grid) bool {
	next := w.Box
	next.X += float64(w.Dir) * WatcherSpeed
	if physics.OverlapsSolid(next, g) {
		return true
	}
	footX := next.X
	if w.Dir > 0 {
		footX = next.Right() - 1e-6
	}
	return !g.Solid(tileOf(footX), tileOf(next.Bottom()))
}

func (w Watcher) EyeY() float64 { return w.Box.Y + WatcherEyeY }

func (w Watcher) eyeX() float64 {
	if w.Dir > 0 {
		return w.Box.Right()
	}
	return w.Box.X
}

func (w Watcher) SightOrigin() geom.Vec { return geom.Vec{X: w.eyeX(), Y: w.EyeY()} }

func (w Watcher) SightEnd(g level.Grid) float64 {
	origin := w.eyeX()
	limit := origin + float64(w.Dir*SightRange*level.TileSize)
	row := tileOf(w.EyeY())
	col := tileOf(origin)
	if w.Dir < 0 {
		col = tileOf(origin - 1e-6)
	}
	for {
		near := float64(col * level.TileSize)
		if w.Dir < 0 {
			near += level.TileSize
		}
		if w.Dir > 0 && near >= limit || w.Dir < 0 && near <= limit {
			return limit
		}
		if !g.InBounds(col, row) || g.Solid(col, row) {
			if w.Dir > 0 {
				return math.Max(near, origin)
			}
			return math.Min(near, origin)
		}
		col += w.Dir
	}
}

func (w Watcher) Sees(target geom.Rect, g level.Grid) bool {
	eye := w.EyeY()
	if eye < target.Y || eye >= target.Bottom() {
		return false
	}
	origin, end := w.eyeX(), w.SightEnd(g)
	if w.Dir > 0 {
		return target.Right() > origin && target.X < end
	}
	return target.X < origin && target.Right() > end
}

func tileOf(v float64) int { return int(math.Floor(v / level.TileSize)) }
