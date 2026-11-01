package camera

import "github.com/carvalhocaio/two-realities/internal/geom"

const (
	DeadZoneW = 48
	DeadZoneH = 32
	Smoothing = 0.15
)

type Camera struct {
	X, Y         float64
	ViewW, ViewH float64
}

func New(viewW, viewH float64) Camera {
	return Camera{ViewW: viewW, ViewH: viewH}
}

func (c *Camera) Snap(target, bounds geom.Rect) {
	cx, cy := center(target)
	c.X = cx - c.ViewW/2
	c.Y = cy - c.ViewH/2
	c.clamp(bounds)
}

func (c *Camera) Follow(target, bounds geom.Rect) {
	cx, cy := center(target)
	wantX := chase(c.X, c.ViewW, DeadZoneW, cx)
	wantY := chase(c.Y, c.ViewH, DeadZoneH, cy)
	c.X += (wantX - c.X) * Smoothing
	c.Y += (wantY - c.Y) * Smoothing
	c.clamp(bounds)
}

func chase(pos, view, zone, target float64) float64 {
	low := pos + (view-zone)/2
	high := low + zone
	switch {
	case target < low:
		return pos - (low - target)
	case target > high:
		return pos + (target - high)
	}
	return pos
}

func (c *Camera) clamp(bounds geom.Rect) {
	c.X = clampAxis(c.X, c.ViewW, bounds.X, bounds.W)
	c.Y = clampAxis(c.Y, c.ViewH, bounds.Y, bounds.H)
}

func clampAxis(pos, view, start, size float64) float64 {
	if size <= view {
		return start - (view-size)/2
	}
	return min(max(pos, start), start+size-view)
}

func center(r geom.Rect) (float64, float64) {
	return r.X + r.W/2, r.Y + r.H/2
}
