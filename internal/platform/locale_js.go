//go:build js

package platform

import "syscall/js"

func Locale() (locale string) {
	defer func() { _ = recover() }()
	nav := js.Global().Get("navigator")
	if !nav.Truthy() {
		return ""
	}
	return nav.Get("language").String()
}
