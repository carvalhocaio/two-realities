package i18n

import (
	"fmt"
	"strings"
)

type Lang int

const (
	EN Lang = iota
	PT
	LangCount
)

var codes = [LangCount]string{EN: "en", PT: "pt"}

func (l Lang) Code() string { return codes[l] }

func (l Lang) Next() Lang { return (l + 1) % LangCount }

func Parse(code string) (Lang, bool) {
	for l, c := range codes {
		if c == code {
			return Lang(l), true
		}
	}
	return EN, false
}

func Detect(locale string) Lang {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if strings.HasPrefix(locale, "pt") {
		return PT
	}
	return EN
}

func (l Lang) T(k Key, args ...any) string {
	text, ok := tables[l][k]
	if !ok {
		text = tables[EN][k]
	}
	if len(args) == 0 {
		return text
	}
	return fmt.Sprintf(text, args...)
}
