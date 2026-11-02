package platform

import "sync/atomic"

type Event int

const (
	Blur Event = iota
	Hidden
	FullscreenExit
)

type Signals struct {
	pending    atomic.Bool
	unfocused  atomic.Bool
	fullscreen atomic.Bool
}

func (s *Signals) Notify(Event) { s.pending.Store(true) }

func (s *Signals) SetFocused(focused bool) {
	wasUnfocused := s.unfocused.Swap(!focused)
	if !focused && !wasUnfocused {
		s.Notify(Blur)
	}
}

func (s *Signals) SetFullscreen(on bool) {
	if s.fullscreen.Swap(on) && !on {
		s.Notify(FullscreenExit)
	}
}

func (s *Signals) Poll() bool { return s.pending.Swap(false) }
