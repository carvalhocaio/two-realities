package level

import (
	"cmp"
	"errors"
	"slices"
)

const (
	tileLamp  = 'L'
	tileRail  = 'O'
	tileGate1 = '1'
	tileGate2 = '2'

	LampCount = 3
	GateCount = LampCount - 1
)

var (
	ErrBossMarkers       = errors.New("lamps, rail and gates need a boss level")
	ErrLampCount         = errors.New("boss level needs exactly 3 lamps")
	ErrRailShape         = errors.New("boss level needs one unbroken rail on a single row")
	ErrRailOutsideGrey   = errors.New("the rail is only allowed in the grey layer")
	ErrGateMissing       = errors.New("boss level needs tiles for every gate")
	ErrBossShiftOverride = errors.New("boss levels take their shift mode from the boss phases")
)

var gateTiles = [GateCount]rune{tileGate1, tileGate2}

type Rail struct {
	Row, From, To int
}

type Boss struct {
	Lamps []Point
	Rail  Rail
	Gates [GateCount][]Point
}

func parseBoss(enabled bool, markers map[rune][]Point) (*Boss, error) {
	if !enabled {
		for _, ch := range append([]rune{tileLamp, tileRail}, gateTiles[:]...) {
			if len(markers[ch]) > 0 {
				return nil, ErrBossMarkers
			}
		}
		return nil, nil
	}
	if len(markers[tileLamp]) != LampCount {
		return nil, ErrLampCount
	}
	rail, err := railOf(markers[tileRail])
	if err != nil {
		return nil, err
	}
	b := &Boss{Lamps: slices.Clone(markers[tileLamp]), Rail: rail}
	slices.SortFunc(b.Lamps, func(a, c Point) int {
		return cmp.Or(cmp.Compare(c.Y, a.Y), cmp.Compare(a.X, c.X))
	})
	for i, ch := range gateTiles {
		if len(markers[ch]) == 0 {
			return nil, ErrGateMissing
		}
		b.Gates[i] = markers[ch]
	}
	return b, nil
}

func railOf(points []Point) (Rail, error) {
	if len(points) == 0 {
		return Rail{}, ErrRailShape
	}
	r := Rail{Row: points[0].Y, From: points[0].X, To: points[0].X}
	for _, p := range points {
		if p.Y != r.Row {
			return Rail{}, ErrRailShape
		}
		r.From, r.To = min(r.From, p.X), max(r.To, p.X)
	}
	if r.To-r.From+1 != len(points) {
		return Rail{}, ErrRailShape
	}
	return r, nil
}
