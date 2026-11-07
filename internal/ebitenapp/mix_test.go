package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/i18n"
)

func TestMixVolume(t *testing.T) {
	tests := []struct {
		name  string
		muted bool
		level int
		want  float64
	}{
		{"full", false, game.VolumeSteps, 0.5},
		{"half the steps is a quarter of the gain", false, game.VolumeSteps / 2, 0.125},
		{"zero is silent", false, 0, 0},
		{"mute wins over any level", true, game.VolumeSteps, 0},
		{"levels above the top are clamped", false, game.VolumeSteps + 5, 0.5},
		{"negative levels are clamped", false, -3, 0},
	}
	for _, tt := range tests {
		if got := mixVolume(0.5, tt.muted, tt.level); got != tt.want {
			t.Errorf("%s: %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestEffectScale(t *testing.T) {
	if shake, flash := effectScale(false); shake != 1 || flash != 1 {
		t.Fatalf("full effects: shake %v flash %v", shake, flash)
	}
	if shake, flash := effectScale(true); shake != 0 || flash != reducedFlash {
		t.Fatalf("reduced effects: shake %v flash %v", shake, flash)
	}
}

func TestSettingsLabelsFit(t *testing.T) {
	for _, lang := range []i18n.Lang{i18n.EN, i18n.PT} {
		for _, label := range []string{
			lang.T(i18n.MenuMusic, volumeLabel(game.VolumeSteps)),
			lang.T(i18n.MenuSounds, volumeLabel(game.VolumeSteps)),
			lang.T(i18n.MenuEffects, onOff(false, lang)),
		} {
			if w := textWidth("- " + label + " -"); w > ScreenWidth-8 {
				t.Errorf("%q is %vpx wide", label, w)
			}
		}
	}
	if volumeLabel(7) != "7/10" {
		t.Fatalf("volume label %q", volumeLabel(7))
	}
}

func TestEpilogueFitsTheScreen(t *testing.T) {
	for _, lang := range []i18n.Lang{i18n.EN, i18n.PT} {
		for _, key := range game.EpilogueLines {
			if w := textWidth(lang.T(key)); w > ScreenWidth-8 {
				t.Errorf("%v %q is %vpx wide", lang, lang.T(key), w)
			}
		}
		if w := scaledTextWidth(lang.T(i18n.EpilogueTitle), subtitleScale); w > ScreenWidth {
			t.Errorf("%v epilogue title is %vpx wide", lang, w)
		}
	}
	if last := epilogueFirstY + (len(game.EpilogueLines)-1)*epilogueStep + 5; last >= epiloguePromptY {
		t.Fatalf("lines end at %d, over the prompt at %d", last, epiloguePromptY)
	}
}
