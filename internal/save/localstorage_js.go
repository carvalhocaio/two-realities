//go:build js

package save

import (
	"fmt"
	"syscall/js"
)

const storageKey = "two-realities/save"

type LocalStorage struct {
	Key string
}

func (l LocalStorage) Load() (d Data, err error) {
	defer recoverJS(&err)
	item := js.Global().Get("localStorage").Call("getItem", l.Key)
	if item.IsNull() || item.IsUndefined() {
		return decode(nil)
	}
	return decode([]byte(item.String()))
}

func (l LocalStorage) Save(d Data) (err error) {
	defer recoverJS(&err)
	raw, err := encode(d)
	if err != nil {
		return err
	}
	js.Global().Get("localStorage").Call("setItem", l.Key, string(raw))
	return nil
}

func recoverJS(err *error) {
	if r := recover(); r != nil {
		*err = fmt.Errorf("localStorage: %v", r)
	}
}

func Default() Storage {
	return LocalStorage{Key: storageKey}
}
