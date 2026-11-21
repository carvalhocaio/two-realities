package ebitenapp

import (
	"testing"
	"time"

	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

func TestCounterPopsOnGainAndDimsOnSpend(t *testing.T) {
	lvl, err := level.Parse(pickupLevel)
	if err != nil {
		t.Fatal(err)
	}
	clock := &rhythm.Manual{}
	s := game.NewSession(clock, lvl, lvl)
	var fx counterFX
	step := func(in input.Intent, n int) {
		for range n {
			clock.Advance(time.Second / 60)
			s.Update(in)
			fx.update(s)
		}
	}
	fx.update(s)
	if lift, c, a := fx.icon(); lift != 0 || c != fragmentColor || a != 1 {
		t.Fatalf("idle icon lift %v color %v alpha %v", lift, c, a)
	}
	step(input.Intent{ShiftPressed: true}, 1)
	for range 30 {
		step(input.Intent{Right: true}, 1)
		if fx.gain > 0 {
			break
		}
	}
	if lift, c, _ := fx.icon(); lift != 1 || c != secretColor || fx.gain != counterGainTicks {
		t.Fatalf("pickup: lift %v color %v gain %d", lift, c, fx.gain)
	}
	step(input.Intent{ShiftPressed: true}, 1)
	if _, _, a := fx.icon(); a != counterSpendAlpha || fx.spend != counterSpendTicks || fx.gain != 0 {
		t.Fatalf("spend: alpha %v spend %d gain %d", a, fx.spend, fx.gain)
	}
	step(input.Intent{}, counterSpendTicks)
	if _, _, a := fx.icon(); a != 1 {
		t.Fatal("spend dim outlived its ticks")
	}
	step(input.Intent{RestartPressed: true}, 1)
	if fx.gain != 0 || fx.spend != 0 || fx.count != s.Fragments() {
		t.Fatalf("restart popped: gain %d spend %d count %d", fx.gain, fx.spend, fx.count)
	}
}

func TestSignalWave(t *testing.T) {
	head := geom.Rect{X: 2, Y: 40, W: signalHeadWidth, H: 2}
	tests := []struct {
		name  string
		tick  int
		ok    bool
		w, y  float64
		alpha float32
	}{
		{"wave starts at the head", 0, true, 2, 38, 1},
		{"halfway it is wider and higher", signalWaveTicks / 2, true, 4, 36, 0.5},
		{"quiet after the wave", signalWaveTicks, false, 0, 0, 0},
		{"repeats every period", signalWavePeriod, true, 2, 38, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, a, ok := signalWave(head, tt.tick)
			if ok != tt.ok || !ok {
				if ok != tt.ok {
					t.Fatalf("ok %v, want %v", ok, tt.ok)
				}
				return
			}
			if r.W != tt.w || r.Y != tt.y || a != tt.alpha || r.X+r.W/2 != head.X+head.W/2 {
				t.Fatalf("wave %+v alpha %v, want w %v y %v alpha %v centered", r, a, tt.w, tt.y, tt.alpha)
			}
		})
	}
}
