package ebitenapp

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 180
	HUDHeight    = 10
	ViewHeight   = ScreenHeight - HUDHeight
	windowScale  = 4
	shakePixels  = 3
	noticeTicks  = 60
)

type App struct {
	flow        *game.Flow
	music       *music
	input       input.Source
	view        *worldView
	title       string
	tick        int
	notice      string
	noticeTimer int
}

func New(levels []level.Level, source input.Source) (*App, error) {
	m, err := newMusic()
	if err != nil {
		return nil, err
	}
	return &App{flow: game.NewFlow(m, levels...), music: m, input: source, view: newWorldView()}, nil
}

func (a *App) Update() error {
	a.tick++
	a.noticeTimer = max(a.noticeTimer-1, 0)
	a.flow.Update(a.input.Read())
	if a.flow.Screen() == game.ScreenPlaying {
		a.music.start()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		a.music.toggleMute()
	}
	handleDevKeys(a)
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
		drawHUD(screen, a.flow.Session(), a.tick, a.hudState())
	}
}

func (a *App) hudState() hudState {
	h := hudState{phase: a.music.Phase(), muted: a.music.muted}
	if a.noticeTimer > 0 {
		h.notice = a.notice
	}
	return h
}

func (a *App) showNotice(msg string) {
	a.notice, a.noticeTimer = msg, noticeTicks
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

func Run(levels []level.Level) error {
	ebiten.SetWindowSize(ScreenWidth*windowScale, ScreenHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)
	app, err := New(levels, Keyboard{})
	if err != nil {
		return err
	}
	return ebiten.RunGame(app)
}
