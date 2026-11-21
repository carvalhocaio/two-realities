package ebitenapp

import (
	"math"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
)

const (
	visorTurnTicks = 4
	visorMinWidth  = 2
	treadStride    = 3
	treadWidth     = 2
	sightPulseLen  = 2
	sightPulseStep = 2
	jitterLevel    = 0.6
	jitterTicks    = 2
)

type watcherFX struct {
	session *game.Session
	level   int
	epoch   int
	dirs    []int
	turns   []int
}

func (f *watcherFX) update(s *game.Session) {
	watchers := s.Watchers()
	if s != f.session || s.LevelIndex() != f.level || s.Epoch() != f.epoch || len(watchers) != len(f.dirs) {
		f.session, f.level, f.epoch = s, s.LevelIndex(), s.Epoch()
		f.dirs, f.turns = f.dirs[:0], f.turns[:0]
		for _, w := range watchers {
			f.dirs = append(f.dirs, w.Dir)
			f.turns = append(f.turns, 0)
		}
		return
	}
	for i, w := range watchers {
		f.turns[i] = max(f.turns[i]-1, 0)
		if w.Dir != f.dirs[i] {
			f.dirs[i], f.turns[i] = w.Dir, visorTurnTicks
		}
	}
}

func (f *watcherFX) turning(i int) int {
	if i >= len(f.turns) {
		return 0
	}
	return f.turns[i]
}

func visorSideX(box geom.Rect, dir int) float64 {
	if dir < 0 {
		return box.X + 1
	}
	return box.Right() - visorWidth - 1
}

func visorRect(w game.Watcher, turning int) geom.Rect {
	y := w.EyeY()
	to := visorSideX(w.Box, w.Dir)
	if turning == 0 {
		return geom.Rect{X: math.Round(to), Y: y, W: visorWidth, H: 1}
	}
	t := 1 - float64(turning)/visorTurnTicks
	from := visorSideX(w.Box, -w.Dir)
	width := math.Round(visorWidth - (visorWidth-visorMinWidth)*math.Sin(math.Pi*t))
	x := from + (to-from)*t + (visorWidth-width)/2
	return geom.Rect{X: math.Round(x), Y: y, W: width, H: 1}
}

func treadNotch(box geom.Rect) geom.Rect {
	x := box.X + 1
	if int(math.Floor(box.X/treadStride))%2 == 1 {
		x = box.Right() - 1 - treadWidth
	}
	return geom.Rect{X: math.Round(x), Y: box.Bottom() - 1, W: treadWidth, H: 1}
}

func sightPulse(x0, x1 float64, dir, tick int) (float64, bool) {
	length := int(x1 - x0)
	if length <= sightPulseLen {
		return 0, false
	}
	off := float64((tick * sightPulseStep) % (length - sightPulseLen))
	if dir < 0 {
		return x1 - sightPulseLen - off, true
	}
	return x0 + off, true
}

func alertJitter(level float64, tick int, reduced bool) float64 {
	if reduced || level <= jitterLevel {
		return 0
	}
	if (tick/jitterTicks)%2 == 0 {
		return -1
	}
	return 1
}
