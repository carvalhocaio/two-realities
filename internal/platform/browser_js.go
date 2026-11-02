//go:build js

package platform

import "syscall/js"

const screenSlack = 2

var listeners []js.Func

func Listen(s *Signals) {
	defer func() { _ = recover() }()
	win := js.Global()
	doc := win.Get("document")

	on(win, "blur", func() { s.Notify(Blur) })
	on(doc, "visibilitychange", func() {
		if doc.Get("visibilityState").String() == "hidden" {
			s.Notify(Hidden)
		}
	})
	track := func() { s.SetFullscreen(isFullscreen(win, doc)) }
	on(doc, "fullscreenchange", track)
	on(win, "resize", track)
	track()
}

func on(target js.Value, event string, handle func()) {
	f := js.FuncOf(func(js.Value, []js.Value) any {
		defer func() { _ = recover() }()
		handle()
		return nil
	})
	listeners = append(listeners, f)
	target.Call("addEventListener", event, f)
}

func isFullscreen(win, doc js.Value) bool {
	if el := doc.Get("fullscreenElement"); el.Truthy() {
		return true
	}
	screen := win.Get("screen")
	return win.Get("innerWidth").Int() >= screen.Get("width").Int()-screenSlack &&
		win.Get("innerHeight").Int() >= screen.Get("height").Int()-screenSlack
}
