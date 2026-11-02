package game

import "github.com/carvalhocaio/two-realities/internal/input"

const ResultShowTicks = 150

type LevelResult struct {
	Ticks   int
	Best    int
	HasBest bool
	NewBest bool
}

type RunSummary struct {
	Ticks        int
	Best         int
	HasBest      bool
	NewBest      bool
	Assisted     bool
	Eligible     bool
	Pressure     int
	NewPressure  bool
	Secrets      int
	SecretsTotal int
}

type runTimer struct {
	level           int
	levelStarted    bool
	levelAssisted   bool
	run             int
	runStarted      bool
	runEligible     bool
	assisted        bool
	newBestRun      bool
	pressureCleared bool
}

func (s *Session) tickTimer(in input.Intent) {
	t := &s.timer
	if !t.levelStarted && in.Acts() {
		t.levelStarted = true
		if s.index == 0 && t.runEligible && !t.runStarted {
			t.runStarted = true
		}
	}
	if t.levelStarted {
		t.level++
	}
	if t.runStarted {
		t.run++
	}
	s.resultTimer = max(s.resultTimer-1, 0)
}

func (s *Session) resetLevelTimer() {
	s.timer.level = 0
	s.timer.levelStarted = false
	s.timer.levelAssisted = s.assist
}

func (s *Session) completeLevel() {
	name := s.level.Name
	res := LevelResult{Ticks: s.timer.level}
	if !s.ranked() {
		res.Best, res.HasBest = s.records.BestLevel(name)
	} else {
		res.Best, res.NewBest = s.records.recordLevel(name, res.Ticks)
		res.HasBest = true
	}
	if s.secretTaken {
		s.records.markSecret(name)
	}
	s.result, s.resultTimer = res, ResultShowTicks
	if s.index == len(s.levels)-1 {
		s.finishRun()
	}
}

func (s *Session) runCounts() bool {
	t := s.timer
	return t.runStarted && t.runEligible && !t.assisted && s.pressure == 0
}

func (s *Session) finishRun() {
	t := &s.timer
	if !t.runStarted || !t.runEligible {
		return
	}
	s.records.unlockPressure()
	if !t.assisted && s.pressure > 0 {
		t.pressureCleared = s.records.clearPressure(s.pressure)
	}
	if s.runCounts() {
		t.newBestRun = s.records.recordRun(t.run)
	}
}

func (s *Session) LevelTicks() int { return s.timer.level }
func (s *Session) RunTicks() int   { return s.timer.run }

func (s *Session) LevelResult() (LevelResult, bool) {
	return s.result, s.resultTimer > 0
}

func (s *Session) Summary() RunSummary {
	best, hasBest := s.records.BestRun()
	found, total := s.records.SecretsFound(s.levels)
	return RunSummary{
		Ticks:        s.timer.run,
		Best:         best,
		HasBest:      hasBest,
		NewBest:      s.timer.newBestRun,
		Assisted:     s.timer.assisted,
		Eligible:     s.timer.runStarted && s.timer.runEligible,
		Pressure:     s.pressure,
		NewPressure:  s.timer.pressureCleared,
		Secrets:      found,
		SecretsTotal: total,
	}
}
