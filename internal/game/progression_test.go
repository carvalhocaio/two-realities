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
	checkpointFirst = "fragments: 1\n[yellow]\n#.........#\n#P.C......#\n###########\n[grey]\n#........E#\n#....Y.^..#\n###########\n"
	checkpointLast  = "fragments: 0\n[yellow]\n#.........#\n#P........#\n###########\n[grey]\n#........E#\n#..Y.C.^..#\n###########\n"
)

func mustParse(t *testing.T, text string) level.Level {
	t.Helper()
	lvl, err := level.Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	return lvl
}

func walkUntilDead(s *Session, in input.Intent) {
	for range maxTicks {
		if s.Death() != Alive {
			break
		}
		s.Update(in)
	}
	for range DeathFreezeTicks {
		s.Update(idle)
	}
}

func TestCheckpoints(t *testing.T) {
	type want struct {
		active    bool
		point     level.Point
		world     world.World
		fragments int
		remaining int
	}
	spawn := level.Point{X: 1, Y: 1}
	tests := []struct {
		name  string
		level string
		play  func(s *Session)
		want  want
	}{
		{"touching in yellow activates", checkpointFirst,
			func(s *Session) { play(s, ticks(right, 10)) },
			want{active: true, point: level.Point{X: 3, Y: 1}, fragments: 1, remaining: 1}},
		{"touching in grey activates", checkpointFirst,
			func(s *Session) { play(s, then(one(shift), ticks(right, 10))) },
			want{active: true, point: level.Point{X: 3, Y: 1}, world: world.Grey, fragments: 1, remaining: 1}},
		{"death respawns at checkpoint with snapshot", checkpointFirst,
			func(s *Session) {
				play(s, then(ticks(right, 10), one(shift)))
				walkUntilDead(s, right)
			},
			want{active: true, point: level.Point{X: 3, Y: 1}, fragments: 1, remaining: 1}},
		{"fragments collected before activation stay collected", checkpointLast,
			func(s *Session) {
				play(s, one(shift))
				walkUntilDead(s, right)
			},
			want{active: true, point: level.Point{X: 5, Y: 1}, fragments: 1}},
		{"death without checkpoint respawns at spawn", checkpointLast,
			func(s *Session) {
				play(s, then(one(shift), ticks(right, 3)))
				s.die(DeathFall)
				play(s, ticks(idle, DeathFreezeTicks))
			},
			want{point: spawn, remaining: 1}},
		{"R discards checkpoint", checkpointFirst,
			func(s *Session) { play(s, then(ticks(right, 10), one(restart))) },
			want{point: spawn, fragments: 1, remaining: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(&rhythm.Manual{}, mustParse(t, tt.level))
			tt.play(s)
			_, active := s.ActiveCheckpoint()
			if active != tt.want.active {
				t.Errorf("checkpoint active = %v, want %v", active, tt.want.active)
			}
			if s.Death() == Alive && !tt.want.active && tt.want.point == spawn && s.PlayerBox() != player.New(spawn).Box {
				t.Errorf("player %+v not at spawn", s.PlayerBox())
			}
			if tt.want.active && s.Active() == world.Yellow && s.PlayerBox() != player.New(tt.want.point).Box && tt.name != "touching in yellow activates" {
				t.Errorf("player %+v not at checkpoint %v", s.PlayerBox(), tt.want.point)
			}
			if s.Active() != tt.want.world || s.Fragments() != tt.want.fragments || len(s.RemainingFragments()) != tt.want.remaining {
				t.Errorf("world=%v fragments=%d remaining=%d, want %v %d %d",
					s.Active(), s.Fragments(), len(s.RemainingFragments()), tt.want.world, tt.want.fragments, tt.want.remaining)
			}
		})
	}
}

func TestEpochChangesOnRespawnAndLevelChange(t *testing.T) {
	lvl := mustParse(t, exitLevel)
	s := NewSession(&rhythm.Manual{}, lvl, lvl, lvl)
	steps := []struct {
		name    string
		act     func()
		changes bool
	}{
		{"walking keeps epoch", func() { s.Update(input.Intent{Left: true}) }, false},
		{"death respawn bumps epoch", func() { s.die(DeathFall); play(s, ticks(idle, DeathFreezeTicks)) }, true},
		{"restart bumps epoch", func() { s.Update(restart) }, true},
		{"next level bumps epoch", func() { play(s, ticks(right, 2)) }, true},
	}
	for _, st := range steps {
		before := s.Epoch()
		st.act()
		if changed := s.Epoch() != before; changed != st.changes {
			t.Fatalf("%s: epoch changed = %v, want %v", st.name, changed, st.changes)
		}
	}
}

func TestLevelProgression(t *testing.T) {
	first := mustParse(t, exitLevel)
	second := mustParse(t, "fragments: 2\n[yellow]\n#....#\n#.PE.#\n######\n[grey]\n#....#\n#....#\n######\n")

	tests := []struct {
		name          string
		levels        []level.Level
		act           func(s *Session)
		wantIndex     int
		wantFragments int
		wantFinished  bool
	}{
		{"starts at first level", []level.Level{first, second}, func(*Session) {}, 0, 0, false},
		{"exit advances", []level.Level{first, second}, func(s *Session) { play(s, ticks(right, 2)) }, 1, 2, false},
		{"last exit finishes", []level.Level{first, second}, func(s *Session) { play(s, ticks(right, 4)) }, 1, 2, true},
		{"single level finishes", []level.Level{first}, func(s *Session) { play(s, ticks(right, 2)) }, 0, 0, true},
		{"restart stays on level", []level.Level{first, second}, func(s *Session) { play(s, then(ticks(right, 2), one(restart))) }, 1, 2, false},
		{"go to level", []level.Level{first, second}, func(s *Session) { s.GoToLevel(1) }, 1, 2, false},
		{"go to level clamps high", []level.Level{first, second}, func(s *Session) { s.GoToLevel(9) }, 1, 2, false},
		{"go to level clamps low", []level.Level{first, second}, func(s *Session) { s.GoToLevel(-1) }, 0, 0, false},
		{"go to level clears finished", []level.Level{first, second}, func(s *Session) { play(s, ticks(right, 4)); s.GoToLevel(0) }, 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(&rhythm.Manual{}, tt.levels...)
			tt.act(s)
			if s.LevelIndex() != tt.wantIndex || s.Fragments() != tt.wantFragments || s.Finished() != tt.wantFinished {
				t.Fatalf("index=%d fragments=%d finished=%v, want %d %d %v",
					s.LevelIndex(), s.Fragments(), s.Finished(), tt.wantIndex, tt.wantFragments, tt.wantFinished)
			}
		})
	}
}

func TestFlow(t *testing.T) {
	lvl := mustParse(t, exitLevel)
	tests := []struct {
		name   string
		inputs []input.Intent
		want   Screen
		level  int
	}{
		{"starts on title", nil, ScreenTitle, -1},
		{"title ignores other keys", []input.Intent{right, restart, {JumpPressed: true}}, ScreenTitle, -1},
		{"shift starts first level", one(shift), ScreenPlaying, 0},
		{"exit moves to second level", then(one(shift), ticks(right, 2)), ScreenPlaying, 1},
		{"last exit shows end screen", then(one(shift), ticks(right, 4)), ScreenEnd, -1},
		{"end screen ignores other keys", then(one(shift), ticks(right, 4), ticks(right, 3)), ScreenEnd, -1},
		{"shift on end returns to title", then(one(shift), ticks(right, 4), one(shift)), ScreenTitle, -1},
		{"new game starts from first level", then(one(shift), ticks(right, 4), one(shift), one(shift)), ScreenPlaying, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlow(&rhythm.Manual{}, lvl, lvl)
			for _, in := range tt.inputs {
				f.Update(in)
			}
			if f.Screen() != tt.want {
				t.Fatalf("screen = %v, want %v", f.Screen(), tt.want)
			}
			if tt.level >= 0 && f.Session().LevelIndex() != tt.level {
				t.Fatalf("level = %d, want %d", f.Session().LevelIndex(), tt.level)
			}
		})
	}
}
