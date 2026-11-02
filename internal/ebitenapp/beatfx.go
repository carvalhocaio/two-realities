package ebitenapp

import (
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	pulseAlpha     = 0.8
	pulseSharpness = 3
	countdownScale = 2
	countdownGap   = 3
)

func drawBeatPulse(dst *ebiten.Image, phase float64, pal palette) {
	strength := math.Pow(1-phase, pulseSharpness)
	fillRect(dst, geom.Rect{Y: HUDHeight - 1, W: ScreenWidth, H: 1}, pal.text, float32(strength*pulseAlpha))
}

func queuedBlink(s *game.Session, tick int) bool {
	return s.ShiftQueued() && (tick/(2*blinkTicks))%2 == 0
}

func shownWorld(s *game.Session, tick int) world.World {
	if queuedBlink(s, tick) {
		return s.Active().Other()
	}
	return s.Active()
}

func drawCountdown(dst *ebiten.Image, s *game.Session) {
	n := s.AutoCountdown()
	if n == 0 {
		return
	}
	msg := strconv.Itoa(n)
	box := s.PlayerBox()
	w := scaledTextWidth(msg, countdownScale)
	y := box.Y - countdownGap - glyphHeight*countdownScale
	drawScaledText(dst, msg, math.Round(box.X+box.W/2-w/2), y, countdownScale, palettes[s.Active()].text, 1)
}
