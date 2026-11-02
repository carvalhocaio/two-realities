package platform

import "testing"

func TestSignals(t *testing.T) {
	tests := []struct {
		name  string
		feed  func(s *Signals)
		polls []bool
	}{
		{"quiet", func(*Signals) {}, []bool{false}},
		{"blur", func(s *Signals) { s.Notify(Blur) }, []bool{true, false}},
		{"hidden", func(s *Signals) { s.Notify(Hidden) }, []bool{true, false}},
		{"several events, one request", func(s *Signals) { s.Notify(Blur); s.Notify(Hidden) }, []bool{true, false}},
		{"focus loss", func(s *Signals) { s.SetFocused(true); s.SetFocused(false) }, []bool{true, false}},
		{"staying unfocused asks once", func(s *Signals) { s.SetFocused(false); s.SetFocused(false); s.SetFocused(false) }, []bool{true, false}},
		{"focused is quiet", func(s *Signals) { s.SetFocused(true); s.SetFocused(true) }, []bool{false}},
		{"fullscreen exit", func(s *Signals) { s.SetFullscreen(true); s.SetFullscreen(false) }, []bool{true, false}},
		{"entering fullscreen is quiet", func(s *Signals) { s.SetFullscreen(false); s.SetFullscreen(true) }, []bool{false}},
		{"windowed resize is quiet", func(s *Signals) { s.SetFullscreen(false); s.SetFullscreen(false) }, []bool{false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Signals
			tt.feed(&s)
			for i, want := range tt.polls {
				if got := s.Poll(); got != want {
					t.Fatalf("poll %d = %v, want %v", i, got, want)
				}
			}
		})
	}
}

func TestFocusLossAgainAfterRefocus(t *testing.T) {
	var s Signals
	s.SetFocused(false)
	s.Poll()
	s.SetFocused(true)
	if s.Poll() {
		t.Fatal("regaining focus must not ask for a pause")
	}
	s.SetFocused(false)
	if !s.Poll() {
		t.Fatal("losing focus again should ask again")
	}
}
