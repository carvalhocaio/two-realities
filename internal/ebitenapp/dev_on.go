//go:build dev

package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/carvalhocaio/two-realities/internal/game"
)

func handleDevKeys(a *App) {
	if a.flow.Screen() != game.ScreenPlaying {
		return
	}
	s := a.flow.Session()
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyBracketRight):
		s.GoToLevel(s.LevelIndex() + 1)
	case inpututil.IsKeyJustPressed(ebiten.KeyBracketLeft):
		s.GoToLevel(s.LevelIndex() - 1)
	}
}
