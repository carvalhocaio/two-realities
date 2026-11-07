package ebitenapp

import (
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/carvalhocaio/two-realities/internal/game"
	tunes "github.com/carvalhocaio/two-realities/internal/music"
)

const sfxVolume = 0.7

var eventSounds = map[game.Event]tunes.Sound{
	game.EventJump:         tunes.SoundJump,
	game.EventLand:         tunes.SoundLand,
	game.EventShiftGrey:    tunes.SoundShiftGrey,
	game.EventShiftYellow:  tunes.SoundShiftYellow,
	game.EventShiftBlocked: tunes.SoundBlocked,
	game.EventShiftQueued:  tunes.SoundQueued,
	game.EventFragment:     tunes.SoundFragment,
	game.EventSecret:       tunes.SoundSecret,
	game.EventCheckpoint:   tunes.SoundCheckpoint,
	game.EventTransmission: tunes.SoundTransmission,
	game.EventSpotted:      tunes.SoundSpotted,
	game.EventDeath:        tunes.SoundDeath,
	game.EventLamp:         tunes.SoundLamp,
	game.EventExit:         tunes.SoundExit,
}

type sfx struct {
	players [tunes.SoundCount]*audio.Player
}

func newSFX(ctx *audio.Context) *sfx {
	s := &sfx{}
	for i := range s.players {
		s.players[i] = ctx.NewPlayerF32FromBytes(stereoBytes(tunes.RenderSound(tunes.Sound(i))))
		s.players[i].SetVolume(sfxVolume)
	}
	return s
}

func stereoBytes(mono []float32) []byte {
	out := make([]byte, len(mono)*bytesPerFrame)
	for i, v := range mono {
		bits := math.Float32bits(v)
		binary.LittleEndian.PutUint32(out[i*bytesPerFrame:], bits)
		binary.LittleEndian.PutUint32(out[i*bytesPerFrame+4:], bits)
	}
	return out
}

func (s *sfx) play(events []game.Event, muted bool) {
	if muted {
		return
	}
	for _, e := range events {
		sound, ok := eventSounds[e]
		if !ok {
			continue
		}
		p := s.players[sound]
		_ = p.Rewind()
		p.Play()
	}
}
