package level

import (
	"errors"
	"strings"
	"testing"

	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const head = "fragments: 0\n"

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{"valid", head + "[yellow]\n#P.#\n####\n[grey]\n#..E\n#..#\n", nil},
		{"valid with blank lines and CRLF", "\r\n" + "fragments: 2\r\n\r\n[yellow]\r\n#P.#\r\n\r\n[grey]\r\n#.E#\r\n", nil},
		{"markers split across layers", head + "[yellow]\n#E.#\n[grey]\n#.P#\n", nil},
		{"fragment in grey", head + "[yellow]\n#PE#\n[grey]\n#YY#\n", nil},
		{"missing header", "[yellow]\n#PE#\n[grey]\n####\n", ErrMissingHeader},
		{"empty file", "", ErrMissingHeader},
		{"header without value", "fragments:\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"header not a number", "fragments: lots\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"header negative", "fragments: -1\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"header unknown key", "lives: 3\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"header without colon", "fragments 3\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"fragment in yellow", head + "[yellow]\n#PEY\n[grey]\n####\n", ErrFragmentOutsideGrey},
		{"spikes in both layers", head + "[yellow]\n#PE^\n[grey]\n#^^#\n", nil},
		{"watcher in grey", head + "[yellow]\n#PE#\n[grey]\n#W.#\n", nil},
		{"watcher in yellow", head + "[yellow]\n#PEW\n[grey]\n####\n", ErrWatcherOutsideGrey},
		{"missing grey", head + "[yellow]\n#PE#\n", ErrMissingSection},
		{"empty section", head + "[yellow]\n#PE#\n[grey]\n", ErrMissingSection},
		{"unknown section", head + "[yellow]\n#PE#\n[blue]\n####\n", ErrUnknownSection},
		{"duplicate section", head + "[yellow]\n#PE#\n[yellow]\n####\n", ErrDuplicateSection},
		{"row before section reads as header", head + "####\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"shift free", head + "shift: free\n[yellow]\n#PE#\n[grey]\n####\n", nil},
		{"shift quantized", head + "shift: quantized\n[yellow]\n#PE#\n[grey]\n####\n", nil},
		{"shift auto", head + "shift: auto\nauto_every: 4\n[yellow]\n#PE#\n[grey]\n####\n", nil},
		{"shift before fragments", "shift: quantized\n" + head + "[yellow]\n#PE#\n[grey]\n####\n", nil},
		{"shift unknown mode", head + "shift: sometimes\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"shift without fragments", "shift: free\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"duplicate key", head + head + "[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"auto without auto_every", head + "shift: auto\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"auto_every zero", head + "shift: auto\nauto_every: 0\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"auto_every not a number", head + "shift: auto\nauto_every: soon\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"auto_every outside auto", head + "shift: quantized\nauto_every: 4\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"auto_every without shift", head + "auto_every: 4\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"fragment in auto", head + "shift: auto\nauto_every: 4\n[yellow]\n#PE#\n[grey]\n#Y.#\n", ErrFragmentInAuto},
		{"start fragments in auto", "fragments: 1\nshift: auto\nauto_every: 4\n[yellow]\n#PE#\n[grey]\n####\n", nil},
		{"transmission in yellow", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\nhello there\n", nil},
		{"transmissions in both layers", head + "[yellow]\n#PTE#\n[grey]\n#T..#\n[transmissions]\nfirst\nsecond\n", nil},
		{"transmission without text", head + "[yellow]\n#PTE#\n[grey]\n#####\n", ErrTransmissionCount},
		{"text without transmission", head + "[yellow]\n#PE#\n[grey]\n####\n[transmissions]\nlost words\n", ErrTransmissionCount},
		{"more texts than markers", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\none\ntwo\n", ErrTransmissionCount},
		{"duplicate transmissions section", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\none\n[transmissions]\n", ErrDuplicateSection},
		{"portuguese transmissions", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\none\n[transmissions.pt]\nação à vista\n", nil},
		{"portuguese before english", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions.pt]\num\n[transmissions]\none\n", nil},
		{"portuguese count mismatch", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\none\n[transmissions.pt]\num\ndois\n", ErrTransmissionCount},
		{"portuguese without english", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions.pt]\num\n", ErrTransmissionCount},
		{"unknown text language", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\none\n[transmissions.fr]\nun\n", ErrUnknownSection},
		{"unsupported accent", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\none\n[transmissions.pt]\nümlaut\n", ErrTransmissionText},
		{"transmission with unsupported character", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\nsmile ;)\n", ErrTransmissionText},
		{"transmission too long", head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\n" + strings.Repeat("a", MaxTransmissionLength+1) + "\n", ErrTransmissionText},
		{"secret in yellow", head + "[yellow]\n#PSE#\n[grey]\n#####\n", nil},
		{"secret in grey", head + "[yellow]\n#P.E#\n[grey]\n#.S.#\n", nil},
		{"two secrets", head + "[yellow]\n#PSE#\n[grey]\n#S..#\n", ErrSecretCount},
		{"secret in auto", head + "shift: auto\nauto_every: 4\n[yellow]\n#PSE#\n[grey]\n#####\n", nil},
		{"act two", head + "act: 2\n[yellow]\n#PE#\n[grey]\n####\n", nil},
		{"act zero", head + "act: 0\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"act four", head + "act: 4\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"act not a number", head + "act: two\n[yellow]\n#PE#\n[grey]\n####\n", ErrInvalidHeader},
		{"ragged rows", head + "[yellow]\n#PE#\n###\n[grey]\n####\n####\n", ErrRaggedRows},
		{"different widths", head + "[yellow]\n#PE#\n[grey]\n#####\n", ErrRaggedRows},
		{"different heights", head + "[yellow]\n#PE#\n[grey]\n####\n####\n", ErrSizeMismatch},
		{"unknown tile", head + "[yellow]\n#PE#\n[grey]\n#x.#\n", ErrUnknownTile},
		{"missing spawn", head + "[yellow]\n#.E#\n[grey]\n####\n", ErrMarkerCount},
		{"missing exit", head + "[yellow]\n#P.#\n[grey]\n####\n", ErrMarkerCount},
		{"two spawns", head + "[yellow]\n#PE#\n[grey]\n#P.#\n", ErrMarkerCount},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.text)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Parse() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseShiftHeader(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		mode      ShiftMode
		autoEvery int
	}{
		{"defaults to free", head, ShiftFree, 0},
		{"free", head + "shift: free\n", ShiftFree, 0},
		{"quantized", head + "shift: quantized\n", ShiftQuantized, 0},
		{"auto", head + "shift:auto\nauto_every: 3\n", ShiftAuto, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl, err := Parse(tt.header + "[yellow]\n#PE#\n[grey]\n####\n")
			if err != nil {
				t.Fatal(err)
			}
			if lvl.Shift != tt.mode || lvl.AutoEvery != tt.autoEvery {
				t.Fatalf("shift = %v every %d, want %v every %d", lvl.Shift, lvl.AutoEvery, tt.mode, tt.autoEvery)
			}
		})
	}
}

func TestParseTransmissionsAndSecret(t *testing.T) {
	lvl, err := Parse(head + "[yellow]\n#P.T.E#\n#######\n[grey]\n#T..S.#\n#.....#\n[transmissions]\nWarm side first\nthen the cold one, it's late?\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []Transmission{
		{Marker{Point{X: 3, Y: 0}, world.Yellow}, [i18n.LangCount]string{i18n.EN: "WARM SIDE FIRST"}},
		{Marker{Point{X: 1, Y: 0}, world.Grey}, [i18n.LangCount]string{i18n.EN: "THEN THE COLD ONE, IT'S LATE?"}},
	}
	if len(lvl.Transmissions) != len(want) {
		t.Fatalf("got %d transmissions, want %d", len(lvl.Transmissions), len(want))
	}
	for i, w := range want {
		if lvl.Transmissions[i] != w {
			t.Errorf("transmission %d = %+v, want %+v", i, lvl.Transmissions[i], w)
		}
	}
	if lvl.Secret == nil || *lvl.Secret != (Marker{Point{X: 4, Y: 0}, world.Grey}) {
		t.Fatalf("secret = %+v", lvl.Secret)
	}
	if lvl.Layer(world.Grey).Solid(4, 0) || lvl.Layer(world.Yellow).Solid(3, 0) {
		t.Fatal("marker tiles must be empty")
	}
}

func TestTransmissionLanguages(t *testing.T) {
	lvl, err := Parse(head + "[yellow]\n#PTTE#\n[grey]\n######\n[transmissions]\nfirst\nsecond\n[transmissions.pt]\nprimeira\nsegunda ação\n")
	if err != nil {
		t.Fatal(err)
	}
	tr := lvl.Transmissions[1]
	if tr.TextIn(i18n.EN) != "SECOND" || tr.TextIn(i18n.PT) != "SEGUNDA AÇÃO" {
		t.Fatalf("texts %q", tr.Text)
	}
	only, err := Parse(head + "[yellow]\n#PTE#\n[grey]\n#####\n[transmissions]\nlonely\n")
	if err != nil {
		t.Fatal(err)
	}
	if got := only.Transmissions[0].TextIn(i18n.PT); got != "LONELY" {
		t.Fatalf("missing portuguese should fall back to english, got %q", got)
	}
}

func TestParseLayout(t *testing.T) {
	lvl, err := Parse("fragments: 3\n[yellow]\n#P.#\n####\n[grey]\n#.YE\n#Y.#\n")
	if err != nil {
		t.Fatal(err)
	}
	checks := []struct {
		name  string
		layer world.World
		x, y  int
		solid bool
	}{
		{"yellow wall", world.Yellow, 0, 0, true},
		{"spawn tile is empty", world.Yellow, 1, 0, false},
		{"yellow floor", world.Yellow, 2, 1, true},
		{"grey gap", world.Grey, 2, 1, false},
		{"exit tile is empty", world.Grey, 3, 0, false},
		{"fragment tile is empty", world.Grey, 2, 0, false},
		{"out of bounds", world.Grey, 9, 9, false},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if got := lvl.Layer(c.layer).Solid(c.x, c.y); got != c.solid {
				t.Fatalf("Solid(%d,%d) = %v, want %v", c.x, c.y, got, c.solid)
			}
		})
	}
	if lvl.Spawn != (Point{1, 0}) || lvl.Exit != (Point{3, 0}) {
		t.Fatalf("markers = %v %v", lvl.Spawn, lvl.Exit)
	}
	if lvl.StartFragments != 3 {
		t.Fatalf("StartFragments = %d, want 3", lvl.StartFragments)
	}
	wantFragments := []Point{{2, 0}, {1, 1}}
	if len(lvl.Fragments) != len(wantFragments) {
		t.Fatalf("Fragments = %v, want %v", lvl.Fragments, wantFragments)
	}
	for i, p := range wantFragments {
		if lvl.Fragments[i] != p {
			t.Fatalf("Fragments = %v, want %v", lvl.Fragments, wantFragments)
		}
	}
}

func TestParseHazards(t *testing.T) {
	lvl, err := Parse(head + "[yellow]\n#P^#\n####\n[grey]\n#W.E\n#.^^\n")
	if err != nil {
		t.Fatal(err)
	}
	checks := []struct {
		name string
		got  []Point
		want []Point
	}{
		{"yellow spikes", lvl.Spikes[world.Yellow], []Point{{2, 0}}},
		{"grey spikes", lvl.Spikes[world.Grey], []Point{{2, 1}, {3, 1}}},
		{"watchers", lvl.Watchers, []Point{{1, 0}}},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if len(c.got) != len(c.want) {
				t.Fatalf("got %v, want %v", c.got, c.want)
			}
			for i := range c.want {
				if c.got[i] != c.want[i] {
					t.Fatalf("got %v, want %v", c.got, c.want)
				}
			}
		})
	}
	if lvl.Layer(world.Yellow).Solid(2, 0) || lvl.Layer(world.Grey).Solid(1, 0) {
		t.Fatal("spike and watcher tiles must not be solid")
	}
}

func TestEmbeddedLevelsLoad(t *testing.T) {
	levels, err := LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	names, _ := Names()
	if len(levels) != len(names) || len(levels) == 0 {
		t.Fatalf("loaded %d levels from %v", len(levels), names)
	}
}
