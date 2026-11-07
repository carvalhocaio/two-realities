package game

import (
	"math/rand/v2"

	"github.com/carvalhocaio/two-realities/internal/geom"
	"github.com/carvalhocaio/two-realities/internal/level"
)

func randomSpot(n int) int { return rand.IntN(n) }

func (s *Session) UseSecretPicker(pick func(n int) int) {
	s.picker = pick
	s.pickSecret()
}

func (s *Session) pickSecret() {
	s.secretSpot = 0
	if n := len(s.level.Secrets); n > 1 {
		s.secretSpot = s.picker(n)
	}
}

func (s *Session) Secret() (level.Marker, bool) {
	if len(s.level.Secrets) == 0 {
		return level.Marker{}, false
	}
	return s.level.Secrets[s.secretSpot], true
}

func (s *Session) SecretBox() (geom.Rect, bool) {
	sec, ok := s.Secret()
	if !ok || s.secretTaken || sec.Layer != s.active {
		return geom.Rect{}, false
	}
	return FragmentBox(sec.Point), true
}

func (s *Session) SecretTaken() bool { return s.secretTaken }

func (s *Session) collectSecret() {
	box, ok := s.SecretBox()
	if !ok || !s.player.Box.Overlaps(box) {
		return
	}
	s.secretTaken = true
	s.secretFeedback = SecretFeedbackTicks
	s.fragments++
	s.restartHint = false
}
