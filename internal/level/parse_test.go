package level

import (
	"errors"
	"testing"

	"github.com/carvalhocaio/two-realities/internal/world"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{"valid", "[yellow]\n#P.#\n####\n[grey]\n#..E\n#..#\n", nil},
		{"valid with blank lines and CRLF", "\r\n[yellow]\r\n#P.#\r\n\r\n[grey]\r\n#.E#\r\n", nil},
		{"markers split across layers", "[yellow]\n#E.#\n[grey]\n#.P#\n", nil},
		{"missing grey", "[yellow]\n#PE#\n", ErrMissingSection},
		{"empty section", "[yellow]\n#PE#\n[grey]\n", ErrMissingSection},
		{"unknown section", "[yellow]\n#PE#\n[blue]\n####\n", ErrUnknownSection},
		{"duplicate section", "[yellow]\n#PE#\n[yellow]\n####\n", ErrDuplicateSection},
		{"row before section", "####\n[yellow]\n#PE#\n[grey]\n####\n", ErrRowOutsideLayer},
		{"ragged rows", "[yellow]\n#PE#\n###\n[grey]\n####\n####\n", ErrRaggedRows},
		{"different widths", "[yellow]\n#PE#\n[grey]\n#####\n", ErrRaggedRows},
		{"different heights", "[yellow]\n#PE#\n[grey]\n####\n####\n", ErrSizeMismatch},
		{"unknown tile", "[yellow]\n#PE#\n[grey]\n#x.#\n", ErrUnknownTile},
		{"missing spawn", "[yellow]\n#.E#\n[grey]\n####\n", ErrMarkerCount},
		{"missing exit", "[yellow]\n#P.#\n[grey]\n####\n", ErrMarkerCount},
		{"two spawns", "[yellow]\n#PE#\n[grey]\n#P.#\n", ErrMarkerCount},
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

func TestParseLayout(t *testing.T) {
	lvl, err := Parse("[yellow]\n#P.#\n####\n[grey]\n#..E\n#..#\n")
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
}

func TestEmbeddedLevelsLoad(t *testing.T) {
	entries, err := files.ReadDir("levels")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()[:len(e.Name())-len(".txt")]
		t.Run(name, func(t *testing.T) {
			if _, err := Load(name); err != nil {
				t.Fatal(err)
			}
		})
	}
}
