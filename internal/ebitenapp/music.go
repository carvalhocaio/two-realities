package ebitenapp

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

const musicVolume = 0.6

type music struct {
	player  *audio.Player
	loop    *loopReader
	muted   bool
	latency time.Duration
}

func newMusic() (*music, error) {
	ctx := audio.NewContext(sampleRate)
	loop := &loopReader{bars: synthBars()}
	player, err := ctx.NewPlayerF32(loop)
	if err != nil {
		return nil, err
	}
	player.SetVolume(musicVolume)
	return &music{player: player, loop: loop, latency: defaultAudioLatency}, nil
}

func (m *music) setIntensity(level int) {
	m.loop.next.Store(int32(min(max(level, 0), intensityLevels-1)))
}

func (m *music) start() {
	if !m.player.IsPlaying() {
		m.player.Play()
	}
}

func (m *music) setMuted(muted bool) {
	if muted == m.muted {
		return
	}
	m.muted = muted
	if muted {
		m.player.SetVolume(0)
		return
	}
	m.player.SetVolume(musicVolume)
}

func (m *music) setOffset(ms int) {
	m.latency = defaultAudioLatency + time.Duration(ms)*time.Millisecond
}

func (m *music) position() rhythm.Position {
	return rhythm.At(m.player.Position()-m.latency, rhythm.BPM)
}

func (m *music) Beat() int      { return m.position().Beat() }
func (m *music) Phase() float64 { return m.position().Phase() }
