package level

import (
	"errors"
	"testing"
)

const bossHead = "fragments: 0\nboss: warden\n"

func bossText(head, yellow, grey string) string {
	return head + "[yellow]\n" + yellow + "[grey]\n" + grey
}

func TestParseBoss(t *testing.T) {
	const (
		yellow = "#........#\n#L.1..2.L#\n#P...L..E#\n"
		grey   = "#.OOOOOO.#\n#........#\n#........#\n"
	)
	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{"valid", bossText(bossHead, yellow, grey), nil},
		{"lamps in grey", bossText(bossHead, "#........#\n#...1..2.#\n#P.......#\n", "#.OOOOOO.#\n#L......L#\n#....L..E#\n"), nil},
		{"unknown boss", bossText(head+"boss: dragon\n", yellow, grey), ErrInvalidHeader},
		{"boss with shift", bossText(bossHead+"shift: free\n", yellow, grey), ErrBossShiftOverride},
		{"boss with auto", bossText(bossHead+"shift: auto\nauto_every: 2\n", yellow, grey), ErrBossShiftOverride},
		{"lamp without boss", bossText(head, "#P.L.E#\n", "#######\n"), ErrBossMarkers},
		{"rail without boss", bossText(head, "#P...E#\n", "#.OO..#\n"), ErrBossMarkers},
		{"gate without boss", bossText(head, "#P.2.E#\n", "#######\n"), ErrBossMarkers},
		{"two lamps", bossText(bossHead, "#........#\n#L.1..2..#\n#P...L..E#\n", grey), ErrLampCount},
		{"four lamps", bossText(bossHead, "#........#\n#L.1..2.L#\n#P..LL..E#\n", grey), ErrLampCount},
		{"no rail", bossText(bossHead, yellow, "#........#\n#........#\n#........#\n"), ErrRailShape},
		{"broken rail", bossText(bossHead, yellow, "#.OO.OOO.#\n#........#\n#........#\n"), ErrRailShape},
		{"rail on two rows", bossText(bossHead, yellow, "#.OOO....#\n#.....OO.#\n#........#\n"), ErrRailShape},
		{"rail in yellow", bossText(bossHead, "#.O......#\n#L.1..2.L#\n#P...L..E#\n", grey), ErrRailOutsideGrey},
		{"missing second gate", bossText(bossHead, "#........#\n#L.1....L#\n#P...L..E#\n", grey), ErrGateMissing},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.text)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestBossLayout(t *testing.T) {
	lvl, err := Parse(bossText(bossHead,
		"#........#\n#L.1..2.L#\n#P...L..E#\n",
		"#.OOOOOO.#\n#..1.....#\n#........#\n"))
	if err != nil {
		t.Fatal(err)
	}
	b := lvl.Boss
	if b == nil {
		t.Fatal("boss missing")
	}
	wantLamps := []Point{{5, 2}, {1, 1}, {8, 1}}
	for i, p := range wantLamps {
		if b.Lamps[i] != p {
			t.Errorf("lamp %d = %v, want %v (bottom to top, then left to right)", i, b.Lamps[i], p)
		}
	}
	if b.Rail != (Rail{Row: 0, From: 2, To: 7}) {
		t.Errorf("rail = %+v", b.Rail)
	}
	if len(b.Gates[0]) != 2 || len(b.Gates[1]) != 1 {
		t.Errorf("gates = %v", b.Gates)
	}
	for _, p := range append(b.Gates[0], b.Lamps...) {
		if lvl.Layer(0).Solid(p.X, p.Y) || lvl.Layer(1).Solid(p.X, p.Y) {
			t.Errorf("%v is parsed as solid; gates close at runtime", p)
		}
	}
}

func TestPlainLevelHasNoBoss(t *testing.T) {
	lvl, err := Parse(bossText(head, "#PE#\n", "####\n"))
	if err != nil || lvl.Boss != nil {
		t.Fatalf("boss = %v, err = %v", lvl.Boss, err)
	}
}
