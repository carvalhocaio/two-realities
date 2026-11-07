//go:build !js

package save

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileErrors(t *testing.T) {
	tests := []struct {
		name  string
		setup func(dir string) string
		load  bool
	}{
		{"load from a directory", func(dir string) string { return dir }, true},
		{"save under a file", func(dir string) string {
			blocker := filepath.Join(dir, "blocker")
			_ = os.WriteFile(blocker, nil, filePerm)
			return filepath.Join(blocker, "save.json")
		}, false},
		{"save with a busy temp path", func(dir string) string {
			path := filepath.Join(dir, "save.json")
			_ = os.Mkdir(path+".tmp", dirPerm)
			return path
		}, false},
		{"save over a directory", func(dir string) string {
			path := filepath.Join(dir, "save.json")
			_ = os.MkdirAll(filepath.Join(path, "full"), dirPerm)
			return path
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := File{Path: tt.setup(t.TempDir())}
			if tt.load {
				d, err := f.Load()
				if err == nil || d.BestLevel == nil || d.Secrets == nil || d.Reached == nil {
					t.Fatalf("load: %+v %v, want normalized empty data and an error", d, err)
				}
				return
			}
			if err := f.Save(Data{BestRun: 1}); err == nil {
				t.Fatal("save should fail")
			}
		})
	}
}

func TestDefaultStorage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	t.Setenv("HOME", dir)
	f, ok := Default().(File)
	if !ok || filepath.Dir(f.Path) != filepath.Join(dir, appDir) {
		t.Fatalf("default storage %#v", Default())
	}
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")
	if _, ok := Default().(*Memory); !ok {
		t.Fatalf("without a config dir, got %#v", Default())
	}
}

type panicOnSave struct{ Memory }

func (*panicOnSave) Save(Data) error { panic("QuotaExceededError") }

func TestFallbackSurvivesAPanickingSave(t *testing.T) {
	f := NewFallback(&panicOnSave{})
	if _, err := f.Load(); err != nil || f.Degraded() {
		t.Fatalf("load %v degraded %v", err, f.Degraded())
	}
	if err := f.Save(Data{BestRun: 5}); err != nil || !f.Degraded() {
		t.Fatalf("save %v degraded %v", err, f.Degraded())
	}
	if d, _ := f.Load(); d.BestRun != 5 {
		t.Fatalf("memory copy lost: %+v", d)
	}
	if !errors.Is(safeSave(&panicOnSave{}, Data{}), ErrUnavailable) {
		t.Fatal("a panic should read as unavailable")
	}
}
