package save

import (
	"errors"
	"testing"
)

type brokenStorage struct {
	loadErr  error
	saveErr  error
	panics   bool
	saved    int
	contents Data
}

func (b *brokenStorage) Load() (Data, error) {
	if b.panics {
		panic("SecurityError: access denied")
	}
	return b.contents, b.loadErr
}

func (b *brokenStorage) Save(d Data) error {
	if b.panics {
		panic("QuotaExceededError")
	}
	if b.saveErr != nil {
		return b.saveErr
	}
	b.saved++
	b.contents = d
	return nil
}

func TestKeyIsPrefixed(t *testing.T) {
	tests := []struct{ name, want string }{
		{"save", "two-realities:save"},
		{"probe", "two-realities:probe"},
	}
	for _, tt := range tests {
		if got := Key(tt.name); got != tt.want {
			t.Errorf("Key(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestFallback(t *testing.T) {
	tests := []struct {
		name         string
		primary      *brokenStorage
		degraded     bool
		primarySaves int
	}{
		{"healthy storage is used", &brokenStorage{}, false, 1},
		{"unavailable on load falls back", &brokenStorage{loadErr: ErrUnavailable}, true, 0},
		{"failing writes fall back", &brokenStorage{saveErr: errors.New("quota")}, true, 0},
		{"panicking storage falls back", &brokenStorage{panics: true}, true, 0},
		{"corrupt data keeps the storage", &brokenStorage{loadErr: errors.New("bad json")}, false, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFallback(tt.primary)
			if _, err := f.Load(); err != nil {
				t.Fatalf("load error leaked: %v", err)
			}
			if err := f.Save(Data{BestRun: 42}); err != nil {
				t.Fatalf("save error leaked: %v", err)
			}
			got, err := f.Load()
			if err != nil || got.BestRun != 42 || got.BestLevel == nil {
				t.Fatalf("session data lost: %+v %v", got, err)
			}
			if f.Degraded() != tt.degraded || tt.primary.saved != tt.primarySaves {
				t.Fatalf("degraded %v saves %d, want %v %d", f.Degraded(), tt.primary.saved, tt.degraded, tt.primarySaves)
			}
		})
	}
}

func TestFallbackStopsTouchingBrokenStorage(t *testing.T) {
	primary := &brokenStorage{saveErr: errors.New("quota")}
	f := NewFallback(primary)
	_ = f.Save(Data{BestRun: 1})
	primary.saveErr = nil
	_ = f.Save(Data{BestRun: 2})
	if primary.saved != 0 {
		t.Fatal("a degraded fallback must stay in memory for the session")
	}
}
