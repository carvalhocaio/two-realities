package i18n

import (
	"strings"
	"testing"

	"github.com/carvalhocaio/two-realities/internal/font"
)

func TestEveryKeyIsTranslated(t *testing.T) {
	for l := range LangCount {
		for k := range KeyCount {
			text, ok := tables[l][k]
			if !ok || text == "" {
				t.Errorf("%s: key %d has no text", l.Code(), k)
			}
			if missing := font.Missing(withoutVerbs(text)); len(missing) > 0 {
				t.Errorf("%s: %q uses glyphs the font lacks: %q", l.Code(), text, string(missing))
			}
		}
	}
}

func withoutVerbs(text string) string {
	return strings.NewReplacer("%d", "0", "%s", "0").Replace(text)
}

func TestDetect(t *testing.T) {
	tests := []struct {
		locale string
		want   Lang
	}{
		{"pt-BR", PT},
		{"pt", PT},
		{"PT_br.UTF-8", PT},
		{"pt_PT", PT},
		{"en-US", EN},
		{"es", EN},
		{"", EN},
	}
	for _, tt := range tests {
		if got := Detect(tt.locale); got != tt.want {
			t.Errorf("Detect(%q) = %v, want %v", tt.locale, got, tt.want)
		}
	}
}

func TestParseRoundTrip(t *testing.T) {
	for l := range LangCount {
		got, ok := Parse(l.Code())
		if !ok || got != l {
			t.Errorf("Parse(%q) = %v %v", l.Code(), got, ok)
		}
	}
	if _, ok := Parse("klingon"); ok {
		t.Error("unknown codes must not parse")
	}
}

func TestFormatting(t *testing.T) {
	if got := PT.T(HudLevel, 3, 8); got != "FASE 3/8" {
		t.Fatalf("got %q", got)
	}
	if got := EN.T(MenuEffects, EN.T(On)); got != "REDUCE EFFECTS: ON" {
		t.Fatalf("got %q", got)
	}
}
