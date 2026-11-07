package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

type countingStore struct {
	save.Memory
	saves int
}

func (c *countingStore) Save(d save.Data) error {
	c.saves++
	return c.Memory.Save(d)
}

func TestRecordsSkipWritesThatChangeNothing(t *testing.T) {
	store := &countingStore{Memory: save.Memory{Stored: save.Data{
		Language:         "pt",
		BestRun:          100,
		PressureUnlocked: true,
		PressureCleared:  3,
		HelpSeen:         true,
		Secrets:          map[string]bool{"a": true},
		Reached:          map[string]bool{"a": true},
		AudioOffsetMs:    20,
	}}}
	r := NewRecords(store)
	for name, write := range map[string]func(){
		"same language":      func() { r.setLanguage(i18n.PT) },
		"same audio offset":  func() { r.setAudioOffset(20) },
		"slower run":         func() { r.recordRun(200) },
		"lower pressure":     func() { r.clearPressure(2) },
		"pressure unlocked":  func() { r.unlockPressure() },
		"known secret":       func() { r.markSecret("a") },
		"known reached":      func() { r.markReached("a") },
		"help already seen":  func() { r.markHelpSeen() },
		"slower level split": func() { r.recordLevel("a", 50); r.recordLevel("a", 60) },
	} {
		before := store.saves
		write()
		if extra := store.saves - before; extra > 1 || (name != "slower level split" && extra != 0) {
			t.Errorf("%s wrote %d times", name, extra)
		}
	}
}

func TestPressureCannotBeChosenWhileLocked(t *testing.T) {
	f := seenFlow(&rhythm.Manual{}, &save.Memory{}, mustParse(t, timedLevel))
	f.Update(moreKey)
	if f.pressure != 0 || f.Pressure() != 0 {
		t.Fatalf("locked pressure moved to %d", f.pressure)
	}
}
