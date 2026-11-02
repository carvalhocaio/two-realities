package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	signalLevel  = "fragments: 0\n[yellow]\n#............................#\n#P.T.....................^..E#\n##############################\n[grey]\n#............................#\n#..T.........................#\n##############################\n[transmissions]\nwarm words\ncold words\n"
	secretLevel  = "fragments: 0\n[yellow]\n#........#\n#P......E#\n##########\n[grey]\n#........#\n#...S....#\n##########\n"
	yellowSecret = "fragments: 0\n[yellow]\n#........#\n#P..S...E#\n##########\n[grey]\n#........#\n#........#\n##########\n"
	timedLevel   = "fragments: 0\n[yellow]\n#......#\n#P....E#\n########\n[grey]\n#......#\n#......#\n########\n"
)

func namedLevel(t *testing.T, name, text string) level.Level {
	t.Helper()
	lvl := mustParse(t, text)
	lvl.Name = name
	return lvl
}

func storedSession(t *testing.T, store save.Storage, levels ...level.Level) *Session {
	t.Helper()
	s := NewSession(&rhythm.Manual{}, levels...)
	s.UseRecords(NewRecords(store))
	return s
}

func TestTransmissions(t *testing.T) {
	tests := []struct {
		name  string
		play  []input.Intent
		open  bool
		text  string
		heard []bool
	}{
		{"untouched stays silent", ticks(idle, 5), false, "", []bool{false, false}},
		{"touching in its layer opens it", ticks(right, 12), true, "WARM WORDS", []bool{true, false}},
		{"other layer marker stays silent in yellow", ticks(right, 12), true, "WARM WORDS", []bool{true, false}},
		{"grey marker needs grey", then(one(shift), ticks(right, 12)), true, "COLD WORDS", []bool{false, true}},
		{"closes after its time", then(ticks(right, 12), ticks(idle, TransmissionTicks+TransmissionFadeTicks)), false, "", []bool{true, false}},
		{"stays open nearby", then(ticks(right, 12), ticks(right, 60)), true, "WARM WORDS", []bool{true, false}},
		{"closes when walking away", then(ticks(right, 12), ticks(right, 130)), false, "", []bool{true, false}},
		{"does not retrigger in the same attempt", then(ticks(right, 12), ticks(idle, TransmissionTicks+TransmissionFadeTicks), ticks(left, 4), ticks(right, 4)), false, "", []bool{true, false}},
		{"restart allows hearing it again", then(ticks(right, 12), one(restart), ticks(right, 12)), true, "WARM WORDS", []bool{true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSession(t, signalLevel)
			play(s, tt.play)
			view, open := s.Transmission()
			if open != tt.open || view.TextIn(i18n.EN) != tt.text {
				t.Fatalf("open %v text %q, want %v %q", open, view.TextIn(i18n.EN), tt.open, tt.text)
			}
			for i, want := range tt.heard {
				if s.Heard(i) != want {
					t.Errorf("heard[%d] = %v, want %v", i, s.Heard(i), want)
				}
			}
		})
	}
}

func TestTransmissionRetriggersAfterDeath(t *testing.T) {
	s := newSession(t, signalLevel)
	play(s, ticks(right, 12))
	walkUntilDead(s, right)
	if s.Heard(0) {
		t.Fatal("death should start a new attempt")
	}
	play(s, ticks(right, 12))
	if view, open := s.Transmission(); !open || view.TextIn(i18n.EN) != "WARM WORDS" {
		t.Fatalf("want the transmission again, got %v %q", open, view.TextIn(i18n.EN))
	}
}

func TestTransmissionFadesBeforeClosing(t *testing.T) {
	s := newSession(t, signalLevel)
	play(s, then(ticks(right, 12), ticks(idle, TransmissionTicks+TransmissionFadeTicks/2)))
	view, open := s.Transmission()
	if !open || view.Fade <= 0 || view.Fade >= 1 {
		t.Fatalf("want a partial fade, got open %v fade %v", open, view.Fade)
	}
}

func TestSecrets(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		play      []input.Intent
		taken     bool
		fragments int
	}{
		{"grey secret ignored while yellow", secretLevel, ticks(right, 30), false, 0},
		{"grey secret taken in grey", secretLevel, then(one(shift), ticks(right, 30)), true, 1},
		{"yellow secret taken in yellow", yellowSecret, ticks(right, 25), true, 1},
		{"yellow secret ignored while grey", yellowSecret, then(one(shift), ticks(right, 25)), false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSession(t, tt.level)
			play(s, tt.play)
			if s.SecretTaken() != tt.taken || s.Fragments() != tt.fragments {
				t.Fatalf("taken %v fragments %d, want %v %d", s.SecretTaken(), s.Fragments(), tt.taken, tt.fragments)
			}
			if _, visible := s.SecretBox(); tt.taken && visible {
				t.Fatal("a taken secret must not be drawn")
			}
		})
	}
}

func TestSecretPersistsOnlyWhenTheLevelIsFinished(t *testing.T) {
	store := &save.Memory{}
	lvl := namedLevel(t, "secret", secretLevel)
	s := storedSession(t, store, lvl, lvl)
	play(s, then(one(shift), ticks(right, 30)))
	if store.Stored.Secrets["secret"] {
		t.Fatal("secret saved before finishing the level")
	}
	play(s, ticks(right, 40))
	if s.LevelIndex() != 1 || !store.Stored.Secrets["secret"] {
		t.Fatalf("level %d, saved %v", s.LevelIndex(), store.Stored.Secrets)
	}
	reloaded := NewRecords(store)
	if found, total := reloaded.SecretsFound([]level.Level{lvl}); found != 1 || total != 1 {
		t.Fatalf("reloaded secrets %d/%d", found, total)
	}
}

func TestTimer(t *testing.T) {
	lvl := namedLevel(t, "timed", timedLevel)
	tests := []struct {
		name    string
		stored  map[string]int
		play    []input.Intent
		ticks   int
		best    int
		newBest bool
	}{
		{"idle does not start", nil, ticks(idle, 50), 0, 0, false},
		{"starts on first input", nil, then(ticks(idle, 50), ticks(right, 3)), 3, 0, false},
		{"first finish sets best", nil, ticks(right, 60), 0, 0, true},
		{"faster replaces best", map[string]int{"timed": 500}, ticks(right, 60), 0, 0, true},
		{"slower keeps best", map[string]int{"timed": 10}, ticks(right, 60), 0, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &save.Memory{Stored: save.Data{BestLevel: tt.stored}}
			s := storedSession(t, store, lvl, lvl)
			play(s, tt.play)
			res, done := s.LevelResult()
			if !done {
				if s.LevelTicks() != tt.ticks {
					t.Fatalf("level ticks %d, want %d", s.LevelTicks(), tt.ticks)
				}
				return
			}
			if res.NewBest != tt.newBest || !res.HasBest {
				t.Fatalf("result %+v, want new best %v", res, tt.newBest)
			}
			if tt.best > 0 && res.Best != tt.best {
				t.Fatalf("best %d, want %d", res.Best, tt.best)
			}
			if store.Stored.BestLevel["timed"] != res.Best {
				t.Fatalf("stored %d, want %d", store.Stored.BestLevel["timed"], res.Best)
			}
		})
	}
}

func TestTimerStopsAtExitAndRunAddsUp(t *testing.T) {
	store := &save.Memory{}
	lvl := namedLevel(t, "timed", timedLevel)
	s := storedSession(t, store, lvl, lvl)
	for s.LevelIndex() == 0 {
		s.Update(right)
	}
	first, _ := s.LevelResult()
	if s.LevelTicks() != 0 || first.Ticks == 0 {
		t.Fatalf("after exit: level %d ticks %d result %+v", s.LevelIndex(), s.LevelTicks(), first)
	}
	play(s, ticks(right, 60))
	sum := s.Summary()
	if !s.Finished() || !sum.NewBest || !sum.Eligible || store.Stored.BestRun != sum.Ticks {
		t.Fatalf("summary %+v stored %d", sum, store.Stored.BestRun)
	}
	if sum.Ticks < 2*first.Ticks {
		t.Fatalf("run %d shorter than two levels of %d", sum.Ticks, first.Ticks)
	}
}

func TestLevelSkipMakesRunIneligible(t *testing.T) {
	store := &save.Memory{}
	lvl := namedLevel(t, "timed", timedLevel)
	s := storedSession(t, store, lvl, lvl)
	s.GoToLevel(1)
	play(s, ticks(right, 60))
	if sum := s.Summary(); sum.Eligible || store.Stored.BestRun != 0 {
		t.Fatalf("summary %+v, stored run %d", sum, store.Stored.BestRun)
	}
}

func TestAssistRunsDoNotSaveBests(t *testing.T) {
	tests := []struct {
		name     string
		assistAt int
	}{
		{"assist from the start", 0},
		{"assist turned on mid run", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &save.Memory{}
			lvl := namedLevel(t, "timed", timedLevel)
			s := storedSession(t, store, lvl, lvl)
			for i := range 2 {
				if i == tt.assistAt {
					s.SetAssist(true)
				}
				play(s, ticks(right, 60))
			}
			sum := s.Summary()
			if !s.Finished() || !sum.Assisted || store.Stored.BestRun != 0 {
				t.Fatalf("summary %+v, stored run %d", sum, store.Stored.BestRun)
			}
			if _, ok := store.Stored.BestLevel["timed"]; ok && tt.assistAt == 0 {
				t.Fatal("assisted level time was saved")
			}
		})
	}
}

func TestAssistAlertTakesTwiceAsLong(t *testing.T) {
	tests := []struct {
		assist bool
		idle   int
		want   DeathCause
	}{
		{false, AlertTicks, DeathCaught},
		{true, AlertTicks, Alive},
		{true, 2*AlertTicks - 2, Alive},
		{true, 2 * AlertTicks, DeathCaught},
	}
	for _, tt := range tests {
		s := newSession(t, watcherAhead)
		s.SetAssist(tt.assist)
		play(s, then(one(shift), ticks(idle, tt.idle)))
		if s.Death() != tt.want {
			t.Errorf("assist %v after %d ticks: death %v, want %v", tt.assist, tt.idle, s.Death(), tt.want)
		}
	}
}

func TestAssistPushesOutOfAutoOverlap(t *testing.T) {
	d := newBeatDriver(t, autoGreyWalled)
	d.s.SetAssist(true)
	before := d.s.PlayerBox()
	d.untilBeat(2)
	if d.s.Death() != Alive || d.s.Active() != world.Grey {
		t.Fatalf("death %v world %v, want alive in grey", d.s.Death(), d.s.Active())
	}
	if d.s.PlayerBox() == before {
		t.Fatal("player was not moved out of the wall")
	}
}

func TestPauseMenu(t *testing.T) {
	pause := input.Intent{PausePressed: true}
	down := input.Intent{DownPressed: true}
	up := input.Intent{UpPressed: true}
	confirm := input.Intent{ConfirmPressed: true}
	mute := input.Intent{MutePressed: true}
	tests := []struct {
		name   string
		inputs []input.Intent
		screen Screen
		assist bool
		muted  bool
	}{
		{"escape pauses", then(one(shift), one(pause)), ScreenPaused, false, false},
		{"escape resumes", then(one(shift), one(pause), one(pause)), ScreenPlaying, false, false},
		{"resume item", then(one(shift), one(pause), one(confirm)), ScreenPlaying, false, false},
		{"assist toggles", then(one(shift), one(pause), one(down), one(confirm)), ScreenPaused, true, false},
		{"mute toggles", then(one(shift), one(pause), one(down), one(down), one(confirm)), ScreenPaused, false, true},
		{"menu wraps upward to title", then(one(shift), one(pause), one(up), one(confirm)), ScreenTitle, false, false},
		{"m mutes anywhere", then(one(mute), one(shift)), ScreenPlaying, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl := mustParse(t, timedLevel)
			f := NewFlow(&rhythm.Manual{}, &save.Memory{}, lvl)
			for _, in := range tt.inputs {
				f.Update(in)
			}
			if f.Screen() != tt.screen || f.Assist() != tt.assist || f.Muted() != tt.muted {
				t.Fatalf("screen %v assist %v muted %v, want %v %v %v", f.Screen(), f.Assist(), f.Muted(), tt.screen, tt.assist, tt.muted)
			}
			if tt.assist && !f.Session().Assist() {
				t.Fatal("assist not applied to the session")
			}
		})
	}
}

func TestPausedSessionDoesNotMove(t *testing.T) {
	lvl := mustParse(t, timedLevel)
	f := NewFlow(&rhythm.Manual{}, &save.Memory{}, lvl)
	f.Update(shift)
	f.Update(input.Intent{PausePressed: true})
	before := f.Session().PlayerBox()
	for range 30 {
		f.Update(right)
	}
	if f.Session().PlayerBox() != before || f.Session().LevelTicks() != 0 {
		t.Fatal("session advanced while paused")
	}
}

func TestAssistCarriesIntoNewRuns(t *testing.T) {
	lvl := mustParse(t, timedLevel)
	f := NewFlow(&rhythm.Manual{}, &save.Memory{}, lvl)
	confirm := input.Intent{ConfirmPressed: true}
	for _, in := range then(one(shift), one(input.Intent{PausePressed: true}), one(input.Intent{DownPressed: true}), one(confirm), one(input.Intent{UpPressed: true}), one(input.Intent{UpPressed: true}), one(confirm), one(shift)) {
		f.Update(in)
	}
	if f.Screen() != ScreenPlaying || !f.Session().Assist() {
		t.Fatalf("screen %v assist %v", f.Screen(), f.Session().Assist())
	}
}
