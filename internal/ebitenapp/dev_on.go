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
	case inpututil.IsKeyJustPressed(ebiten.KeyMinus):
		a.music.nudgeLatency(-1)
		a.showNotice(fmt.Sprintf("LAT %dMS", a.music.latency.Milliseconds()))
	case inpututil.IsKeyJustPressed(ebiten.KeyEqual):
		a.music.nudgeLatency(1)
		a.showNotice(fmt.Sprintf("LAT %dMS", a.music.latency.Milliseconds()))
	}
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
