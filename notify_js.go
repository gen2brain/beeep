// +build js

package beeep

import (
	"github.com/gopherjs/gopherjs/js"
)

// Notify sends desktop notification.
//
// On Web, in Firefox it just works, in Chrome you must call it from some "user gesture" like `onclick`,
// and you must use TLS certificate, it doesn't work with plain http.
func Notify(title, message, appIcon string) (err error) {
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

	n := js.Global.Get("window").Get("Notification")

	if n.Get("permission").String() == "granted" {
		n.New(title, map[string]interface{}{
			"body": message,
			"icon": appIcon,
		})
	} else {
		n.Call("requestPermission", func(permission string) {
			if permission == "granted" {
				n.New(title, map[string]interface{}{
					"body": message,
					"icon": appIcon,
				})
			}
		})
	}

	return
}
