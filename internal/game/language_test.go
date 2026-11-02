package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

var langKey = input.Intent{LanguagePressed: true}

func toLanguageItem() []input.Intent {
	return then(one(shift), one(pauseKey), ticks(downKey, int(MenuLanguage)))
}

func TestLanguage(t *testing.T) {
	tests := []struct {
		name   string
		saved  string
		locale string
		inputs []input.Intent
		want   i18n.Lang
	}{
		{"english browser", "", "en-US", nil, i18n.EN},
		{"brazilian browser", "", "pt-BR", nil, i18n.PT},
		{"saved choice beats the browser", "en", "pt-BR", nil, i18n.EN},
		{"saved portuguese", "pt", "en-US", nil, i18n.PT},
		{"garbage saved value falls back to detection", "xx", "pt-BR", nil, i18n.PT},
		{"l on the title switches", "", "en-US", one(langKey), i18n.PT},
		{"l twice switches back", "", "en-US", then(one(langKey), one(langKey)), i18n.EN},
		{"l is ignored while playing", "", "en-US", then(one(shift), one(langKey)), i18n.EN},
		{"menu item confirms", "", "en-US", then(toLanguageItem(), one(input.Intent{ConfirmPressed: true})), i18n.PT},
		{"menu item arrows", "", "en-US", then(toLanguageItem(), one(moreKey)), i18n.PT},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlow(&rhythm.Manual{}, &save.Memory{Stored: save.Data{Language: tt.saved}}, mustParse(t, timedLevel))
			f.UseLocale(tt.locale)
			for _, in := range tt.inputs {
				f.Update(in)
			}
			if f.Language() != tt.want {
				t.Fatalf("language %v, want %v", f.Language(), tt.want)
			}
		})
	}
}

func TestLanguageChoicePersists(t *testing.T) {
	store := &save.Memory{}
	f := NewFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
	f.UseLocale("en-US")
	f.Update(langKey)
	again := NewFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
	again.UseLocale("en-US")
	if again.Language() != i18n.PT || store.Stored.Language != "pt" {
		t.Fatalf("reloaded %v stored %q", again.Language(), store.Stored.Language)
	}
}
