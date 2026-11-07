package game

import (
	"slices"

	"github.com/carvalhocaio/two-realities/internal/i18n"
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
	ScreenHelp
	ScreenEpilogue
)

type MenuItem int

const (
	MenuResume MenuItem = iota
	MenuHelp
	MenuAssist
	MenuSkip
	MenuMusic
	MenuSounds
	MenuEffects
	MenuAudioOffset
	MenuLanguage
	MenuTitle
	MenuCount
)

const (
	AudioOffsetStep  = 10
	AudioOffsetLimit = 500
)

type Flow struct {
	clock    rhythm.Clock
	levels   []level.Level
	records  *Records
	screen   Screen
	session  *Session
	menu     MenuItem
	assist   bool
	muted    bool
	started  bool
	offset   int
	lang     i18n.Lang
	pressure int
	startAt  int
	helpBack Screen
	music    int
	sound    int
	reduced  bool
	epilogue epilogueState
}

func NewFlow(clock rhythm.Clock, store save.Storage, levels ...level.Level) *Flow {
	records := NewRecords(store)
	f := &Flow{clock: clock, levels: levels, records: records, offset: clampOffset(records.AudioOffset()), music: records.MusicVolume(), sound: records.SoundVolume()}
	f.reduced, _ = records.ReducedEffects()
	f.toTitle()
	return f
}

func (f *Flow) toTitle() {
	f.session = nil
	f.screen = ScreenTitle
	f.startAt = 0
	if !f.records.CampaignDone() {
		f.startAt = f.records.Furthest(f.levels)
	}
}

func (f *Flow) StartLevel() int   { return f.startAt }
func (f *Flow) StartChoices() int { return f.records.Furthest(f.levels) + 1 }

func (f *Flow) chooseStart(in input.Intent) {
	switch {
	case in.UpPressed:
		f.startAt = min(f.startAt+1, f.StartChoices()-1)
	case in.DownPressed:
		f.startAt = max(f.startAt-1, 0)
	}
}

func (f *Flow) MenuItems() []MenuItem {
	items := make([]MenuItem, 0, MenuCount)
	for m := range MenuCount {
		if m != MenuSkip || f.canSkip() {
			items = append(items, m)
		}
	}
	return items
}

func (f *Flow) canSkip() bool {
	return f.assist && f.session != nil && f.session.LevelIndex() < f.session.LevelCount()-1
}

func (f *Flow) moveMenu(step int) {
	items := f.MenuItems()
	i := max(slices.Index(items, f.menu), 0)
	f.menu = items[(i+step+len(items))%len(items)]
}

func (f *Flow) Screen() Screen      { return f.screen }
func (f *Flow) Session() *Session   { return f.session }
func (f *Flow) LevelCount() int     { return len(f.levels) }
func (f *Flow) Menu() MenuItem      { return f.menu }
func (f *Flow) Assist() bool        { return f.assist }
func (f *Flow) Muted() bool         { return f.muted }
func (f *Flow) AudioOffsetMs() int  { return f.offset }
func (f *Flow) Language() i18n.Lang { return f.lang }

func (f *Flow) UseLocale(locale string) {
	if saved, ok := f.records.Language(); ok {
		f.lang = saved
		return
	}
	f.lang = i18n.Detect(locale)
}

func (f *Flow) switchLanguage() {
	f.lang = f.lang.Next()
	f.records.setLanguage(f.lang)
}
func (f *Flow) Records() *Records { return f.records }

func (f *Flow) Update(in input.Intent) {
	if in.MutePressed {
		f.muted = !f.muted
	}
	switch f.screen {
	case ScreenTitle:
		if in.LanguagePressed {
			f.switchLanguage()
		}
		f.choosePressure(in)
		f.chooseStart(in)
		switch {
		case in.HelpPressed:
			f.openHelp(ScreenTitle)
		case in.StartPressed() && !f.records.HelpSeen():
			f.openHelp(ScreenPlaying)
		case in.StartPressed():
			f.start()
		}
	case ScreenPlaying:
		if in.PausePressed || in.FocusLost {
			f.screen, f.menu = ScreenPaused, MenuResume
			return
		}
		f.session.Update(in)
		if f.session.Finished() {
			f.screen, f.epilogue = ScreenEpilogue, epilogueState{}
		}
	case ScreenPaused:
		f.updateMenu(in)
	case ScreenHelp:
		f.updateHelp(in)
	case ScreenEpilogue:
		f.updateEpilogue(in)
	case ScreenEnd:
		if in.StartPressed() {
			f.toTitle()
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
	f.session.SetPressure(f.Pressure())
	if f.startAt > 0 {
		f.session.GoToLevel(f.startAt)
	}
	f.screen = ScreenPlaying
}

func (f *Flow) updateMenu(in input.Intent) {
	switch {
	case in.ConfirmPressed:
		f.confirm()
	case in.PausePressed:
		f.resume()
	case in.UpPressed:
		f.moveMenu(-1)
	case in.DownPressed:
		f.moveMenu(1)
	case f.adjustVolume(in):
	case f.menu == MenuAudioOffset && in.DecreasePressed:
		f.setOffset(f.offset - AudioOffsetStep)
	case f.menu == MenuAudioOffset && in.IncreasePressed:
		f.setOffset(f.offset + AudioOffsetStep)
	case f.menu == MenuLanguage && (in.DecreasePressed || in.IncreasePressed):
		f.switchLanguage()
	}
}

func (f *Flow) resume() {
	f.screen = ScreenPlaying
	f.session.Resume()
}

func (f *Flow) setOffset(ms int) {
	f.offset = clampOffset(ms)
	f.records.setAudioOffset(f.offset)
}

func clampOffset(ms int) int {
	return max(-AudioOffsetLimit, min(AudioOffsetLimit, ms))
}

func (f *Flow) confirm() {
	switch f.menu {
	case MenuResume:
		f.resume()
	case MenuHelp:
		f.openHelp(ScreenPaused)
	case MenuAssist:
		f.assist = !f.assist
		f.session.SetAssist(f.assist)
	case MenuEffects:
		f.toggleEffects()
	case MenuLanguage:
		f.switchLanguage()
	case MenuSkip:
		f.session.GoToLevel(f.session.LevelIndex() + 1)
		f.resume()
	case MenuTitle:
		f.toTitle()
	}
}

func (f *Flow) Pressure() int {
	if !f.records.PressureUnlocked() {
		return 0
	}
	return f.pressure
}

func (f *Flow) choosePressure(in input.Intent) {
	if !f.records.PressureUnlocked() {
		return
	}
	switch {
	case in.DecreasePressed:
		f.pressure = max(f.pressure-1, 0)
	case in.IncreasePressed:
		f.pressure = min(f.pressure+1, MaxPressure)
	}
}
