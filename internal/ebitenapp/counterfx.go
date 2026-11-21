package ebitenapp

import (
	"image/color"

	"github.com/carvalhocaio/two-realities/internal/game"
)

const (
	counterGainTicks  = 6
	counterSpendTicks = 4
	counterSpendAlpha = 0.35
)

type counterFX struct {
	session *game.Session
	level   int
	epoch   int
	count   int
	gain    int
	spend   int
}

func (f *counterFX) update(s *game.Session) {
	if s != f.session || s.LevelIndex() != f.level || s.Epoch() != f.epoch {
		*f = counterFX{session: s, level: s.LevelIndex(), epoch: s.Epoch(), count: s.Fragments()}
		return
	}
	f.gain, f.spend = max(f.gain-1, 0), max(f.spend-1, 0)
	switch n := s.Fragments(); {
	case n > f.count:
		f.gain, f.spend = counterGainTicks, 0
	case n < f.count:
		f.spend, f.gain = counterSpendTicks, 0
	}
	f.count = s.Fragments()
}

func (f *counterFX) icon() (lift float64, c color.Color, alpha float32) {
	switch {
	case f.gain > 0:
		return 1, secretColor, 1
	case f.spend > 0:
		return 0, fragmentColor, counterSpendAlpha
	}
	return 0, fragmentColor, 1
}
