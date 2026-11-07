package input

type Device int

const (
	DeviceKeyboard Device = iota
	DevicePad
)

func (i Intent) Any() bool {
	return i.Acts() || i.RestartPressed || i.PausePressed || i.UpPressed || i.DownPressed ||
		i.ConfirmPressed || i.MutePressed || i.DecreasePressed || i.IncreasePressed ||
		i.LanguagePressed || i.HelpPressed
}

func LastDevice(prev Device, keyboard, pad Intent) Device {
	switch kb, p := keyboard.Any(), pad.Any(); {
	case p && !kb:
		return DevicePad
	case kb && !p:
		return DeviceKeyboard
	}
	return prev
}
