//go:build dev

package ebitenapp

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/save"
)

const TitleShot = -1

type Shot struct {
	Name  string
	Level int
	Every int
	Steps []Step
}

type Step struct {
	Ticks int
	In    input.Intent
}

type scripted struct{ in input.Intent }

func (s *scripted) Read() input.Intent { return s.in }

type shotRunner struct {
	app     *App
	src     *scripted
	shots   []Shot
	dir     string
	shot    int
	step    int
	inStep  int
	tick    int
	entered bool
	capture string
	err     error
}

func RunShots(levels []level.Level, shots []Shot, dir, locale string) error {
	ebiten.SetWindowSize(ScreenWidth*windowScale, ScreenHeight*windowScale)
	ebiten.SetTPS(60)
	src := &scripted{}
	app, err := New(levels, &save.Memory{}, src)
	if err != nil {
		return err
	}
	app.flow.UseLocale(locale)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	r := &shotRunner{app: app, src: src, shots: shots, dir: dir}
	if err := ebiten.RunGame(r); err != nil && !errors.Is(err, ebiten.Termination) {
		return err
	}
	return r.err
}

func (r *shotRunner) Update() error {
	if r.err != nil {
		return r.err
	}
	if r.shot >= len(r.shots) {
		return ebiten.Termination
	}
	sh := r.shots[r.shot]
	if !r.entered {
		r.enter(sh)
		return nil
	}
	r.src.in = r.intent(sh)
	if err := r.app.Update(); err != nil {
		return err
	}
	if sh.Every > 0 && r.tick%sh.Every == 0 {
		r.capture = fmt.Sprintf("%s_%04d.png", sh.Name, r.tick)
	}
	r.advance(sh)
	return nil
}

func (r *shotRunner) enter(sh Shot) {
	r.entered = true
	if sh.Level == TitleShot {
		return
	}
	if r.app.flow.Screen() == game.ScreenTitle {
		r.src.in = input.Intent{ShiftPressed: true}
		_ = r.app.Update()
	}
	r.app.flow.Session().GoToLevel(sh.Level)
}

func (r *shotRunner) intent(sh Shot) input.Intent {
	in := sh.Steps[r.step].In
	if r.inStep > 0 {
		return input.Intent{Left: in.Left, Right: in.Right}
	}
	return in
}

func (r *shotRunner) advance(sh Shot) {
	r.tick++
	r.inStep++
	if r.inStep < sh.Steps[r.step].Ticks {
		return
	}
	r.step, r.inStep = r.step+1, 0
	if r.step < len(sh.Steps) {
		return
	}
	r.shot, r.step, r.tick, r.entered = r.shot+1, 0, 0, false
}

func (r *shotRunner) Draw(screen *ebiten.Image) {
	r.app.Draw(screen)
	if r.capture == "" {
		return
	}
	img := image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))
	screen.ReadPixels(img.Pix)
	r.err = writePNG(filepath.Join(r.dir, r.capture), img)
	r.capture = ""
}

func (r *shotRunner) Layout(w, h int) (int, int) { return r.app.Layout(w, h) }

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
