package ebitenapp

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

const (
	musicVolume  = 0.6
	latencyStep  = 10 * time.Millisecond
	latencyLimit = 500 * time.Millisecond
)

type music struct {
	player  *audio.Player
	muted   bool
	latency time.Duration
}

func newMusic() (*music, error) {
	ctx := audio.NewContext(sampleRate)
	player, err := ctx.NewPlayerF32(&loopReader{data: synthBar()})
	if err != nil {
		return nil, err
	}
	player.SetVolume(musicVolume)
	return &music{player: player, latency: defaultAudioLatency}, nil
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

func (m *music) nudgeLatency(steps int) {
	m.latency = max(-latencyLimit, min(latencyLimit, m.latency+time.Duration(steps)*latencyStep))
}

func (m *music) position() rhythm.Position {
	return rhythm.At(m.player.Position()-m.latency, rhythm.BPM)
}

func (m *music) Beat() int      { return m.position().Beat() }
func (m *music) Phase() float64 { return m.position().Phase() }
