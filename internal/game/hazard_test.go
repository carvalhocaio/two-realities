package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/player"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	yellowSpikes   = "fragments: 0\n[yellow]\n#....#\n#P^..#\n######\n[grey]\n#...E#\n#....#\n######\n"
	greySpikes     = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#.^..#\n######\n"
	spikeOnSpawn   = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#^...#\n######\n"
	fragmentSpikes = "fragments: 0\n[yellow]\n#....#\n#P...#\n######\n[grey]\n#...E#\n#.Y^.#\n######\n"
	watcherBehind  = "fragments: 1\n[yellow]\n#.................#\n#....P............#\n###################\n[grey]\n#...............E.#\n#........W........#\n###################\n"
	watcherAhead   = "fragments: 1\n[yellow]\n#.................#\n#.......P.........#\n###################\n[grey]\n#...............E.#\n#W................#\n###################\n"
)

func TestSpikes(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		inputs []input.Intent
		want   DeathCause
	}{
		{"yellow spikes kill in yellow", yellowSpikes, ticks(right, 5), DeathSpikes},
		{"yellow spikes are harmless in grey", yellowSpikes, then(one(shift), ticks(right, 5)), Alive},
		{"grey spikes are harmless in yellow", greySpikes, ticks(right, 5), Alive},
		{"grey spikes kill in grey", greySpikes, then(one(shift), ticks(right, 5)), DeathSpikes},
		{"shifting onto spikes kills", spikeOnSpawn, one(shift), DeathSpikes},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSession(t, tt.level)
			play(s, tt.inputs)
			if s.Death() != tt.want {
				t.Fatalf("death = %v, want %v", s.Death(), tt.want)
			}
		})
	}
}

func TestDeathFreezeAndReset(t *testing.T) {
	s := newSession(t, fragmentSpikes)
	play(s, then(one(shift), ticks(right, 3)))
	if s.Fragments() != 1 {
		t.Fatalf("fragments = %d, want 1 before dying", s.Fragments())
	}
	for range 10 {
		if s.Death() != Alive {
			break
		}
		s.Update(right)
	}
	if s.Death() != DeathSpikes {
		t.Fatalf("death = %v, want spikes", s.Death())
	}

	frozen := s.PlayerBox()
	play(s, then(one(restart), ticks(right, DeathFreezeTicks-2)))
	if s.Death() != DeathSpikes || s.PlayerBox() != frozen {
		t.Fatalf("expected frozen death, got %v at %+v", s.Death(), s.PlayerBox())
	}

	play(s, one(idle))
	got := observe(s)
	want := state{world: world.Yellow, remaining: 1, atSpawn: true}
	if s.Death() != Alive || got != want {
		t.Fatalf("after freeze: death=%v\n got %+v\nwant %+v", s.Death(), got, want)
	}
}

func TestWatcherAlert(t *testing.T) {
	jump := input.Intent{JumpPressed: true}
	tests := []struct {
		name      string
		level     string
		inputs    []input.Intent
		wantAlert int
		wantDeath DeathCause
		wantWorld world.World
	}{
		{"cannot see yellow", watcherAhead, ticks(idle, 60), 0, Alive, world.Yellow},
		{"cannot see behind", watcherBehind, then(one(shift), ticks(idle, 60)), 0, Alive, world.Grey},
		{"alert fills while seen", watcherAhead, then(one(shift), ticks(idle, 9)), 10, Alive, world.Grey},
		{"caught when meter is full", watcherAhead, then(one(shift), ticks(idle, AlertTicks-1)), AlertTicks, DeathCaught, world.Grey},
		{"leaving sight resets meter", watcherAhead, then(one(shift), ticks(idle, 9), one(jump), ticks(idle, 3)), 0, Alive, world.Grey},
		{"dodge by shifting back to yellow", watcherAhead, then(one(shift), ticks(idle, AlertTicks-2), one(shift)), 0, Alive, world.Yellow},
		{"dodge holds while yellow", watcherAhead, then(one(shift), ticks(idle, AlertTicks-2), one(shift), ticks(idle, 120)), 0, Alive, world.Yellow},
		{"dodge too late", watcherAhead, then(one(shift), ticks(idle, AlertTicks-1), one(shift)), AlertTicks, DeathCaught, world.Grey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSession(t, tt.level)
			play(s, tt.inputs)
			w := s.Watchers()[0]
			if w.Alert != tt.wantAlert || s.Death() != tt.wantDeath || s.Active() != tt.wantWorld {
				t.Fatalf("alert=%d death=%v world=%v, want alert=%d death=%v world=%v",
					w.Alert, s.Death(), s.Active(), tt.wantAlert, tt.wantDeath, tt.wantWorld)
			}
		})
	}
}

func TestWatchersPatrolTheSameInBothWorlds(t *testing.T) {
	yellow := newSession(t, watcherBehind)
	grey := newSession(t, watcherBehind)
	grey.Update(shift)
	yellow.Update(idle)
	for range 200 {
		yellow.Update(idle)
		grey.Update(idle)
	}
	if yellow.Watchers()[0].Box != grey.Watchers()[0].Box {
		t.Fatalf("yellow %+v vs grey %+v", yellow.Watchers()[0].Box, grey.Watchers()[0].Box)
	}
}

func TestCaughtResetsLevel(t *testing.T) {
	s := newSession(t, watcherAhead)
	start := NewWatcher(s.Level().Watchers[0])
	play(s, then(one(shift), ticks(idle, AlertTicks-1+DeathFreezeTicks)))
	w := s.Watchers()[0]
	if s.Death() != Alive || s.Active() != world.Yellow || s.Fragments() != 1 {
		t.Fatalf("death=%v world=%v fragments=%d", s.Death(), s.Active(), s.Fragments())
	}
	if w.Box != start.Box || w.Alert != 0 || s.PlayerBox() != player.New(s.Level().Spawn).Box {
		t.Fatalf("watcher %+v player %+v not reset", w, s.PlayerBox())
	}
}

func TestLevelProgression(t *testing.T) {
	parse := func(text string) level.Level {
		lvl, err := level.Parse(text)
		if err != nil {
			t.Fatal(err)
		}
		return lvl
	}
	first := parse(exitLevel)
	second := parse("fragments: 2\n[yellow]\n#....#\n#.PE.#\n######\n[grey]\n#....#\n#....#\n######\n")

	tests := []struct {
		name          string
		levels        []level.Level
		inputs        []input.Intent
		wantIndex     int
		wantFragments int
	}{
		{"starts at first level", []level.Level{first, second}, nil, 0, 0},
		{"exit advances", []level.Level{first, second}, ticks(right, 2), 1, 2},
		{"last level loops to first", []level.Level{first, second}, ticks(right, 4), 0, 0},
		{"single level loops to itself", []level.Level{first}, ticks(right, 2), 0, 0},
		{"restart stays on level", []level.Level{first, second}, then(ticks(right, 2), one(restart)), 1, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSession(tt.levels...)
			play(s, tt.inputs)
			if s.LevelIndex() != tt.wantIndex || s.Fragments() != tt.wantFragments {
				t.Fatalf("index=%d fragments=%d, want index=%d fragments=%d",
					s.LevelIndex(), s.Fragments(), tt.wantIndex, tt.wantFragments)
			}
			if s.PlayerBox() != player.New(s.Level().Spawn).Box {
				t.Fatalf("player not at spawn: %+v", s.PlayerBox())
			}
		})
	}
}
