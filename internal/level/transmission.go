package level

import (
	"errors"
	"fmt"
	"strings"

	"github.com/carvalhocaio/two-realities/internal/font"
	"github.com/carvalhocaio/two-realities/internal/i18n"
	"github.com/carvalhocaio/two-realities/internal/world"
)

const MaxTransmissionLength = 110

var ErrTransmissionText = errors.New("invalid transmission text")

type Transmission struct {
	Marker
	Text [i18n.LangCount]string
}

func (t Transmission) TextIn(l i18n.Lang) string {
	if t.Text[l] != "" {
		return t.Text[l]
	}
	return t.Text[i18n.EN]
}

func markersIn(layer world.World, points []Point) []Marker {
	out := make([]Marker, len(points))
	for i, p := range points {
		out[i] = Marker{Point: p, Layer: layer}
	}
	return out
}

func pairTransmissions(signals []Marker, texts textBlocks) ([]Transmission, error) {
	if len(texts[i18n.EN]) != len(signals) {
		return nil, fmt.Errorf("%w: %d markers, %d lines", ErrTransmissionCount, len(signals), len(texts[i18n.EN]))
	}
	out := make([]Transmission, len(signals))
	for i, m := range signals {
		out[i].Marker = m
	}
	for lang, lines := range texts {
		if len(lines) != len(signals) {
			return nil, fmt.Errorf("%w: %s has %d lines for %d markers", ErrTransmissionCount, lang.Code(), len(lines), len(signals))
		}
		for i, line := range lines {
			text := strings.ToUpper(line)
			if err := checkText(text); err != nil {
				return nil, err
			}
			out[i].Text[lang] = text
		}
	}
	return out, nil
}

func checkText(text string) error {
	if len([]rune(text)) > MaxTransmissionLength {
		return fmt.Errorf("%w: longer than %d characters", ErrTransmissionText, MaxTransmissionLength)
	}
	if missing := font.Missing(text); len(missing) > 0 {
		return fmt.Errorf("%w: unsupported characters %q", ErrTransmissionText, string(missing))
	}
	return nil
}
