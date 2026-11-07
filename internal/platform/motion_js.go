//go:build js

package platform

import "syscall/js"

func ReducedMotion() (reduced bool) {
	defer func() { _ = recover() }()
	win := js.Global()
	if !win.Get("matchMedia").Truthy() {
		return false
	}
	return win.Call("matchMedia", "(prefers-reduced-motion: reduce)").Get("matches").Bool()
}
