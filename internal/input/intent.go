package input

type Intent struct {
	Left, Right     bool
	JumpPressed     bool
	ShiftPressed    bool
	RestartPressed  bool
	PausePressed    bool
	UpPressed       bool
	DownPressed     bool
	ConfirmPressed  bool
	MutePressed     bool
	FocusLost       bool
	DecreasePressed bool
	IncreasePressed bool
	LanguagePressed bool
	HelpPressed     bool
}

func (i Intent) StartPressed() bool { return i.ShiftPressed || i.ConfirmPressed }

func (i Intent) Acts() bool {
	return i.Left || i.Right || i.JumpPressed || i.ShiftPressed
}

type Source interface {
	Read() Intent
}
