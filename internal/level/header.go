package level

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	fragmentsKey = "fragments"
	shiftKey     = "shift"
	autoEveryKey = "auto_every"
	bossKey      = "boss"
	wardenBoss   = "warden"
)

var (
	ErrMissingHeader = errors.New("missing header line before layers")
	ErrInvalidHeader = errors.New("invalid header")
)

type ShiftMode int

const (
	ShiftFree ShiftMode = iota
	ShiftQuantized
	ShiftAuto
)

var shiftModes = map[string]ShiftMode{
	"free":      ShiftFree,
	"quantized": ShiftQuantized,
	"auto":      ShiftAuto,
}

func (m ShiftMode) String() string {
	for name, mode := range shiftModes {
		if mode == m {
			return name
		}
	}
	return "unknown"
}

type header struct {
	fragments int
	shift     ShiftMode
	autoEvery int
	boss      bool
}

func splitHeader(lines []string) (header, []string, error) {
	fields := map[string]string{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") {
			if len(fields) == 0 {
				return header{}, nil, ErrMissingHeader
			}
			h, err := parseHeader(fields)
			return h, lines[i:], err
		}
		key, value, ok := strings.Cut(line, ":")
		key = strings.TrimSpace(key)
		if !ok {
			return header{}, nil, fmt.Errorf("%w: %q", ErrInvalidHeader, line)
		}
		if _, dup := fields[key]; dup {
			return header{}, nil, fmt.Errorf("%w: duplicate %s", ErrInvalidHeader, key)
		}
		fields[key] = strings.TrimSpace(value)
	}
	return header{}, nil, ErrMissingHeader
}

func parseHeader(fields map[string]string) (header, error) {
	var h header
	for key, value := range fields {
		var err error
		switch key {
		case fragmentsKey:
			h.fragments, err = nonNegative(key, value)
		case shiftKey:
			mode, ok := shiftModes[value]
			if !ok {
				err = fmt.Errorf("%w: %s must be free, quantized or auto", ErrInvalidHeader, key)
			}
			h.shift = mode
		case autoEveryKey:
			h.autoEvery, err = nonNegative(key, value)
			if err == nil && h.autoEvery == 0 {
				err = fmt.Errorf("%w: %s must be positive", ErrInvalidHeader, key)
			}
		case bossKey:
			h.boss = value == wardenBoss
			if !h.boss {
				err = fmt.Errorf("%w: unknown boss %q", ErrInvalidHeader, value)
			}
		default:
			err = fmt.Errorf("%w: unknown key %q", ErrInvalidHeader, key)
		}
		if err != nil {
			return header{}, err
		}
	}
	if _, ok := fields[fragmentsKey]; !ok {
		return header{}, fmt.Errorf("%w: %s is required", ErrInvalidHeader, fragmentsKey)
	}
	_, hasShift := fields[shiftKey]
	if h.boss && hasShift {
		return header{}, ErrBossShiftOverride
	}
	if (h.shift == ShiftAuto) != (h.autoEvery > 0) {
		return header{}, fmt.Errorf("%w: %s is required in auto mode and only there", ErrInvalidHeader, autoEveryKey)
	}
	return h, nil
}

func nonNegative(key, value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("%w: %s must be a non-negative integer", ErrInvalidHeader, key)
	}
	return n, nil
}
