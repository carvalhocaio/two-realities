package input

import (
	"reflect"
	"testing"
)

func pad(buttons ...PadButton) PadState {
	p := PadState{Connected: true}
	for _, b := range buttons {
		p.Buttons[b] = true
	}
	return p
}

func stick(x, y float64) PadState { return PadState{Connected: true, StickX: x, StickY: y} }

func TestPadIntent(t *testing.T) {
	tests := []struct {
		name      string
		prev, cur PadState
		want      Intent
	}{
		{"idle", PadState{}, pad(), Intent{}},
		{"a jumps and confirms", pad(), pad(PadA), Intent{JumpPressed: true, ConfirmPressed: true}},
		{"holding a does not repeat", pad(PadA), pad(PadA), Intent{}},
		{"b shifts", pad(), pad(PadB), Intent{ShiftPressed: true}},
		{"x shifts", pad(), pad(PadX), Intent{ShiftPressed: true}},
		{"lb shifts", pad(), pad(PadLB), Intent{ShiftPressed: true}},
		{"rb shifts", pad(), pad(PadRB), Intent{ShiftPressed: true}},
		{"a second shift button while one is held", pad(PadB), pad(PadB, PadRB), Intent{ShiftPressed: true}},
		{"y opens help", pad(), pad(PadY), Intent{HelpPressed: true}},
		{"back restarts", pad(), pad(PadBack), Intent{RestartPressed: true}},
		{"start pauses", pad(), pad(PadStart), Intent{PausePressed: true}},
		{"dpad up", pad(), pad(PadUp), Intent{UpPressed: true}},
		{"dpad down", pad(), pad(PadDown), Intent{DownPressed: true}},
		{"dpad left walks and decreases", pad(), pad(PadLeft), Intent{Left: true, DecreasePressed: true}},
		{"dpad right held walks only", pad(PadRight), pad(PadRight), Intent{Right: true}},
		{"stick right past the dead zone", stick(0, 0), stick(0.9, 0), Intent{Right: true, IncreasePressed: true}},
		{"stick left past the dead zone", stick(0, 0), stick(-0.9, 0), Intent{Left: true, DecreasePressed: true}},
		{"stick inside the dead zone", stick(0, 0), stick(0.3, -0.3), Intent{}},
		{"stick up for menus", stick(0, 0), stick(0, -0.8), Intent{UpPressed: true}},
		{"stick down for menus", stick(0, 0), stick(0, 0.8), Intent{DownPressed: true}},
		{"stick held down does not repeat", stick(0, 0.8), stick(0, 0.9), Intent{}},
		{"releasing gives nothing", pad(PadA, PadLeft), pad(), Intent{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadIntent(tt.prev, tt.cur); got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestCombinePads(t *testing.T) {
	a := pad(PadA)
	a.StickX = 0.2
	b := stick(-0.7, 0.1)
	b.Buttons[PadStart] = true
	got := a.Combine(b)
	if !got.Connected || !got.Buttons[PadA] || !got.Buttons[PadStart] || got.StickX != -0.7 || got.StickY != 0.1 {
		t.Fatalf("combined %+v", got)
	}
	if none := (PadState{}).Combine(PadState{}); none.Connected {
		t.Fatal("no pads means not connected")
	}
}

func TestMergeKeepsEveryField(t *testing.T) {
	v := reflect.ValueOf(&Intent{}).Elem()
	for i := range v.NumField() {
		var one Intent
		reflect.ValueOf(&one).Elem().Field(i).SetBool(true)
		name := v.Type().Field(i).Name
		if Merge(one, Intent{}) != one || Merge(Intent{}, one) != one {
			t.Errorf("merge drops %s", name)
		}
	}
	if Merge(Intent{}, Intent{}) != (Intent{}) {
		t.Fatal("merge of nothing is nothing")
	}
}

func TestStartPressed(t *testing.T) {
	for in, want := range map[Intent]bool{
		{}:                     false,
		{ShiftPressed: true}:   true,
		{ConfirmPressed: true}: true,
		{JumpPressed: true}:    false,
	} {
		if in.StartPressed() != want {
			t.Errorf("%+v start = %v", in, !want)
		}
	}
}
