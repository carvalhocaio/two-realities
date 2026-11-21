package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
)

const levelFadeTicks = 14

type levelFade struct {
	from    *ebiten.Image
	session *game.Session
	level   int
	start   int
	running bool
}

func (f *levelFade) begin(s *game.Session, last *ebiten.Image, tick int) {
	advanced := s == f.session && s.LevelIndex() != f.level
	f.session, f.level = s, s.LevelIndex()
	if !advanced {
		return
	}
	if f.from == nil {
		f.from = ebiten.NewImage(ScreenWidth, ViewHeight)
	}
	f.from.Clear()
	f.from.DrawImage(last, nil)
	f.start, f.running = tick, true
}

func (f *levelFade) draw(screen *ebiten.Image, tick int) {
	if !f.running {
		return
	}
	alpha := levelFadeAlpha(tick - f.start)
	if alpha <= 0 {
		f.running = false
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, HUDHeight)
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(f.from, op)
}

func levelFadeAlpha(age int) float32 {
	if age < 0 || age >= levelFadeTicks {
		return 0
	}
	return 1 - float32(age)/levelFadeTicks
}
