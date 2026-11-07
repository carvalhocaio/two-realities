package game

import "github.com/carvalhocaio/two-realities/internal/input"

func (f *Flow) openHelp(back Screen) {
	f.helpBack = back
	f.screen = ScreenHelp
}

func (f *Flow) updateHelp(in input.Intent) {
	if !in.ShiftPressed && !in.ConfirmPressed && !in.PausePressed && !in.HelpPressed {
		return
	}
	f.records.markHelpSeen()
	if f.helpBack == ScreenPlaying {
		f.start()
		return
	}
	f.screen = f.helpBack
}
