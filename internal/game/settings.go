package game

import "github.com/carvalhocaio/two-realities/internal/input"

const (
	VolumeSteps    = 10
	effectsReduced = "reduced"
	effectsFull    = "full"
)

func clampVolume(v int) int { return max(0, min(VolumeSteps, v)) }

func (f *Flow) MusicVolume() int    { return f.music }
func (f *Flow) SoundVolume() int    { return f.sound }
func (f *Flow) ReduceEffects() bool { return f.reduced }

func (f *Flow) UseReducedMotion(preferred bool) {
	if _, set := f.records.ReducedEffects(); !set {
		f.reduced = preferred
	}
}

func (f *Flow) adjustVolume(in input.Intent) bool {
	step := 0
	switch {
	case in.DecreasePressed:
		step = -1
	case in.IncreasePressed:
		step = 1
	default:
		return false
	}
	switch f.menu {
	case MenuMusic:
		f.music = clampVolume(f.music + step)
	case MenuSounds:
		f.sound = clampVolume(f.sound + step)
	default:
		return false
	}
	f.records.setVolumes(f.music, f.sound)
	return true
}

func (f *Flow) toggleEffects() {
	f.reduced = !f.reduced
	f.records.setReducedEffects(f.reduced)
}
