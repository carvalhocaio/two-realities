package level

import (
	"embed"
	"fmt"
)

//go:embed levels/*.txt
var files embed.FS

func Load(name string) (Level, error) {
	data, err := files.ReadFile("levels/" + name + ".txt")
	if err != nil {
		return Level{}, fmt.Errorf("load level %s: %w", name, err)
	}
	lvl, err := Parse(string(data))
	if err != nil {
		return Level{}, fmt.Errorf("parse level %s: %w", name, err)
	}
	return lvl, nil
}
