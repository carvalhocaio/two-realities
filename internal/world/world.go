package world

type World int

const (
	Yellow World = iota
	Grey
	Count = 2
)

func (w World) Other() World {
	if w == Yellow {
		return Grey
	}
	return Yellow
}

func (w World) String() string {
	if w == Yellow {
		return "yellow"
	}
	return "grey"
}
