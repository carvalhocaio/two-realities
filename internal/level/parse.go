package level

import (
	"errors"
	"fmt"
	"strings"

	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	tileSolid = '#'
	tileEmpty = '.'
	tileSpawn = 'P'
	tileExit  = 'E'
)

var (
	ErrUnknownSection   = errors.New("unknown section")
	ErrDuplicateSection = errors.New("duplicate section")
	ErrRowOutsideLayer  = errors.New("row outside of a section")
	ErrMissingSection   = errors.New("missing section")
	ErrRaggedRows       = errors.New("rows have different lengths")
	ErrSizeMismatch     = errors.New("layers have different sizes")
	ErrUnknownTile      = errors.New("unknown tile")
	ErrMarkerCount      = errors.New("level needs exactly one spawn and one exit")
)

var sectionNames = map[string]world.World{
	"[yellow]": world.Yellow,
	"[grey]":   world.Grey,
}

func Parse(text string) (Level, error) {
	rows, err := splitSections(text)
	if err != nil {
		return Level{}, err
	}
	if err := checkSizes(rows); err != nil {
		return Level{}, err
	}

	var lvl Level
	markers := map[rune][]Point{}
	for w := range world.Count {
		grid, err := parseGrid(rows[w], markers)
		if err != nil {
			return Level{}, fmt.Errorf("%s layer: %w", world.World(w), err)
		}
		lvl.Layers[w] = grid
	}

	if len(markers[tileSpawn]) != 1 || len(markers[tileExit]) != 1 {
		return Level{}, ErrMarkerCount
	}
	lvl.Spawn = markers[tileSpawn][0]
	lvl.Exit = markers[tileExit][0]
	return lvl, nil
}

func splitSections(text string) ([world.Count][]string, error) {
	var rows [world.Count][]string
	seen := map[world.World]bool{}
	current := world.World(-1)

	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") {
			w, ok := sectionNames[line]
			if !ok {
				return rows, fmt.Errorf("%w: %s", ErrUnknownSection, line)
			}
			if seen[w] {
				return rows, fmt.Errorf("%w: %s", ErrDuplicateSection, line)
			}
			seen[w] = true
			current = w
			continue
		}
		if current < 0 {
			return rows, ErrRowOutsideLayer
		}
		rows[current] = append(rows[current], line)
	}

	for w := range world.Count {
		if len(rows[w]) == 0 {
			return rows, fmt.Errorf("%w: %s", ErrMissingSection, world.World(w))
		}
	}
	return rows, nil
}

func checkSizes(rows [world.Count][]string) error {
	width := len(rows[0][0])
	for _, layer := range rows {
		for _, row := range layer {
			if len(row) != width {
				return ErrRaggedRows
			}
		}
	}
	if len(rows[world.Yellow]) != len(rows[world.Grey]) {
		return ErrSizeMismatch
	}
	return nil
}

func parseGrid(rows []string, markers map[rune][]Point) (Grid, error) {
	width, height := len(rows[0]), len(rows)
	solid := make([]bool, width*height)
	for y, row := range rows {
		for x, ch := range row {
			switch ch {
			case tileSolid:
				solid[y*width+x] = true
			case tileEmpty:
			case tileSpawn, tileExit:
				markers[ch] = append(markers[ch], Point{X: x, Y: y})
			default:
				return Grid{}, fmt.Errorf("%w %q at %d,%d", ErrUnknownTile, ch, x, y)
			}
		}
	}
	return NewGrid(width, height, solid), nil
}
