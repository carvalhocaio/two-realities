package game

import (
	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/save"
)

type Records struct {
	store save.Storage
	data  save.Data
}

func NewRecords(store save.Storage) *Records {
	data, _ := store.Load()
	return &Records{store: store, data: data.Normalized()}
}

func (r *Records) Degraded() bool {
	d, ok := r.store.(interface{ Degraded() bool })
	return ok && d.Degraded()
}

func (r *Records) BestLevel(name string) (int, bool) {
	t, ok := r.data.BestLevel[name]
	return t, ok
}

func (r *Records) BestRun() (int, bool) {
	return r.data.BestRun, r.data.BestRun > 0
}

func (r *Records) SecretFound(name string) bool { return r.data.Secrets[name] }

func (r *Records) SecretsFound(levels []level.Level) (found, total int) {
	for _, l := range levels {
		if len(l.Secrets) == 0 {
			continue
		}
		total++
		if r.SecretFound(l.Name) {
			found++
		}
	}
	return found, total
}

func (r *Records) Language() (i18n.Lang, bool) { return i18n.Parse(r.data.Language) }

func (r *Records) setLanguage(l i18n.Lang) {
	if r.data.Language == l.Code() {
		return
	}
	r.data.Language = l.Code()
	r.persist()
}

func (r *Records) AudioOffset() int { return r.data.AudioOffsetMs }

func (r *Records) setAudioOffset(ms int) {
	if r.data.AudioOffsetMs == ms {
		return
	}
	r.data.AudioOffsetMs = ms
	r.persist()
}

func (r *Records) recordLevel(name string, ticks int) (best int, improved bool) {
	if old, ok := r.data.BestLevel[name]; ok && old <= ticks {
		return old, false
	}
	r.data.BestLevel[name] = ticks
	r.persist()
	return ticks, true
}

func (r *Records) recordRun(ticks int) bool {
	if r.data.BestRun > 0 && r.data.BestRun <= ticks {
		return false
	}
	r.data.BestRun = ticks
	r.persist()
	return true
}

func (r *Records) PressureUnlocked() bool { return r.data.PressureUnlocked || r.data.BestRun > 0 }
func (r *Records) PressureCleared() int   { return r.data.PressureCleared }

func (r *Records) unlockPressure() {
	if r.data.PressureUnlocked {
		return
	}
	r.data.PressureUnlocked = true
	r.persist()
}

func (r *Records) clearPressure(p int) bool {
	if p <= r.data.PressureCleared {
		return false
	}
	r.data.PressureCleared = p
	r.persist()
	return true
}

func (r *Records) Reached(levels []level.Level, i int) bool {
	switch {
	case i == 0, r.CampaignDone(), r.data.Reached[levels[i].Name]:
		return true
	}
	_, cleared := r.data.BestLevel[levels[i-1].Name]
	return cleared
}

func (r *Records) Furthest(levels []level.Level) int {
	for i := len(levels) - 1; i > 0; i-- {
		if r.Reached(levels, i) {
			return i
		}
	}
	return 0
}

func (r *Records) CampaignDone() bool { return r.PressureUnlocked() }

func (r *Records) HelpSeen() bool {
	return r.data.HelpSeen || r.PressureUnlocked() || len(r.data.BestLevel) > 0 || len(r.data.Reached) > 0
}

func (r *Records) markHelpSeen() {
	if r.data.HelpSeen {
		return
	}
	r.data.HelpSeen = true
	r.persist()
}

func (r *Records) markReached(name string) {
	if r.data.Reached[name] {
		return
	}
	r.data.Reached[name] = true
	r.persist()
}

func (r *Records) markSecret(name string) {
	if r.data.Secrets[name] {
		return
	}
	r.data.Secrets[name] = true
	r.persist()
}

func (r *Records) persist() {
	_ = r.store.Save(r.data)
}
