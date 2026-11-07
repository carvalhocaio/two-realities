package ebitenapp

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	tunes "github.com/carvalhocaio/two-realities/internal/music"
	"github.com/carvalhocaio/two-realities/internal/platform"
	"github.com/carvalhocaio/two-realities/internal/save"
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
	sfx         *sfx
	input       input.Source
	view        *worldView
	textbox     textbox
	title       string
	tick        int
	notice      string
	noticeTimer int
	showTimer   bool
}

func New(levels []level.Level, store save.Storage, source input.Source) (*App, error) {
	ctx := audio.NewContext(tunes.SampleRate)
	m, err := newMusic(ctx)
	if err != nil {
		return nil, err
	}
	flow := game.NewFlow(m, store, levels...)
	flow.UseLocale(platform.Locale())
	return &App{flow: flow, music: m, sfx: newSFX(ctx), input: source, view: newWorldView()}, nil
}

func (a *App) Update() error {
	a.tick++
	a.noticeTimer = max(a.noticeTimer-1, 0)
	a.flow.Update(a.input.Read())
	if a.flow.Screen() == game.ScreenPlaying {
		a.music.start()
	}
	a.music.setMuted(a.flow.Muted())
	if s := a.flow.Session(); s != nil {
		a.sfx.play(s.Events(), a.flow.Muted())
	}
	a.music.setOffset(a.flow.AudioOffsetMs())
	a.music.setTrack(musicTrack(a.flow))
	if anyJustPressed(timerKeys) {
		a.showTimer = !a.showTimer
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
		drawTitleScreen(screen, a.flow.Language(), a.flow.StorageNotice(), pressureChoice{a.flow.Records().PressureUnlocked(), a.flow.Pressure()}, a.tick)
	case game.ScreenEnd:
		drawEndScreen(screen, a.flow.Session().Summary(), a.flow.Language(), a.tick)
	default:
		s := a.flow.Session()
		a.view.draw(screen, s, a.music.Phase(), a.tick)
		drawHUD(screen, s, a.tick, a.hudState())
		a.textbox.draw(screen, s, a.view.onScreen(s.PlayerBox()), a.flow.Language())
		if a.flow.Screen() == game.ScreenPaused {
			drawPauseMenu(screen, a.flow, a.music.Phase())
		}
	}
}

func (a *App) hudState() hudState {
	h := hudState{phase: a.music.Phase(), muted: a.music.muted, timer: a.showTimer, lang: a.flow.Language()}
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

func musicTrack(f *game.Flow) tunes.Track {
	if f.Screen() != game.ScreenPlaying && f.Screen() != game.ScreenPaused {
		return tunes.ActOne
	}
	s := f.Session()
	if w, ok := s.Warden(); ok && w.Lit < level.LampCount {
		return tunes.WardenOne + tunes.Track(w.Lit)
	}
	return tunes.ActOne + tunes.Track(min(max(s.Level().Act, 1), 3)-1)
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

func Run(levels []level.Level, store save.Storage) error {
	ebiten.SetWindowSize(ScreenWidth*windowScale, ScreenHeight*windowScale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)
	app, err := New(levels, store, NewKeyboard())
	if err != nil {
		return err
	}
	return ebiten.RunGame(app)
}
