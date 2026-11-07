//go:build dev

package main

import (
	"flag"
	"log"

	"github.com/carvalhocaio/two-realities/internal/ebitenapp"
	"github.com/carvalhocaio/two-realities/internal/input"
	"github.com/carvalhocaio/two-realities/internal/level"
)

type step = ebitenapp.Step

func hold(ticks int, in input.Intent) step { return step{Ticks: ticks, In: in} }

var (
	idle  = input.Intent{}
	right = input.Intent{Right: true}
	jump  = input.Intent{Right: true, JumpPressed: true}
	shift = input.Intent{ShiftPressed: true}
	help  = input.Intent{HelpPressed: true}
)

var shots = []ebitenapp.Shot{
	{Name: "title", Level: ebitenapp.TitleShot, Every: 30, Steps: []step{hold(60, idle)}},
	{Name: "help", Level: ebitenapp.TitleShot, Every: 20, Steps: []step{hold(1, help), hold(20, idle)}},
	{Name: "first_shift", Level: 0, Every: 1, Steps: []step{
		hold(20, idle), hold(30, right), hold(50, jump), hold(1, shift), hold(1, right), hold(150, right), hold(45, idle),
	}},
	{Name: "crossfire", Level: 6, Every: 10, Steps: []step{hold(60, idle), hold(1, shift), hold(240, idle)}},
	{Name: "exit_yellow", Level: 1, Every: 20, Steps: []step{
		hold(20, idle), hold(42, right), hold(1, shift), hold(125, right), hold(1, shift), hold(95, right), hold(60, idle),
	}},
	{Name: "warden", Level: 12, Every: 6, Steps: []step{hold(30, idle), hold(1, shift), hold(400, idle)}},
}

func main() {
	dir := flag.String("out", "shots", "output directory")
	locale := flag.String("lang", "en", "locale")
	flag.Parse()
	levels, err := level.LoadAll()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebitenapp.RunShots(levels, shots, *dir, *locale); err != nil {
		log.Fatal(err)
	}
}
