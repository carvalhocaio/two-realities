package input

import (
	"reflect"
	"testing"
)

func TestLastDevice(t *testing.T) {
	press := Intent{JumpPressed: true}
	tests := []struct {
		name          string
		prev          Device
		keyboard, pad Intent
		want          Device
	}{
		{"nothing keeps the keyboard", DeviceKeyboard, Intent{}, Intent{}, DeviceKeyboard},
		{"nothing keeps the pad", DevicePad, Intent{}, Intent{}, DevicePad},
		{"a pad press switches to the pad", DeviceKeyboard, Intent{}, press, DevicePad},
		{"a key press switches back", DevicePad, press, Intent{}, DeviceKeyboard},
		{"holding a direction counts", DeviceKeyboard, Intent{}, Intent{Right: true}, DevicePad},
		{"both at once keep the last one", DevicePad, press, press, DevicePad},
		{"focus loss is not a device", DeviceKeyboard, Intent{}, Intent{FocusLost: true}, DeviceKeyboard},
	}
	for _, tt := range tests {
		if got := LastDevice(tt.prev, tt.keyboard, tt.pad); got != tt.want {
			t.Errorf("%s: %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAnyCoversEveryPress(t *testing.T) {
	v := reflect.ValueOf(Intent{})
	for i := range v.NumField() {
		name := v.Type().Field(i).Name
		var in Intent
		reflect.ValueOf(&in).Elem().Field(i).SetBool(true)
		if in.Any() != (name != "FocusLost") {
			t.Errorf("Any() with only %s = %v", name, in.Any())
		}
	}
}
