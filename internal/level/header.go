package level

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const fragmentsKey = "fragments"

var (
	ErrMissingHeader = errors.New("missing header line before layers")
	ErrInvalidHeader = errors.New("invalid header")
)

type header struct {
	fragments int
}

func splitHeader(lines []string) (header, []string, error) {
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") {
			return header{}, nil, ErrMissingHeader
		}
		h, err := parseHeader(line)
		return h, lines[i+1:], err
	}
	return header{}, nil, ErrMissingHeader
}

func parseHeader(line string) (header, error) {
	key, value, ok := strings.Cut(line, ":")
	if !ok || strings.TrimSpace(key) != fragmentsKey {
		return header{}, fmt.Errorf("%w: %q", ErrInvalidHeader, line)
	}
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || n < 0 {
		return header{}, fmt.Errorf("%w: %s must be a non-negative integer", ErrInvalidHeader, fragmentsKey)
	}
	return header{fragments: n}, nil
}
