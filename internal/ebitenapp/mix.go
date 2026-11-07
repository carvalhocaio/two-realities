package ebitenapp

import "github.com/carvalhocaio/two-realities/internal/game"

const reducedFlash = 0.25

func mixVolume(base float64, muted bool, level int) float64 {
	if muted {
		return 0
	}
	x := float64(max(0, min(game.VolumeSteps, level))) / game.VolumeSteps
	return base * x * x
}

func effectScale(reduced bool) (shake, flash float64) {
	if reduced {
		return 0, reducedFlash
	}
	return 1, 1
}
