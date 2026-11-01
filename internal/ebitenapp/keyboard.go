package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/carvalhocaio/two-realities/internal/input"
)

var (
	leftKeys    = []ebiten.Key{ebiten.KeyA, ebiten.KeyArrowLeft}
	rightKeys   = []ebiten.Key{ebiten.KeyD, ebiten.KeyArrowRight}
	jumpKeys    = []ebiten.Key{ebiten.KeySpace, ebiten.KeyW, ebiten.KeyArrowUp}
	shiftKeys   = []ebiten.Key{ebiten.KeyShiftLeft, ebiten.KeyShiftRight}
	restartKeys = []ebiten.Key{ebiten.KeyR}
)

type Keyboard struct{}

var _ input.Source = Keyboard{}

func (Keyboard) Read() input.Intent {
	return input.Intent{
		Left:           anyPressed(leftKeys),
		Right:          anyPressed(rightKeys),
		JumpPressed:    anyJustPressed(jumpKeys),
		ShiftPressed:   anyJustPressed(shiftKeys),
		RestartPressed: anyJustPressed(restartKeys),
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
