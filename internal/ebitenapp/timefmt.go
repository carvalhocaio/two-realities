package ebitenapp

import "fmt"

const ticksPerSecond = 60

func formatTicks(ticks int) string {
	centis := ticks * 100 / ticksPerSecond
	return fmt.Sprintf("%d:%02d.%02d", centis/6000, centis/100%60, centis%100)
}
