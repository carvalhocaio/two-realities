package ebitenapp

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	tunes "github.com/carvalhocaio/two-realities/internal/music"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

func TestHelpCardFitsTheScreen(t *testing.T) {
	for _, lang := range []i18n.Lang{i18n.EN, i18n.PT} {
		for _, row := range helpRows {
			msg := lang.T(row.key)
			if w := helpTextX + textWidth(msg); w > ScreenWidth-helpIconX {
				t.Errorf("%v %q is %vpx wide", lang, msg, w)
			}
		}
		for _, key := range []i18n.Key{i18n.HelpTitle, i18n.HelpPrompt, i18n.TitleLanguage} {
			if w := scaledTextWidth(lang.T(key), subtitleScale); key == i18n.HelpTitle && w > ScreenWidth {
				t.Errorf("%v title is %vpx wide", lang, w)
			}
			if w := textWidth(lang.T(key)); w > ScreenWidth-8 {
				t.Errorf("%v %q is %vpx wide", lang, lang.T(key), w)
			}
		}
	}
	if last := helpFirstY + (len(helpRows)-1)*helpStep + 10; last >= helpPromptY {
		t.Errorf("rows end at %d, over the prompt at %d", last, helpPromptY)
	}
}

func TestMusicTrackOnTheHelpCard(t *testing.T) {
	lvl, err := level.Parse("fragments: 0\nact: 2\n[yellow]\n#....#\n#P..E#\n######\n[grey]\n#....#\n#....#\n######\n")
	if err != nil {
		t.Fatal(err)
	}
	f := game.NewFlow(&rhythm.Manual{}, &save.Memory{}, lvl)
	f.Update(input.Intent{HelpPressed: true})
	if got := musicTrack(f); got != tunes.ActOne {
		t.Fatalf("help from the title plays %v", got)
	}
	f.Update(input.Intent{HelpPressed: true})
	f.Update(input.Intent{ShiftPressed: true})
	f.Update(input.Intent{ShiftPressed: true})
	f.Update(input.Intent{PausePressed: true})
	f.Update(input.Intent{DownPressed: true})
	f.Update(input.Intent{ConfirmPressed: true})
	if f.Screen() != game.ScreenHelp {
		t.Fatalf("screen %v, want help", f.Screen())
	}
	if got := musicTrack(f); got != tunes.ActTwo {
		t.Fatalf("help over a paused act 2 level plays %v, want act 2", got)
	}
}

func TestPauseMenuFitsWithEveryItem(t *testing.T) {
	if bottom := pauseFirstY + int(game.MenuCount)*pauseStep + pauseHintGap + 5; bottom > ScreenHeight {
		t.Fatalf("pause hint ends at %d with all %d items", bottom, game.MenuCount)
	}
	if pauseFirstY < pauseTitleY+5*subtitleScale+4 {
		t.Fatal("first pause item overlaps the title")
	}
	for _, lang := range []i18n.Lang{i18n.EN, i18n.PT} {
		if w := textWidth(lang.T(i18n.OffsetHint)); w > ScreenWidth {
			t.Errorf("%v offset hint is %vpx wide", lang, w)
		}
	}
}

func TestTitlePadLineFits(t *testing.T) {
	for _, lang := range []i18n.Lang{i18n.EN, i18n.PT} {
		if w := textWidth(lang.T(i18n.TitleControlsPad)); w > ScreenWidth-8 {
			t.Errorf("%v pad line is %vpx wide", lang, w)
		}
	}
}

func TestEveryPadButtonHasItsOwnMapping(t *testing.T) {
	seen := map[ebiten.StandardGamepadButton]bool{}
	for b, sb := range padButtons {
		if seen[sb] {
			t.Errorf("pad button %d reuses %v", b, sb)
		}
		seen[sb] = true
	}
}

func TestStartLabelShowsTheWholeCampaign(t *testing.T) {
	tests := []struct {
		lang i18n.Lang
		s    startChoice
		want string
	}{
		{i18n.EN, startChoice{level: 4, count: 5, total: 13}, "START AT LEVEL 5/13 - UP/DOWN"},
		{i18n.PT, startChoice{level: 4, count: 5, total: 13}, "COMEÇAR NA FASE 5/13 - CIMA/BAIXO"},
		{i18n.EN, startChoice{level: 0, count: 13, total: 13}, "START AT LEVEL 1/13 - UP/DOWN"},
	}
	for _, tt := range tests {
		if got := startLabel(tt.lang, tt.s); got != tt.want {
			t.Errorf("%v %+v: %q, want %q", tt.lang, tt.s, got, tt.want)
		}
	}
}
