package ebitenapp

import (
	"math"
	"testing"
)

func TestShakeOffset(t *testing.T) {
	if x, y := shakeOffset(0, 7); x != 0 || y != 0 {
		t.Fatalf("no strength shook by %v,%v", x, y)
	}
	seen := map[[2]float64]bool{}
	for tick := range 30 {
		x, y := shakeOffset(1, tick)
		if math.Abs(x) > shakePixels || math.Abs(y) > shakePixels {
			t.Fatalf("tick %d offset %v,%v beyond %v px", tick, x, y, float64(shakePixels))
		}
		if x2, y2 := shakeOffset(1, tick); x2 != x || y2 != y {
			t.Fatalf("tick %d gave %v,%v then %v,%v", tick, x, y, x2, y2)
		}
		seen[[2]float64{x, y}] = true
	}
	if len(seen) < 25 {
		t.Fatalf("only %d distinct offsets in 30 ticks", len(seen))
	}
	if x, _ := shakeOffset(0.5, 1); math.Abs(x) > shakePixels/2.0 {
		t.Fatalf("half strength offset %v beyond %v px", x, shakePixels/2.0)
	}
}
