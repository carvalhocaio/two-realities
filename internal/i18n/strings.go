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
	TitleControlsPad
	WorldYellow
	WorldGrey
	HudLevel
	HudRestart
	HudSecret
	HudShiftHint
	HudMuted
	HudWarden
	ResultTime
	ResultTimeBest
	ResultNewBest
	Paused
	MenuResume
	MenuAssist
	MenuMute
	MenuAudioOffset
	MenuLanguage
	MenuSkip
	MenuHelp
	HelpTitle
	HelpWorlds
	HelpFragments
	HelpSpikes
	HelpWatchers
	HelpBeat
	HelpRestart
	HelpPrompt
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
	TitlePressure
	TitlePressureOff
	TitleStartLevel
	PressureFasterWatchers
	PressureNoCheckpoints
	PressureBeatGhost
	PressureRestartLevel
	PressureRestartAct
	HudPressure
	EndPressure
	EndPressureCleared
	KeyCount
)

var tables = [LangCount]map[Key]string{
	EN: {
		Title:                  "TWO REALITIES",
		TitlePrompt:            "PRESS SHIFT TO START",
		TitleControlsMove:      "ARROWS OR A D: MOVE    SPACE: JUMP",
		TitleControlsWorld:     "SHIFT: SWAP WORLDS    R: RESTART",
		TitleControlsExtra:     "M: MUTE    T: TIMER    P OR ESC: PAUSE",
		TitleLanguage:          "L: PORTUGUÊS    H: HOW TO PLAY",
		StorageNotice:          "PROGRESS WON'T BE SAVED IN THIS BROWSER",
		TitleControlsPad:       "GAMEPAD: A JUMP  B/X/LB/RB SHIFT  START PAUSE  Y HELP",
		WorldYellow:            "YELLOW",
		WorldGrey:              "GREY",
		HudLevel:               "LV %d/%d",
		HudRestart:             "R TO RESTART",
		HudSecret:              "SECRET FOUND!",
		HudShiftHint:           "SHIFT: SWAP WORLDS",
		HudMuted:               "MUTE",
		HudWarden:              "THE WARDEN",
		ResultTime:             "TIME %s",
		ResultTimeBest:         "TIME %s  BEST %s",
		ResultNewBest:          "NEW BEST %s",
		Paused:                 "PAUSED",
		MenuResume:             "RESUME",
		MenuAssist:             "ASSIST MODE: %s",
		MenuMute:               "MUTE: %s",
		MenuAudioOffset:        "AUDIO OFFSET: %s",
		MenuLanguage:           "LANGUAGE: ENGLISH",
		MenuSkip:               "SKIP LEVEL",
		MenuHelp:               "HOW TO PLAY",
		HelpTitle:              "HOW TO PLAY",
		HelpWorlds:             "SHIFT SWAPS WORLDS. DOTTED TILES ARE THE OTHER ONE: NOT SOLID.",
		HelpFragments:          "GOING GREY IS FREE. COMING BACK COSTS 1 FRAGMENT, FOUND IN GREY.",
		HelpSpikes:             "SPIKES ONLY HURT IN THEIR OWN WORLD.",
		HelpWatchers:           "WATCHERS SEE YOU ONLY IN GREY. SHIFT BEFORE THE ALERT FILLS.",
		HelpBeat:               "LATER, SHIFT WAITS FOR THE BEAT, THEN SWAPS ON ITS OWN.",
		HelpRestart:            "OUT OF FRAGMENTS? R RESTARTS THE LEVEL.",
		HelpPrompt:             "PRESS SHIFT TO CONTINUE",
		MenuTitle:              "BACK TO TITLE",
		On:                     "ON",
		Off:                    "OFF",
		OffsetHint:             "LEFT/RIGHT OR -/= TO ADJUST, MATCH THE DOT TO THE KICK",
		EndTitle:               "BOTH WORLDS CROSSED",
		EndThanks:              "THANKS FOR PLAYING",
		EndRun:                 "RUN %s",
		EndAssist:              "ASSIST",
		EndPractice:            "PRACTICE",
		EndNewBest:             "NEW BEST!",
		EndBestRun:             "BEST RUN %s",
		EndNoBest:              "BEST RUN -",
		EndSecrets:             "SECRETS %d/%d",
		EndPrompt:              "PRESS SHIFT FOR TITLE",
		TitlePressure:          "PRESSURE %d - LEFT/RIGHT",
		TitlePressureOff:       "PRESSURE OFF - LEFT/RIGHT",
		TitleStartLevel:        "START AT LEVEL %d/%d - UP/DOWN",
		PressureFasterWatchers: "+ FASTER WATCHERS",
		PressureNoCheckpoints:  "+ NO CHECKPOINTS",
		PressureBeatGhost:      "+ THE OTHER WORLD ONLY SHOWS ON THE BEAT",
		PressureRestartLevel:   "+ DEATH RESTARTS THE LEVEL",
		PressureRestartAct:     "+ DEATH RESTARTS THE ACT",
		HudPressure:            "P%d",
		EndPressure:            "PRESSURE %d",
		EndPressureCleared:     "PRESSURE %d CLEARED!",
	},
	PT: {
		Title:                  "TWO REALITIES",
		TitlePrompt:            "APERTE SHIFT PARA COMEÇAR",
		TitleControlsMove:      "SETAS OU A D: ANDAR    ESPAÇO: PULAR",
		TitleControlsWorld:     "SHIFT: TROCAR DE MUNDO    R: REINICIAR",
		TitleControlsExtra:     "M: SOM    T: TEMPO    P OU ESC: PAUSA",
		TitleLanguage:          "L: ENGLISH    H: COMO JOGAR",
		StorageNotice:          "SEU PROGRESSO NÃO SERÁ SALVO NESTE NAVEGADOR",
		TitleControlsPad:       "CONTROLE: A PULA  B/X/LB/RB TROCA  START PAUSA  Y AJUDA",
		WorldYellow:            "AMARELO",
		WorldGrey:              "CINZA",
		HudLevel:               "FASE %d/%d",
		HudRestart:             "R PARA REINICIAR",
		HudSecret:              "SEGREDO ENCONTRADO!",
		HudShiftHint:           "SHIFT: TROCAR DE MUNDO",
		HudMuted:               "MUDO",
		HudWarden:              "O CARCEREIRO",
		ResultTime:             "TEMPO %s",
		ResultTimeBest:         "TEMPO %s  RECORDE %s",
		ResultNewBest:          "NOVO RECORDE %s",
		Paused:                 "PAUSA",
		MenuResume:             "CONTINUAR",
		MenuAssist:             "MODO ASSISTIDO: %s",
		MenuMute:               "MUDO: %s",
		MenuAudioOffset:        "ATRASO DO ÁUDIO: %s",
		MenuLanguage:           "IDIOMA: PORTUGUÊS",
		MenuSkip:               "PULAR FASE",
		MenuHelp:               "COMO JOGAR",
		HelpTitle:              "COMO JOGAR",
		HelpWorlds:             "SHIFT TROCA DE MUNDO. O PONTILHADO É O OUTRO: NÃO É SÓLIDO.",
		HelpFragments:          "IR AO CINZA É GRÁTIS. VOLTAR CUSTA 1 FRAGMENTO, ACHADO NO CINZA.",
		HelpSpikes:             "ESPINHOS SÓ FEREM NO PRÓPRIO MUNDO.",
		HelpWatchers:           "VIGIAS SÓ TE VEEM NO CINZA. TROQUE ANTES DO ALERTA ENCHER.",
		HelpBeat:               "DEPOIS, A TROCA ESPERA A BATIDA, E ENTÃO TROCA SOZINHA.",
		HelpRestart:            "SEM FRAGMENTOS? R REINICIA A FASE.",
		HelpPrompt:             "APERTE SHIFT PARA CONTINUAR",
		MenuTitle:              "VOLTAR AO INÍCIO",
		On:                     "SIM",
		Off:                    "NÃO",
		OffsetHint:             "ESQUERDA/DIREITA OU -/= PARA AJUSTAR O PONTO À BATIDA",
		EndTitle:               "OS DOIS MUNDOS FICARAM PARA TRÁS",
		EndThanks:              "OBRIGADO POR JOGAR",
		EndRun:                 "PARTIDA %s",
		EndAssist:              "ASSISTIDO",
		EndPractice:            "TREINO",
		EndNewBest:             "NOVO RECORDE!",
		EndBestRun:             "RECORDE %s",
		EndNoBest:              "RECORDE -",
		EndSecrets:             "SEGREDOS %d/%d",
		EndPrompt:              "APERTE SHIFT PARA VOLTAR",
		TitlePressure:          "PRESSÃO %d - ESQUERDA/DIREITA",
		TitlePressureOff:       "PRESSÃO DESLIGADA - ESQUERDA/DIREITA",
		TitleStartLevel:        "COMEÇAR NA FASE %d/%d - CIMA/BAIXO",
		PressureFasterWatchers: "+ VIGIAS MAIS RÁPIDOS",
		PressureNoCheckpoints:  "+ SEM CHECKPOINTS",
		PressureBeatGhost:      "+ O OUTRO MUNDO SÓ APARECE NA BATIDA",
		PressureRestartLevel:   "+ MORRER REINICIA A FASE",
		PressureRestartAct:     "+ MORRER VOLTA AO INÍCIO DO ATO",
		HudPressure:            "P%d",
		EndPressure:            "PRESSÃO %d",
		EndPressureCleared:     "PRESSÃO %d VENCIDA!",
	},
}
