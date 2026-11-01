package input

type Intent struct {
	Left, Right    bool
	JumpPressed    bool
	ShiftPressed   bool
	RestartPressed bool
}

type Source interface {
	Read() Intent
}
