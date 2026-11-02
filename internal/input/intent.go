package input

type Intent struct {
	Left, Right    bool
	JumpPressed    bool
	ShiftPressed   bool
	RestartPressed bool
	PausePressed   bool
	UpPressed      bool
	DownPressed    bool
	ConfirmPressed bool
	MutePressed    bool
	FocusLost      bool
}

func (i Intent) Acts() bool {
	return i.Left || i.Right || i.JumpPressed || i.ShiftPressed
}

type Source interface {
	Read() Intent
}
