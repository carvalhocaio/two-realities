//go:build !js

package save

import (
	"os"
	"path/filepath"
	"testing"
)

func sample() Data {
	return Data{
		BestLevel: map[string]int{"01_first_shift": 600},
		BestRun:   9000,
		Secrets:   map[string]bool{"03_spikes": true},
	}
}

func TestStorageRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		store Storage
	}{
		{"memory", &Memory{}},
		{"file", File{Path: filepath.Join(t.TempDir(), "nested", "save.json")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			empty, err := tt.store.Load()
			if err != nil || empty.BestLevel == nil || empty.Secrets == nil || empty.BestRun != 0 {
				t.Fatalf("empty load = %+v, %v", empty, err)
			}
			if err := tt.store.Save(sample()); err != nil {
				t.Fatal(err)
			}
			got, err := tt.store.Load()
			if err != nil {
				t.Fatal(err)
			}
			if got.BestLevel["01_first_shift"] != 600 || got.BestRun != 9000 || !got.Secrets["03_spikes"] {
				t.Fatalf("loaded %+v", got)
			}
		})
	}
}

func TestFileWithGarbageLoadsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "save.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := File{Path: path}.Load()
	if err == nil {
		t.Fatal("want a decode error")
	}
	if got.BestLevel == nil || got.Secrets == nil {
		t.Fatalf("want usable empty data, got %+v", got)
	}
}
