package level

const TileSize = 10

type Point struct {
	X, Y int
}

type Grid struct {
	width, height int
	solid         []bool
}

func NewGrid(width, height int, solid []bool) Grid {
	return Grid{width: width, height: height, solid: solid}
}

func (g Grid) Width() int  { return g.width }
func (g Grid) Height() int { return g.height }

func (g Grid) InBounds(x, y int) bool {
	return x >= 0 && x < g.width && y >= 0 && y < g.height
}

func (g Grid) Solid(x, y int) bool {
	return g.InBounds(x, y) && g.solid[y*g.width+x]
}

func (g Grid) PixelHeight() float64 { return float64(g.height * TileSize) }
