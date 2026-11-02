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

func Load(name string) (Level, error) {
	data, err := files.ReadFile(levelsDir + "/" + name + ".txt")
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

func Names() ([]string, error) {
	entries, err := fs.ReadDir(files, levelsDir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, strings.TrimSuffix(e.Name(), ".txt"))
	}
	sort.Strings(names)
	return names, nil
}

func LoadAll() ([]Level, error) {
	names, err := Names()
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return nil, ErrNoLevels
	}
	levels := make([]Level, 0, len(names))
	for _, name := range names {
		lvl, err := Load(name)
		if err != nil {
			return nil, err
		}
		levels = append(levels, lvl)
	}
	return levels, nil
}
