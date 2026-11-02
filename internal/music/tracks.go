package music

type Track int

const (
	ActOne Track = iota
	ActTwo
	ActThree
	WardenOne
	WardenTwo
	WardenThree
	TrackCount
)

type chord struct {
	root  int
	major bool
}

func (c chord) tones() [3]int {
	third := 3
	if c.major {
		third = 4
	}
	return [3]int{c.root, c.root + third, c.root + 7}
}

type note struct {
	beat, length float64
	pitch        int
}

type song struct {
	chords                 [BarsPerLoop]chord
	kickGain, hatGain      float64
	kickEvery              int
	hats16, snare, rimFour bool
	bass                   bassStyle
	pad                    bool
	arp                    int
	lead                   []note
}

type bassStyle int

const (
	bassHalves bassStyle = iota
	bassEighths
	bassOctaves
)

var (
	bassVoice = voice{wave: saw, gain: 0.3, attack: 0.005, decay: 0.12, sustain: 0.55, release: 0.05, cutoff: 320}
	padVoice  = voice{wave: saw, gain: 0.05, attack: 0.4, decay: 0.6, sustain: 0.8, release: 0.5, cutoff: 1100, detune: 0.006}
	arpVoice  = voice{wave: triangle, gain: 0.11, attack: 0.003, decay: 0.12, sustain: 0.2, release: 0.05}
	leadVoice = voice{wave: square, gain: 0.06, attack: 0.01, decay: 0.2, sustain: 0.5, release: 0.12, cutoff: 2200}
)

var songs = [TrackCount]song{
	ActOne: {
		chords:   [BarsPerLoop]chord{{45, false}, {41, true}, {48, true}, {43, true}},
		kickGain: 0.45, kickEvery: 2, hatGain: 0.07,
		bass: bassHalves, pad: true, arp: 2,
	},
	ActTwo: {
		chords:   [BarsPerLoop]chord{{38, false}, {46, true}, {43, false}, {45, true}},
		kickGain: 0.6, kickEvery: 1, hatGain: 0.08, rimFour: true,
		bass: bassEighths, pad: true, arp: 2,
	},
	ActThree: {
		chords:   [BarsPerLoop]chord{{40, false}, {36, true}, {43, true}, {38, true}},
		kickGain: 0.7, kickEvery: 1, hatGain: 0.06, hats16: true, snare: true,
		bass: bassOctaves, pad: true, arp: 4,
		lead: []note{
			{0, 2, 64}, {2, 1, 67}, {3, 1, 71},
			{4, 2, 72}, {6, 1, 71}, {7, 1, 67},
			{8, 2, 67}, {10, 1, 71}, {11, 1, 74},
			{12, 1.5, 74}, {13.5, 0.5, 72}, {14, 2, 71},
		},
	},
	WardenOne: {
		chords:   [BarsPerLoop]chord{{36, false}, {44, true}, {41, false}, {43, true}},
		kickGain: 0.7, kickEvery: 1, hatGain: 0.08,
		bass: bassEighths, pad: true,
	},
	WardenTwo: {
		chords:   [BarsPerLoop]chord{{36, false}, {44, true}, {41, false}, {43, true}},
		kickGain: 0.7, kickEvery: 1, hatGain: 0.06, hats16: true,
		bass: bassEighths, pad: true, arp: 4,
	},
	WardenThree: {
		chords:   [BarsPerLoop]chord{{36, false}, {44, true}, {41, false}, {43, true}},
		kickGain: 0.75, kickEvery: 1, hatGain: 0.06, hats16: true, snare: true,
		bass: bassOctaves, pad: true, arp: 4,
		lead: []note{
			{0, 1, 72}, {1, 1, 75}, {2, 2, 79},
			{4, 2, 80}, {6, 1, 79}, {7, 1, 75},
			{8, 2, 77}, {10, 1, 80}, {11, 1, 77},
			{12, 3, 79}, {15, 1, 71},
		},
	},
}

func Render(t Track, bpm int) []float32 {
	s := songs[t]
	beat := BeatFrames(bpm)
	b := make(buffer, BarFrames(bpm)*BarsPerLoop)
	for bar, c := range s.chords {
		start := bar * BeatsPerBar * beat
		s.drums(b, start, beat)
		s.bassLine(b, start, beat, c)
		if s.pad {
			for _, n := range c.tones() {
				padVoice.play(b, start, BeatsPerBar*beat, n+12)
			}
		}
		s.arpeggio(b, start, beat, c)
	}
	for _, n := range s.lead {
		leadVoice.play(b, int(n.beat*float64(beat)), int(n.length*float64(beat)), n.pitch)
	}
	return normalize(b)
}

func (s song) drums(b buffer, start, beat int) {
	for i := range BeatsPerBar {
		at := start + i*beat
		if i%s.kickEvery == 0 {
			kick(b, at, s.kickGain)
		}
		hat(b, at+beat/2, s.hatGain)
		if s.hats16 {
			hat(b, at+beat/4, s.hatGain/2)
			hat(b, at+3*beat/4, s.hatGain/2)
		}
		if s.snare && i%2 == 1 {
			snare(b, at, 0.3)
		}
	}
	if s.rimFour {
		snare(b, start+3*beat, 0.12)
	}
}

func (s song) bassLine(b buffer, start, beat int, c chord) {
	switch s.bass {
	case bassHalves:
		for i := 0; i < BeatsPerBar; i += 2 {
			bassVoice.play(b, start+i*beat, 2*beat-beat/8, c.root)
		}
	case bassEighths:
		for i := range 2 * BeatsPerBar {
			bassVoice.play(b, start+i*beat/2, beat/2-beat/16, c.root)
		}
	case bassOctaves:
		for i := range 2 * BeatsPerBar {
			bassVoice.play(b, start+i*beat/2, beat/2-beat/16, c.root+12*(i%2))
		}
	}
}

func (s song) arpeggio(b buffer, start, beat int, c chord) {
	if s.arp == 0 {
		return
	}
	tones := c.tones()
	step := beat / s.arp
	for i := range BeatsPerBar * s.arp {
		pitch := tones[i%3] + 24
		if i%6 >= 3 {
			pitch += 12
		}
		arpVoice.play(b, start+i*step, step, pitch)
	}
}
