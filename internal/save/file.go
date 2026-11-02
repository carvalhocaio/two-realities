//go:build !js

package save

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	appDir   = "two-realities"
	fileName = "save.json"
	filePerm = 0o644
	dirPerm  = 0o755
)

type File struct {
	Path string
}

func (f File) Load() (Data, error) {
	raw, err := os.ReadFile(f.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return decode(nil)
	}
	if err != nil {
		empty, _ := decode(nil)
		return empty, err
	}
	return decode(raw)
}

func (f File) Save(d Data) error {
	raw, err := encode(d)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(f.Path), dirPerm); err != nil {
		return err
	}
	tmp := f.Path + ".tmp"
	if err := os.WriteFile(tmp, raw, filePerm); err != nil {
		return err
	}
	return os.Rename(tmp, f.Path)
}

func Default() Storage {
	dir, err := os.UserConfigDir()
	if err != nil {
		return &Memory{}
	}
	return File{Path: filepath.Join(dir, appDir, fileName)}
}
