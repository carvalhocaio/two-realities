package geom

type Vec struct {
	X, Y float64
}

type Rect struct {
	X, Y, W, H float64
}

func (r Rect) Right() float64  { return r.X + r.W }
func (r Rect) Bottom() float64 { return r.Y + r.H }

func (r Rect) Overlaps(o Rect) bool {
	return r.X < o.Right() && o.X < r.Right() && r.Y < o.Bottom() && o.Y < r.Bottom()
}
