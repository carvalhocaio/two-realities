package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	openLevel      = "fragments: 1\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#....#\n######\n"
	brokeLevel     = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#....#\n######\n"
	blockedLevel   = "fragments: 1\n[yellow]\n#...E#\n#P...#\n######\n[grey]\n#....#\n##...#\n######\n"
	yellowWall     = "fragments: 1\n[yellow]\n#....#\n#P.#E#\n######\n[grey]\n#....#\n#....#\n######\n"
	brokeWall      = "fragments: 0\n[yellow]\n#....#\n#P.#E#\n######\n[grey]\n#....#\n#....#\n######\n"
	fragmentLevel  = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#.Y..#\n######\n"
	fragmentToExit = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#....#\n#.Y.E#\n######\n"
	exitLevel      = "fragments: 0\n[yellow]\n#....#\n#PE..#\n######\n[grey]\n#....#\n#....#\n######\n"
	pitLevel       = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#....#\n#....#\n"
)

var (
	idle    = input.Intent{}
	shift   = input.Intent{ShiftPressed: true}
	right   = input.Intent{Right: true}
	left    = input.Intent{Left: true}
	restart = input.Intent{RestartPressed: true}
)

func ticks(in input.Intent, n int) []input.Intent {
	out := make([]input.Intent, n)
	for i := range out {
		out[i] = in
	}
	return out
}

func then(parts ...[]input.Intent) []input.Intent {
	var out []input.Intent
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func one(in input.Intent) []input.Intent { return []input.Intent{in} }

func newSession(t *testing.T, text string) *Session {
	t.Helper()
	lvl, err := level.Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	return NewSession(&rhythm.Manual{}, lvl, lvl)
}

func play(s *Session, inputs []input.Intent) {
	for _, in := range inputs {
		s.Update(in)
	}
}

type state struct {
	world     world.World
	fragments int
	remaining int
	feedback  Feedback
	hint      bool
	atSpawn   bool
}

func observe(s *Session) state {
	return state{
		world:     s.Active(),
		fragments: s.Fragments(),
		remaining: len(s.RemainingFragments()),
		feedback:  s.Feedback(),
		hint:      s.RestartHint(),
		atSpawn:   s.PlayerBox() == player.New(s.Level().Spawn).Box,
	}
}

func TestSession(t *testing.T) {
	full := Feedback{}
	half := BlockedFeedbackTicks / 2
	pickup := then(one(shift), ticks(right, 5))

	tests := []struct {
		name   string
		level  string
		inputs []input.Intent
		want   state
	}{
		{"starts in yellow with header fragments", openLevel, nil,
			state{world: world.Yellow, fragments: 1, atSpawn: true}},
		{"yellow to grey is free", openLevel, one(shift),
			state{world: world.Grey, fragments: 1, atSpawn: true}},
		{"grey to yellow costs one", openLevel, ticks(shift, 2),
			state{world: world.Yellow, fragments: 0, atSpawn: true}},
		{"grey to yellow blocked without fragments", brokeLevel, ticks(shift, 2),
			state{world: world.Grey, feedback: Feedback{CounterFlash: 1}, hint: true, atSpawn: true}},
		{"counter flash fades", brokeLevel, then(ticks(shift, 2), ticks(idle, half)),
			state{world: world.Grey, feedback: Feedback{CounterFlash: 0.5}, hint: true, atSpawn: true}},
		{"hint stays after flash ends", brokeLevel, then(ticks(shift, 2), ticks(idle, BlockedFeedbackTicks)),
			state{world: world.Grey, feedback: full, hint: true, atSpawn: true}},
		{"yellow to grey blocked by solid", blockedLevel, one(shift),
			state{world: world.Yellow, fragments: 1, feedback: Feedback{Shake: 1}, atSpawn: true}},
		{"shake fades", blockedLevel, then(one(shift), ticks(idle, half)),
			state{world: world.Yellow, fragments: 1, feedback: Feedback{Shake: 0.5}, atSpawn: true}},
		{"grey to yellow blocked by solid keeps fragment", yellowWall, then(one(shift), ticks(right, 9), one(shift)),
			state{world: world.Grey, fragments: 1, feedback: Feedback{Shake: 1}}},
		{"no fragments wins over solid", brokeWall, then(one(shift), ticks(right, 9), one(shift)),
			state{world: world.Grey, feedback: Feedback{CounterFlash: 1}, hint: true}},
		{"no pickup while yellow", fragmentLevel, ticks(right, 5),
			state{world: world.Yellow, remaining: 1}},
		{"pickup while grey", fragmentLevel, pickup,
			state{world: world.Grey, fragments: 1}},
		{"fragment picked up once", fragmentLevel, then(pickup, ticks(left, 5), ticks(right, 5)),
			state{world: world.Grey, fragments: 1}},
		{"spending collected fragment", fragmentLevel, then(pickup, one(shift)),
			state{world: world.Yellow}},
		{"pickup clears hint", fragmentLevel, then(ticks(shift, 2), ticks(right, 5)),
			state{world: world.Grey, fragments: 1, feedback: Feedback{CounterFlash: 7.0 / 12}}},
		{"restart key resets fragments", fragmentLevel, then(pickup, one(restart)),
			state{world: world.Yellow, remaining: 1, atSpawn: true}},
		{"restart key clears hint", brokeLevel, then(ticks(shift, 2), one(restart)),
			state{world: world.Yellow, atSpawn: true}},
		{"reaching exit resets fragments", fragmentToExit, then(one(shift), ticks(right, 16)),
			state{world: world.Yellow, remaining: 1, atSpawn: true}},
		{"reaching exit restarts", exitLevel, ticks(right, 2),
			state{world: world.Yellow, atSpawn: true}},
		{"falling out dies then restarts", pitLevel, then(one(shift), ticks(idle, 20+DeathFreezeTicks)),
			state{world: world.Yellow, atSpawn: true}},
		{"moving away leaves spawn", openLevel, ticks(right, 3),
			state{world: world.Yellow, fragments: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSession(t, tt.level)
			play(s, tt.inputs)
			if got := observe(s); got != tt.want {
				t.Fatalf("\n got %+v\nwant %+v", got, tt.want)
			}
		})
	}
}

func TestTryShiftOutcome(t *testing.T) {
	tests := []struct {
		name  string
		level string
		setup []input.Intent
		want  ShiftOutcome
	}{
		{"free shift into grey", openLevel, nil, ShiftDone},
		{"paid shift into yellow", openLevel, one(shift), ShiftDone},
		{"solid tile in grey", blockedLevel, nil, ShiftBlockedSolid},
		{"solid tile in yellow", yellowWall, then(one(shift), ticks(right, 9)), ShiftBlockedSolid},
		{"no fragments", brokeLevel, one(shift), ShiftBlockedNoFragments},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSession(t, tt.level)
			play(s, tt.setup)
			if got := s.TryShift(); got != tt.want {
				t.Fatalf("TryShift() = %v, want %v", got, tt.want)
			}
		})
	}
}
