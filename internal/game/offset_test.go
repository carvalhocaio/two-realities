package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

var (
	lessKey = input.Intent{DecreasePressed: true}
	moreKey = input.Intent{IncreasePressed: true}
)

func toOffsetItem() []input.Intent {
	return then(one(shift), one(pauseKey), ticks(downKey, menuDowns(MenuAudioOffset)))
}

func TestAudioOffsetMenu(t *testing.T) {
	tests := []struct {
		name   string
		stored int
		inputs []input.Intent
		want   int
	}{
		{"starts at zero", 0, nil, 0},
		{"loads the saved value", 70, nil, 70},
		{"clamps a bad saved value", 9000, nil, AudioOffsetLimit},
		{"increase", 0, then(toOffsetItem(), ticks(moreKey, 3)), 3 * AudioOffsetStep},
		{"decrease below zero", 0, then(toOffsetItem(), ticks(lessKey, 2)), -2 * AudioOffsetStep},
		{"clamps at the top", 0, then(toOffsetItem(), ticks(moreKey, 80)), AudioOffsetLimit},
		{"clamps at the bottom", 0, then(toOffsetItem(), ticks(lessKey, 80)), -AudioOffsetLimit},
		{"other items ignore it", 0, then(one(shift), one(pauseKey), ticks(moreKey, 3)), 0},
		{"gameplay ignores it", 0, then(one(shift), ticks(moreKey, 3)), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &save.Memory{Stored: save.Data{AudioOffsetMs: tt.stored}}
			f := seenFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
			for _, in := range tt.inputs {
				f.Update(in)
			}
			if f.AudioOffsetMs() != tt.want {
				t.Fatalf("offset %d, want %d", f.AudioOffsetMs(), tt.want)
			}
		})
	}
}

func TestAudioOffsetPersists(t *testing.T) {
	store := &save.Memory{}
	f := seenFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
	for _, in := range then(toOffsetItem(), ticks(moreKey, 4)) {
		f.Update(in)
	}
	again := seenFlow(&rhythm.Manual{}, store, mustParse(t, timedLevel))
	if again.AudioOffsetMs() != 4*AudioOffsetStep {
		t.Fatalf("reloaded offset %d, want %d", again.AudioOffsetMs(), 4*AudioOffsetStep)
	}
}

func menuDowns(item MenuItem) int {
	if item > MenuSkip {
		return int(item) - 1
	}
	return int(item)
}
