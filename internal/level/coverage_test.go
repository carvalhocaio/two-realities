package level

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/carvalhocaio/two-realities/internal/world"
)

const tinyLevel = "fragments: 0\n[yellow]\n#....#\n#P..E#\n######\n[grey]\n#....#\n#....#\n######\n"

func TestLoadFromAFileSystem(t *testing.T) {
	tests := []struct {
		name  string
		fsys  fstest.MapFS
		want  []string
		check func(error) bool
	}{
		{"sorted levels", fstest.MapFS{"levels/b.txt": {Data: []byte(tinyLevel)}, "levels/a.txt": {Data: []byte(tinyLevel)}}, []string{"a", "b"}, nil},
		{"no levels directory", fstest.MapFS{}, nil, func(err error) bool { return errors.Is(err, fs.ErrNotExist) }},
		{"empty directory", fstest.MapFS{"levels": {Mode: fs.ModeDir}}, nil, func(err error) bool { return errors.Is(err, ErrNoLevels) }},
		{"broken level", fstest.MapFS{"levels/a.txt": {Data: []byte("fragments: 0\n")}}, nil, func(err error) bool { return err != nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			levels, err := loadAll(tt.fsys)
			if tt.check != nil {
				if !tt.check(err) {
					t.Fatalf("error %v", err)
				}
				return
			}
			if err != nil || len(levels) != len(tt.want) {
				t.Fatalf("levels %d, err %v", len(levels), err)
			}
			for i, name := range tt.want {
				if levels[i].Name != name {
					t.Fatalf("level %d named %q, want %q", i, levels[i].Name, name)
				}
			}
		})
	}
}

func TestLoadUnknownLevel(t *testing.T) {
	if _, err := Load("99_nope"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("error %v", err)
	}
}

func TestGridAndLevelGeometry(t *testing.T) {
	lvl, err := Parse(tinyLevel)
	if err != nil {
		t.Fatal(err)
	}
	g := lvl.Layer(world.Grey)
	if g.PixelWidth() != 6*TileSize || g.PixelHeight() != 3*TileSize {
		t.Fatalf("pixel size %vx%v", g.PixelWidth(), g.PixelHeight())
	}
	if b := lvl.Bounds(); b.W != 60 || b.H != 30 || b.X != 0 || b.Y != 0 {
		t.Fatalf("bounds %+v", b)
	}
	closed := g.WithSolid([]Point{{X: 2, Y: 1}, {X: 99, Y: 99}})
	if !closed.Solid(2, 1) || g.Solid(2, 1) {
		t.Fatal("WithSolid should fill a copy and ignore points out of bounds")
	}
}

func TestShiftModeNames(t *testing.T) {
	for mode, want := range map[ShiftMode]string{ShiftFree: "free", ShiftQuantized: "quantized", ShiftAuto: "auto", ShiftMode(42): "unknown"} {
		if mode.String() != want {
			t.Errorf("%d is %q, want %q", mode, mode.String(), want)
		}
	}
}

func TestParseSectionErrors(t *testing.T) {
	tests := []struct {
		name string
		text string
		want error
	}{
		{"row before any section", "fragments: 0\n#P.E#\n[yellow]\n#....#\n[grey]\n#....#\n", ErrInvalidHeader},
		{"duplicate layer", "fragments: 0\n[yellow]\n#P.E#\n[yellow]\n#....#\n[grey]\n#....#\n", ErrDuplicateSection},
		{"unknown section", "fragments: 0\n[purple]\n#P.E#\n", ErrUnknownSection},
		{"missing grey", "fragments: 0\n[yellow]\n#P.E#\n", ErrMissingSection},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Parse(tt.text); !errors.Is(err, tt.want) {
				t.Fatalf("error %v, want %v", err, tt.want)
			}
		})
	}
}

func TestTileRect(t *testing.T) {
	r := TileRect(Point{X: 3, Y: 2})
	if r.X != 3*TileSize || r.Y != 2*TileSize || r.W != TileSize || r.H != TileSize {
		t.Fatalf("tile rect %+v", r)
	}
}
