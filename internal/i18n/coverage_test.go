package i18n

import "testing"

func TestNextCyclesLanguages(t *testing.T) {
	if EN.Next() != PT || PT.Next() != EN {
		t.Fatal("languages should cycle EN -> PT -> EN")
	}
}

func TestMissingKeyFallsBackToEnglish(t *testing.T) {
	const fake Key = 1 << 20
	tables[EN][fake] = "ONLY ENGLISH"
	defer delete(tables[EN], fake)
	if got := PT.T(fake); got != "ONLY ENGLISH" {
		t.Fatalf("got %q", got)
	}
}
