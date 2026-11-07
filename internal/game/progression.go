package game

func (s *Session) LevelIndex() int { return s.index }
func (s *Session) LevelCount() int { return len(s.levels) }

func (s *Session) GoToLevel(i int) {
	s.finished = false
	s.timer.runEligible = false
	s.enterLevel(min(max(i, 0), len(s.levels)-1))
}

func (s *Session) enterLevel(i int) {
	s.index = i
	s.level = s.levels[i]
	s.resetLevelTimer()
	s.pickSecret()
	s.shifted, s.shiftHint = false, false
	s.Restart()
}

func (s *Session) nextLevel() {
	if s.index == len(s.levels)-1 {
		s.finished = true
		return
	}
	s.enterLevel(s.index + 1)
}
