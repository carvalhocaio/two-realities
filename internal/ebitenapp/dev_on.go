//go:build dev

package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/carvalhocaio/two-realities/internal/game"
)

func handleDevKeys(f *game.Flow) {
	if f.Screen() != game.ScreenPlaying {
		return
	}
	s := f.Session()
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyN):
		s.GoToLevel(s.LevelIndex() + 1)
	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		s.GoToLevel(s.LevelIndex() - 1)
	}
}
