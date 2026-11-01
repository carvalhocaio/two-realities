package main

import (
	"log"

	"github.com/carvalhocaio/two-realities/internal/ebitenapp"
	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/level"
)

const firstLevel = "01_test"

func main() {
	lvl, err := level.Load(firstLevel)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebitenapp.Run(game.NewSession(lvl)); err != nil {
		log.Fatal(err)
	}
}
