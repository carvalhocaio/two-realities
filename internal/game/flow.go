package game

import (
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

type Screen int

const (
	ScreenTitle Screen = iota
	ScreenPlaying
	ScreenPaused
	ScreenEnd
)

type MenuItem int

const (
	MenuResume MenuItem = iota
	MenuAssist
	MenuMute
	MenuTitle
	MenuCount
)

type Flow struct {
	clock   rhythm.Clock
	levels  []level.Level
	records *Records
	screen  Screen
	session *Session
	menu    MenuItem
	assist  bool
	muted   bool
	started bool
}

func NewFlow(clock rhythm.Clock, store save.Storage, levels ...level.Level) *Flow {
	return &Flow{clock: clock, levels: levels, records: NewRecords(store), screen: ScreenTitle}
}

func (f *Flow) Screen() Screen    { return f.screen }
func (f *Flow) Session() *Session { return f.session }
func (f *Flow) LevelCount() int   { return len(f.levels) }
func (f *Flow) Menu() MenuItem    { return f.menu }
func (f *Flow) Assist() bool      { return f.assist }
func (f *Flow) Muted() bool       { return f.muted }
func (f *Flow) Records() *Records { return f.records }

func (f *Flow) Update(in input.Intent) {
	if in.MutePressed {
		f.muted = !f.muted
	}
	switch f.screen {
	case ScreenTitle:
		if in.ShiftPressed {
			f.start()
		}
	case ScreenPlaying:
		if in.PausePressed || in.FocusLost {
			f.screen, f.menu = ScreenPaused, MenuResume
			return
		}
		f.session.Update(in)
		if f.session.Finished() {
			f.screen = ScreenEnd
		}
	case ScreenPaused:
		f.updateMenu(in)
	case ScreenEnd:
		if in.ShiftPressed {
			f.session = nil
			f.screen = ScreenTitle
		}
	}
}

func (f *Flow) StorageNotice() bool {
	return !f.started && f.records.Degraded()
}

func (f *Flow) start() {
	f.started = true
	f.session = NewSession(f.clock, f.levels...)
	f.session.UseRecords(f.records)
	f.session.SetAssist(f.assist)
	f.screen = ScreenPlaying
}

func (f *Flow) updateMenu(in input.Intent) {
	switch {
	case in.ConfirmPressed:
		f.confirm()
	case in.PausePressed:
		f.screen = ScreenPlaying
	case in.UpPressed:
		f.menu = (f.menu + MenuCount - 1) % MenuCount
	case in.DownPressed:
		f.menu = (f.menu + 1) % MenuCount
	}
}

func (f *Flow) confirm() {
	switch f.menu {
	case MenuResume:
		f.screen = ScreenPlaying
	case MenuAssist:
		f.assist = !f.assist
		f.session.SetAssist(f.assist)
	case MenuMute:
		f.muted = !f.muted
	case MenuTitle:
		f.session = nil
		f.screen = ScreenTitle
	}
}
