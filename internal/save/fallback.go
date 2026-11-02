package save

import (
	"errors"
	"sync"
)

const KeyPrefix = "two-realities:"

var ErrUnavailable = errors.New("storage unavailable")

func Key(name string) string { return KeyPrefix + name }

type Fallback struct {
	mu       sync.Mutex
	primary  Storage
	memory   Memory
	degraded bool
}

func NewFallback(primary Storage) *Fallback {
	return &Fallback{primary: primary}
}

func (f *Fallback) Load() (Data, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.degraded {
		return f.memory.Load()
	}
	d, err := safeLoad(f.primary)
	if errors.Is(err, ErrUnavailable) {
		f.degraded = true
		return f.memory.Load()
	}
	d = d.Normalized()
	_ = f.memory.Save(d)
	return d, nil
}

func (f *Fallback) Save(d Data) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	_ = f.memory.Save(d)
	if f.degraded {
		return nil
	}
	if err := safeSave(f.primary, d); err != nil {
		f.degraded = true
	}
	return nil
}

func (f *Fallback) Degraded() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.degraded
}

func safeLoad(s Storage) (d Data, err error) {
	defer func() {
		if r := recover(); r != nil {
			d, err = Data{}, ErrUnavailable
		}
	}()
	return s.Load()
}

func safeSave(s Storage, d Data) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrUnavailable
		}
	}()
	return s.Save(d)
}
