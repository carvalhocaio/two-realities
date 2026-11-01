package main

import (
	"log"

	"github.com/carvalhocaio/two-realities/internal/ebitenapp"
	"github.com/carvalhocaio/two-realities/internal/game"
	"github.com/carvalhocaio/two-realities/internal/level"
)

func main() {
	levels, err := level.LoadAll()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebitenapp.Run(game.NewFlow(levels...)); err != nil {
		log.Fatal(err)
	}
}
