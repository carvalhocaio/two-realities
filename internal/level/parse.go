package level

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/carvalhocaio/two-realities/internal/i18n"
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

	noLayer = world.World(-1)

	MaxSecretSpots = 3
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
	ErrSecretCount         = errors.New("level has too many secret spots")
)

var sectionNames = map[string]world.World{
	"[yellow]": world.Yellow,
	"[grey]":   world.Grey,
}

var textSections = map[string]i18n.Lang{
	"[transmissions]":    i18n.EN,
	"[transmissions.pt]": i18n.PT,
}

type textBlocks map[i18n.Lang][]string

var sectionHeader = regexp.MustCompile(`^\[[a-z][a-z.]*\]$`)

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

	lvl := Level{StartFragments: head.fragments, Shift: head.shift, AutoEvery: head.autoEvery, Act: head.act}
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
	if lvl.Boss, err = parseBoss(head.boss, markers); err != nil {
		return Level{}, err
	}
	if len(secrets) > MaxSecretSpots {
		return Level{}, ErrSecretCount
	}
	lvl.Secrets = secrets
	lvl.Transmissions, err = pairTransmissions(signals, texts)
	if err != nil {
		return Level{}, err
	}
	return lvl, nil
}

func splitSections(lines []string) ([world.Count][]string, textBlocks, error) {
	var rows [world.Count][]string
	texts := textBlocks{}
	seen := map[string]bool{}
	current := noLayer
	inText, textLang := false, i18n.EN

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lang, isText := textSections[line]
		_, isLayer := sectionNames[line]
		if isText || isLayer || sectionHeader.MatchString(line) || (!inText && strings.HasPrefix(line, "[")) {
			if !isText && !isLayer {
				return rows, nil, fmt.Errorf("%w: %s", ErrUnknownSection, line)
			}
			if seen[line] {
				return rows, nil, fmt.Errorf("%w: %s", ErrDuplicateSection, line)
			}
			seen[line] = true
			inText, textLang, current = isText, lang, sectionNames[line]
			if isText {
				texts[lang] = []string{}
			}
			continue
		}
		if inText {
			texts[textLang] = append(texts[textLang], line)
			continue
		}
		if current == noLayer {
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
	tileRail:     ErrRailOutsideGrey,
}

func isMarker(ch rune) bool {
	switch ch {
	case tileSpawn, tileExit, tileSpike, tileFragment, tileWatcher, tileCheckpoint, tileSignal, tileSecret,
		tileLamp, tileRail, tileGate1, tileGate2:
		return true
	}
	return false
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
			default:
				if !isMarker(ch) {
					return Grid{}, nil, fmt.Errorf("%w %q at %d,%d", ErrUnknownTile, ch, x, y)
				}
				if err, ok := greyOnly[ch]; ok && layer != world.Grey {
					return Grid{}, nil, fmt.Errorf("%w (at %d,%d)", err, x, y)
				}
				markers[ch] = append(markers[ch], Point{X: x, Y: y})
			}
		}
	}
	return NewGrid(width, height, solid), markers, nil
}
