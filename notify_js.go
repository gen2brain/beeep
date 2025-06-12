//go:build js

package beeep

import (
	"encoding/base64"
	"fmt"
	"syscall/js"
)

// Notify sends desktop notification.
// The icon can be string with a path/url to png file or png []byte data. Stock icon names can also be used where supported.
//
// On the Web it uses the Notification API, in Firefox it just works, in Chrome you must call it from some "user gesture"
// like `onclick`, and you must use TLS certificate, it doesn't work with plain http.
func Notify(title, message string, icon any) (err error) {
	defer func() {
		e := recover()

		if e == nil {
			return
		}

		if e, ok := e.(*js.Error); ok {
			err = e
		} else {
			panic(e)
		}
	}()

	var img string
	switch i := icon.(type) {
	case string:
		img = i
	case []byte:
		img = fmt.Sprintf("data:image/png;base64, %s", base64.StdEncoding.EncodeToString(i))
	default:
		return fmt.Errorf("unsupported argument: %T", icon)
	}

	n := js.Global().Get("Notification")

	opts := js.Global().Get("Object").Invoke()
	opts.Set("body", message)
	opts.Set("icon", img)

	if n.Get("permission").String() == "granted" {
		n.New(js.ValueOf(title), opts)
	} else {
		var f js.Func
		f = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if args[0].String() == "granted" {
				n.New(js.ValueOf(title), opts)
			}
			f.Release()
			return nil
		})

		n.Call("requestPermission", f)
	}

	return
}
