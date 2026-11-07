package font

import (
	"testing"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		r      rune
		ok     bool
		accent bool
	}{
		{' ', true, false},
		{'A', true, false},
		{'7', true, false},
		{'Ã', true, true},
		{'Ç', true, true},
		{'@', false, false},
		{'a', false, false},
	}
	for _, tt := range tests {
		g, ok := Lookup(tt.r)
		if ok != tt.ok || (g.Above != "" || g.Below != "") != tt.accent {
			t.Errorf("%q: ok %v glyph %+v", tt.r, ok, g)
		}
	}
	if g, _ := Lookup('É'); g.Rows != base['E'] {
		t.Error("accented letters reuse the base glyph")
	}
}

func TestMissing(t *testing.T) {
	if got := Missing("OLÁ, MUNDO!"); len(got) != 0 {
		t.Fatalf("missing %q", string(got))
	}
	if got := Missing("A@B#c"); string(got) != "@#c" {
		t.Fatalf("missing %q, want @#c", string(got))
	}
}

func TestUpper(t *testing.T) {
	for in, want := range map[string]string{
		"press shift": "PRESS SHIFT",
		"ação":        "AÇÃO",
		"3x faster":   "3x FASTER",
		"":            "",
	} {
		if got := Upper(in); got != want {
			t.Errorf("Upper(%q) = %q, want %q", in, got, want)
		}
	}
}
