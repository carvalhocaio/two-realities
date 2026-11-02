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
}

func newWorldView() *worldView {
	return &worldView{cam: camera.New(ScreenWidth, ViewHeight), level: -1}
}

func (v *worldView) follow(s *game.Session) {
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

func (v *worldView) draw(screen *ebiten.Image, s *game.Session, tick int) {
	v.prepare(s)
	drawScene(v.canvas, s, v.ghost(s), tick)
	drawCountdown(v.canvas, s)

	blocked := s.Feedback().Shake
	dying := deathStrength(s)
	sx, sy := shakeOffset(max(blocked, dying))

	screen.Fill(palettes[s.Active()].background)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(math.Round(-v.cam.X+sx), math.Round(-v.cam.Y+sy)+HUDHeight)
	screen.DrawImage(v.canvas, op)

	viewport := geom.Rect{Y: HUDHeight, W: ScreenWidth, H: ViewHeight}
	drawFlash(screen, viewport, dangerColor, blocked*dangerColorAlpha)
	drawFlash(screen, viewport, deathFlash, dying*deathFlashAlpha)
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
