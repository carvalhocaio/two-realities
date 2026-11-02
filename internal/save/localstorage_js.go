//go:build js

package save

import (
	"fmt"
	"syscall/js"
)

const (
	saveName  = "save"
	probeName = "probe"
	legacyKey = "two-realities/save"
)

type LocalStorage struct {
	Name string
}

func (l LocalStorage) Load() (d Data, err error) {
	defer recoverJS(&err)
	store, err := localStorage()
	if err != nil {
		return Data{}, err
	}
	store.Call("setItem", Key(probeName), "1")
	store.Call("removeItem", Key(probeName))
	item := store.Call("getItem", Key(l.Name))
	if !present(item) && l.Name == saveName {
		item = store.Call("getItem", legacyKey)
		if present(item) {
			store.Call("setItem", Key(l.Name), item)
			store.Call("removeItem", legacyKey)
		}
	}
	if !present(item) {
		return decode(nil)
	}
	return decode([]byte(item.String()))
}

func (l LocalStorage) Save(d Data) (err error) {
	defer recoverJS(&err)
	store, err := localStorage()
	if err != nil {
		return err
	}
	raw, err := encode(d)
	if err != nil {
		return err
	}
	store.Call("setItem", Key(l.Name), string(raw))
	return nil
}

func present(v js.Value) bool { return !v.IsNull() && !v.IsUndefined() }

func localStorage() (js.Value, error) {
	store := js.Global().Get("Reflect").Call("get", js.Global(), "localStorage")
	if !store.Truthy() {
		return js.Value{}, ErrUnavailable
	}
	return store, nil
}

func recoverJS(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("%w: %v", ErrUnavailable, r)
	}
}

func Default() Storage {
	return NewFallback(LocalStorage{Name: saveName})
}
