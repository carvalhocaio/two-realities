package level

import (
	"errors"
	"fmt"
	"strings"

	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	tileSolid      = '#'
	tileEmpty      = '.'
	tileSpawn      = 'P'
	tileExit       = 'E'
	tileFragment   = 'Y'
	tileSpike      = '^'
	tileWatcher    = 'W'
	tileCheckpoint = 'C'
	tileSignal     = 'T'
	tileSecret     = 'S'

	transmissionsSection = "[transmissions]"
	textSection          = world.World(-2)
)

var (
	ErrFragmentOutsideGrey = errors.New("fragments are only allowed in the grey layer")
	ErrWatcherOutsideGrey  = errors.New("watchers are only allowed in the grey layer")
	ErrUnknownSection      = errors.New("unknown section")
	ErrDuplicateSection    = errors.New("duplicate section")
	ErrRowOutsideLayer     = errors.New("row outside of a section")
	ErrMissingSection      = errors.New("missing section")
	ErrRaggedRows          = errors.New("rows have different lengths")
	ErrSizeMismatch        = errors.New("layers have different sizes")
	ErrUnknownTile         = errors.New("unknown tile")
	ErrMarkerCount         = errors.New("level needs exactly one spawn and one exit")
	ErrFragmentInAuto      = errors.New("fragments are not allowed in auto shift mode")
	ErrTransmissionCount   = errors.New("each T marker needs exactly one transmission line")
	ErrSecretCount         = errors.New("level has more than one secret")
)

var sectionNames = map[string]world.World{
	"[yellow]": world.Yellow,
	"[grey]":   world.Grey,
}

func Parse(text string) (Level, error) {
	head, body, err := splitHeader(strings.Split(text, "\n"))
	if err != nil {
		return Level{}, err
	}
	rows, texts, err := splitSections(body)
	if err != nil {
		return Level{}, err
	}
	if err := checkSizes(rows); err != nil {
		return Level{}, err
	}

	lvl := Level{StartFragments: head.fragments, Shift: head.shift, AutoEvery: head.autoEvery}
	markers := map[rune][]Point{}
	var signals, secrets []Marker
	for w := range world.Count {
		layer := world.World(w)
		grid, layerMarkers, err := parseGrid(layer, rows[w])
		if err != nil {
			return Level{}, fmt.Errorf("%s layer: %w", layer, err)
		}
		lvl.Layers[w] = grid
		lvl.Spikes[w] = layerMarkers[tileSpike]
		signals = append(signals, markersIn(layer, layerMarkers[tileSignal])...)
		secrets = append(secrets, markersIn(layer, layerMarkers[tileSecret])...)
		for ch, points := range layerMarkers {
			markers[ch] = append(markers[ch], points...)
		}
	}

	if len(markers[tileSpawn]) != 1 || len(markers[tileExit]) != 1 {
		return Level{}, ErrMarkerCount
	}
	lvl.Spawn = markers[tileSpawn][0]
	lvl.Exit = markers[tileExit][0]
	lvl.Fragments = markers[tileFragment]
	lvl.Watchers = markers[tileWatcher]
	lvl.Checkpoints = markers[tileCheckpoint]
	if lvl.Shift == ShiftAuto && len(lvl.Fragments) > 0 {
		return Level{}, ErrFragmentInAuto
	}
	if len(secrets) > 1 {
		return Level{}, ErrSecretCount
	}
	if len(secrets) == 1 {
		lvl.Secret = &secrets[0]
	}
	lvl.Transmissions, err = pairTransmissions(signals, texts)
	if err != nil {
		return Level{}, err
	}
	return lvl, nil
}

func splitSections(lines []string) ([world.Count][]string, []string, error) {
	var rows [world.Count][]string
	var texts []string
	seen := map[world.World]bool{}
	current := world.World(-1)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if current == textSection {
			texts = append(texts, line)
			continue
		}
		if strings.HasPrefix(line, "[") {
			w, ok := sectionNames[line]
			if line == transmissionsSection {
				w, ok = textSection, true
			}
			if !ok {
				return rows, nil, fmt.Errorf("%w: %s", ErrUnknownSection, line)
			}
			if seen[w] {
				return rows, nil, fmt.Errorf("%w: %s", ErrDuplicateSection, line)
			}
			seen[w] = true
			current = w
			continue
		}
		if current < 0 {
			return rows, nil, ErrRowOutsideLayer
		}
		rows[current] = append(rows[current], line)
	}

	for w := range world.Count {
		if len(rows[w]) == 0 {
			return rows, nil, fmt.Errorf("%w: %s", ErrMissingSection, world.World(w))
		}
	}
	return rows, texts, nil
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

var greyOnly = map[rune]error{
	tileFragment: ErrFragmentOutsideGrey,
	tileWatcher:  ErrWatcherOutsideGrey,
}

func parseGrid(layer world.World, rows []string) (Grid, map[rune][]Point, error) {
	width, height := len(rows[0]), len(rows)
	solid := make([]bool, width*height)
	markers := map[rune][]Point{}
	for y, row := range rows {
		for x, ch := range row {
			switch ch {
			case tileSolid:
				solid[y*width+x] = true
			case tileEmpty:
			case tileSpawn, tileExit, tileSpike, tileFragment, tileWatcher, tileCheckpoint, tileSignal, tileSecret:
				if err, ok := greyOnly[ch]; ok && layer != world.Grey {
					return Grid{}, nil, fmt.Errorf("%w (at %d,%d)", err, x, y)
				}
				markers[ch] = append(markers[ch], Point{X: x, Y: y})
			default:
				return Grid{}, nil, fmt.Errorf("%w %q at %d,%d", ErrUnknownTile, ch, x, y)
			}
		}
	}
	return NewGrid(width, height, solid), markers, nil
}
