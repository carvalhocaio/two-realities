package game

func (s *Session) LevelIndex() int { return s.index }
func (s *Session) LevelCount() int { return len(s.levels) }

func (s *Session) enterLevel(i int) {
	s.index = i
	s.level = s.levels[i]
	s.Restart()
}

func (s *Session) nextLevel() {
	s.enterLevel((s.index + 1) % len(s.levels))
}
