package input

import "testing"

func TestActs(t *testing.T) {
	tests := []struct {
		in   Intent
		want bool
	}{
		{Intent{}, false},
		{Intent{Left: true}, true},
		{Intent{Right: true}, true},
		{Intent{JumpPressed: true}, true},
		{Intent{ShiftPressed: true}, true},
		{Intent{PausePressed: true, MutePressed: true, RestartPressed: true, HelpPressed: true}, false},
	}
	for _, tt := range tests {
		if tt.in.Acts() != tt.want {
			t.Errorf("%+v acts = %v, want %v", tt.in, !tt.want, tt.want)
		}
	}
}
