package game

import (
	"testing"

	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/rhythm"
	"github.com/carvalhocaio/two-realities/internal/save"
)

type pad struct {
	f    *Flow
	prev input.PadState
}

func (p *pad) press(b input.PadButton) {
	var cur input.PadState
	cur.Connected, cur.Buttons[b] = true, true
	p.send(cur)
	p.send(input.PadState{Connected: true})
}

func (p *pad) stick(x, y float64) {
	p.send(input.PadState{Connected: true, StickX: x, StickY: y})
	p.send(input.PadState{Connected: true})
}

func (p *pad) send(cur input.PadState) {
	p.f.Update(input.PadIntent(p.prev, cur))
	p.prev = cur
}

func TestAGamepadAloneCanPlay(t *testing.T) {
	p := &pad{f: NewFlow(&rhythm.Manual{}, &save.Memory{}, mustParse(t, timedLevel))}
	p.press(input.PadA)
	if p.f.Screen() != ScreenHelp {
		t.Fatalf("A on a new save: screen %v, want the help card", p.f.Screen())
	}
	p.press(input.PadA)
	if p.f.Screen() != ScreenPlaying {
		t.Fatalf("A closes the card into the run: screen %v", p.f.Screen())
	}
	p.press(input.PadB)
	if p.f.Session().Active() == 0 {
		t.Fatal("B should shift worlds")
	}
	p.press(input.PadStart)
	if p.f.Screen() != ScreenPaused {
		t.Fatalf("start pauses: screen %v", p.f.Screen())
	}
	p.stick(0, 0.9)
	if p.f.Menu() != MenuHelp {
		t.Fatalf("stick down moves the menu: %v", p.f.Menu())
	}
	p.stick(0, 0.2)
	if p.f.Menu() != MenuHelp {
		t.Fatal("a push inside the dead zone should not move the menu")
	}
	p.press(input.PadDown)
	p.press(input.PadA)
	if !p.f.Assist() {
		t.Fatal("A confirms the selected item")
	}
	p.press(input.PadStart)
	if p.f.Screen() != ScreenPlaying {
		t.Fatalf("start resumes: screen %v", p.f.Screen())
	}
}

func TestConfirmStartsAndLeavesTheEndScreen(t *testing.T) {
	confirm := input.Intent{ConfirmPressed: true}
	f := seenFlow(&rhythm.Manual{}, &save.Memory{}, mustParse(t, exitLevel))
	f.Update(confirm)
	if f.Screen() != ScreenPlaying {
		t.Fatalf("confirm on the title: screen %v", f.Screen())
	}
	for range 20 {
		f.Update(right)
	}
	for f.Screen() == ScreenEpilogue {
		f.Update(confirm)
	}
	if f.Screen() != ScreenEnd {
		t.Fatalf("screen %v, want the end screen", f.Screen())
	}
	f.Update(confirm)
	if f.Screen() != ScreenTitle {
		t.Fatalf("confirm on the end screen: screen %v", f.Screen())
	}
}
