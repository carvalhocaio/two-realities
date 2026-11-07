package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

func skipEpilogue() []input.Intent { return ticks(shift, 2*len(EpilogueLines)) }

func passEpilogue(t *testing.T, f *Flow) {
	t.Helper()
	for range 2*len(EpilogueLines) + 1 {
		if f.Screen() != ScreenEpilogue {
			return
		}
		f.Update(shift)
	}
	t.Fatal("the epilogue never ended")
}

func epilogueFlow(t *testing.T) *Flow {
	t.Helper()
	f := seenFlow(&rhythm.Manual{}, &save.Memory{}, mustParse(t, exitLevel))
	f.Update(shift)
	for range 20 {
		if f.Screen() == ScreenEpilogue {
			return f
		}
		f.Update(right)
	}
	t.Fatalf("screen %v, want the epilogue", f.Screen())
	return nil
}

func TestEpilogueTypesEachLine(t *testing.T) {
	f := epilogueFlow(t)
	first := []rune(f.Language().T(EpilogueLines[0]))
	if lines, done := f.Epilogue(); len(lines) != 1 || lines[0] != "" || done {
		t.Fatalf("start: %q done %v", lines, done)
	}
	for range 3 * EpilogueTicksPerRune {
		f.Update(idle)
	}
	if lines, _ := f.Epilogue(); lines[0] != string(first[:3]) {
		t.Fatalf("after 3 runes: %q", lines[0])
	}
	for range len(first) * EpilogueTicksPerRune * 2 {
		f.Update(idle)
	}
	if lines, done := f.Epilogue(); lines[0] != string(first) || !done {
		t.Fatalf("typed out: %q done %v", lines[0], done)
	}
}

func TestEpilogueConfirmCompletesThenAdvances(t *testing.T) {
	f := epilogueFlow(t)
	f.Update(shift)
	lines, done := f.Epilogue()
	if len(lines) != 1 || lines[0] != f.Language().T(EpilogueLines[0]) || !done {
		t.Fatalf("confirm should complete the line: %q done %v", lines, done)
	}
	f.Update(shift)
	if lines, done := f.Epilogue(); len(lines) != 2 || lines[1] != "" || done {
		t.Fatalf("second confirm should start the next line: %q done %v", lines, done)
	}
	for range 2*len(EpilogueLines) - 3 {
		f.Update(shift)
	}
	if f.Screen() != ScreenEpilogue {
		t.Fatalf("ended early on %v", f.Screen())
	}
	if lines, done := f.Epilogue(); len(lines) != len(EpilogueLines) || !done {
		t.Fatalf("last line: %d lines done %v", len(lines), done)
	}
	f.Update(shift)
	if f.Screen() != ScreenEnd {
		t.Fatalf("screen %v after the last line, want end", f.Screen())
	}
	if lines, _ := f.Epilogue(); len(lines) != len(EpilogueLines) {
		t.Fatalf("after the end %d lines", len(lines))
	}
}

func TestEpilogueIgnoresOtherKeys(t *testing.T) {
	f := epilogueFlow(t)
	for _, in := range []input.Intent{right, left, {JumpPressed: true}, {PausePressed: true}, {FocusLost: true}, restart, {HelpPressed: true}} {
		f.Update(in)
		if f.Screen() != ScreenEpilogue {
			t.Fatalf("%+v left the epilogue for %v", in, f.Screen())
		}
	}
	if lines, _ := f.Epilogue(); len(lines) != 1 {
		t.Fatal("other keys should not advance the lines")
	}
}

func TestEpilogueFollowsTheLanguage(t *testing.T) {
	f := epilogueFlow(t)
	f.lang = i18n.PT
	f.Update(shift)
	if lines, _ := f.Epilogue(); lines[0] != i18n.PT.T(EpilogueLines[0]) {
		t.Fatalf("portuguese line %q", lines[0])
	}
}
