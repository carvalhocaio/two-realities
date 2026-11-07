package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/i18n"
)

func TestPadHints(t *testing.T) {
	for _, lang := range []i18n.Lang{i18n.EN, i18n.PT} {
		for key, pad := range padHints {
			if hint(lang, key, false) != lang.T(key) {
				t.Errorf("%v keyboard hint for %v changed", lang, key)
			}
			text := hint(lang, key, true)
			if text != lang.T(pad) || text == lang.T(key) {
				t.Errorf("%v pad hint for %v is %q", lang, key, text)
			}
			if w := helpTextX + textWidth(text); w > ScreenWidth-helpIconX {
				t.Errorf("%v %q is %vpx wide", lang, text, w)
			}
		}
	}
	if hint(i18n.EN, i18n.Title, true) != i18n.EN.T(i18n.Title) {
		t.Fatal("keys without a pad version stay the same")
	}
}
