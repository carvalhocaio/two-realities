package input

type PadButton int

const (
	PadA PadButton = iota
	PadB
	PadX
	PadY
	PadLB
	PadRB
	PadBack
	PadStart
	PadUp
	PadDown
	PadLeft
	PadRight
	PadButtonCount
)

const StickDeadZone = 0.5

type PadState struct {
	Connected      bool
	Buttons        [PadButtonCount]bool
	StickX, StickY float64
}

func (p PadState) Held(b PadButton) bool {
	switch b {
	case PadLeft:
		return p.Buttons[b] || p.StickX < -StickDeadZone
	case PadRight:
		return p.Buttons[b] || p.StickX > StickDeadZone
	case PadUp:
		return p.Buttons[b] || p.StickY < -StickDeadZone
	case PadDown:
		return p.Buttons[b] || p.StickY > StickDeadZone
	}
	return p.Buttons[b]
}

func (p PadState) Combine(o PadState) PadState {
	out := PadState{Connected: p.Connected || o.Connected, StickX: strongest(p.StickX, o.StickX), StickY: strongest(p.StickY, o.StickY)}
	for b := range PadButtonCount {
		out.Buttons[b] = p.Buttons[b] || o.Buttons[b]
	}
	return out
}

func strongest(a, b float64) float64 {
	if b*b > a*a {
		return b
	}
	return a
}

func PadIntent(prev, cur PadState) Intent {
	pressed := func(buttons ...PadButton) bool {
		for _, b := range buttons {
			if cur.Held(b) && !prev.Held(b) {
				return true
			}
		}
		return false
	}
	return Intent{
		Left:            cur.Held(PadLeft),
		Right:           cur.Held(PadRight),
		JumpPressed:     pressed(PadA),
		ShiftPressed:    pressed(PadB, PadX, PadLB, PadRB),
		RestartPressed:  pressed(PadBack),
		PausePressed:    pressed(PadStart),
		UpPressed:       pressed(PadUp),
		DownPressed:     pressed(PadDown),
		ConfirmPressed:  pressed(PadA),
		DecreasePressed: pressed(PadLeft),
		IncreasePressed: pressed(PadRight),
		HelpPressed:     pressed(PadY),
	}
}

func Merge(a, b Intent) Intent {
	return Intent{
		Left:            a.Left || b.Left,
		Right:           a.Right || b.Right,
		JumpPressed:     a.JumpPressed || b.JumpPressed,
		ShiftPressed:    a.ShiftPressed || b.ShiftPressed,
		RestartPressed:  a.RestartPressed || b.RestartPressed,
		PausePressed:    a.PausePressed || b.PausePressed,
		UpPressed:       a.UpPressed || b.UpPressed,
		DownPressed:     a.DownPressed || b.DownPressed,
		ConfirmPressed:  a.ConfirmPressed || b.ConfirmPressed,
		MutePressed:     a.MutePressed || b.MutePressed,
		FocusLost:       a.FocusLost || b.FocusLost,
		DecreasePressed: a.DecreasePressed || b.DecreasePressed,
		IncreasePressed: a.IncreasePressed || b.IncreasePressed,
		LanguagePressed: a.LanguagePressed || b.LanguagePressed,
		HelpPressed:     a.HelpPressed || b.HelpPressed,
	}
}
