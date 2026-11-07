package ebitenapp

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/camera"
	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/world"
)

type worldView struct {
	cam     camera.Camera
	epoch   int
	session *game.Session
	canvas  *ebiten.Image
	ghosts  [world.Count]*ebiten.Image
	level   int
	wipe    *wipe
	towers  []towerLayer
	fx      playerFX
	reduced bool
}

func newWorldView() *worldView {
	return &worldView{cam: camera.New(ScreenWidth, ViewHeight), level: -1, wipe: newWipe()}
}

func (v *worldView) follow(s *game.Session) {
	v.fx.update(s)
	bounds := s.Level().Bounds()
	if s != v.session || s.Epoch() != v.epoch {
		v.session, v.epoch = s, s.Epoch()
		v.cam.Snap(s.PlayerBox(), bounds)
		return
	}
	v.cam.Follow(s.PlayerBox(), bounds)
}

func (v *worldView) prepare(s *game.Session) {
	if v.level == s.LevelIndex() && v.canvas != nil {
		return
	}
	v.level = s.LevelIndex()
	b := s.Level().Bounds()
	v.canvas = ebiten.NewImage(int(b.W), int(b.H))
	v.ghosts = [world.Count]*ebiten.Image{}
	v.towers = newBackdrop(v.level, b.W, b.H)
}

func (v *worldView) ghost(s *game.Session) *ebiten.Image {
	active := s.Active()
	if v.ghosts[active] == nil {
		inactive := active.Other()
		lvl := s.Level()
		v.ghosts[active] = newGhostImage(lvl.Layer(inactive), lvl.Spikes[inactive], palettes[inactive].tile, palettes[inactive].tile)
	}
	return v.ghosts[active]
}

func (v *worldView) draw(screen *ebiten.Image, s *game.Session, phase float64, tick int) {
	v.prepare(s)
	v.canvas.Fill(palettes[s.Active()].background)
	if s.Active() == world.Grey {
		drawBackdrop(v.canvas, v.towers, v.cam.X)
		drawVignette(v.canvas, math.Round(v.cam.X), math.Round(v.cam.Y))
	}
	drawScene(v.canvas, s, v.ghost(s), ghostAlpha(s, phase), &v.fx, tick)
	drawCountdown(v.canvas, s)

	blocked := s.Feedback().Shake
	dying := deathStrength(s)
	shake, flash := effectScale(v.reduced)
	sx, sy := shakeOffset(max(blocked, dying) * shake)

	screen.Fill(palettes[s.Active()].background)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(math.Round(-v.cam.X+sx), math.Round(-v.cam.Y+sy)+HUDHeight)
	screen.DrawImage(v.canvas, op)

	viewport := geom.Rect{Y: HUDHeight, W: ScreenWidth, H: ViewHeight}
	drawFlash(screen, viewport, dangerColor, blocked*dangerColorAlpha*flash)
	drawFlash(screen, viewport, deathFlash, dying*deathFlashAlpha*flash)

	v.wipe.begin(s.Active(), s.Epoch(), tick)
	v.wipe.remember(screen)
	if !v.reduced {
		v.wipe.draw(screen, tick)
	}
}

func deathStrength(s *game.Session) float64 {
	if s.Death() == game.Alive {
		return 0
	}
	return 1 - s.DeathProgress()
}

func (v *worldView) onScreen(r geom.Rect) geom.Rect {
	r.X -= math.Round(v.cam.X)
	r.Y += HUDHeight - math.Round(v.cam.Y)
	return r
}

func ghostAlpha(s *game.Session, phase float64) float32 {
	if !s.Has(game.BeatGhost) {
		return 1
	}
	return float32(max(0, 1-phase*beatGhostFade))
}
