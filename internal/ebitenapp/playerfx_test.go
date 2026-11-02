package ebitenapp

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

const dropLevel = "fragments: 0\n[yellow]\n#######\n#P...E#\n#.....#\n#.....#\n#.....#\n#######\n[grey]\n#######\n#.....#\n#.....#\n#.....#\n#.....#\n#######\n"

func TestLandingSquashesAndRaisesDust(t *testing.T) {
	lvl, err := level.Parse(dropLevel)
	if err != nil {
		t.Fatal(err)
	}
	clock := &rhythm.Manual{}
	s := game.NewSession(clock, lvl)
	var fx playerFX
	landed := false
	for range 120 {
		clock.Advance(time.Second / 60)
		s.Update(input.Intent{})
		fx.update(s)
		if fx.squash == squashTicks {
			landed = true
			break
		}
	}
	if !landed || len(fx.dust) != dustPerLanding {
		t.Fatalf("landing: squash=%d dust=%d", fx.squash, len(fx.dust))
	}
	body := fx.body(s.PlayerBox())
	if body.W <= s.PlayerBox().W || body.H >= s.PlayerBox().H || body.Bottom() != s.PlayerBox().Bottom() {
		t.Fatalf("squashed body %+v from %+v", body, s.PlayerBox())
	}
	for range dustTicks {
		s.Update(input.Intent{Left: true})
		fx.update(s)
	}
	if fx.squash != 0 || len(fx.dust) != 0 || fx.facing != -1 {
		t.Fatalf("after settling: squash=%d dust=%d facing=%d", fx.squash, len(fx.dust), fx.facing)
	}
}
