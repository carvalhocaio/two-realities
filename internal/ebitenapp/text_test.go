package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/font"
)

func TestRunePrefix(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"HELLO", 0, ""},
		{"HELLO", -3, ""},
		{"HELLO", 2, "HE"},
		{"HELLO", 5, "HELLO"},
		{"HELLO", 9, "HELLO"},
		{"AÇÃO", 2, "AÇ"},
		{"AÇÃO", 3, "AÇÃ"},
		{"", 4, ""},
	}
	for _, tt := range tests {
		if got := runePrefix(tt.s, tt.n); got != tt.want {
			t.Fatalf("runePrefix(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
		}
	}
	if n := testing.AllocsPerRun(100, func() { runePrefix("TRANSMISSÃO", 7) }); n != 0 {
		t.Fatalf("runePrefix allocated %v times", n)
	}
}

func TestTextWidthCountsRunes(t *testing.T) {
	if got, want := scaledTextWidth("AÇÃO", 2), float64(4*(glyphWidth+glyphSpacing)-glyphSpacing)*2; got != want {
		t.Fatalf("width %v, want %v", got, want)
	}
	if scaledTextWidth("", 3) != 0 {
		t.Fatal("empty text has width")
	}
}

func TestGlyphPixelsKeepAccentAndCedillaRows(t *testing.T) {
	lit := func(pix []byte, x, y int) bool { return pix[4*(y*glyphWidth+x)+3] == 0xff }
	for _, ch := range []rune{'A', 'Á', 'Ç'} {
		g, ok := font.Lookup(ch)
		if !ok {
			t.Fatalf("%q missing from the font", ch)
		}
		pix := glyphPixels(g)
		rows := append(append([]string{g.Above}, g.Rows[:]...), g.Below)
		for y := range glyphCellHeight {
			for x := range glyphWidth {
				want := y < len(rows) && x < len(rows[y]) && rows[y][x] == '#'
				if lit(pix, x, y) != want {
					t.Fatalf("%q pixel %d,%d lit %v, want %v", ch, x, y, lit(pix, x, y), want)
				}
			}
		}
	}
}

func TestTextboxWrapsOncePerText(t *testing.T) {
	var tb textbox
	first := tb.wrapped("THE WARM FLOOR GIVES UP HALFWAY", 60)
	if again := tb.wrapped("THE WARM FLOOR GIVES UP HALFWAY", 60); &again[0] != &first[0] {
		t.Fatal("same text wrapped again")
	}
	if n := testing.AllocsPerRun(100, func() { tb.wrapped("THE WARM FLOOR GIVES UP HALFWAY", 60) }); n != 0 {
		t.Fatalf("cached wrap allocated %v times", n)
	}
	if other := tb.wrapped("O CHÃO QUENTE CEDE NO MEIO", 60); other[0] == first[0] {
		t.Fatal("new text kept the old lines")
	}
	narrow := tb.wrapped("O CHÃO QUENTE CEDE NO MEIO", 30)
	if len(narrow) <= 1 {
		t.Fatalf("width change did not rewrap: %v", narrow)
	}
}
