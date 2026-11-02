package level

import (
	"errors"
	"fmt"
	"strings"

	"github.com/carvalhocaio/two-realities/internal/world"
)

const (
	MaxTransmissionLength = 110
	TextCharset           = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 .,-:/!?'"
)

var ErrTransmissionText = errors.New("invalid transmission text")

func markersIn(layer world.World, points []Point) []Marker {
	out := make([]Marker, len(points))
	for i, p := range points {
		out[i] = Marker{Point: p, Layer: layer}
	}
	return out
}

func pairTransmissions(signals []Marker, texts []string) ([]Transmission, error) {
	if len(signals) != len(texts) {
		return nil, fmt.Errorf("%w: %d markers, %d lines", ErrTransmissionCount, len(signals), len(texts))
	}
	out := make([]Transmission, len(signals))
	for i, m := range signals {
		text := strings.ToUpper(texts[i])
		if err := checkText(text); err != nil {
			return nil, err
		}
		out[i] = Transmission{Marker: m, Text: text}
	}
	return out, nil
}

func checkText(text string) error {
	if len(text) > MaxTransmissionLength {
		return fmt.Errorf("%w: longer than %d characters", ErrTransmissionText, MaxTransmissionLength)
	}
	for _, ch := range text {
		if !strings.ContainsRune(TextCharset, ch) {
			return fmt.Errorf("%w: unsupported character %q", ErrTransmissionText, ch)
		}
	}
	return nil
}
