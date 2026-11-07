package game

import (
	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/input"
)

const EpilogueTicksPerRune = 2

var EpilogueLines = []i18n.Key{i18n.Epilogue1, i18n.Epilogue2, i18n.Epilogue3, i18n.Epilogue4}

type epilogueState struct {
	line, ticks int
}

func (f *Flow) epilogueText(line int) []rune { return []rune(f.lang.T(EpilogueLines[line])) }

func (f *Flow) lineTicks() int { return len(f.epilogueText(f.epilogue.line)) * EpilogueTicksPerRune }

func (f *Flow) updateEpilogue(in input.Intent) {
	e := &f.epilogue
	if !in.StartPressed() {
		e.ticks = min(e.ticks+1, f.lineTicks())
		return
	}
	if e.ticks < f.lineTicks() {
		e.ticks = f.lineTicks()
		return
	}
	e.line, e.ticks = e.line+1, 0
	if e.line == len(EpilogueLines) {
		f.screen = ScreenEnd
	}
}

func (f *Flow) Epilogue() (lines []string, lineDone bool) {
	e := f.epilogue
	for i := range min(e.line, len(EpilogueLines)) {
		lines = append(lines, string(f.epilogueText(i)))
	}
	if e.line >= len(EpilogueLines) {
		return lines, true
	}
	text := f.epilogueText(e.line)
	return append(lines, string(text[:e.ticks/EpilogueTicksPerRune])), e.ticks >= f.lineTicks()
}
