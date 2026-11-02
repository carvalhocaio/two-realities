package game

import "github.com/carvalhocaio/two-realities/internal/geom"

func (s *Session) SecretBox() (geom.Rect, bool) {
	sec := s.level.Secret
	if sec == nil || s.secretTaken || sec.Layer != s.active {
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
	s.fragments++
	s.restartHint = false
}
