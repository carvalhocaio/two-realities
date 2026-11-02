package i18n

type Key int

const (
	Title Key = iota
	TitlePrompt
	TitleControlsMove
	TitleControlsWorld
	TitleControlsExtra
	TitleLanguage
	StorageNotice
	WorldYellow
	WorldGrey
	HudLevel
	HudRestart
	HudMuted
	ResultTime
	ResultTimeBest
	ResultNewBest
	Paused
	MenuResume
	MenuAssist
	MenuMute
	MenuAudioOffset
	MenuLanguage
	MenuTitle
	On
	Off
	OffsetHint
	EndTitle
	EndThanks
	EndRun
	EndAssist
	EndPractice
	EndNewBest
	EndBestRun
	EndNoBest
	EndSecrets
	EndPrompt
	KeyCount
)

var tables = [LangCount]map[Key]string{
	EN: {
		Title:              "TWO REALITIES",
		TitlePrompt:        "PRESS SHIFT TO START",
		TitleControlsMove:  "ARROWS OR A D: MOVE    SPACE: JUMP",
		TitleControlsWorld: "SHIFT: SWAP WORLDS    R: RESTART",
		TitleControlsExtra: "M: MUTE    T: TIMER    P OR ESC: PAUSE",
		TitleLanguage:      "L: PORTUGUÊS",
		StorageNotice:      "PROGRESS WON'T BE SAVED IN THIS BROWSER",
		WorldYellow:        "YELLOW",
		WorldGrey:          "GREY",
		HudLevel:           "LV %d/%d",
		HudRestart:         "R TO RESTART",
		HudMuted:           "MUTE",
		ResultTime:         "TIME %s",
		ResultTimeBest:     "TIME %s  BEST %s",
		ResultNewBest:      "NEW BEST %s",
		Paused:             "PAUSED",
		MenuResume:         "RESUME",
		MenuAssist:         "ASSIST MODE: %s",
		MenuMute:           "MUTE: %s",
		MenuAudioOffset:    "AUDIO OFFSET: %s",
		MenuLanguage:       "LANGUAGE: ENGLISH",
		MenuTitle:          "BACK TO TITLE",
		On:                 "ON",
		Off:                "OFF",
		OffsetHint:         "LEFT/RIGHT OR -/= TO ADJUST, MATCH THE DOT TO THE KICK",
		EndTitle:           "BOTH WORLDS CROSSED",
		EndThanks:          "THANKS FOR PLAYING",
		EndRun:             "RUN %s",
		EndAssist:          "ASSIST",
		EndPractice:        "PRACTICE",
		EndNewBest:         "NEW BEST!",
		EndBestRun:         "BEST RUN %s",
		EndNoBest:          "BEST RUN -",
		EndSecrets:         "SECRETS %d/%d",
		EndPrompt:          "PRESS SHIFT FOR TITLE",
	},
	PT: {
		Title:              "TWO REALITIES",
		TitlePrompt:        "APERTE SHIFT PARA COMEÇAR",
		TitleControlsMove:  "SETAS OU A D: ANDAR    ESPAÇO: PULAR",
		TitleControlsWorld: "SHIFT: TROCAR DE MUNDO    R: REINICIAR",
		TitleControlsExtra: "M: SOM    T: TEMPO    P OU ESC: PAUSA",
		TitleLanguage:      "L: ENGLISH",
		StorageNotice:      "SEU PROGRESSO NÃO SERÁ SALVO NESTE NAVEGADOR",
		WorldYellow:        "AMARELO",
		WorldGrey:          "CINZA",
		HudLevel:           "FASE %d/%d",
		HudRestart:         "R PARA REINICIAR",
		HudMuted:           "MUDO",
		ResultTime:         "TEMPO %s",
		ResultTimeBest:     "TEMPO %s  RECORDE %s",
		ResultNewBest:      "NOVO RECORDE %s",
		Paused:             "PAUSA",
		MenuResume:         "CONTINUAR",
		MenuAssist:         "MODO ASSISTIDO: %s",
		MenuMute:           "MUDO: %s",
		MenuAudioOffset:    "ATRASO DO ÁUDIO: %s",
		MenuLanguage:       "IDIOMA: PORTUGUÊS",
		MenuTitle:          "VOLTAR AO INÍCIO",
		On:                 "SIM",
		Off:                "NÃO",
		OffsetHint:         "ESQUERDA/DIREITA OU -/= PARA AJUSTAR O PONTO À BATIDA",
		EndTitle:           "OS DOIS MUNDOS FICARAM PARA TRÁS",
		EndThanks:          "OBRIGADO POR JOGAR",
		EndRun:             "PARTIDA %s",
		EndAssist:          "ASSISTIDO",
		EndPractice:        "TREINO",
		EndNewBest:         "NOVO RECORDE!",
		EndBestRun:         "RECORDE %s",
		EndNoBest:          "RECORDE -",
		EndSecrets:         "SEGREDOS %d/%d",
		EndPrompt:          "APERTE SHIFT PARA VOLTAR",
	},
}
