package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/platform"
)

var (
	leftKeys    = []ebiten.Key{ebiten.KeyA, ebiten.KeyArrowLeft}
	rightKeys   = []ebiten.Key{ebiten.KeyD, ebiten.KeyArrowRight}
	jumpKeys    = []ebiten.Key{ebiten.KeySpace, ebiten.KeyW, ebiten.KeyArrowUp}
	shiftKeys   = []ebiten.Key{ebiten.KeyShiftLeft, ebiten.KeyShiftRight}
	restartKeys = []ebiten.Key{ebiten.KeyR}
	pauseKeys   = []ebiten.Key{ebiten.KeyEscape, ebiten.KeyP, ebiten.KeyEnter}
	upKeys      = []ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyW}
	downKeys    = []ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyS}
	confirmKeys = []ebiten.Key{ebiten.KeyEnter, ebiten.KeySpace, ebiten.KeyShiftLeft, ebiten.KeyShiftRight}
	muteKeys    = []ebiten.Key{ebiten.KeyM}
	timerKeys   = []ebiten.Key{ebiten.KeyT}
	langKeys    = []ebiten.Key{ebiten.KeyL}
	helpKeys    = []ebiten.Key{ebiten.KeyH}
	lessKeys    = []ebiten.Key{ebiten.KeyMinus, ebiten.KeyNumpadSubtract, ebiten.KeyArrowLeft, ebiten.KeyA}
	moreKeys    = []ebiten.Key{ebiten.KeyEqual, ebiten.KeyNumpadAdd, ebiten.KeyArrowRight, ebiten.KeyD}
)

type Keyboard struct {
	signals *platform.Signals
}

var _ input.Source = Keyboard{}

func NewKeyboard() Keyboard {
	k := Keyboard{signals: &platform.Signals{}}
	platform.Listen(k.signals)
	return k
}

func (k Keyboard) Read() input.Intent {
	k.signals.SetFocused(ebiten.IsFocused())
	return input.Intent{
		Left:            anyPressed(leftKeys),
		Right:           anyPressed(rightKeys),
		JumpPressed:     anyJustPressed(jumpKeys),
		ShiftPressed:    anyJustPressed(shiftKeys),
		RestartPressed:  anyJustPressed(restartKeys),
		PausePressed:    anyJustPressed(pauseKeys),
		UpPressed:       anyJustPressed(upKeys),
		DownPressed:     anyJustPressed(downKeys),
		ConfirmPressed:  anyJustPressed(confirmKeys),
		MutePressed:     anyJustPressed(muteKeys),
		FocusLost:       k.signals.Poll(),
		DecreasePressed: anyJustPressed(lessKeys),
		IncreasePressed: anyJustPressed(moreKeys),
		LanguagePressed: anyJustPressed(langKeys),
		HelpPressed:     anyJustPressed(helpKeys),
	}
}

func anyPressed(keys []ebiten.Key) bool {
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

func anyJustPressed(keys []ebiten.Key) bool {
	for _, k := range keys {
		if inpututil.IsKeyJustPressed(k) {
			return true
		}
	}
	return false
}
