package level

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

const levelsDir = "levels"

//go:embed levels/*.txt
var files embed.FS

var ErrNoLevels = errors.New("no embedded levels")

func Load(name string) (Level, error) { return load(files, name) }

func Names() ([]string, error) { return names(files) }

func LoadAll() ([]Level, error) { return loadAll(files) }

func load(fsys fs.FS, name string) (Level, error) {
	data, err := fs.ReadFile(fsys, levelsDir+"/"+name+".txt")
	if err != nil {
		return Level{}, fmt.Errorf("load level %s: %w", name, err)
	}
	lvl, err := Parse(string(data))
	if err != nil {
		return Level{}, fmt.Errorf("parse level %s: %w", name, err)
	}
	lvl.Name = name
	return lvl, nil
}

func names(fsys fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(fsys, levelsDir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, strings.TrimSuffix(e.Name(), ".txt"))
	}
	sort.Strings(out)
	return out, nil
}

func loadAll(fsys fs.FS) ([]Level, error) {
	list, err := names(fsys)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNoLevels
	}
	levels := make([]Level, 0, len(list))
	for _, name := range list {
		lvl, err := load(fsys, name)
		if err != nil {
			return nil, err
		}
		levels = append(levels, lvl)
	}
	return levels, nil
}
