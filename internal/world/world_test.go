package world

import "testing"

func TestWorld(t *testing.T) {
	if Yellow.Other() != Grey || Grey.Other() != Yellow {
		t.Fatal("Other should swap worlds")
	}
	if Yellow.String() != "yellow" || Grey.String() != "grey" {
		t.Fatalf("names %q %q", Yellow, Grey)
	}
}
