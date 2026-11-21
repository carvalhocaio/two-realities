package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/i18n"
)

func TestTitleWorldsFlankTheTitle(t *testing.T) {
	for _, lang := range []i18n.Lang{i18n.EN, i18n.PT} {
		w := scaledTextWidth(lang.T(i18n.Title), titleScale)
		textLeft, textRight := (ScreenWidth-w)/2, (ScreenWidth+w)/2
		left := textLeft - titleBlockGap - titleBlock
		right := textRight + titleBlockGap
		if left < 0 || right+titleBlock > ScreenWidth {
			t.Fatalf("%v: blocks at %v and %v leave the screen", lang, left, right)
		}
	}
	if titleShifted(0) || !titleShifted(titleSwap) || titleShifted(2*titleSwap) {
		t.Fatal("title blocks should swap every titleSwap ticks")
	}
}
