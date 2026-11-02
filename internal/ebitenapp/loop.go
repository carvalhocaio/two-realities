package ebitenapp

import (
	"encoding/binary"
	"math"
	"sync/atomic"
)

const bytesPerFrame = 8

type loopReader struct {
	tracks [][]float32
	bar    int
	frames int
	cur    int
	next   atomic.Int32
	pos    int
}

func newLoopReader(tracks [][]float32, barFrames int) *loopReader {
	return &loopReader{tracks: tracks, frames: barFrames}
}

func (l *loopReader) Read(p []byte) (int, error) {
	n := 0
	for n+bytesPerFrame <= len(p) {
		track := l.tracks[l.cur]
		v := track[(l.bar*l.frames+l.pos)%len(track)]
		bits := math.Float32bits(v)
		binary.LittleEndian.PutUint32(p[n:], bits)
		binary.LittleEndian.PutUint32(p[n+4:], bits)
		n += bytesPerFrame
		l.pos++
		if l.pos == l.frames {
			l.pos = 0
			l.bar++
			l.cur = int(l.next.Load())
		}
	}
	return n, nil
}
