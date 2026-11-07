//go:build !js

package platform

import "testing"

func TestLocale(t *testing.T) {
	tests := []struct {
		name                string
		all, messages, lang string
		want                string
	}{
		{"lc_all wins", "pt_BR.UTF-8", "en_US", "en_US", "pt_BR.UTF-8"},
		{"then lc_messages", "", "pt_PT", "en_US", "pt_PT"},
		{"then lang", "", "", "en_GB", "en_GB"},
		{"nothing set", "", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LC_ALL", tt.all)
			t.Setenv("LC_MESSAGES", tt.messages)
			t.Setenv("LANG", tt.lang)
			if got := Locale(); got != tt.want {
				t.Fatalf("locale %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListenIsANoOpOnDesktop(t *testing.T) {
	s := &Signals{}
	Listen(s)
	if s.Poll() {
		t.Fatal("desktop listen should not raise signals")
	}
}
