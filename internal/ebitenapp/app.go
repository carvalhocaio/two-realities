package ebitenapp

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 180
	windowScale  = 4
	shakePixels  = 3
)

type App struct {
	session *game.Session
	input   input.Source
	canvas  *ebiten.Image
	shown   world.World
}

func New(session *game.Session, source input.Source) *App {
	return &App{
		session: session,
		input:   source,
		canvas:  ebiten.NewImage(ScreenWidth, ScreenHeight),
		shown:   -1,
	}
}

func (a *App) Update() error {
	a.session.Update(a.input.Read())
	if a.session.Active() != a.shown {
		a.shown = a.session.Active()
		ebiten.SetWindowTitle(worldLabel(a.shown))
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	drawScene(a.canvas, a.session)
	drawFlash(a.canvas, a.session.Feedback())

	screen.Fill(palettes[a.session.Active()].background)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(shakeOffset(a.session.Feedback()))
	screen.DrawImage(a.canvas, op)
}

func (a *App) Layout(int, int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func shakeOffset(strength float64) (float64, float64) {
	if strength <= 0 {
		return 0, 0
	}
	amp := shakePixels * strength
	return (rand.Float64()*2 - 1) * amp, (rand.Float64()*2 - 1) * amp
}

func worldLabel(w world.World) string {
	return "Two Realities — " + w.String()
}

func Run(session *game.Session) error {
	ebiten.SetWindowSize(ScreenWidth*windowScale, ScreenHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)
	return ebiten.RunGame(New(session, Keyboard{}))
}
