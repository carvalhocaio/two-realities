package ebitenapp

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/geom"
)

var deadBox = geom.Rect{X: 50, Y: 40, W: 6, H: 9}

func TestShardsBurstOutAndFall(t *testing.T) {
	n := shardCount(deadBox)
	if n < 6 {
		t.Fatalf("only %d shards", n)
	}
	for i := range n {
		if start := shard(deadBox, i, 0); !start.Overlaps(deadBox) {
			t.Fatalf("shard %d starts outside the body: %+v", i, start)
		}
	}
	left, right := shard(deadBox, 0, 0.5), shard(deadBox, 2, 0.5)
	if left.X >= deadBox.X || right.Right() <= deadBox.Right() {
		t.Fatalf("outer shards did not spread: %+v %+v", left, right)
	}
	if mid := shard(deadBox, 1, 0.3); mid.Y >= deadBox.Y {
		t.Fatalf("top shard did not burst up: %+v", mid)
	}
	for i := range n {
		end := shard(deadBox, i, 1)
		if end.Y <= deadBox.Y || end.Bottom() > deadBox.Bottom() {
			t.Fatalf("shard %d should end on the floor it died on: %+v", i, end)
		}
	}
}

func TestBeamDeathDissolvesUpward(t *testing.T) {
	if b := dissolvedBody(deadBox, 0); b != deadBox {
		t.Fatalf("body %+v at the start, want %+v", b, deadBox)
	}
	if b := dissolvedBody(deadBox, 0.5); b.H >= deadBox.H || b.Y != deadBox.Y {
		t.Fatalf("half dissolved body %+v", b)
	}
	if b := dissolvedBody(deadBox, 1); b.H != 0 {
		t.Fatalf("body left at the end: %+v", b)
	}
	for i := range ashMotes {
		start, later := ashMote(deadBox, i, 0), ashMote(deadBox, i, 0.5)
		if !start.Overlaps(deadBox) || later.Y >= start.Y {
			t.Fatalf("mote %d: %+v then %+v", i, start, later)
		}
	}
}

func TestCaughtSlicesSplitBothWays(t *testing.T) {
	if sliceOffset(0, 0) != 0 || sliceOffset(1, 0) != 0 {
		t.Fatal("slices offset before the death started")
	}
	a, b := sliceOffset(0, 0.4), sliceOffset(1, 0.4)
	if a <= 0 || b >= 0 {
		t.Fatalf("slices should split both ways: %v %v", a, b)
	}
}

func TestFallSparksRiseFromTheBottomEdge(t *testing.T) {
	bounds := geom.Rect{W: 200, H: 100}
	fallen := geom.Rect{X: 80, Y: 130, W: 6, H: 9}
	for i := range fallSparks {
		start, later := fallSpark(fallen, bounds, i, 0), fallSpark(fallen, bounds, i, 0.5)
		if start.Y != bounds.Bottom()-1 || later.Y >= start.Y {
			t.Fatalf("spark %d: %+v then %+v", i, start, later)
		}
		if end := fallSpark(fallen, bounds, i, 1); end.Y < bounds.Bottom()-1-fallRiseMax {
			t.Fatalf("spark %d rose past %dpx: %+v", i, fallRiseMax, end)
		}
		if start.X < fallen.X-fallSpread*2 || start.X > fallen.Right()+fallSpread*2 {
			t.Fatalf("spark %d far from where the player fell: %+v", i, start)
		}
	}
}
