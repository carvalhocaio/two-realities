package ebitenapp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/carvalhocaio/two-realities/internal/input"
)

var padButtons = [input.PadButtonCount]ebiten.StandardGamepadButton{
	input.PadA:     ebiten.StandardGamepadButtonRightBottom,
	input.PadB:     ebiten.StandardGamepadButtonRightRight,
	input.PadX:     ebiten.StandardGamepadButtonRightLeft,
	input.PadY:     ebiten.StandardGamepadButtonRightTop,
	input.PadLB:    ebiten.StandardGamepadButtonFrontTopLeft,
	input.PadRB:    ebiten.StandardGamepadButtonFrontTopRight,
	input.PadBack:  ebiten.StandardGamepadButtonCenterLeft,
	input.PadStart: ebiten.StandardGamepadButtonCenterRight,
	input.PadUp:    ebiten.StandardGamepadButtonLeftTop,
	input.PadDown:  ebiten.StandardGamepadButtonLeftBottom,
	input.PadLeft:  ebiten.StandardGamepadButtonLeftLeft,
	input.PadRight: ebiten.StandardGamepadButtonLeftRight,
}

type Controls struct {
	keyboard Keyboard
	ids      []ebiten.GamepadID
	prev     input.PadState
	device   input.Device
}

var _ input.Source = (*Controls)(nil)

func NewControls() *Controls { return &Controls{keyboard: NewKeyboard()} }

func (c *Controls) Read() input.Intent {
	cur := c.readPads()
	kb, pad := c.keyboard.Read(), input.PadIntent(c.prev, cur)
	c.prev, c.device = cur, input.LastDevice(c.device, kb, pad)
	return input.Merge(kb, pad)
}

func (c *Controls) PadConnected() bool { return c.prev.Connected }
func (c *Controls) UsingPad() bool     { return c.device == input.DevicePad }

func (c *Controls) readPads() input.PadState {
	var all input.PadState
	c.ids = ebiten.AppendGamepadIDs(c.ids[:0])
	for _, id := range c.ids {
		if !ebiten.IsStandardGamepadLayoutAvailable(id) {
			continue
		}
		p := input.PadState{
			Connected: true,
			StickX:    ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal),
			StickY:    ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical),
		}
		for b, sb := range padButtons {
			p.Buttons[b] = ebiten.IsStandardGamepadButtonPressed(id, sb)
		}
		all = all.Combine(p)
	}
	return all
}
