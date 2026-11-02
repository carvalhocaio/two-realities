package ebitenapp

import (
	"encoding/binary"
	"math"
	"testing"
)

func TestLoopSwitchesTracksOnlyOnBarLines(t *testing.T) {
	const bar = 4
	a := []float32{1, 1, 1, 1, 2, 2, 2, 2}
	b := []float32{5, 5, 5, 5, 6, 6, 6, 6}
	l := newLoopReader([][]float32{a, b}, bar)
	read := func(frames int) []float32 {
		p := make([]byte, frames*bytesPerFrame)
		n, _ := l.Read(p)
		out := make([]float32, n/bytesPerFrame)
		for i := range out {
			out[i] = math.Float32frombits(binary.LittleEndian.Uint32(p[i*bytesPerFrame:]))
		}
		return out
	}
	got := read(2)
	l.next.Store(1)
	got = append(got, read(8)...)
	want := []float32{1, 1, 1, 1, 6, 6, 6, 6, 5, 5}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("frames %v, want %v", got, want)
		}
	}
}
