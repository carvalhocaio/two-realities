package player

type jumpWindow struct {
	coyote, buffer int
}

func (w *jumpWindow) tick(grounded, pressed bool) bool {
	w.coyote = refresh(w.coyote, grounded, CoyoteTicks)
	w.buffer = refresh(w.buffer, pressed, BufferTicks)
	if w.coyote > 0 && w.buffer > 0 {
		*w = jumpWindow{}
		return true
	}
	return false
}

func refresh(timer int, active bool, full int) int {
	if active {
		return full
	}
	return max(timer-1, 0)
}
