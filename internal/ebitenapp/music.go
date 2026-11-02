package ebitenapp

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"

	tunes "github.com/carvalhocaio/two-realities/internal/music"
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
	ctx := audio.NewContext(tunes.SampleRate)
	tracks := make([][]float32, tunes.TrackCount)
	for i := range tracks {
		tracks[i] = tunes.Render(tunes.Track(i), rhythm.BPM)
	}
	loop := newLoopReader(tracks, tunes.BarFrames(rhythm.BPM))
	player, err := ctx.NewPlayerF32(loop)
	if err != nil {
		return nil, err
	}
	player.SetVolume(musicVolume)
	return &music{player: player, loop: loop, latency: defaultAudioLatency}, nil
}

func (m *music) setTrack(t tunes.Track) {
	m.loop.next.Store(int32(t))
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
