package level

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/world"
)

func standsOnYellow(l Level, p Point) bool {
	y := l.Layer(world.Yellow)
	return !y.Solid(p.X, p.Y) && y.Solid(p.X, p.Y+1)
}

func TestEmbeddedLevelInvariants(t *testing.T) {
	names, err := Names()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			lvl, err := Load(name)
			if err != nil {
				t.Fatal(err)
			}
			yellow, grey := lvl.Layer(world.Yellow), lvl.Layer(world.Grey)
			checks := []struct {
				what string
				ok   bool
			}{
				{"layers have the same size", yellow.Width() == grey.Width() && yellow.Height() == grey.Height()},
				{"spawn stands on yellow ground", standsOnYellow(lvl, lvl.Spawn)},
				{"exit is inside the level", yellow.InBounds(lvl.Exit.X, lvl.Exit.Y)},
				{"exit is reachable in some world", !yellow.Solid(lvl.Exit.X, lvl.Exit.Y) || !grey.Solid(lvl.Exit.X, lvl.Exit.Y)},
			}
			for _, p := range lvl.Checkpoints {
				checks = append(checks, struct {
					what string
					ok   bool
				}{"checkpoint stands on yellow ground", standsOnYellow(lvl, p)})
			}
			for _, p := range append(append([]Point{}, lvl.Fragments...), lvl.Watchers...) {
				checks = append(checks, struct {
					what string
					ok   bool
				}{"grey-only marker sits in open grey space", !grey.Solid(p.X, p.Y)})
			}
			for _, c := range checks {
				if !c.ok {
					t.Errorf("%s", c.what)
				}
			}
		})
	}
}
