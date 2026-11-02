//go:build dev

package ebitenapp

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/carvalhocaio/two-realities/internal/game"
)

func handleDevKeys(a *App) {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyBracketLeft):
		a.music.nudgeLatency(-1)
		a.showNotice(fmt.Sprintf("LAT %dMS", a.music.latency.Milliseconds()))
	case inpututil.IsKeyJustPressed(ebiten.KeyBracketRight):
		a.music.nudgeLatency(1)
		a.showNotice(fmt.Sprintf("LAT %dMS", a.music.latency.Milliseconds()))
	}
	if a.flow.Screen() != game.ScreenPlaying {
		return
	}
	s := a.flow.Session()
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyN):
		s.GoToLevel(s.LevelIndex() + 1)
	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		s.GoToLevel(s.LevelIndex() - 1)
	}
}
