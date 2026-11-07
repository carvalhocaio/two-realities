package game

import "github.com/carvalhocaio/two-realities/internal/level"

type snapshot struct {
	index     int
	point     level.Point
	fragments int
	collected []bool
	secret    bool
	lit       int
}

func (s *Session) startState() snapshot {
	if s.checkpoint != nil {
		return *s.checkpoint
	}
	return snapshot{
		index:     -1,
		point:     s.level.Spawn,
		fragments: s.level.StartFragments,
		collected: make([]bool, len(s.level.Fragments)),
	}
}

func (s *Session) activateCheckpoints() {
	if !s.CheckpointsActive() {
		return
	}
	for i, p := range s.level.Checkpoints {
		if s.checkpoint != nil && s.checkpoint.index == i {
			continue
		}
		if s.player.Box.Overlaps(level.TileRect(p)) {
			s.checkpoint = &snapshot{
				index:     i,
				point:     p,
				fragments: s.fragments,
				collected: append([]bool(nil), s.collected...),
				secret:    s.secretTaken,
			}
			s.emit(EventCheckpoint)
		}
	}
}

func (s *Session) ActiveCheckpoint() (int, bool) {
	if s.checkpoint == nil || s.checkpoint.index < 0 {
		return 0, false
	}
	return s.checkpoint.index, true
}
