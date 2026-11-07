package game

import (
	"slices"
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
)

const (
	checkpointLevel = "fragments: 0\n[yellow]\n#.....#\n#P.C..#\n#######\n[grey]\n#....E#\n#.....#\n#######\n"
	jumpLevel       = "fragments: 0\n[yellow]\n#....#\n#....#\n#....#\n#....#\n#P...#\n######\n[grey]\n#...E#\n#....#\n#....#\n#....#\n#....#\n######\n"
	radioLevel      = "fragments: 0\n[yellow]\n#.....#\n#P.T..#\n#######\n[grey]\n#....E#\n#.....#\n#######\n[transmissions]\nHELLO\n"
)

func recordEvents(s *Session, inputs []input.Intent) []Event {
	var out []Event
	for _, in := range inputs {
		s.Update(in)
		out = append(out, s.Events()...)
	}
	return out
}

func TestSessionEvents(t *testing.T) {
	jump := input.Intent{JumpPressed: true}
	tests := []struct {
		name   string
		level  string
		inputs []input.Intent
		want   []Event
	}{
		{"walking is silent", openLevel, ticks(right, 3), nil},
		{"jump and land", jumpLevel, then(one(jump), ticks(idle, 40)), []Event{EventJump, EventLand}},
		{"shift into grey", openLevel, one(shift), []Event{EventShiftGrey}},
		{"shift back to yellow", openLevel, then(one(shift), one(shift)), []Event{EventShiftGrey, EventShiftYellow}},
		{"blocked by a solid tile", blockedLevel, one(shift), []Event{EventShiftBlocked}},
		{"blocked without fragments", brokeLevel, then(one(shift), one(shift)), []Event{EventShiftGrey, EventShiftBlocked}},
		{"fragment", fragmentLevel, then(one(shift), ticks(right, 8)), []Event{EventShiftGrey, EventFragment}},
		{"checkpoint", checkpointLevel, ticks(right, 15), []Event{EventCheckpoint}},
		{"transmission", radioLevel, ticks(right, 15), []Event{EventTransmission}},
		{"exit", exitLevel, ticks(right, 10), []Event{EventExit}},
		{"fall", pitLevel, then(one(shift), ticks(idle, 20)), []Event{EventShiftGrey, EventDeath}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvl, err := level.Parse(tt.level)
			if err != nil {
				t.Fatal(err)
			}
			got := recordEvents(NewSession(&rhythm.Manual{}, lvl), tt.inputs)
			if !slices.Equal(got, tt.want) {
				t.Fatalf("events %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventsDrainOnce(t *testing.T) {
	s := newSession(t, openLevel)
	s.Update(shift)
	if len(s.Events()) != 1 || len(s.Events()) != 0 {
		t.Fatal("events should be returned once")
	}
}
