package ebitenapp

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/input"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 180
	HUDHeight    = 10
	ViewHeight   = ScreenHeight - HUDHeight
	windowScale  = 4
	shakePixels  = 3
)

type App struct {
	flow  *game.Flow
	input input.Source
	view  *worldView
	title string
	tick  int
}

func New(flow *game.Flow, source input.Source) *App {
	return &App{flow: flow, input: source, view: newWorldView()}
}

func (a *App) Update() error {
	a.tick++
	a.flow.Update(a.input.Read())
	handleDevKeys(a.flow)
	if a.flow.Screen() == game.ScreenPlaying {
		a.view.follow(a.flow.Session())
	}
	a.setTitle(windowTitle(a.flow))
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	switch a.flow.Screen() {
	case game.ScreenTitle:
		drawTitleScreen(screen, a.tick)
	case game.ScreenEnd:
		drawEndScreen(screen, a.tick)
	default:
		a.view.draw(screen, a.flow.Session(), a.tick)
		drawHUD(screen, a.flow.Session(), a.tick)
	}
}

func (a *App) Layout(int, int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (a *App) setTitle(title string) {
	if title != a.title {
		a.title = title
		ebiten.SetWindowTitle(title)
	}
}

func windowTitle(f *game.Flow) string {
	if f.Screen() != game.ScreenPlaying {
		return "Two Realities"
	}
	return "Two Realities — " + f.Session().Active().String()
}

func shakeOffset(strength float64) (float64, float64) {
	if strength <= 0 {
		return 0, 0
	}
	amp := shakePixels * strength
	return (rand.Float64()*2 - 1) * amp, (rand.Float64()*2 - 1) * amp
}

func Run(flow *game.Flow) error {
	ebiten.SetWindowSize(ScreenWidth*windowScale, ScreenHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)
	return ebiten.RunGame(New(flow, Keyboard{}))
}
