package game

import (
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

type Screen int

const (
	ScreenTitle Screen = iota
	ScreenPlaying
	ScreenEnd
)

type Flow struct {
	clock   rhythm.Clock
	levels  []level.Level
	screen  Screen
	session *Session
}

func NewFlow(clock rhythm.Clock, levels ...level.Level) *Flow {
	return &Flow{clock: clock, levels: levels, screen: ScreenTitle}
}

func (f *Flow) Screen() Screen    { return f.screen }
func (f *Flow) Session() *Session { return f.session }
func (f *Flow) LevelCount() int   { return len(f.levels) }

func (f *Flow) Update(in input.Intent) {
	switch f.screen {
	case ScreenTitle:
		if in.ShiftPressed {
			f.session = NewSession(f.clock, f.levels...)
			f.screen = ScreenPlaying
		}
	case ScreenPlaying:
		f.session.Update(in)
		if f.session.Finished() {
			f.screen = ScreenEnd
		}
	case ScreenEnd:
		if in.ShiftPressed {
			f.session = nil
			f.screen = ScreenTitle
		}
	}
}
