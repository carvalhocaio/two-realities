package ebitenapp

import "github.com/carvalhocaio/two-realities/internal/i18n"

var padHints = map[i18n.Key]i18n.Key{
	i18n.HudRestart:   i18n.HudRestartPad,
	i18n.HudShiftHint: i18n.HudShiftHintPad,
	i18n.TitlePrompt:  i18n.TitlePromptPad,
	i18n.HelpPrompt:   i18n.HelpPromptPad,
	i18n.EndPrompt:    i18n.EndPromptPad,
	i18n.HelpWorlds:   i18n.HelpWorldsPad,
	i18n.HelpRestart:  i18n.HelpRestartPad,
}

func hint(lang i18n.Lang, k i18n.Key, pad bool) string {
	if p, ok := padHints[k]; ok && pad {
		return lang.T(p)
	}
	return lang.T(k)
}

type padState struct {
	connected, using bool
}
