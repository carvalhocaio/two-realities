package main

import (
	"log"

	"github.com/carvalhocaio/two-realities/internal/ebitenapp"
	"github.com/carvalhocaio/two-realities/internal/level"
	"github.com/carvalhocaio/two-realities/internal/save"
)

func main() {
	levels, err := level.LoadAll()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebitenapp.Run(levels, save.Default()); err != nil {
		log.Fatal(err)
	}
}
